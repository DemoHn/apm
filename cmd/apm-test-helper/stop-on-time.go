package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// args[0]: timeout in millisecond
func stopTimeoutHandler(args []string) {
	var exitCode = 0
	timeout, _ := strconv.Atoi(args[0])
	if len(args) > 1 {
		exitCode, _ = strconv.Atoi(args[1])
	}

	fmt.Printf("will finish in %dms\n", timeout)
	<-time.After(time.Millisecond * time.Duration(timeout))

	os.Exit(exitCode)
}
