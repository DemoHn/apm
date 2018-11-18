package cli

import "net/rpc"

const (
	sockFile = "/tmp/apm.sock"
)

func sendRequest(method string, input interface{}, output interface{}) error {
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
