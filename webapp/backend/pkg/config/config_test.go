package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_MergeConfigMap(t *testing.T) {
	//setup
	testConfig := configuration{
		Viper: viper.New(),
	}
	testConfig.Set("user.dashboard_display", "hello")
	testConfig.SetDefault("user.layout", "hello")

	mergeSettings := map[string]interface{}{
		"user": map[string]interface{}{
			"dashboard_display": "dashboard_display",
			"layout":            "layout",
		},
	}
	//test
	err := testConfig.MergeConfigMap(mergeSettings)

	//verify
	require.NoError(t, err)

	// if using Set, the MergeConfigMap functionality will not override
	// if using SetDefault, the MergeConfigMap will override correctly
	require.Equal(t, "hello", testConfig.GetString("user.dashboard_display"))
	require.Equal(t, "layout", testConfig.GetString("user.layout"))

}
