package daemon

import (
	"github.com/DemoHn/apm/mod/logger"
	"github.com/sevlyar/go-daemon"
)

// StartDaemon - start the main daemon
func StartDaemon(debugMode bool) {
	// init logger
	log := logger.Init(debugMode)

	// init context
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Error("[apm] Unable to run daemon: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	err2 := daemonHandler(debugMode)
	if err2 != nil {
		log.Error("[apm] error to start:", err)
	}
}
