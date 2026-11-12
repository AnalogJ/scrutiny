package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestSendTestNotification_DoesNotLeakNotifyURLSecrets(t *testing.T) {
	// fake credentials, sent to a closed local port
	appConfig, err := config.Create()
	require.NoError(t, err)
	appConfig.Set("notify.urls", []string{"ntfy://fakeuser:FAKEntfyPASSWORD@127.0.0.1:1/faketopic"})

	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/health/notify", func(c *gin.Context) {
		c.Set("LOGGER", logrus.NewEntry(logger))
		c.Set("CONFIG", appConfig)
		SendTestNotification(c)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/health/notify", nil))

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "notification 1/1 via ntfy")

	var logs []string
	for _, entry := range hook.AllEntries() {
		line, err := entry.String()
		require.NoError(t, err)
		logs = append(logs, line)
	}
	for _, secret := range []string{"fakeuser", "FAKEntfyPASSWORD", "faketopic"} {
		require.NotContains(t, recorder.Body.String(), secret)
		require.NotContains(t, strings.Join(logs, "\n"), secret)
	}
}
