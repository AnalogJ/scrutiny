package config

import (
	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
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
	c.SetDefault("web.listen.port", "8080")
	c.SetDefault("web.listen.host", "0.0.0.0")
	c.SetDefault("web.src.frontend.path", "/scrutiny/web")

	c.SetDefault("web.database.location", "/scrutiny/config/scrutiny.db")

	c.SetDefault("disks.include", []string{})
	c.SetDefault("disks.exclude", []string{})

	c.SetDefault("notify.metric.script", "/scrutiny/config/notify-metrics.sh")
	c.SetDefault("notify.long.script", "/scrutiny/config/notify-long-test.sh")
	c.SetDefault("notify.short.script", "/scrutiny/config/notify-short-test.sh")

	c.SetDefault("collect.metric.enable", true)
	c.SetDefault("collect.metric.command", "-a -o on -S on")
	c.SetDefault("collect.long.enable", true)
	c.SetDefault("collect.long.command", "-a -o on -S on")
	c.SetDefault("collect.short.enable", true)
	c.SetDefault("collect.short.command", "-a -o on -S on")

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
		log.Printf("No configuration file found at %v. Skipping", configFilePath)
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

	////deserialize Questions
	//questionsMap := map[string]Question{}
	//err := c.UnmarshalKey("questions", &questionsMap)
	//
	//if err != nil {
	//	log.Printf("questions could not be deserialized correctly. %v", err)
	//	return err
	//}
	//
	//for _, v := range questionsMap {
	//
	//	typeContent, ok := v.Schema["type"].(string)
	//	if !ok || len(typeContent) == 0 {
	//		return errors.QuestionSyntaxError("`type` is required for questions")
	//	}
	//}
	//
	//

	return nil
}
