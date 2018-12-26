package cli

import (
	"github.com/DemoHn/apm/mod/logger"
	"github.com/DemoHn/apm/mod/master"
	"github.com/urfave/cli"
)

var restartFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "instance ID",
	},
}

func restartHandler(c *cli.Context) error {
	var resp master.RestartInstanceResponse
	var err error

	log := logger.Init(false)
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
