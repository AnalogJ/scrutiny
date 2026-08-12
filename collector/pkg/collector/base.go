package collector

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/sirupsen/logrus"
)

var httpClient = &http.Client{Timeout: 60 * time.Second}

type BaseCollector struct {
	logger *logrus.Entry
}

func (c *BaseCollector) postJson(url string, body interface{}, target interface{}) error {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	r, err := httpClient.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (c *BaseCollector) LogSmartctlExitStatus(exitStatus collector.SmartctlExitStatus) {
	for _, description := range exitStatus.Descriptions() {
		c.logger.Errorln(description)
	}
}
