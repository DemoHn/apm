package cli

import (
	"fmt"

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

	req := &master.StopInstanceRequest{
		ID: c.Int("id"),
	}

	err := sendRequest("Tower.StopInstance", req, &resp)
	if err != nil {
		return err
	}

	if resp.IsSuccess == true {
		fmt.Printf("[apm] stop instance success - ID = %d, code = %d", resp.InstanceID, resp.ExitCode)
	} else {
		fmt.Printf("[apm] stop instance error - error = %s", resp.Error)
	}
	return nil
}
