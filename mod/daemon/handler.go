package daemon

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DemoHn/apm/mod/config"
	"github.com/DemoHn/apm/mod/logger"
	"github.com/DemoHn/apm/mod/master"
)

// handle daemon
func daemonHandler(debugMode bool) error {
	var err error
	// signal
	quit := make(chan os.Signal)
	// get config instance
	configN := config.Get()
	log := logger.Get()

	// create & init master
	masterN := master.New(debugMode)

	// sockFile
	sockFile, _ := configN.FindString("global.sockFile")
	err = masterN.Init(sockFile)
	if err != nil {
		return err
	}

	go func() {
		err = masterN.Listen()
		if err != nil {
			log.Error("[apm] server encounters an error:", err)
			// send quit signal
			quit <- os.Interrupt
		}
	}()

	log.Info("[apm] listening to server")
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// wait for quit signal
	<-quit

	log.Info("[apm] going to teardown")
	err = masterN.Teardown()
	if err != nil {
		return err
	}

	return nil
}
