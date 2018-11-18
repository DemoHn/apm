package cli

import (
	"github.com/urfave/cli"
)

// Parse -
func Parse(args []string) error {
	app := cli.NewApp()
	// init app properties
	app.Name = "apm"
	app.Usage = "The process manager that controls the process flow"
	app.Version = "1.0.0"

	// define command
	app.Commands = []cli.Command{
		{
			Name:   "daemon",
			Usage:  "start the daemon",
			Flags:  daemonFlags,
			Action: daemonHandler,
		},
		{
			Name:   "start",
			Usage:  "create & start the instance",
			Flags:  startFlags,
			Action: startHandler,
		},
	}

	// parse and go
	return app.Run(args)
}
