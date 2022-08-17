package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"zion-test/excel"
)

func main() {
	app := &cli.App{
		Name:   "zion test",
		Usage:  "",
		Action: start,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config.json",
				Usage: "configuration file",
			},
			&cli.StringFlag{
				Name:  "excel",
				Value: "testcase.xlsx",
				Usage: "test case excel file",
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Start error", "err", err)
	}
}

func start(c *cli.Context) error {
	excelPath := c.String("excel")
	if len(excelPath) == 0 {
		log.Fatal("excel path empty")
	}

	excel.ParseExcel(excelPath)

	return nil
}
