package cli

import (
	"io/ioutil"
	"strconv"
	"syscall"

	"github.com/DemoHn/apm/mod/config"
	"github.com/DemoHn/apm/mod/logger"
	"github.com/urfave/cli"
)

var killFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "force,f",
		Usage: "kill apm daemon by force",
	},
}

func killHandler(c *cli.Context) error {
	// get param
	var force = c.Bool("force")

	// read config
	configN := config.Init(nil)
	log := logger.Init(false)

	pidFile, err := configN.FindString("global.pidFile")
	if err != nil {
		return err
	}

	// read pidFile
	pidData, errR := ioutil.ReadFile(pidFile)
	if errR != nil {
		return errR
	}

	// parse to int
	pid, errT := strconv.Atoi(string(pidData))
	if errT != nil {
		return errT
	}

	// send quit signal (if exists)
	if force {
		return syscall.Kill(pid, syscall.SIGKILL)
	}

	errK := syscall.Kill(pid, syscall.SIGTERM)
	if errK != nil {
		return errK
	}

	log.Infof("[apm] kill PID:%d success", pid)
	return nil
}
