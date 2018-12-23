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
	quit := make(chan os.Signal)

	// get config instance
	configN := config.Get()
	log := logger.Init(debugMode)

	// create master object
	masterN := master.New(debugMode)

	// get sockFile configlet
	var sockFile string
	if sockFile, err = configN.FindString("global.sockFile"); err != nil {
		return err
	}

	// init master
	if err = masterN.Init(sockFile); err != nil {
		return err
	}

	go func() {
		if err = masterN.Listen(); err != nil {
			log.Errorf("[apm] daemon encounters an error on listening '%s'", sockFile)

			quit <- os.Interrupt
		}
	}()

	log.Infof("[apm] daemon start listening to '%s'", sockFile)

	// wait for quit signal
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 'quit' channel will receive data from two sources:
	// 1. masterN.Listen() quits with error
	// 2. parent process receives a SIGTERM signal
	// Thus if err != nil, `quit` must happens on condition #1
	if err != nil {
		return err
	}

	log.Info("[apm] going to teardown")
	return masterN.Teardown()
}
