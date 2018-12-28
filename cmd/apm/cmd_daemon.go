package main

import (
	"github.com/DemoHn/apm/infra/logger"
	"github.com/DemoHn/apm/mod/daemon"
	"github.com/urfave/cli"
)

// DaemonCmd - handling daemon
func DaemonCmd(name string) cli.Command {
	return cli.Command{
		Name:   name,
		Usage:  "start the daemon",
		Flags:  daemonFlags,
		Action: daemonHandler,
	}
}

var daemonFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "debug,d",
		Usage: "debug mode",
	},
	cli.BoolFlag{
		Name:  "foreground,fg",
		Usage: "start the daemon on foreground",
	},
}

func daemonHandler(c *cli.Context) error {
	var err error

	// logger with debugMode = false
	log := logger.Get()
	debugMode := c.Bool("debug")
	fg := c.Bool("foreground")

	if fg {
		err = daemon.StartForeground(debugMode)
	} else {
		err = daemon.Start(debugMode)
	}
	// create & init master
	if err == nil {
		log.Info("[apm] start apm daemon succeed")
	}
	return err
}
