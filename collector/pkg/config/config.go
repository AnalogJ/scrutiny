package config

import (
	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/collector/pkg/errors"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
	"os"
)

// When initializing this class the following methods must be called:
// Config.New
// Config.Init
// This is done automatically when created via the Factory.
type configuration struct {
	*viper.Viper
}

//Viper uses the following precedence order. Each item takes precedence over the item below it:
// explicit call to Set
// flag
// env
// config
// key/value store
// default

func (c *configuration) Init() error {
	c.Viper = viper.New()
	//set defaults
	c.SetDefault("host.id", "")

	c.SetDefault("devices", []string{})

	c.SetDefault("log.level", "INFO")
	c.SetDefault("log.file", "")

	c.SetDefault("api.endpoint", "http://localhost:8080")

	//c.SetDefault("collect.short.command", "-a -o on -S on")

	//if you want to load a non-standard location system config file (~/drawbridge.yml), use ReadConfig
	c.SetConfigType("yaml")
	//c.SetConfigName("drawbridge")
	//c.AddConfigPath("$HOME/")

	//CLI options will be added via the `Set()` function
	return nil
}

func (c *configuration) ReadConfig(configFilePath string) error {
	configFilePath, err := utils.ExpandPath(configFilePath)
	if err != nil {
		return err
	}

	if !utils.FileExists(configFilePath) {
		log.Printf("No configuration file found at %v. Using Defaults.", configFilePath)
		return errors.ConfigFileMissingError("The configuration file could not be found.")
	}

	//validate config file contents
	//err = c.ValidateConfigFile(configFilePath)
	//if err != nil {
	//	log.Printf("Config file at `%v` is invalid: %s", configFilePath, err)
	//	return err
	//}

	log.Printf("Loading configuration file: %s", configFilePath)

	config_data, err := os.Open(configFilePath)
	if err != nil {
		log.Printf("Error reading configuration file: %s", err)
		return err
	}

	err = c.MergeConfig(config_data)
	if err != nil {
		return err
	}

	return c.ValidateConfig()
}

// This function ensures that the merged config works correctly.
func (c *configuration) ValidateConfig() error {

	//TODO:
	// check that device prefix matches OS
	// check that schema of config file is valid

	return nil
}

func (c *configuration) GetScanOverrides() []models.ScanOverride {
	// we have to support 2 types of device types.
	// - simple device type (device_type: 'sat')
	// and list of device types (type: \n- 3ware,0 \n- 3ware,1 \n- 3ware,2)
	// GetString will return "" if this is a list of device types.

	overrides := []models.ScanOverride{}
	c.UnmarshalKey("devices", &overrides, func(c *mapstructure.DecoderConfig) { c.WeaklyTypedInput = true })
	return overrides
}
