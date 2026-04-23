package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/analogj/scrutiny/collector/pkg/common/shell"
	"github.com/analogj/scrutiny/collector/pkg/config"
	"github.com/analogj/scrutiny/collector/pkg/detect"
	"github.com/analogj/scrutiny/collector/pkg/errors"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
)

type MetricsCollector struct {
	config config.Interface
	BaseCollector
	apiEndpoint *url.URL
	shell       shell.Interface
}

func CreateMetricsCollector(appConfig config.Interface, logger *logrus.Entry, apiEndpoint string) (MetricsCollector, error) {
	apiEndpointUrl, err := url.Parse(apiEndpoint)
	if err != nil {
		return MetricsCollector{}, err
	}

	sc := MetricsCollector{
		config:      appConfig,
		apiEndpoint: apiEndpointUrl,
		BaseCollector: BaseCollector{
			logger: logger,
		},
		shell: shell.Create(),
	}

	return sc, nil
}

func (mc *MetricsCollector) Run() error {
	err := mc.Validate()
	if err != nil {
		return err
	}

	apiEndpoint, _ := url.Parse(mc.apiEndpoint.String())
	apiEndpoint, _ = apiEndpoint.Parse("api/devices/register") //this acts like filepath.Join()

	deviceRespWrapper := new(models.DeviceWrapper)

	deviceDetector := detect.Detect{
		Logger: mc.logger,
		Config: mc.config,
	}
	scannedDevices, err := deviceDetector.Scan()
	if err != nil {
		mc.reportScanError(err)
		return err
	}

	rawDetectedStorageDevices, err := deviceDetector.Info(scannedDevices)
	if err != nil {
		mc.logger.Errorf("An error occurred while retrieving device info: %v", err)
		// report errors for devices that failed to get info (no ScrutinyUUID)
		for _, device := range rawDetectedStorageDevices {
			if device.ScrutinyUUID.IsNil() {
				mc.reportDeviceError(device, err)
			}
		}
		// continue with devices that succeeded
	}

	// Ignore any device without a Scrutiny UUID. This should never happen...
	detectedStorageDevices := make([]models.Device, 0, len(rawDetectedStorageDevices))
	for _, device := range rawDetectedStorageDevices {
		if device.ScrutinyUUID.IsNil() {
			mc.logger.Errorf("Device %s has no scrutiny UUID; skipping (no data association possible).", device.DeviceName)
			mc.logger.Debugf("Raw detected device: model=%q serial=%q wwn=%q ScrutinyUUID=%s",
				device.ModelName, device.SerialNumber, device.WWN, device.ScrutinyUUID)
			continue
		}
		detectedStorageDevices = append(detectedStorageDevices, device)
	}

	mc.logger.Infof("Sending %d/%d detected devices to API for filtering & validation",
		len(detectedStorageDevices), len(rawDetectedStorageDevices))
	jsonObj, _ := json.Marshal(detectedStorageDevices)
	mc.logger.Debugf("Detected devices: %v", string(jsonObj))
	err = mc.postJson(apiEndpoint.String(), models.DeviceWrapper{
		Data: detectedStorageDevices,
	}, &deviceRespWrapper)
	if err != nil {
		return err
	}

	if !deviceRespWrapper.Success {
		mc.logger.Errorln("An error occurred while retrieving filtered devices")
		mc.logger.Debugln(deviceRespWrapper)
		return errors.ApiServerCommunicationError("An error occurred while retrieving filtered devices")
	} else {
		mc.logger.Debugln(deviceRespWrapper)
		//var wg sync.WaitGroup
		for _, device := range deviceRespWrapper.Data {
			// execute collection in parallel go-routines
			//wg.Add(1)
			//go mc.Collect(&wg, device.WWN, device.DeviceName, device.DeviceType)
			mc.Collect(device.ScrutinyUUID, device.DeviceName, device.DeviceType)

			if mc.config.GetInt("commands.metrics_smartctl_wait") > 0 {
				time.Sleep(time.Duration(mc.config.GetInt("commands.metrics_smartctl_wait")) * time.Second)
			}
		}

		//mc.logger.Infoln("Main: Waiting for workers to finish")
		//wg.Wait()
		mc.logger.Infoln("Main: Completed")
	}

	return nil
}

func (mc *MetricsCollector) Validate() error {
	mc.logger.Infoln("Verifying required tools")
	_, lookErr := exec.LookPath(mc.config.GetString("commands.metrics_smartctl_bin"))

	if lookErr != nil {
		return errors.DependencyMissingError(fmt.Sprintf("%s binary is missing", mc.config.GetString("commands.metrics_smartctl_bin")))
	}

	return nil
}

