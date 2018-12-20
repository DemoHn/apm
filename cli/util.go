package cli

import (
	"net/rpc"

	"github.com/DemoHn/apm/mod/config"
)

func sendRequest(method string, input interface{}, output interface{}) error {
	configN := config.Init(nil)
	sockFile, _ := configN.FindString("global.sockFile")

	client, err := rpc.DialHTTP("unix", sockFile)
	if err != nil {
		return err
	}

	err2 := client.Call(method, input, output)
	if err2 != nil {
		return err2
	}

	return nil
}
