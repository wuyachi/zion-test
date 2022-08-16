package main

import (
	"github.com/polynetwork/bridge-common/log"
)

func Run() (err error) {
	var cases []*Case
	cs := make(chan *Case)
	res := make(chan *Case)
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
			c := <- res
			log.Info("Ran case", "index", c.index, "err", c.err)
		}
	} ()

	runCases(cs, res)
	return
}

func runCases(cs, res chan *Case) {
	for i := 0; i < CONFIG.ChainCount; i++ {
		go func(index int) {
			chain := &Chain{index, cs, res}
			log.Info("Launching chain", "index", index)
			chain.Run()
		} (i)
	}
}

type Chain struct {
	index int
	cs, res chan *Case
}

func (c *Chain) Run() (err error) {
	for {
		cs := <- c.cs
		if cs == nil {
			break
		}
		c.Start()
		ctx := &Context{}
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

