package main

import (
	"net/rpc"
	"time"

	"github.com/DemoHn/apm/infra/config"
	"github.com/DemoHn/apm/mod/daemon"
)

func sendRequest(method string, input interface{}, output interface{}) error {
	var err error
	configN := config.Get()

	var sockFile string
	if sockFile, err = configN.FindString("global.sockFile"); err != nil {
		return err
	}
	// start daemon to ensure server is running
	if err = daemon.Start(false); err != nil {
		return err
	}
	if err = daemon.PingTimeout(100*time.Millisecond, 5*time.Second); err != nil {
		return err
	}

	var client *rpc.Client
	if client, err = rpc.DialHTTP("unix", sockFile); err != nil {
		return err
	}

	return client.Call(method, input, output)
}
