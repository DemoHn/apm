package main

import (
	"os"

	"github.com/DemoHn/apm/cli"
)

func main() {
	err := cli.Parse(os.Args)
	if err != nil {
		panic(err)
	}
}
