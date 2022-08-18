package main

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/polynetwork/bridge-common/chains/eth"
	"github.com/polynetwork/bridge-common/log"
)

type Context struct {
	nodes *eth.SDK
	sync.RWMutex
	height uint64
}

func (ctx *Context) Till(height uint64) {
	pass := false
	for {
		ctx.RLock()
		pass = ctx.height >= height
		ctx.Unlock()
		if pass {
			return
		}
		time.Sleep(time.Millisecond * 300)
	}
}

type Case struct {
	index   int
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
			case <- exit:
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
		 err error
		 index int
	}
	res := make(chan result, len(c.actions))
	for i, item := range c.plan {
		log.Info("Running plan", "index", i, "action_count", len(item.actions), "at", item.start)
		for _, action := range item.actions {
			go func(a Action) {
				ctx.Till(a.StartAt())
				res <- result{ a.Run(ctx), a.Index() }
			}(action)
		}
	}

	for j := 0; j < len(c.actions); j++ {
		r := <- res
		if r.err != nil {
			return fmt.Errorf("%w case %v action %v", r.err, c.index, r.index)
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
}

type ActionBase struct {
	Epoch        uint64
	Block        uint64
	ShouldBefore uint64
	index 		 int
}

func (a *ActionBase) StartAt() uint64 { return a.Block + (a.Epoch - 1) * uint64(CONFIG.BlocksPerEpoch) }
func (a *ActionBase) Before() uint64  { return a.ShouldBefore + (a.Epoch - 1) * uint64(CONFIG.BlocksPerEpoch) }
func (a *ActionBase) SetIndex(index int) { a.index = index }
func (a *ActionBase) Index() int { return a.index }

type SendTx struct {
	ActionBase
	Tx            *types.Transaction
	ShouldSucceed bool
}

func (a *SendTx) Run(ctx *Context) (err error) {
	err = ctx.nodes.Node().SendTransaction(context.Background(), a.Tx)
	if err != nil { return }
	for i := 0; i < 10; i++ {
		height, _, pending, err := ctx.nodes.Node().Confirm(a.Tx.Hash(), 1, 10)
		if err != nil { return err }
		if !pending { return fmt.Errorf("possible tx lost") }
		if height > 0 {
			if height <= a.Before() {
				rec, err := ctx.nodes.Node().TransactionReceipt(context.Background(), a.Tx.Hash())
				if err != nil { return err }
				if (rec.Status == 1) == a.ShouldSucceed { return nil }
				return fmt.Errorf("transaction status error, status %v, wanted %v", rec.Status == 1, a.ShouldSucceed )
			}
			return fmt.Errorf("tx packed too late, height %v, expected before %v", height, a.Before())
		}
		time.Sleep(time.Second)
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
	if err != nil { return }
    err = Assert(output, a.Assertions)
	return
}
