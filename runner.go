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

func parseCases(path string) (cases []*Case, err error) {
	rawCases, err := ParseExcel(path)
	if err != nil {
		log.Error("parse case file failed", "err", err)
		return
	}

	for _, rawCase := range rawCases {
		c, e := rawCase.Pack()
		if e != nil {
			err = fmt.Errorf("pack rawCase failed. err=%s", e)
			cases = append(cases, &c)
		}
	}
	return
}

func dumpResult(cases []*Case) (err error) {
	return
}

func Run() (err error) {
	cases, err := parseCases(CONFIG.Input)
	if err != nil {
		return
	}

	log.Info("Parsed cases", "count", len(cases))

	cs := make(chan *Case)
	res := make(chan *Case)
	go func() {
		for i, c := range cases {
			// c.index = i
			cs <- c
			log.Info("Running case", "index", i, "action_count", len(c.actions))
		}
		// Signal to stop chains
		for i := 0; i < CONFIG.ChainCount; i++ {
			cs <- nil
		}
	}()
	done := make(chan bool)
	go func() {
		for i := 0; i < len(cases); i++ {
			c := <-res
			log.Info("Ran case", "index", c.index, "err", c.err)
		}
		close(done)
	}()

	runCases(cs, res)
	<-done
	return dumpResult(cases)
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
	index   int
	bin     string
	cs, res chan *Case
	nodes   int
	port    int
	sdk     *eth.SDK
}

func (c *Chain) Run() (err error) {
	for {
		cs := <-c.cs
		if cs == nil {
			break
		}
		c.Start(cs.index)
		ctx := &Context{nodes: c.sdk}
		cs.err = cs.Run(ctx)
		c.res <- cs
		c.Stop(cs.index)
	}
	return
}

func (c *Chain) Start(caseIndex int) {
	err := runCmd(CONFIG.StartScript, c.bin, CONFIG.ChainDir, fmt.Sprint(c.index), fmt.Sprint(CONFIG.NodesPerChain), fmt.Sprint(CONFIG.NodesPortStart+(c.index*100)),
		CONFIG.CheckBin, fmt.Sprint(caseIndex),
	)
	if err != nil {
		log.Fatal("Failed to start chain", "index", c.index, "err", err)
	}
	time.Sleep(time.Second * 10)
	var urls []string
	for i := 0; i < c.nodes; i++ {
		urls = append(urls, fmt.Sprintf("http://127.0.0.1:%v", CONFIG.NodesPortStart+(c.index*100)+i))
	}
	c.sdk, err = eth.WithOptions(0, urls, time.Minute, 1)
	if err != nil {
		log.Fatal("Failed to create eth client", "index", c.index, "err", err)
	}
	height, err := c.sdk.Node().GetLatestHeight()
	if err != nil {
		log.Fatal("Failed to get chain height", "index", c.index, "err", err)
	}
	log.Info("Chain started", "index", c.index, "height", height)
}

func (c *Chain) Stop(caseIndex int) {
	if c.sdk != nil {
		c.sdk.Stop()
		c.sdk = nil
	}
	err := runCmd(CONFIG.StopScript, c.bin, CONFIG.ChainDir, fmt.Sprint(c.index), fmt.Sprint(CONFIG.NodesPerChain), fmt.Sprint(CONFIG.NodesPortStart+(c.index*100)),
		CONFIG.CheckBin, fmt.Sprint(caseIndex),
	)
	if err != nil {
		log.Fatal("Failed to stop chain", "index", c.index, "err", err)
	}
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
	chain.Start(0)
	return
}

func stopChain(ctx *cli.Context) (err error) {
	chain := &Chain{0, CONFIG.Bin, nil, nil, CONFIG.NodesPerChain, CONFIG.NodesPortStart, nil}
	chain.Stop(0)
	return
}
