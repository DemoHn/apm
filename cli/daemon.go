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

	debugMode := c.Bool("debug")
	// create & init master
	daemon.StartDaemon(debugMode)
	return nil
}
