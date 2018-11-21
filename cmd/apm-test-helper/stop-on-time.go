package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// args[0]: timeout in millisecond
func stopTimeoutHandler(args []string) {
	timeout, _ := strconv.Atoi(args[0])
	fmt.Printf("will finish in %dms\n", timeout)
	<-time.After(time.Millisecond * time.Duration(timeout))

	os.Exit(0)
}
