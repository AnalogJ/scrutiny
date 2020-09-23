package detect

import (
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/jaypipes/ghw"
	"strings"
)

func DevicePrefix() string {
	return "/dev/"
}

func (d *Detect) Start() ([]models.Device, error) {
	// call the base/common functionality to get a list of devicess
	detectedDevices, err := d.smartctlScan()
	if err != nil {
		return nil, err
	}

	//smartctl --scan doesn't seem to detect mac nvme drives, lets see if we can detect them manually.
	missingDevices, err := d.findMissingDevices(detectedDevices) //we dont care about the error here, just continue retrieving device info.
	if err == nil {
		detectedDevices = append(detectedDevices, missingDevices...)
	}

	//inflate device info for detected devices.
	for ndx, _ := range detectedDevices {
		d.smartCtlInfo(&detectedDevices[ndx]) //ignore errors.
	}

	return detectedDevices, nil
}

func (d *Detect) findMissingDevices(detectedDevices []models.Device) ([]models.Device, error) {

	missingDevices := []models.Device{}

	block, err := ghw.Block()
	if err != nil {
		d.Logger.Errorf("Error getting block storage info: %v", err)
		return nil, err
	}

	for _, disk := range block.Disks {

		// ignore optical drives and floppy disks
		if disk.DriveType == ghw.DRIVE_TYPE_FDD || disk.DriveType == ghw.DRIVE_TYPE_ODD {
			d.Logger.Debugf(" => Ignore: Optical or floppy disk - (found %s)\n", disk.DriveType.String())
			continue
		}

		// ignore removable disks
		if disk.IsRemovable {
			d.Logger.Debugf(" => Ignore: Removable disk (%v)\n", disk.IsRemovable)
			continue
		}

		// ignore virtual disks & mobile phone storage devices
		if disk.StorageController == ghw.STORAGE_CONTROLLER_VIRTIO || disk.StorageController == ghw.STORAGE_CONTROLLER_MMC {
			d.Logger.Debugf(" => Ignore: Virtual/multi-media storage controller - (found %s)\n", disk.StorageController.String())
			continue
		}

		// Skip unknown storage controllers, not usually S.M.A.R.T compatible.
		if disk.StorageController == ghw.STORAGE_CONTROLLER_UNKNOWN {
			d.Logger.Debugf(" => Ignore: Unknown storage controller - (found %s)\n", disk.StorageController.String())
			continue
		}

		//check if device is already detected.
		alreadyDetected := false
		diskName := strings.TrimPrefix(disk.Name, DevicePrefix())
		for _, detectedDevice := range detectedDevices {

			if detectedDevice.DeviceName == diskName {
				alreadyDetected = true
				break
			}
		}
		if !alreadyDetected {
			missingDevices = append(missingDevices, models.Device{
				DeviceName: diskName,
				DeviceType: "",
			})
		}
	}
	return missingDevices, nil
}

//WWN values NVMe and SCSI
func (d *Detect) wwnFallback(detectedDevice *models.Device) {
	block, err := ghw.Block()
	if err == nil {
		for _, disk := range block.Disks {
			if disk.Name == detectedDevice.DeviceName {
				d.Logger.Debugf("Found matching block device. WWN: %s", disk.WWN)
				detectedDevice.WWN = disk.WWN
				break
			}
		}
	}

	//no WWN found, or could not open Block devices. Either way, fallback to serial number
	if len(detectedDevice.WWN) == 0 {
		d.Logger.Debugf("WWN is empty, falling back to serial number: %s", detectedDevice.SerialNumber)
		detectedDevice.WWN = detectedDevice.SerialNumber
	}

	//wwn must always be lowercase.
	detectedDevice.WWN = strings.ToLower(detectedDevice.WWN)
}
