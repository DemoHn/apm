package main

import (
	"os"
	"os/signal"
	"strconv"
)

func normalRunHandler(args []string) {
	exitCode := 0
	if len(args) > 0 {
		exitCode, _ = strconv.Atoi(args[0])
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	<-c
	os.Exit(exitCode)
}
