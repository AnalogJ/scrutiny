package collector

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApiEndpointParse(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost:8080/")

	url1, _ := baseURL.Parse("d/e")
	require.Equal(t, "http://localhost:8080/d/e", url1.String())

	url2, _ := baseURL.Parse("/d/e")
	require.Equal(t, "http://localhost:8080/d/e", url2.String())
}

func TestApiEndpointParse_WithBasepathWithoutTrailingSlash(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost:8080/scrutiny")

	// This testcase is unexpected and can cause issues. We need to ensure the apiEndpoint always has a trailing slash.
	url1, _ := baseURL.Parse("d/e")
	require.Equal(t, "http://localhost:8080/d/e", url1.String())

	url2, _ := baseURL.Parse("/d/e")
	require.Equal(t, "http://localhost:8080/d/e", url2.String())
}

func TestApiEndpointParse_WithBasepathWithTrailingSlash(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost:8080/scrutiny/")

	url1, _ := baseURL.Parse("d/e")
	require.Equal(t, "http://localhost:8080/scrutiny/d/e", url1.String())

	url2, _ := baseURL.Parse("/d/e")
	require.Equal(t, "http://localhost:8080/d/e", url2.String())
}
