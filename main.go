package main

import (
	"os"

	"github.com/polynetwork/bridge-common/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "zion-test",
		Usage:  "zion test framework",
		Action: start,
		Before: Init,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config.json",
			},
			&cli.StringFlag{
				Name:  "excel",
				Value: "testcase.xlsx",
				Usage: "test case excel file",
			},
		},

		Commands: cli.Commands{
			{
				Name:   "dump",
				Action: dump,
			},
			{
				Name:   "spawn",
				Action: runChain,
			},
			{
				Name:   "stop",
				Action: stopChain,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Start failure", "err", err)
	}
}

func start(ctx *cli.Context) (err error) {
	err = Run()
	return
}

func Init(ctx *cli.Context) (err error) {
	log.Init(nil)
	err = NewConfig(ctx.String("config"))
	if err != nil {
		return
	}

	if len(CONFIG.Input) == 0 {
		CONFIG.Input = ctx.String("excel")
	}

	err = Setup()
	return
}
