package main

import (
	"encoding/json"
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestLoggableSettings_RemovesSecrets(t *testing.T) {
	// fake credentials
	appConfig, err := config.Create()
	require.NoError(t, err)
	appConfig.Set("web.influxdb.token", "FAKEinfluxTOKEN")
	appConfig.Set("web.influxdb.init_password", "FAKEinfluxPASSWORD")
	appConfig.Set("notify.urls", []string{"ntfy://:FAKEntfyTOKEN@ntfy.example.invalid/topic"})

	settings := loggableSettings(appConfig)

	settingsData, err := json.Marshal(settings)
	require.NoError(t, err)
	for _, secret := range []string{"FAKEinfluxTOKEN", "FAKEinfluxPASSWORD", "FAKEntfyTOKEN"} {
		require.NotContains(t, string(settingsData), secret)
	}

	influxdbSettings := settings["web"].(map[string]interface{})["influxdb"].(map[string]interface{})
	require.Equal(t, "[redacted]", influxdbSettings["token"])
	require.Equal(t, "[redacted]", influxdbSettings["init_password"])
	require.Equal(t, "localhost", influxdbSettings["host"])
	require.Equal(t, []string{"notification 1/1 via ntfy"}, settings["notify"].(map[string]interface{})["urls"])

	// the config itself keeps the secrets
	require.Equal(t, "FAKEinfluxTOKEN", appConfig.GetString("web.influxdb.token"))
	require.Equal(t, "FAKEinfluxPASSWORD", appConfig.GetString("web.influxdb.init_password"))
	require.Equal(t, []string{"ntfy://:FAKEntfyTOKEN@ntfy.example.invalid/topic"}, appConfig.GetStringSlice("notify.urls"))
}

func TestLoggableSettings_RemovesSecretsFromEnvironment(t *testing.T) {
	t.Setenv("SCRUTINY_WEB_INFLUXDB_TOKEN", "FAKEenvTOKEN")
	appConfig, err := config.Create()
	require.NoError(t, err)
	require.Equal(t, "FAKEenvTOKEN", appConfig.GetString("web.influxdb.token"))

	settingsData, err := json.Marshal(loggableSettings(appConfig))
	require.NoError(t, err)
	require.NotContains(t, string(settingsData), "FAKEenvTOKEN")
}
