package common

import (
	"sort"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/polynetwork/bridge-common/chains/eth"
	"github.com/polynetwork/bridge-common/log"
)

type Context struct {
	nodes *eth.SDK
}

type Case struct {
	index   int
	err     error
	actions []Action
	plan    []*ActionItem
}

func (c *Case) Run(ctx *Context) (err error) {
	// Sort actions
	bp := make(map[uint64][]Action)
	for _, a := range c.actions {
		bp[a.StartAt()] = append(bp[a.StartAt()], a)
	}
	c.plan = make([]*ActionItem, 0, len(bp))
	for b, actions := range bp {
		c.plan = append(c.plan, &ActionItem{b, actions})
	}
	sort.Slice(c.plan, func(i, j int) bool { return c.plan[i].start < c.plan[j].start })

	for i, item := range c.plan {
		log.Info("Running plan", "index", i, "action_count", len(item.actions), "at", item.start)
		res := make(chan error, len(item.actions))
		for _, action := range item.actions {
			go func(a Action) {
				res <- a.Run(ctx)
			}(action)
		}
		for j := 0; j < len(item.actions); j++ {
			err = <-res
			if err != nil {
				return err
			}
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
}

type ActionBase struct {
	Epoch        uint64
	Block        uint64
	ShouldBefore uint64
}

func (a *ActionBase) StartAt() uint64 { return a.Block }
func (a *ActionBase) Before() uint64  { return a.ShouldBefore }

type SendTx struct {
	ActionBase
	Tx            types.Transaction
	ShouldSucceed bool
}

func (a *SendTx) Run(ctx *Context) error {
	return nil
}

type Query struct {
	ActionBase
	Request    ethereum.CallMsg
	Assertions []Assertion
}

func (a *Query) Run(ctx *Context) error {
	return nil
}
