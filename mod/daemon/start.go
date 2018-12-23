package daemon

import (
	"fmt"
	"os"

	"github.com/DemoHn/apm/mod/config"
	"github.com/DemoHn/apm/mod/logger"
	"github.com/sevlyar/go-daemon"
)

// Start - start the main daemon
func Start(debugMode bool) error {
	var err error

	// init gObjects
	configN := config.Init(nil)
	log := logger.Init(debugMode)

	// fetch globalDir
	var globalDir string
	if globalDir, err = configN.FindString("global.dir"); err != nil {
		return addTag("cfg", err)
	}

	// make directory
	err = os.MkdirAll(globalDir, os.ModePerm)
	if err != nil {
		return err
	}

	// init logFile & pidFile
	var pidFile, logFile string
	if pidFile, err = configN.FindString("global.pidFile"); err != nil {
		return addTag("cfg", err)
	}
	if logFile, err = configN.FindString("global.logFile"); err != nil {
		return addTag("cfg", err)
	}

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

	var p *os.Process
	p, err = cntxt.Reborn()
	if err != nil {
		log.Infof("[apm] daemon has started")
		return nil
	}
	// if fork process succeed, let the parent process
	// go and run the folllowing logic in the child process
	if p != nil {
		return nil
	}
	defer cntxt.Release()

	// CHILD PROCESS
	return daemonHandler(debugMode)
}

// internal function
func addTag(tag string, err error) error {
	return fmt.Errorf("[%s]: %s", tag, err.Error())
}
