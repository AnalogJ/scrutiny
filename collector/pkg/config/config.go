package config

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/collector/pkg/errors"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// When initializing this class the following methods must be called:
// Config.New
// Config.Init
// This is done automatically when created via the Factory.
type configuration struct {
	*viper.Viper

	deviceOverrides []models.ScanOverride
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

	c.SetDefault("commands.metrics_smartctl_bin", "smartctl")
	c.SetDefault("commands.metrics_scan_args", "--scan --json")
	c.SetDefault("commands.metrics_info_args", "--info --json")
	c.SetDefault("commands.metrics_smart_args", "--xall --json")
	c.SetDefault("commands.metrics_smartctl_wait", 0)

	//configure env variable parsing.
	c.SetEnvPrefix("COLLECTOR")
	c.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	c.AutomaticEnv()

	//c.SetDefault("collect.short.command", "-a -o on -S on")

	c.SetDefault("allow_listed_devices", []string{})

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

	// check that the collector commands are valid
	commandArgStrings := map[string]string{
		"commands.metrics_scan_args":  c.GetString("commands.metrics_scan_args"),
		"commands.metrics_info_args":  c.GetString("commands.metrics_info_args"),
		"commands.metrics_smart_args": c.GetString("commands.metrics_smart_args"),
	}

	errorStrings := []string{}
	for configKey, commandArgString := range commandArgStrings {
		args := strings.Split(commandArgString, " ")
		//ensure that the args string contains `--json` or `-j` flag
		containsJsonFlag := false
		containsDeviceFlag := false
		for _, flag := range args {
			if strings.HasPrefix(flag, "--json") || strings.HasPrefix(flag, "-j") {
				containsJsonFlag = true
			}
			if strings.HasPrefix(flag, "--device") || strings.HasPrefix(flag, "-d") {
				containsDeviceFlag = true
			}
		}

		if !containsJsonFlag {
			errorStrings = append(errorStrings, fmt.Sprintf("configuration key '%s' is missing '--json' flag", configKey))
		}

		if containsDeviceFlag {
			errorStrings = append(errorStrings, fmt.Sprintf("configuration key '%s' must not contain '--device' or '-d' flag", configKey))
		}
	}
	//sort(errorStrings)
	sort.Strings(errorStrings)

	if len(errorStrings) == 0 {
		return nil
	} else {
		return errors.ConfigValidationError(strings.Join(errorStrings, ", "))
	}
}

func (c *configuration) GetDeviceOverrides() []models.ScanOverride {
	// we have to support 2 types of device types.
	// - simple device type (device_type: 'sat')
	// and list of device types (type: \n- 3ware,0 \n- 3ware,1 \n- 3ware,2)
	// GetString will return "" if this is a list of device types.

	if c.deviceOverrides == nil {
		overrides := []models.ScanOverride{}
		c.UnmarshalKey("devices", &overrides, func(c *mapstructure.DecoderConfig) { c.WeaklyTypedInput = true })
		c.deviceOverrides = overrides
	}

	return c.deviceOverrides
}

func (c *configuration) GetCommandMetricsInfoArgs(deviceName string) string {
	overrides := c.GetDeviceOverrides()

	for _, deviceOverrides := range overrides {
		if strings.EqualFold(deviceName, deviceOverrides.Device) {
			//found matching device
			if len(deviceOverrides.Commands.MetricsInfoArgs) > 0 {
				return deviceOverrides.Commands.MetricsInfoArgs
			} else {
				return c.GetString("commands.metrics_info_args")
			}
		}
	}
	return c.GetString("commands.metrics_info_args")
}

func (c *configuration) GetCommandMetricsSmartArgs(deviceName string) string {
	overrides := c.GetDeviceOverrides()

	for _, deviceOverrides := range overrides {
		if strings.EqualFold(deviceName, deviceOverrides.Device) {
			//found matching device
			if len(deviceOverrides.Commands.MetricsSmartArgs) > 0 {
				return deviceOverrides.Commands.MetricsSmartArgs
			} else {
				return c.GetString("commands.metrics_smart_args")
			}
		}
	}
	return c.GetString("commands.metrics_smart_args")
}

func (c *configuration) IsAllowlistedDevice(deviceName string) bool {
	allowList := c.GetStringSlice("allow_listed_devices")
	if len(allowList) == 0 {
		return true
	}

	for _, item := range allowList {
		if item == deviceName {
			return true
		}
	}

	return false
}
