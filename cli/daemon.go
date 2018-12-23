package cli

import (
	"github.com/DemoHn/apm/mod/daemon"
	"github.com/urfave/cli"
)

var daemonFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "debug,d",
		Usage: "debug mode",
	},
}

func daemonHandler(c *cli.Context) error {
	var err error

	debugMode := c.Bool("debug")
	// create & init master
	err = daemon.Start(debugMode)
	return err
}
