package collector

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/analogj/scrutiny/collector/pkg/models"
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

// https://github.com/smartmontools/smartmontools/blob/b1bb7d73c37c16ddddf73906ac9e9e9ad673d481/src/smartctl.h#L16-L47
func (c *BaseCollector) LogSmartctlExitCode(exitStatus models.SmartctlExitStatus) {
	if exitStatus.HasFailCmd() {
		c.logger.Errorln("smartctl could not parse commandline")
	}
	if exitStatus.HasFailDev() {
		c.logger.Errorln("smartctl could not open device")
	}
	if exitStatus.HasFailSmart() {
		c.logger.Errorln("smartctl detected a checksum error")
	}
	if exitStatus.HasFailStatus() {
		c.logger.Errorln("smartctl detected a failing disk")
	}
	if exitStatus.HasFailAttr() {
		c.logger.Errorln("smartctl detected a disk in pre-fail")
	}
	if exitStatus.HasFailAge() {
		c.logger.Errorln("smartctl detected a disk close to failure")
	}
	if exitStatus.HasFailErr() {
		c.logger.Errorln("smartctl detected an error log with errors")
	}
	if exitStatus.HasFailLog() {
		c.logger.Errorln("smartctl detected a self test log with errors")
	}
}
