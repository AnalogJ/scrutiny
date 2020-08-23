package collector

import (
	"bytes"
	"fmt"
	"github.com/analogj/scrutiny/collector/pkg/errors"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/sirupsen/logrus"
	"net/url"
	"os/exec"
	"strings"
	"sync"
)

type MetricsCollector struct {
	BaseCollector
	apiEndpoint *url.URL
}

func CreateMetricsCollector(logger *logrus.Entry, apiEndpoint string) (MetricsCollector, error) {
	apiEndpointUrl, err := url.Parse(apiEndpoint)
	if err != nil {
		return MetricsCollector{}, err
	}

	sc := MetricsCollector{
		apiEndpoint: apiEndpointUrl,
		BaseCollector: BaseCollector{
			logger: logger,
		},
	}

	return sc, nil
}

func (mc *MetricsCollector) Run() error {
	err := mc.Validate()
	if err != nil {
		return err
	}

	apiEndpoint, _ := url.Parse(mc.apiEndpoint.String())
	apiEndpoint.Path = "/api/devices/register"

	deviceRespWrapper := new(models.DeviceWrapper)
	detectedStorageDevices, err := mc.DetectStorageDevices()
	if err != nil {
		return err
	}

	mc.logger.Infoln("Sending detected devices to API, for filtering & validation")
	err = mc.postJson(apiEndpoint.String(), models.DeviceWrapper{
		Data: detectedStorageDevices,
	}, &deviceRespWrapper)
	if err != nil {
		return err
	}

	if !deviceRespWrapper.Success {
		mc.logger.Errorln("An error occurred while retrieving filtered devices")
		return errors.ApiServerCommunicationError("An error occurred while retrieving filtered devices")
	} else {
		mc.logger.Debugln(deviceRespWrapper)
		var wg sync.WaitGroup

		for _, device := range deviceRespWrapper.Data {
			// execute collection in parallel go-routines
			wg.Add(1)
			go mc.Collect(&wg, device.WWN, device.DeviceName)
		}

		mc.logger.Infoln("Main: Waiting for workers to finish")
		wg.Wait()
		mc.logger.Infoln("Main: Completed")
	}

	return nil
}

func (mc *MetricsCollector) Validate() error {
	mc.logger.Infoln("Verifying required tools")
	_, lookErr := exec.LookPath("smartctl")

	if lookErr != nil {
		return errors.DependencyMissingError("smartctl is missing")
	}

	return nil
}

func (mc *MetricsCollector) Collect(wg *sync.WaitGroup, deviceWWN string, deviceName string) {
	defer wg.Done()
	mc.logger.Infof("Collecting smartctl results for %s\n", deviceName)

	result, err := mc.ExecCmd("smartctl", []string{"-a", "-j", fmt.Sprintf("/dev/%s", deviceName)}, "", nil)
	resultBytes := []byte(result)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// smartctl command exited with an error, we should still push the data to the API server
			mc.logger.Errorf("smartctl returned an error code (%d) while processing %s\n", exitError.ExitCode(), deviceName)
			mc.Publish(deviceWWN, resultBytes)
		} else {
			mc.logger.Errorf("error while attempting to execute smartctl: %s\n", deviceName)
			mc.logger.Errorf("ERROR MESSAGE: %v", err)
			mc.logger.Errorf("IGNORING RESULT: %v", result)
		}
		return
	} else {
		//successful run, pass the results directly to webapp backend for parsing and processing.
		mc.Publish(deviceWWN, resultBytes)
	}
}

func (mc *MetricsCollector) Publish(deviceWWN string, payload []byte) error {
	mc.logger.Infof("Publishing smartctl results for %s\n", deviceWWN)

	apiEndpoint, _ := url.Parse(mc.apiEndpoint.String())
	apiEndpoint.Path = fmt.Sprintf("/api/device/%s/smart", strings.ToLower(deviceWWN))

	resp, err := httpClient.Post(apiEndpoint.String(), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
