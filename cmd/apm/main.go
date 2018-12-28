package main

import (
	"os"

	"github.com/DemoHn/apm/cli"
	"github.com/DemoHn/apm/infra"
)

func main() {
	var err error
	_, log := infra.Init(nil, false)

	if err = cli.Parse(os.Args); err != nil {
		log.Errorf("%s", err.Error())
	}
}
