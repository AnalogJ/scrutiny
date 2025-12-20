package collector

import (
	"net/url"

	"github.com/sirupsen/logrus"
)

type SelfTestCollector struct {
	BaseCollector

	apiEndpoint *url.URL
	logger      *logrus.Entry
}

// CreateSelfTestCollector creates a new SelfTestCollector with a default 60-second timeout
// TODO: accept config.Interface to use configurable timeout like MetricsCollector
func CreateSelfTestCollector(logger *logrus.Entry, apiEndpoint string) (SelfTestCollector, error) {
	apiEndpointUrl, err := url.Parse(apiEndpoint)
	if err != nil {
		return SelfTestCollector{}, err
	}

	stc := SelfTestCollector{
		BaseCollector: BaseCollector{
			logger:     logger,
			httpClient: NewHTTPClient(60), // Default timeout, will use config when refactored
		},
		apiEndpoint: apiEndpointUrl,
		logger:      logger,
	}

	return stc, nil
}

func (sc *SelfTestCollector) Run() error {
	return nil
}
