package cli

import (
	"net/rpc"

	"github.com/DemoHn/apm/mod/config"
)

func sendRequest(method string, input interface{}, output interface{}) error {
	var err error
	configN := config.Init(nil)

	var sockFile string
	if sockFile, err = configN.FindString("global.sockFile"); err != nil {
		return err
	}

	var client *rpc.Client
	if client, err = rpc.DialHTTP("unix", sockFile); err != nil {
		return err
	}

	return client.Call(method, input, output)
}
