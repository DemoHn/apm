package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DemoHn/apm/mod/master"
	"github.com/urfave/cli"
)

var daemonFlags = []cli.Flag{}

func daemonHandler(c *cli.Context) error {
	var err error
	// signal
	quit := make(chan os.Signal)

	// create & init master
	m := master.New()
	err = m.Init(sockFile)
	if err != nil {
		return err
	}

	go func() {
		err = m.Listen()
		if err != nil {
			fmt.Println("[apm] server encounters an error:", err)
			// send quit signal
			quit <- os.Interrupt
		}
	}()

	fmt.Println("[apm] listening to server")
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// wait for quit signal
	<-quit

	fmt.Println("[apm] going to shutdown")
	err = m.Shutdown()
	if err != nil {
		return err
	}

	return nil
}
