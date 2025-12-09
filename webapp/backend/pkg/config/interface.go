package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Create mock using:
// mockgen -source=webapp/backend/pkg/config/interface.go -destination=webapp/backend/pkg/config/mock/mock_config.go
type Interface interface {
	Init() error
	ReadConfig(configFilePath string, logger *logrus.Entry) error
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
	GetString(key string) string
	GetStringSlice(key string) []string
	GetIntSlice(key string) []int
	UnmarshalKey(key string, rawVal interface{}, decoderOpts ...viper.DecoderConfigOption) error
}
