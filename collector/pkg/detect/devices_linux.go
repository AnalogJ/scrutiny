package detect

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/analogj/scrutiny/collector/pkg/common/shell"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/jaypipes/ghw"
)

func DevicePrefix() string {
	return "/dev/"
}

func (d *Detect) Start() ([]models.Device, error) {
	d.Shell = shell.Create()
	// call the base/common functionality to get a list of devices
	detectedDevices, err := d.SmartctlScan()
	if err != nil {
		return nil, err
	}

	// inflate device info for detected devices.
	for ndx := range detectedDevices {
		d.SmartCtlInfo(&detectedDevices[ndx])   // ignore errors.
		populateUdevInfo(&detectedDevices[ndx]) // ignore errors.
	}

	return detectedDevices, nil
}

// WWN values NVMe and SCSI
func (d *Detect) wwnFallback(detectedDevice *models.Device) {
	block, err := ghw.Block()
	if err == nil {
		for _, disk := range block.Disks {
			if disk.Name == detectedDevice.DeviceName && strings.ToLower(disk.WWN) != "unknown" {
				d.Logger.Debugf("Found matching block device. WWN: %s", disk.WWN)
				detectedDevice.WWN = disk.WWN
				break
			}
		}
	}

	// no WWN found, or could not open Block devices. Either way, fallback to serial number
	if len(detectedDevice.WWN) == 0 {
		d.Logger.Debugf("WWN is empty, falling back to serial number: %s", detectedDevice.SerialNumber)
		detectedDevice.WWN = detectedDevice.SerialNumber
	}

	// wwn must always be lowercase.
	detectedDevice.WWN = strings.ToLower(detectedDevice.WWN)
}

// as discussed in
// - https://github.com/AnalogJ/scrutiny/issues/225
// - https://github.com/jaypipes/ghw/issues/59#issue-361915216
// udev exposes its data in a standardized way under /run/udev/data/....
func populateUdevInfo(detectedDevice *models.Device) error {
	// Get device major:minor numbers
	// `cat /sys/class/block/sda/dev`
	devNo, err := ioutil.ReadFile(filepath.Join("/sys/class/block/", detectedDevice.DeviceName, "dev"))
	if err != nil {
		return err
	}

	// Look up block device in udev runtime database
	// `cat /run/udev/data/b8:0`
	udevID := "b" + strings.TrimSpace(string(devNo))
	udevBytes, err := ioutil.ReadFile(filepath.Join("/run/udev/data/", udevID))
	if err != nil {
		return err
	}

	deviceMountPaths := []string{}
	udevInfo := make(map[string]string)
	for _, udevLine := range strings.Split(string(udevBytes), "\n") {
		if strings.HasPrefix(udevLine, "E:") {
			if s := strings.SplitN(udevLine[2:], "=", 2); len(s) == 2 {
				udevInfo[s[0]] = s[1]
			}
		} else if strings.HasPrefix(udevLine, "S:") {
			deviceMountPaths = append(deviceMountPaths, udevLine[2:])
		}
	}

	// Set additional device information.
	if deviceLabel, exists := udevInfo["ID_FS_LABEL"]; exists {
		detectedDevice.DeviceLabel = deviceLabel
	}
	if deviceUUID, exists := udevInfo["ID_FS_UUID"]; exists {
		detectedDevice.DeviceUUID = deviceUUID
	}
	if deviceSerialID, exists := udevInfo["ID_SERIAL"]; exists {
		detectedDevice.DeviceSerialID = fmt.Sprintf("%s-%s", udevInfo["ID_BUS"], deviceSerialID)
	}

	return nil
}
