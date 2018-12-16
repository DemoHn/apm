package config

import (
	"os"
	"regexp"

	"github.com/DemoHn/apm/util/configparser"
)

// Config - config object
type Config = configparser.Config

var gConfig *Config

// default config
var defaultConfig = map[string]interface{}{
	// put global.dir at beginning!
	"global.dir":      "${HOME}/.apm",
	"global.sockFile": "$(global.dir)/apm.sock",
	"global.pidFile":  "$(global.dir)/apm.pid",
	"global.logFile":  "$(global.dir)/apm.log",
}

// Init - init config from config dir (yaml)
func Init(configDir *string) *Config {
	if gConfig == nil {
		config := configparser.New("yaml")
		// set macro parser
		config.SetMacroParser(func(key string, item interface{}) interface{} {
			// replace $HOME -> <homeDir>
			re := regexp.MustCompile("\\$\\{HOME\\}")
			// if item is string to replace
			if itemStr, ok := item.(string); ok {
				return re.ReplaceAllString(itemStr, os.Getenv("HOME"))
			}

			return item
		})

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
