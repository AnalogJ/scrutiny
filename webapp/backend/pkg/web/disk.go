package web

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/jaypipes/ghw"
)

func RetrieveStorageDevices() ([]db.Device, error) {

	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
		return nil, err
	}

	approvedDisks := []db.Device{}
	for _, disk := range block.Disks {
		//TODO: always allow if in approved list
		fmt.Printf(" %v\n", disk)

		// ignore optical drives and floppy disks
		if disk.DriveType == ghw.DRIVE_TYPE_FDD || disk.DriveType == ghw.DRIVE_TYPE_ODD {
			fmt.Printf(" => Ignore: Optical or floppy disk - (found %s)\n", disk.DriveType.String())
			continue
		}

		// ignore removable disks
		if disk.IsRemovable {
			fmt.Printf(" => Ignore: Removable disk (%v)\n", disk.IsRemovable)
			continue
		}

		// ignore virtual disks & mobile phone storage devices
		if disk.StorageController == ghw.STORAGE_CONTROLLER_VIRTIO || disk.StorageController == ghw.STORAGE_CONTROLLER_MMC {
			fmt.Printf(" => Ignore: Virtual/multi-media storage controller - (found %s)\n", disk.StorageController.String())
			continue
		}

		// ignore NVMe devices (not currently supported) TBA
		if disk.StorageController == ghw.STORAGE_CONTROLLER_NVME {
			fmt.Printf(" => Ignore: NVMe storage controller - (found %s)\n", disk.StorageController.String())
			continue
		}

		// Skip unknown storage controllers, not usually S.M.A.R.T compatible.
		if disk.StorageController == ghw.STORAGE_CONTROLLER_UNKNOWN {
			fmt.Printf(" => Ignore: Unknown storage controller - (found %s)\n", disk.StorageController.String())
			continue
		}

		//TODO: remove if in excluded list

		diskModel := db.Device{
			WWN:           disk.WWN,
			Manufacturer:  disk.Vendor,
			ModelName:     disk.Model,
			InterfaceType: disk.StorageController.String(),
			//InterfaceSpeed: string
			SerialNumber: disk.SerialNumber,
			Capacity:     int64(disk.SizeBytes),
			//Firmware       string
			//RotationSpeed  int

			DeviceName: disk.Name,
		}
		if len(diskModel.WWN) == 0 {
			//(macOS and some other os's) do not provide a WWN, so we're going to fallback to
			//diskname as identifier if WWN is not present
			diskModel.WWN = disk.Name
		}

		approvedDisks = append(approvedDisks, diskModel)
	}

	return approvedDisks, nil
}
