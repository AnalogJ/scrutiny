package collector

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var httpClient = &http.Client{Timeout: 60 * time.Second}

type BaseCollector struct {
	logger *logrus.Entry
}

func (c *BaseCollector) getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (c *BaseCollector) postJson(url string, body, target interface{}) error {
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

// http://www.linuxguide.it/command_line/linux-manpage/do.php?file=smartctl#sect7
func (c *BaseCollector) LogSmartctlExitCode(exitCode int) {
	if exitCode&0x01 != 0 {
		c.logger.Errorln("smartctl could not parse commandline")
	} else if exitCode&0x02 != 0 {
		c.logger.Errorln("smartctl could not open device")
	} else if exitCode&0x04 != 0 {
		c.logger.Errorln("smartctl detected a checksum error")
	} else if exitCode&0x08 != 0 {
		c.logger.Errorln("smartctl detected a failing disk ")
	} else if exitCode&0x10 != 0 {
		c.logger.Errorln("smartctl detected a disk in pre-fail")
	} else if exitCode&0x20 != 0 {
		c.logger.Errorln("smartctl detected a disk close to failure")
	} else if exitCode&0x40 != 0 {
		c.logger.Errorln("smartctl detected a error log with errors")
	} else if exitCode&0x80 != 0 {
		c.logger.Errorln("smartctl detected a self test log with errors")
	}
}
