package main

import (
	"os"

	"github.com/DemoHn/apm/cli"
	"github.com/DemoHn/apm/infra/logger"
)

func main() {
	var err error
	log := logger.Init(false)

	if err = cli.Parse(os.Args); err != nil {
		log.Errorf("%s", err.Error())
	}
}
