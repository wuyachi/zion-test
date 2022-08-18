package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/polynetwork/bridge-common/chains/eth"
	"github.com/polynetwork/bridge-common/log"
	"github.com/urfave/cli/v2"
)

func Run() (err error) {
	var cases []*Case
	cs := make(chan *Case)
	res := make(chan *Case)
	go func() {
		for i, c := range cases {
			// c.index = i
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

func runCases(cs, res chan *Case) {
	wg := &sync.WaitGroup{}
	for i := 0; i < CONFIG.ChainCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			chain := &Chain{index, CONFIG.Bin, cs, res, CONFIG.NodesPerChain, CONFIG.NodesPortStart, nil}
			log.Info("Launching chain", "index", index)
			chain.Run()
		}(i)
	}
	wg.Wait()
}

type Chain struct {
	index int
	bin string
	cs, res chan *Case
	nodes int
	port int
	sdk *eth.SDK
}

func (c *Chain) Run() (err error) {
	for {
		cs := <-c.cs
		if cs == nil {
			break
		}
		c.Start()
		ctx := &Context{nodes: c.sdk}
		cs.err = cs.Run(ctx)
		c.res <- cs
		c.Stop()
	}
	return
}

func (c *Chain) Start() {
	err := runCmd(CONFIG.StartScript, c.bin, CONFIG.ChainDir, fmt.Sprint(c.index), fmt.Sprint(CONFIG.NodesPerChain), fmt.Sprint(CONFIG.NodesPortStart + (c.index * 100)))
	if err != nil {
		log.Fatal("Failed to start chain", "index", c.index, "err", err)
	}
	time.Sleep(time.Second * 30)
	var urls []string
	for i := 0; i < c.nodes; i++ {
		urls = append(urls, fmt.Sprintf("http://127.0.0.1:%v", CONFIG.NodesPortStart + (c.index * 100) + i))
	}
	c.sdk, err = eth.WithOptions(0, urls, time.Minute, 1)
	if err != nil { 
		log.Fatal("Failed to create eth client", "index", c.index, "err", err)
	}
}

func (c *Chain) Stop() {
	err := runCmd(CONFIG.StopScript, c.bin, CONFIG.ChainDir, fmt.Sprint(c.index), fmt.Sprint(CONFIG.NodesPerChain), fmt.Sprint(CONFIG.NodesPortStart + (c.index * 100)))
	if err != nil {
		log.Fatal("Failed to stop chain", "index", c.index, "err", err)
	}
	c.sdk = nil
}

func runCmd(bin string, args ...string) (err error) {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}


func runChain(ctx *cli.Context) (err error) {
	chain := &Chain{0, CONFIG.Bin, nil, nil, CONFIG.NodesPerChain, CONFIG.NodesPortStart, nil}
	chain.Start()
	return
}

func stopChain(ctx *cli.Context) (err error) {
	chain := &Chain{0, CONFIG.Bin, nil, nil, CONFIG.NodesPerChain, CONFIG.NodesPortStart, nil}
	chain.Stop()
	return
}