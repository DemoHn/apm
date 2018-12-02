package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DemoHn/apm/mod/logger"
	"github.com/DemoHn/apm/mod/master"
	"github.com/urfave/cli"
)

var log *logger.Logger
var daemonFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "debug,d",
		Usage: "debug mode",
	},
}

func daemonHandler(c *cli.Context) error {
	var err error
	// signal
	quit := make(chan os.Signal)
	debugMode := c.Bool("debug")
	// create & init master
	m := master.New(debugMode)
	// init logger
	log = logger.Get()
	err = m.Init(sockFile)
	if err != nil {
		return err
	}

	go func() {
		err = m.Listen()
		if err != nil {
			log.Error("[apm] server encounters an error:", err)
			// send quit signal
			quit <- os.Interrupt
		}
	}()

	log.Infof("[apm] listening to server: unix(%s)", sockFile)
	signal.Notify(quit, syscall.SIGINT)
	// wait for quit signal
	sig := <-quit

	log.Infof("[apm] going to shutdown due to recieve signal: %s", sig.String())
	err = m.Shutdown()
	if err != nil {
		return err
	}

	return nil
}
