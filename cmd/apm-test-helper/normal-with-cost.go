package main

import (
	"os"
	"os/signal"
	"strconv"
	"time"
)

func normalWithCostHandler(args []string) {
	exitCode := 0
	if len(args) > 0 {
		exitCode, _ = strconv.Atoi(args[0])
	}

	go func() {
		j := 2
		for {
			for i := 1; i < 100000; i++ {
				j = i * j
			}
			time.Sleep(20 * time.Microsecond)
		}

	}()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	<-c
	os.Exit(exitCode)
}
