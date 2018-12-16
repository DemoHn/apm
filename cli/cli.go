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
		{
			Name:   "stop",
			Usage:  "stop the instance of assigned ID",
			Flags:  stopFlags,
			Action: stopHandler,
		},
		{
			Name:   "list",
			Usage:  "list current status of instance",
			Flags:  listFlags,
			Action: listHandler,
		},
		{
			Name:   "kill",
			Usage:  "kill apm daemon",
			Flags:  killFlags,
			Action: killHandler,
		},
	}

	// parse and go
	return app.Run(args)
}
