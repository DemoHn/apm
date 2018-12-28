package cli

import (
	"github.com/DemoHn/apm/infra/logger"
	"github.com/DemoHn/apm/mod/master"
	"github.com/urfave/cli"
)

var stopFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "instance ID",
	},
}

func stopHandler(c *cli.Context) error {
	var resp master.StopInstanceResponse
	var err error

	log := logger.Init(false)
	req := &master.StopInstanceRequest{
		ID: c.Int("id"),
	}

	if err = sendRequest("Tower.StopInstance", req, &resp); err != nil {
		return err
	}

	if resp.IsSuccess == true {
		log.Infof("[apm] stop instance success - ID = %d, code = %d", resp.InstanceID, resp.ExitCode)
	} else {
		log.Infof("[apm] stop instance error - error = %s", resp.Error)
	}
	return nil
}
