package collector

import (
	"bytes"
	"fmt"
	"github.com/analogj/scrutiny/collector/pkg/errors"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type MetricsCollector struct {
	BaseCollector

	apiEndpoint *url.URL
	logger      *logrus.Entry
}

func CreateMetricsCollector(logger *logrus.Entry, apiEndpoint string) (MetricsCollector, error) {
	apiEndpointUrl, err := url.Parse(apiEndpoint)
	if err != nil {
		return MetricsCollector{}, err
	}

	sc := MetricsCollector{
		apiEndpoint: apiEndpointUrl,
		logger:      logger,
	}

	return sc, nil
}

func (mc *MetricsCollector) Run() error {
	err := mc.Validate()
	if err != nil {
		return err
	}

	apiEndpoint, _ := url.Parse(mc.apiEndpoint.String())
	apiEndpoint.Path = "/api/devices"

	deviceRespWrapper := new(models.DeviceRespWrapper)

	fmt.Println("Getting devices")
	err = mc.getJson(apiEndpoint.String(), &deviceRespWrapper)
	if err != nil {
		return err
	}

	if !deviceRespWrapper.Success {
		//TODO print error payload
		fmt.Println("An error occurred while retrieving devices")
	} else {
		fmt.Println(deviceRespWrapper)
		var wg sync.WaitGroup

		for _, device := range deviceRespWrapper.Data {
			// execute collection in parallel go-routines
			wg.Add(1)
			go mc.Collect(&wg, device.WWN, device.DeviceName)
		}

		fmt.Println("Main: Waiting for workers to finish")
		wg.Wait()
		fmt.Println("Main: Completed")
	}

	return nil
}

func (mc *MetricsCollector) Validate() error {
	fmt.Println("Verifying required tools")
	_, lookErr := exec.LookPath("smartctl")

	if lookErr != nil {
		return errors.DependencyMissingError("smartctl is missing")
	}

	return nil
}

func (mc *MetricsCollector) Collect(wg *sync.WaitGroup, deviceWWN string, deviceName string) {
	defer wg.Done()
	fmt.Printf("Collecting smartctl results for %s\n", deviceName)

	result, err := mc.execCmd("smartctl", []string{"-a", "-j", fmt.Sprintf("/dev/%s", deviceName)}, "", nil)
	resultBytes := []byte(result)
	if err != nil {
		fmt.Printf("error while retrieving data from smartctl %s\n", deviceName)
		fmt.Printf("ERROR MESSAGE: %v", err)
		fmt.Printf("RESULT: %v", result)
		// TODO: error while retrieving data from smartctl.
		// TODO: we should pass this data on to scrutiny API for recording.
		return
	} else {
		//successful run, pass the results directly to webapp backend for parsing and processing.
		mc.Publish(deviceWWN, resultBytes)
	}
}

func (mc *MetricsCollector) Publish(deviceWWN string, payload []byte) error {
	fmt.Printf("Publishing smartctl results for %s\n", deviceWWN)

	apiEndpoint, _ := url.Parse(mc.apiEndpoint.String())
	apiEndpoint.Path = fmt.Sprintf("/api/device/%s/smart", strings.ToLower(deviceWWN))

	resp, err := httpClient.Post(apiEndpoint.String(), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
