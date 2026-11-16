package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
	"github.com/go-viper/mapstructure/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/spf13/viper"
)

const DB_USER_SETTINGS_SUBKEY = "user"

// DeviceConfig is an entry in the `devices` list
type DeviceConfig struct {
	ScrutinyUUID  uuid.UUID `mapstructure:"scrutiny_uuid"`
	Notifications struct {
		// nil means not set, which defaults to the top level `notify` key. Pointers are needed because viper
		// defaults are per key and can't reach into list elements, so they couldn't tell unset from false or 0.
		OnMissedUpload  *bool `mapstructure:"on_missed_upload"`
		OnSmartctlError *bool `mapstructure:"on_smartctl_error"`
		// 0 means estimate it from the device's uploads
		UploadPeriod *time.Duration `mapstructure:"upload_period"`
		// no notifications of any kind are sent for the device
		Disabled bool `mapstructure:"disabled"`
	} `mapstructure:"notifications"`
}

func (c *configuration) GetDevices() ([]DeviceConfig, error) {
	var devices []DeviceConfig
	// keep viper's default hooks, which parse durations
	err := c.UnmarshalKey("devices", &devices, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.TextUnmarshallerHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)))
	return devices, err
}

func (c *configuration) GetDeviceConfig(scrutinyUUID uuid.UUID) *DeviceConfig {
	// ValidateConfig already rejected invalid devices
	devices, _ := c.GetDevices()
	for i := range devices {
		if devices[i].ScrutinyUUID == scrutinyUUID {
			return &devices[i]
		}
	}
	return nil
}

// When initializing this class the following methods must be called:
// Config.New
// Config.Init
// This is done automatically when created via the Factory.
type configuration struct {
	*viper.Viper

	// viper isn't safe for concurrent use, and request handlers read and write the config at runtime
	// (e.g. GET /api/settings calls SetDefault). Methods taking the lock must not call each other.
	mu sync.RWMutex
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
	c.SetDefault("notify.on_smartctl_error", true)
	c.SetDefault("notify.on_missed_upload", true)
	c.SetDefault("notify.upload_period", "0s")

	c.SetDefault("web.influxdb.scheme", "http")
	c.SetDefault("web.influxdb.host", "localhost")
	c.SetDefault("web.influxdb.port", "8086")
	c.SetDefault("web.influxdb.org", "scrutiny")
	c.SetDefault("web.influxdb.bucket", "metrics")
	c.SetDefault("web.influxdb.init_username", "admin")
	c.SetDefault("web.influxdb.init_password", "password12345")
	c.SetDefault("web.influxdb.token", "scrutiny-default-admin-token")
	c.SetDefault("web.influxdb.tls.insecure_skip_verify", false)
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
	// the sub config shares maps with this one, so hold the lock while reading it
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.Sub(key).AllKeys()
}

func (c *configuration) Sub(key string) Interface {
	c.mu.RLock()
	defer c.mu.RUnlock()
	config := configuration{
		Viper: c.Viper.Sub(key),
	}
	return &config
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

	// GetDuration would silently turn an invalid duration into 0
	if period, err := time.ParseDuration(c.GetString("notify.upload_period")); err != nil {
		return errors.ConfigValidationError(fmt.Sprintf("`notify.upload_period` is invalid: %v", err))
	} else if period < 0 {
		return errors.ConfigValidationError("`notify.upload_period` is negative")
	}

	devices, err := c.GetDevices()
	if err != nil {
		return errors.ConfigValidationError(fmt.Sprintf("`devices` is invalid: %v", err))
	}
	for _, device := range devices {
		if device.ScrutinyUUID.IsNil() {
			return errors.ConfigValidationError("`scrutiny_uuid` is missing from a device")
		}
		if device.Notifications.UploadPeriod != nil && *device.Notifications.UploadPeriod < 0 {
			return errors.ConfigValidationError(fmt.Sprintf("`upload_period` of device %s is negative", device.ScrutinyUUID))
		}
	}

	return nil
}

// The methods below wrap viper's so they can be used concurrently.

func (c *configuration) MergeConfig(in io.Reader) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Viper.MergeConfig(in)
}

func (c *configuration) MergeConfigMap(cfg map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Viper.MergeConfigMap(cfg)
}

func (c *configuration) WriteConfig() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.WriteConfig()
}

func (c *configuration) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Viper.Set(key, value)
}

func (c *configuration) SetDefault(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Viper.SetDefault(key, value)
}

func (c *configuration) AllSettings() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.AllSettings()
}

func (c *configuration) AllKeys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.AllKeys()
}

func (c *configuration) IsSet(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.IsSet(key)
}

func (c *configuration) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.Get(key)
}

func (c *configuration) GetBool(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.GetBool(key)
}

func (c *configuration) GetInt(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.GetInt(key)
}

func (c *configuration) GetInt64(key string) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.GetInt64(key)
}

func (c *configuration) GetDuration(key string) time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.GetDuration(key)
}

func (c *configuration) GetString(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.GetString(key)
}

func (c *configuration) GetStringSlice(key string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.GetStringSlice(key)
}

func (c *configuration) UnmarshalKey(key string, rawVal interface{}, decoderOpts ...viper.DecoderConfigOption) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Viper.UnmarshalKey(key, rawVal, decoderOpts...)
}
