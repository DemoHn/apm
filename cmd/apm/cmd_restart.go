package main

import (
	"github.com/DemoHn/apm/infra/logger"
	"github.com/DemoHn/apm/mod/master"
	"github.com/urfave/cli"
)

// RestartCmd -
func RestartCmd(name string) cli.Command {
	return cli.Command{
		Name:   "restart",
		Usage:  "restart the instance of assigned ID",
		Flags:  restartFlags,
		Action: restartHandler,
	}
}

var restartFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "instance ID",
	},
}

func restartHandler(c *cli.Context) error {
	var resp master.RestartInstanceResponse
	var err error

	log := logger.Get()
	req := &master.RestartInstanceRequest{
		ID: c.Int("id"),
	}

	if err = sendRequest("Tower.RestartInstance", req, &resp); err != nil {
		return err
	}

	if resp.IsSuccess == true {
		log.Infof("[apm] restart instance success - ID = %d", resp.InstanceID)
	} else {
		log.Infof("[apm] restart instance error - error = %s", resp.Error)
	}
	return nil
}
