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

func CreateSelfTestCollector(logger *logrus.Entry, apiEndpoint string) (SelfTestCollector, error) {
	apiEndpointUrl, err := url.Parse(apiEndpoint)
	if err != nil {
		return SelfTestCollector{}, err
	}

	stc := SelfTestCollector{
		apiEndpoint: apiEndpointUrl,
		logger:      logger,
	}

	return stc, nil
}

func (sc *SelfTestCollector) Run() error {
	return nil
}
