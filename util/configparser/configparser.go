package configparser

import (
	"fmt"
	"io/ioutil"
	"os"
)

// MacroParserFunc - macro parser function
type MacroParserFunc func(key string, item interface{}) interface{}

// Config type
type Config struct {
	configType  string
	configItem  map[string]interface{}
	macroParser MacroParserFunc
}

// New - initialize a new config instance
func New(configType string) *Config {
	return &Config{
		configType: configType,
		macroParser: func(key string, item interface{}) interface{} {
			// default macroParser, just return the original value
			return item
		},
	}
}

// SetMacroParser - set macro parser to replace macros like "$HOME$", "$TMP$", etc.
func (config *Config) SetMacroParser(parser MacroParserFunc) {
	config.macroParser = parser
}

// Load - parse & load from a file
// currently only supports "yaml"
func (config *Config) Load(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	// read all
	data, errRead := ioutil.ReadAll(fd)
	if errRead != nil {
		return errRead
	}
	// then parse config
	errParse := config.parseConfig(data)
	if errParse != nil {
		return errParse
	}
	// finish
	return nil
}

// LoadFromData - instead of loading config from a file,
// you can directly inject data
func (config *Config) LoadFromData(data []byte) error {
	errParse := config.parseConfig(data)

	if errParse != nil {
		return errParse
	}

	return nil
}

// LoadDefault - load config from default value
func (config *Config) LoadDefault(defaultConf map[string]interface{}) {
	for k, v := range defaultConf {
		config.configItem[k] = v
	}
}

// Find - find config value from key.
func (config *Config) Find(key string) (result interface{}, err error) {
	if value, ok := config.configItem[key]; ok {
		result = value
		err = nil
		return
	}

	err = fmt.Errorf("Config key `%s` not found", key)
	return
}

// parseConfig - parse config from data
func (config *Config) parseConfig(data []byte) (err error) {
	switch config.configType {
	case "yaml":
		{
			config.configItem, err = ParseYamlConfig(data, config.macroParser)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("Config Type `%s` not found", config.configType)
}
