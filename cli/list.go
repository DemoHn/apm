package cli

import (
	"fmt"

	"github.com/DemoHn/apm/mod/master"
	"github.com/urfave/cli"
)

var listFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "instance ID. Use this to show the information of this one",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "Filter shown instances by name",
	},
}

func listHandler(c *cli.Context) error {
	var resp master.ListInstanceResponse
	// get id
	var id int = c.Int("id")
	var rid *int
	if id != 0 {
		rid = &id
	}
	var name string = c.String("name")
	var rname *string
	if name != "" {
		rname = &name
	}
	req := &master.ListInstanceRequest{
		ID:   rid,
		Name: rname,
	}
	err := sendRequest("Tower.ListInstance", req, &resp)
	if err != nil {
		return err
	}
	// just print output
	// TODO: more human-readable response
	fmt.Println("resp", resp)
	return nil
}
