package notify

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

// fails every request with an error that echoes the request url, like a careless service could. Requests to
// 127.0.0.1 go to the real transport, to get Go's own errors from a closed local port.
type urlEchoingTransport struct {
	original http.RoundTripper
	requests atomic.Int64
}

func (s *urlEchoingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Hostname() == "127.0.0.1" {
		return s.original.RoundTrip(req)
	}
	s.requests.Add(1)
	return nil, fmt.Errorf("stub transport refused %s", req.URL.String())
}

// stubDefaultTransport keeps the discord, telegram and generic requests (which use http.DefaultTransport) off the network
func stubDefaultTransport(t *testing.T) *urlEchoingTransport {
	original := http.DefaultTransport
	stub := &urlEchoingTransport{original: original}
	http.DefaultTransport = stub
	t.Cleanup(func() { http.DefaultTransport = original })
	return stub
}

// all tokens are fakes. ntfy (which uses its own transport) and the http webhook target a closed local port.
var notifyURLSecretCases = []struct {
	name    string
	url     string
	scheme  string
	secrets []string
	// sent through the stubbed http.DefaultTransport
	stubbed bool
}{
	{
		name:    "ntfy",
		url:     "ntfy://fakeuser:FAKEntfyPASSWORD@127.0.0.1:1/faketopic?priority=high",
		scheme:  "ntfy",
		secrets: []string{"fakeuser", "FAKEntfyPASSWORD", "faketopic"},
	},
	{
		name:    "discord",
		url:     "discord://FAKEdiscordTOKEN@987654321",
		scheme:  "discord",
		secrets: []string{"FAKEdiscordTOKEN", "987654321"},
		stubbed: true,
	},
	{
		name:    "telegram",
		url:     "telegram://123456789:FAKEtelegramTOKEN@telegram?chats=@fakechannel",
		scheme:  "telegram",
		secrets: []string{"123456789", "FAKEtelegramTOKEN"},
		stubbed: true,
	},
	{
		// shoutrrr before v0.22.0 echoed invalid telegram tokens in its error
		name:    "telegram invalid token",
		url:     "telegram://FAKEinvalidTOKEN@telegram?chats=@fakechannel",
		scheme:  "telegram",
		secrets: []string{"FAKEinvalidTOKEN"},
	},
	{
		name:    "shoutrrr generic webhook",
		url:     "generic://hooks.example.invalid/notify?token=FAKEgenericTOKEN",
		scheme:  "generic",
		secrets: []string{"FAKEgenericTOKEN"},
		stubbed: true,
	},
	{
		name: "http webhook",
		// Go's http client masks the password in its errors, but not the user, path or query
		url:     "https://fakeuser:FAKEwebhookPASSWORD@127.0.0.1:1/notify/FAKEpathTOKEN?token=FAKEqueryTOKEN",
		scheme:  "https",
		secrets: []string{"fakeuser", "FAKEwebhookPASSWORD", "FAKEpathTOKEN", "FAKEqueryTOKEN"},
	},
}

func requireNoSecrets(t *testing.T, text string, secrets []string) {
	t.Helper()
	for _, secret := range secrets {
		require.NotContains(t, text, secret)
	}
}

func logOutput(hook *test.Hook) string {
	var lines []string
	for _, entry := range hook.AllEntries() {
		line, err := entry.String()
		if err != nil {
			line = entry.Message
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func newRedactTestNotify(t *testing.T, urls []string) (Notify, *test.Hook) {
	appConfig, err := config.Create()
	require.NoError(t, err)
	appConfig.Set("notify.urls", urls)

	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)
	return New(logger, appConfig, models.Device{DeviceName: "/dev/sda", SerialNumber: "FAKESERIAL"}, true), hook
}

func TestNotifyURLLabel(t *testing.T) {
	for _, tc := range notifyURLSecretCases {
		t.Run(tc.name, func(t *testing.T) {
			label := NotifyURLLabel(1, 3, tc.url)
			require.Equal(t, "notification 2/3 via "+tc.scheme, label)
		})
	}

	require.Equal(t, "notification 1/1 via unknown", NotifyURLLabel(0, 1, "FAKEnotAURL"))
	require.Equal(t, "notification 1/1 via unknown", NotifyURLLabel(0, 1, "not a scheme://FAKEtoken"))
	require.Equal(t, "notification 1/1 via generic+https", NotifyURLLabel(0, 1, "generic+https://host/?token=FAKEtoken"))
}

func TestSend_DoesNotLeakNotifyURLSecrets(t *testing.T) {
	stub := stubDefaultTransport(t)

	for _, tc := range notifyURLSecretCases {
		t.Run(tc.name, func(t *testing.T) {
			notify, hook := newRedactTestNotify(t, []string{tc.url})
			requestsBefore := stub.requests.Load()

			err := notify.Send()

			if tc.stubbed {
				require.Greater(t, stub.requests.Load(), requestsBefore)
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), "notification 1/1 via "+tc.scheme)
			requireNoSecrets(t, err.Error(), tc.secrets)

			output := logOutput(hook)
			require.Contains(t, output, "Sending notification 1/1 via "+tc.scheme)
			requireNoSecrets(t, output, tc.secrets)
		})
	}
}

func TestSend_LogsOnlyLabelsOfConfiguredURLs(t *testing.T) {
	stubDefaultTransport(t)

	var urls, secrets []string
	for _, tc := range notifyURLSecretCases {
		urls = append(urls, tc.url)
		secrets = append(secrets, tc.secrets...)
	}
	notify, hook := newRedactTestNotify(t, urls)

	err := notify.Send()

	require.Error(t, err)
	requireNoSecrets(t, err.Error(), secrets)
	output := logOutput(hook)
	require.Contains(t, output, "notification 1/6 via ntfy notification 2/6 via discord")
	requireNoSecrets(t, output, secrets)
}
