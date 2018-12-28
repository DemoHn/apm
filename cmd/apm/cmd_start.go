package main

import (
	"github.com/DemoHn/apm/infra/logger"
	"github.com/DemoHn/apm/mod/master"
	"github.com/urfave/cli"
)

// StartCmd -
func StartCmd(name string) cli.Command {
	return cli.Command{
		Name:   name,
		Usage:  "create & start the instance",
		Flags:  startFlags,
		Action: startHandler,
	}
}

var startFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "cmd, c",
		Usage: "program command to execute",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "instance name",
	},
	cli.IntFlag{
		Name:  "id",
		Usage: "existing instance ID to start",
	},
}

func startHandler(c *cli.Context) error {
	var resp master.StartInstanceResponse
	var id int = c.Int("id")
	var rid *int
	log := logger.Init(false)

	if id != 0 {
		rid = &id
	}
	req := &master.StartInstanceRequest{
		Command: c.String("cmd"),
		Name:    c.String("name"),
		ID:      rid,
	}

	err := sendRequest("Tower.StartInstance", req, &resp)
	if err != nil {
		return err
	}
	if resp.IsSuccess == true {
		log.Infof("[apm] start instance success - ID = %d, pid = %d", resp.InstanceID, resp.PID)
	} else {
		log.Infof("[apm] start instance failed - error: %s", resp.Error)
	}
	return nil
}
