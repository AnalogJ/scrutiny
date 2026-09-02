package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
		mc.ReportError(models.Device{}, err)
		return err
	}

	rawDetectedStorageDevices, infoErrors := deviceDetector.Info(scannedDevices)
	for _, infoError := range infoErrors {
		mc.ReportError(infoError.Device, infoError.Err)
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

	device := models.Device{ScrutinyUUID: scrutiny_uuid, DeviceName: deviceName, DeviceType: deviceType}

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
		exitError, isExitError := err.(*exec.ExitError)
		if !isExitError {
			mc.logger.Errorf("error while attempting to execute smartctl: %s\n", deviceName)
			mc.logger.Errorf("ERROR MESSAGE: %v", err)
			mc.logger.Errorf("IGNORING RESULT: %v", result)
			mc.ReportError(device, err)
			return
		}

		// a process killed by a signal has no exit status to interpret; ExitCode() reports -1,
		// which would otherwise decode as every failure bit being set.
		if exitError.ExitCode() < 0 {
			mc.logger.Errorf("smartctl was terminated before it could finish processing %s: %v", deviceName, exitError)
			mc.ReportError(device, exitError)
			return
		}

		exitStatus := collector.SmartctlExitStatus(exitError.ExitCode())
		mc.logger.Errorf("smartctl returned an error code (%d) while processing %s\n", exitError.ExitCode(), deviceName)
		mc.LogSmartctlExitStatus(exitStatus)

		if exitStatus.IsFatal() {
			if isLowPowerExit(resultBytes) {
				mc.logger.Infof("%s is in a low power mode, skipping collection\n", deviceName)
				return
			}

			mc.logger.Errorf("smartctl output for %s is incomplete, not publishing results\n", deviceName)
			mc.ReportError(device, fmt.Errorf("smartctl exited with status %d: %s", exitError.ExitCode(), exitStatus))
			return
		}
		// the remaining exit status bits describe problems with the disk rather than with smartctl,
		// so the results are still worth publishing.
	}

	mc.Publish(scrutiny_uuid, resultBytes)
}

// isLowPowerExit checks smartctl's own output rather than the configured arguments: `-n`/`--nocheck`
// being present says nothing about whether this particular device was asleep, and a device that has
// genuinely gone away shares the same exit status.
func isLowPowerExit(result []byte) bool {
	var smartData collector.SmartInfo
	if err := json.Unmarshal(result, &smartData); err != nil {
		return false
	}
	return smartData.IsLowPowerExit()
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

	//the server rejects results it cannot store, and the collector would otherwise report success
	if resp.StatusCode != http.StatusOK {
		mc.logger.Errorf("The API server rejected the SMART data for device (%s): %s", scrutinyUuid, resp.Status)
		return errors.ApiServerCommunicationError(fmt.Sprintf("The API server rejected the SMART data for device (%s)", scrutinyUuid))
	}

	return nil
}

// ReportError tells the API server that the collector could not gather data, so that it can notify
// the user. device may be empty when the failure was not specific to one device.
func (mc *MetricsCollector) ReportError(device models.Device, collectorErr error) {
	payload := collector.CollectorError{
		HostId:     mc.config.GetString("host.id"),
		DeviceName: device.DeviceName,
		DeviceType: device.DeviceType,
		Error:      collectorErr.Error(),
	}
	if !device.ScrutinyUUID.IsNil() {
		payload.ScrutinyUUID = device.ScrutinyUUID.String()
	}
	mc.logger.Debugf("Reporting collector error to API: %v", payload)

	apiEndpoint, _ := url.Parse(mc.apiEndpoint.String())
	apiEndpoint, _ = apiEndpoint.Parse("api/collector/error")

	response := models.DeviceWrapper{}
	if err := mc.postJson(apiEndpoint.String(), payload, &response); err != nil {
		mc.logger.Errorf("An error occurred while reporting a collector error to the API: %v", err)
	}
}
