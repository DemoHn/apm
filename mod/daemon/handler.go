package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DemoHn/apm/mod/config"
	"github.com/DemoHn/apm/mod/master"
)

// handle daemon
func daemonHandler(debugMode bool) error {
	var err error
	// signal
	quit := make(chan os.Signal)
	// get config instance
	configN := config.Get()
	// create & init master
	m := master.New(debugMode)

	// sockFile
	sockFile, _ := configN.FindString("global.sockFile")
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
