package config

import (
	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
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
	c.SetDefault("web.listen.port", "8080")
	c.SetDefault("web.listen.host", "0.0.0.0")
	c.SetDefault("web.listen.basepath", "")
	c.SetDefault("web.src.frontend.path", "/opt/scrutiny/web")
	c.SetDefault("web.database.location", "/opt/scrutiny/config/scrutiny.db")

	c.SetDefault("log.level", "INFO")
	c.SetDefault("log.file", "")

	c.SetDefault("notify.urls", []string{})

	c.SetDefault("web.influxdb.scheme", "http")
	c.SetDefault("web.influxdb.host", "localhost")
	c.SetDefault("web.influxdb.port", "8086")
	c.SetDefault("web.influxdb.org", "scrutiny")
	c.SetDefault("web.influxdb.bucket", "metrics")
	c.SetDefault("web.influxdb.init_username", "admin")
	c.SetDefault("web.influxdb.init_password", "password12345")
	c.SetDefault("web.influxdb.token", "scrutiny-default-admin-token")
	c.SetDefault("web.influxdb.retention_policy", true)

	//c.SetDefault("disks.include", []string{})
	//c.SetDefault("disks.exclude", []string{})

	//if you want to load a non-standard location system config file (~/drawbridge.yml), use ReadConfig
	c.SetConfigType("yaml")
	//c.SetConfigName("drawbridge")
	//c.AddConfigPath("$HOME/")

	//configure env variable parsing.
	c.SetEnvPrefix("SCRUTINY")
	c.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	c.AutomaticEnv()

	//CLI options will be added via the `Set()` function
	return c.ValidateConfig()
}

func (c *configuration) SubKeys(key string) []string {
	return c.Sub(key).AllKeys()
}

func (c *configuration) ReadConfig(configFilePath string) error {
	//make sure that we specify that this is the correct config path (for eventual WriteConfig() calls)
	c.SetConfigFile(configFilePath)

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

	//the following keys are deprecated, and no longer supported
	/*
		- notify.filter_attributes (replaced by metrics.status.filter_attributes SETTING)
		- notify.level (replaced by metrics.notify.level and metrics.status.threshold SETTING)
	*/
	//TODO add docs and upgrade doc.
	if c.IsSet("notify.filter_attributes") {
		return errors.ConfigValidationError("`notify.filter_attributes` configuration option is deprecated. Replaced by option in Dashboard Settings page")
	}
	if c.IsSet("notify.level") {
		return errors.ConfigValidationError("`notify.level` configuration option is deprecated. Replaced by option in Dashboard Settings page")
	}

	return nil
}
