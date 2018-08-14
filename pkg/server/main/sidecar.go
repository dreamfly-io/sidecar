package main

import (
	"time"
	"os"
	"github.com/urfave/cli"


	_ "github.com/dreamfly-io/sidecar/pkg/proxy/network/listener"
	_ "github.com/dreamfly-io/sidecar/pkg/proxy/network/connection"
	"github.com/dreamfly-io/sidecar/pkg/server/starter"
	"github.com/dreamfly-io/sidecar/pkg/util/log"
	"github.com/dreamfly-io/sidecar/pkg/server/config"
)

func main() {
	app := cli.NewApp()
	app.Name = "sidecar"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Copyright = "(c) 2018 DreamFly.IO"
	app.Usage = "Just for prototype."

	//commands
	app.Commands = []cli.Command{
		cmdStart,
	}

	//action
	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)

		c.App.Setup()
		return nil
	}

	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Run(os.Args)
}

var (
	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start sidecar",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				Value:  "configs/config.json",
			},
		},
		Action: func(c *cli.Context) error {

			configPath := c.String("config")

			log.StartLogger.Debug(configPath)

			conf := config.Load(configPath)
			starter.Start(conf)
			return nil
		},
	}
)


