package main

import (
	"os"
	"os/signal"
)

func normalRunHandler(args []string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	<-c
	os.Exit(0)
}