// func (mc *MetricsCollector) Collect(wg *sync.WaitGroup, deviceWWN string, deviceName string, deviceType string) {
func (mc *MetricsCollector) Collect(scrutiny_uuid uuid.UUID, deviceName string, deviceType string) {
	//defer wg.Done()
	// Run() filters out devices with nil ScrutinyUUIDs before calling Collect, so this should never
	// happen; guarded here in case Collect is called from elsewhere in the future.
	if scrutiny_uuid.IsNil() {
		mc.logger.Errorf("Device %s has no scrutiny UUID; skipping collection (no data association possible).", deviceName)
		return
	}
	mc.logger.Infof("Collecting smartctl results for %s\n", deviceName)

	fullDeviceName := fmt.Sprintf("%s%s", detect.DevicePrefix(), deviceName)
	args := strings.Split(mc.config.GetCommandMetricsSmartArgs(fullDeviceName), " ")
	//only include the device type if its a non-standard one, or the user set it in the config file.
	//In some cases ata drives are detected as scsi in docker, and metadata is lost, so a scanned scsi/ata type is dropped.
	if len(deviceType) > 0 &&
		((deviceType != "scsi" && deviceType != "ata") || mc.config.HasDeviceTypeOverride(fullDeviceName)) {
		args = append(args, "--device", deviceType)
	}
	args = append(args, fullDeviceName)

	result, err := mc.shell.Command(mc.logger, mc.config.GetString("commands.metrics_smartctl_bin"), args, "", os.Environ())
	resultBytes := []byte(result)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitStatus := models.NewSmartctlExitStatus(exitError.ExitCode())
			mc.logger.Errorf("smartctl returned an error code (%d) while processing %s\n", exitError.ExitCode(), deviceName)
			mc.LogSmartctlExitCode(exitStatus)

			// check for fatal errors that mean we should not publish
			if exitStatus.HasFailCmd() || exitStatus.HasFailSmart() {
				mc.logger.Errorf("Fatal smartctl error for %s, not publishing results\n", deviceName)
				mc.reportDeviceError(models.Device{ScrutinyUUID: scrutiny_uuid, DeviceName: deviceName}, err)
				return
			}

			if exitStatus.HasFailDev() {
				// if the only issue is FailDev and we passed "-n standby", the drive is just in standby — not a real error
				hasStandbyArg := false
				for i, arg := range args {
					if arg == "-n" && i+1 < len(args) && args[i+1] == "standby" {
						hasStandbyArg = true
						break
					}
				}
				if !hasStandbyArg {
					mc.logger.Errorf("Fatal smartctl error (device open failed) for %s, not publishing results\n", deviceName)
					mc.reportDeviceError(models.Device{ScrutinyUUID: scrutiny_uuid, DeviceName: deviceName}, err)
					return
				}
			}

			mc.Publish(scrutiny_uuid, resultBytes)
		} else {
			mc.logger.Errorf("error while attempting to execute smartctl: %s\n", deviceName)
			mc.logger.Errorf("ERROR MESSAGE: %v", err)
			mc.logger.Errorf("IGNORING RESULT: %v", result)
		}
		return
	} else {
		//successful run, pass the results directly to webapp backend for parsing and processing.
		mc.Publish(scrutiny_uuid, resultBytes)
	}
}

func (mc *MetricsCollector) Publish(scrutinyUuid uuid.UUID, payload []byte) error {
	mc.logger.Infof("Publishing smartctl results for %s\n", scrutinyUuid)

	apiEndpoint, _ := url.Parse(mc.apiEndpoint.String())
	apiEndpoint, _ = apiEndpoint.Parse(fmt.Sprintf("api/device/%s/smart", scrutinyUuid.String()))

	resp, err := httpClient.Post(apiEndpoint.String(), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		mc.logger.Errorf("An error occurred while publishing SMART data for device (%s): %v", scrutinyUuid, err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (mc *MetricsCollector) reportScanError(scanErr error) {
	mc.logger.Debugf("Reporting scan error to API: %v", scanErr)

	apiEndpoint, _ := url.Parse(mc.apiEndpoint.String())
	apiEndpoint, _ = apiEndpoint.Parse("api/collector_scan_error")

	payload := collector.CollectorError{
		Error:  scanErr.Error(),
		HostId: mc.config.GetString("host.id"),
	}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		mc.logger.Errorf("An error occurred while marshalling scan error payload: %v", err)
		return
	}

	resp, err := httpClient.Post(apiEndpoint.String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		mc.logger.Errorf("An error occurred while reporting scan error to API: %v", err)
		return
	}
	defer resp.Body.Close()
}

func (mc *MetricsCollector) reportDeviceError(device models.Device, deviceErr error) {
	scrutinyUuid := device.ScrutinyUUID
	if scrutinyUuid.IsNil() {
		scrutinyUuid = detect.GenerateScrutinyUUID("", "", device.DeviceName)
	}

	mc.logger.Debugf("Reporting device error for %s (%s) to API: %v", device.DeviceName, scrutinyUuid, deviceErr)

	apiEndpoint, _ := url.Parse(mc.apiEndpoint.String())
	apiEndpoint, _ = apiEndpoint.Parse(fmt.Sprintf("api/device/%s/collector_error", scrutinyUuid.String()))

	payload := collector.CollectorError{
		Error:  deviceErr.Error(),
		HostId: mc.config.GetString("host.id"),
	}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		mc.logger.Errorf("An error occurred while marshalling device error payload: %v", err)
		return
	}

	resp, err := httpClient.Post(apiEndpoint.String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		mc.logger.Errorf("An error occurred while reporting device error to API: %v", err)
		return
	}
	defer resp.Body.Close()
}
