package main

import (
	"io/ioutil"
	"strconv"
	"syscall"

	"github.com/DemoHn/apm/infra/config"
	"github.com/DemoHn/apm/infra/logger"
	"github.com/urfave/cli"
)

// KillCmd -
func KillCmd(name string) cli.Command {
	return cli.Command{
		Name:   "kill",
		Usage:  "kill apm daemon",
		Flags:  killFlags,
		Action: killHandler,
	}
}

var killFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "force,f",
		Usage: "kill apm daemon by force",
	},
}

func killHandler(c *cli.Context) error {
	var err error
	// get param
	var force = c.Bool("force")

	// read config
	configN := config.Get()
	log := logger.Get()

	var pidFile string
	if pidFile, err = configN.FindString("global.pidFile"); err != nil {
		return err
	}

	// read pidFile
	var pidData []byte
	if pidData, err = ioutil.ReadFile(pidFile); err != nil {
		return err
	}

	// parse to int
	var pid int
	if pid, err = strconv.Atoi(string(pidData)); err != nil {
		return err
	}

	// send quit signal (if exists)
	if force {
		return syscall.Kill(pid, syscall.SIGKILL)
	}

	if err = syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return err
	}

	log.Infof("[apm] kill apm daemon (PID:%d) success", pid)
	return nil
}
