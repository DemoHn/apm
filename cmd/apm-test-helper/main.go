package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) > 1 {
		command := os.Args[1]
		args := os.Args[2:]
		switch command {
		case "normal-run":
			normalRunHandler(args)
		case "stop-on-time":
			stopTimeoutHandler(args)
		}
	} else {
		fmt.Println("[apm-test-helper] no command received")
	}

}
