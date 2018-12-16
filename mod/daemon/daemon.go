package daemon

import (
	"os"

	"github.com/DemoHn/apm/mod/config"
	"github.com/DemoHn/apm/mod/logger"
	"github.com/sevlyar/go-daemon"
)

// StartDaemon - start the main daemon
func StartDaemon(debugMode bool) {
	// init config
	configN := config.Init(nil)
	// init logger
	log := logger.Init(debugMode)

	// init globalDir
	globalDir, _ := configN.FindString("global.dir")
	errM := os.MkdirAll(globalDir, os.ModePerm)
	if errM != nil {
		log.Error(errM)
	}

	// init logFile & pidFile
	pidFile, _ := configN.FindString("global.pidFile")
	logFile, _ := configN.FindString("global.logFile")
	// init context
	cntxt := &daemon.Context{
		PidFileName: pidFile,
		PidFilePerm: 0644,
		LogFileName: logFile,
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
