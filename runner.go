package main

import (
	"main/common"
	"sync"

	"github.com/polynetwork/bridge-common/log"
)

func Run() (err error) {
	var cases []*common.Case
	cs := make(chan *common.Case)
	res := make(chan *common.Case)
	go func() {
		for i, c := range cases {
			c.index = i
			cs <- c
			log.Info("Running case", "index", i, "action_count", len(c.actions))
		}
		// Signal to stop chains
		cs <- nil
	}()
	go func() {
		for i := 0; i < len(cases); i++ {
			c := <-res
			log.Info("Ran case", "index", c.index, "err", c.err)
		}
	}()

	runCases(cs, res)
	return
}

func runCases(cs, res chan *common.Case) {
	wg := &sync.WaitGroup{}
	for i := 0; i < CONFIG.ChainCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			chain := &Chain{index, cs, res, CONFIG.Bin, CONFIG.NodesPerChain, CONFIG.NodesPortStart}
			log.Info("Launching chain", "index", index)
			chain.Run()
		}(i)
	}
	wg.Wait()
}

type Chain struct {
	index   int
	cs, res chan *common.Case
	bin     string
	nodes   int
	port    int
}

func (c *Chain) Run() (err error) {
	for {
		cs := <-c.cs
		if cs == nil {
			break
		}
		c.Start()
		ctx := &common.Context{}
		cs.err = cs.Run(ctx)
		c.res <- cs
		c.Stop()
	}
	return
}

func (c *Chain) Start() {
}

func (c *Chain) Stop() {
}
