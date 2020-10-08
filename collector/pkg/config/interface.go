package config

import (
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/spf13/viper"
)

// Create mock using:
// mockgen -source=collector/pkg/config/interface.go -destination=collector/pkg/config/mock/mock_config.go
type Interface interface {
	Init() error
	ReadConfig(configFilePath string) error
	Set(key string, value interface{})
	SetDefault(key string, value interface{})

	AllSettings() map[string]interface{}
	IsSet(key string) bool
	Get(key string) interface{}
	GetBool(key string) bool
	GetInt(key string) int
	GetString(key string) string
	GetStringSlice(key string) []string
	UnmarshalKey(key string, rawVal interface{}, decoderOpts ...viper.DecoderConfigOption) error

	GetScanOverrides() []models.ScanOverride
}
