package cli

import (
	"fmt"

	"github.com/DemoHn/apm/mod/master"
	"github.com/urfave/cli"
)

var startFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "cmd, c",
		Usage: "program command to execute",
	},
}

func startHandler(c *cli.Context) error {
	var resp master.StartInstanceResponse

	req := &master.StartInstanceRequest{
		Command: c.String("cmd"),
	}

	err := sendRequest("Tower.StartInstance", req, &resp)
	if err != nil {
		return err
	}

	if resp.IsSuccess == true {
		fmt.Printf("[apm] start instance success. ID = %d, pid = %d", resp.InstanceID, resp.PID)
	} else {
		fmt.Printf("[apm] start instance failed - error: %v", resp.Error)
	}

	return nil
}
