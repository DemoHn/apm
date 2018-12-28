package infra

import (
	"github.com/DemoHn/apm/infra/config"
	"github.com/DemoHn/apm/infra/logger"
)

// Init - init infra helpers - including config, logger
func Init(configFile *string, debugMode bool) (*config.Config, *logger.Logger) {
	configN := config.Init(configFile)
	log := logger.Init(debugMode)

	return configN, log
}
