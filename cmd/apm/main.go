package main

import (
	"os"

	"github.com/DemoHn/apm/infra"
	"github.com/urfave/cli"
)

func main() {
	var err error
	_, log := infra.Init(nil, false)

	if err = parseArgs(os.Args); err != nil {
		log.Errorf("%s", err.Error())
	}
}

func parseArgs(args []string) error {
	app := cli.NewApp()
	// init app properties
	app.Name = "apm"
	app.Usage = "The process manager that controls the process flow"
	app.Version = "1.0.0"

	// define command
	app.Commands = []cli.Command{
		DaemonCmd("daemon"),
		StartCmd("start"),
		StopCmd("stop"),
		ListCmd("list"),
		KillCmd("kill"),
	}

	// parse and go
	return app.Run(args)
}
