package detect

import (
	"encoding/json"
	"fmt"
	"github.com/analogj/scrutiny/collector/pkg/common"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/denisbrodbeck/machineid"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Detect struct {
	Logger *logrus.Entry
}

//private/common functions

// This function calls smartctl --scan which can be used to detect storage devices.
// It has a couple of issues however:
// - --scan does not return any results on mac
//
// To handle these issues, we have OS specific wrapper functions that update/modify these detected devices.
// models.Device returned from this function only contain the minimum data for smartctl to execute: device type and device name (device file).
func (d *Detect) smartctlScan() ([]models.Device, error) {
	//we use smartctl to detect all the drives available.
	detectedDeviceConnJson, err := common.ExecCmd("smartctl", []string{"--scan", "-j"}, "", os.Environ())
	if err != nil {
		d.Logger.Errorf("Error scanning for devices: %v", err)
		return nil, err
	}

	var detectedDeviceConns models.Scan
	err = json.Unmarshal([]byte(detectedDeviceConnJson), &detectedDeviceConns)
	if err != nil {
		d.Logger.Errorf("Error decoding detected devices: %v", err)
		return nil, err
	}

	detectedDevices := []models.Device{}

	for _, detectedDevice := range detectedDeviceConns.Devices {
		detectedDevices = append(detectedDevices, models.Device{
			DeviceType: detectedDevice.Type,
			DeviceName: strings.TrimPrefix(detectedDevice.Name, DevicePrefix()),
		})
	}

	return detectedDevices, nil
}

//updates a device model with information from smartctl --scan
// It has a couple of issues however:
// - WWN is provided as component data, rather than a "string". We'll have to generate the WWN value ourselves
// - WWN from smartctl only provided for ATA protocol drives, NVMe and SCSI drives do not include WWN.
func (d *Detect) smartCtlInfo(device *models.Device) error {

	args := []string{"--info", "-j"}
	//only include the device type if its a non-standard one. In some cases ata drives are detected as scsi in docker, and metadata is lost.
	if len(device.DeviceType) > 0 && device.DeviceType != "scsi" && device.DeviceType != "ata" {
		args = append(args, "-d", device.DeviceType)
	}
	args = append(args, fmt.Sprintf("%s%s", DevicePrefix(), device.DeviceName))

	availableDeviceInfoJson, err := common.ExecCmd("smartctl", args, "", os.Environ())
	if err != nil {
		d.Logger.Errorf("Could not retrieve device information for %s: %v", device.DeviceName, err)
		return err
	}

	var availableDeviceInfo collector.SmartInfo
	err = json.Unmarshal([]byte(availableDeviceInfoJson), &availableDeviceInfo)
	if err != nil {
		d.Logger.Errorf("Could not decode device information for %s: %v", device.DeviceName, err)
		return err
	}

	//DeviceType and DeviceName are already populated.
	//WWN
	//InterfaceType:
	device.ModelName = availableDeviceInfo.ModelName
	device.InterfaceSpeed = availableDeviceInfo.InterfaceSpeed.Current.String
	device.SerialNumber = availableDeviceInfo.SerialNumber
	device.Firmware = availableDeviceInfo.FirmwareVersion
	device.RotationSpeed = availableDeviceInfo.RotationRate
	device.Capacity = availableDeviceInfo.UserCapacity.Bytes
	device.FormFactor = availableDeviceInfo.FormFactor.Name
	device.DeviceProtocol = availableDeviceInfo.Device.Protocol
	if len(availableDeviceInfo.Vendor) > 0 {
		device.Manufacturer = availableDeviceInfo.Vendor
	}

	//populate WWN is possible if present
	if availableDeviceInfo.Wwn.Naa != 0 { //valid values are 1-6 (5 is what we handle correctly)
		d.Logger.Info("Generating WWN")
		wwn := Wwn{
			Naa: availableDeviceInfo.Wwn.Naa,
			Oui: availableDeviceInfo.Wwn.Oui,
			Id:  availableDeviceInfo.Wwn.ID,
		}
		device.WWN = wwn.ToString()
	} else {
		d.Logger.Info("Using WWN Fallback")
		d.wwnFallback(device)
	}

	return nil
}

//uses https://github.com/denisbrodbeck/machineid to get a OS specific unique machine ID.
func (d *Detect) getMachineId() (string, error) {
	return machineid.ProtectedID("scrutiny")
}
