package config

import (
	"github.com/DemoHn/apm/util/configparser"
)

// Config - config object
type Config = configparser.Config

var gConfig *Config

// default config
var defaultConfig = map[string]interface{}{
	"global.dir":      "${HOME}/.apm",
	"global.sockFile": "$(global.dir)/apm.sock",
	"global.pidFile":  "$(global.dir)/apm.pid",
	"global.logFile":  "$(global.dir)/apm.log",
}

// Init - init config from config dir (yaml)
func Init(configDir *string) *Config {
	if gConfig == nil {
		config := configparser.New("yaml")
		config.LoadDefault(defaultConfig)

		if configDir != nil {
			config.Load(*configDir)
		}

		gConfig = config
	}

	return gConfig
}

// Get - get global instance of config
func Get() *Config {
	return gConfig
}