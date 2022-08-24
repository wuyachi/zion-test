package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"io/ioutil"
	"math/big"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/polynetwork/bridge-common/chains/eth"
	"github.com/polynetwork/bridge-common/log"
)

type Context struct {
	nodes *eth.SDK
	sync.RWMutex
	height   uint64
	checkUrl string
}

func (ctx *Context) Till(height uint64) {
	pass := false
	for {
		ctx.RLock()
		pass = ctx.height >= height
		ctx.RUnlock()
		if pass {
			return
		}
		time.Sleep(time.Millisecond * 300)
	}
}

type Case struct {
	index   int64
	err     error
	actions []Action
	plan    []*ActionItem
}

func (c *Case) Run(ctx *Context) (err error) {
	exit := make(chan struct{})
	defer func() {
		close(exit)
	}()

	go func() {
		for {
			select {
			case <-exit:
				return
			default:
				height, err := ctx.nodes.Node().GetLatestHeight()
				if err != nil {
					log.Error("Get chain latest height failed", "err", err)
				} else {
					ctx.Lock()
					ctx.height = height
					ctx.Unlock()
				}
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()

	// Sort actions
	bp := make(map[uint64][]Action)
	for i, a := range c.actions {
		a.SetIndex(i)
		bp[a.StartAt()] = append(bp[a.StartAt()], a)
	}
	c.plan = make([]*ActionItem, 0, len(bp))
	for b, actions := range bp {
		c.plan = append(c.plan, &ActionItem{b, actions})
	}
	sort.Slice(c.plan, func(i, j int) bool { return c.plan[i].start < c.plan[j].start })

	type result struct {
		err   error
		index int
	}
	res := make(chan result, len(c.actions))
	for i, item := range c.plan {
		log.Info("Scheduling plan", "index", i, "action_count", len(item.actions), "at", item.start)
		for _, action := range item.actions {
			go func(a Action) {
				ctx.Till(a.StartAt())
				log.Info("Running case action", "case_index", c.index, "action_index", a.Index())
				res <- result{a.Run(ctx), a.Index()}
			}(action)
		}
	}

	for j := 0; j < len(c.actions); j++ {
		log.Info("Waiting case actions result", "case_index", c.index, "progress", j+1, "total", len(c.actions))
		r := <-res
		c.actions[r.index].SetError(r.err)
		if r.err != nil {
			err = fmt.Errorf("action failure, err: %v, action_index: %v", r.err, r.index)
			log.Error("Run case action failed", "case", c.index, "action", r.index, "err", err)
		}
	}
	return
}

type ActionItem struct {
	start   uint64
	actions []Action
}

type Action interface {
	Run(*Context) error
	StartAt() uint64
	Before() uint64
	SetIndex(int)
	Index() int
	Error() error
	SetError(err error)
}

type ActionBase struct {
	Epoch        uint64
	Block        uint64
	ShouldBefore uint64
	index        int
	err          error
}

func (a *ActionBase) StartAt() uint64 { return a.Block + a.Epoch*uint64(CONFIG.BlocksPerEpoch) }
func (a *ActionBase) Before() uint64 {
	return a.ShouldBefore + a.Epoch*uint64(CONFIG.BlocksPerEpoch)
}
func (a *ActionBase) SetIndex(index int) { a.index = index }
func (a *ActionBase) Index() int         { return a.index }
func (a *ActionBase) SetError(err error) { a.err = err }
func (a *ActionBase) Error() error       { return a.err }

type SendTx struct {
	ActionBase
	Tx            *types.Transaction
	ShouldSucceed bool
}

func (a *SendTx) Run(ctx *Context) (err error) {
	err = ctx.nodes.Node().SendTransaction(context.Background(), a.Tx)
	if err != nil {
		return
	}
	currentHeight, _ := ctx.nodes.Node().GetLatestHeight()
	log.Info("Sent tx", "current_height", currentHeight, "hash", a.Tx.Hash(), "index", a.Index())
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 2)
		height, _, pending, err := ctx.nodes.Node().Confirm(a.Tx.Hash(), 1, 10)
		if err != nil {
			return err
		}

		if height > 0 {
			if height <= a.Before() {
				rec, err := ctx.nodes.Node().TransactionReceipt(context.Background(), a.Tx.Hash())
				if err != nil {
					return err
				}
				if (rec.Status == 1) == a.ShouldSucceed {
					return nil
				}
				return fmt.Errorf("transaction status error, status %v, wanted %v", rec.Status == 1, a.ShouldSucceed)
			}
			return fmt.Errorf("tx packed too late, height %v, expected before %v", height, a.Before())
		} else if !pending {
			return fmt.Errorf("possible tx lost %s", a.Tx.Hash())
		}
	}
	return nil
}

type Query struct {
	ActionBase
	Request    ethereum.CallMsg
	Assertions []Assertion
}

func (a *Query) Run(ctx *Context) (err error) {
	output, err := ctx.nodes.Node().CallContract(context.Background(), a.Request, big.NewInt(int64(a.StartAt())))
	if err != nil {
		return
	}
	err = Assert(output, a.Assertions)
	return
}

type CheckBalance struct {
	ActionBase
	Addresses []common.Address
}

func (a *CheckBalance) Run(ctx *Context) (err error) {
	expectedBalances, err := a.CheckBalances(ctx, a.Addresses)
	if err != nil {
		return
	}
	maxDelta := new(big.Int).Mul(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), big.NewInt(1))

	for i, address := range a.Addresses {
		balance, err := ctx.nodes.Node().BalanceAt(context.Background(), address, big.NewInt(int64(a.StartAt())))
		if err != nil {
			return err
		}
		delta := new(big.Int).Abs(new(big.Int).Sub(balance, expectedBalances[i]))
		if delta.Cmp(maxDelta) == 1 {
			return fmt.Errorf("balance check failure, balance %s, expected %s, delta %s", balance, expectedBalances[i], delta)
		}
	}
	return
}

type CheckBalanceReq struct {
	Id        string   `json:"Id"`
	Addresses []string `json:"Addresses"`
	EndHeight uint64   `json:"EndHeight"`
}

type CheckBalanceRsp struct {
	Action string `json:"action"`
	Desc   string `json:"desc"`
	Error  int    `json:"error"`
	Result struct {
		Amount []string `json:"Amount"`
		Id     string   `json:"Id"`
	} `json:"result"`
}

func (a *CheckBalance) CheckBalances(ctx *Context, addresses []common.Address) (balances []*big.Int, err error) {
	addressArr := make([]string, 0)
	for _, address := range addresses {
		addressArr = append(addressArr, address.String())
	}
	req := &CheckBalanceReq{Addresses: addressArr, EndHeight: a.StartAt()}
	rsp := &CheckBalanceRsp{}
	balances = make([]*big.Int, 0)
	err = PostJsonFor(ctx.checkUrl, req, rsp)
	if err != nil {
		return
	}
	for _, amount := range rsp.Result.Amount {
		balance := big.NewInt(int64(math.MustParseUint64(amount)))
		balances = append(balances, balance)
	}
	return
}

func PostJsonFor(url string, payload interface{}, result interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, result)
	if err != nil {
		log.Error("PostJson response", "Body", string(respBody))
	} else {
		log.Debug("PostJson response", "Body", string(respBody))
	}
	return err
}
