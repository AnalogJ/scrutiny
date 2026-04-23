package detect

import (
	"github.com/analogj/scrutiny/collector/pkg/common/shell"
	"github.com/analogj/scrutiny/collector/pkg/models"
)

func DevicePrefix() string {
	return ""
}

func (d *Detect) Scan() ([]models.Device, error) {
	d.Shell = shell.Create()
	// call the base/common functionality to get a list of devices
	return d.SmartctlScan()
}

func (d *Detect) Info(detectedDevices []models.Device) ([]models.Device, error) {
	//inflate device info for detected devices.
	var firstErr error
	for ndx := range detectedDevices {
		err := d.SmartCtlInfo(&detectedDevices[ndx])
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return detectedDevices, firstErr
}

// WWN values NVMe and SCSI
func (d *Detect) wwnFallback(detectedDevice *models.Device) {
	// No fallback on windows
}
