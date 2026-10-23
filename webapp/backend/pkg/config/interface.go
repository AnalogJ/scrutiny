package config

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/spf13/viper"
)

// Create mock using:
// mockgen -source=webapp/backend/pkg/config/interface.go -destination=webapp/backend/pkg/config/mock/mock_config.go
type Interface interface {
	Init() error
	ReadConfig(configFilePath string) error
	WriteConfig() error
	Set(key string, value interface{})
	SetDefault(key string, value interface{})
	MergeConfigMap(cfg map[string]interface{}) error

	Sub(key string) Interface
	AllSettings() map[string]interface{}
	AllKeys() []string
	SubKeys(key string) []string
	IsSet(key string) bool
	Get(key string) interface{}
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	GetDuration(key string) time.Duration
	GetString(key string) string
	GetStringSlice(key string) []string
	UnmarshalKey(key string, rawVal interface{}, decoderOpts ...viper.DecoderConfigOption) error

	GetDevices() ([]DeviceConfig, error)
	// GetDeviceConfig returns the `devices` entry for the scrutiny UUID, or nil if there is none
	GetDeviceConfig(scrutinyUUID uuid.UUID) *DeviceConfig
}
