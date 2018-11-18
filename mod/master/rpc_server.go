package master

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

// rpc_server: all about RPC server

const (
	unixNetwork string = "unix"
)

type rpcServer struct {
	sockFile   string
	httpServer *http.Server
}

func (m *Master) initRPC(sockFile string) error {
	var err error
	tower := &Tower{
		master: m,
	}
	err = rpc.Register(tower)
	if err != nil {
		return err
	}

	rpc.HandleHTTP()
	m.rpc = &rpcServer{
		sockFile:   sockFile,
		httpServer: &http.Server{},
	}
	return nil
}

func (m *Master) listen() error {
	var l net.Listener
	var err error
	if m.rpc == nil {
		return fmt.Errorf("Listen to server failed - is master.rpc initialized?")
	}

	l, err = net.Listen(unixNetwork, m.rpc.sockFile)
	if err != nil {
		return err
	}

	err = m.rpc.httpServer.Serve(l)
	if err != nil {
		return err
	}

	return nil
}

func (m *Master) shutdown() error {
	// TODO - shutdown logic
	return nil
}
