package collector

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/jaypipes/ghw"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type BaseCollector struct {
	logger *logrus.Entry
}

func (c *BaseCollector) DetectStorageDevices() ([]models.Device, error) {

	//availableDisksJson, err := c.ExecCmd("smartctl", []string{"-j", "--scan"}, "", os.Environ())
	//if err != nil {
	//	c.logger.Errorf("Error getting block storage info: %v", err)
	//	return nil, err
	//}
	//
	//var smartctlScan models.Scan
	//err = json.Unmarshal([]byte(availableDisksJson), &smartctlScan)
	//if err != nil {
	//	return nil, err
	//}

	block, err := ghw.Block()
	if err != nil {
		c.logger.Errorf("Error getting block storage info: %v", err)
		return nil, err
	}

	approvedDisks := []models.Device{}
	for _, disk := range block.Disks {

		// ignore optical drives and floppy disks
		if disk.DriveType == ghw.DRIVE_TYPE_FDD || disk.DriveType == ghw.DRIVE_TYPE_ODD {
			c.logger.Debugf(" => Ignore: Optical or floppy disk - (found %s)\n", disk.DriveType.String())
			continue
		}

		// ignore removable disks
		if disk.IsRemovable {
			c.logger.Debugf(" => Ignore: Removable disk (%v)\n", disk.IsRemovable)
			continue
		}

		// ignore virtual disks & mobile phone storage devices
		if disk.StorageController == ghw.STORAGE_CONTROLLER_VIRTIO || disk.StorageController == ghw.STORAGE_CONTROLLER_MMC {
			c.logger.Debugf(" => Ignore: Virtual/multi-media storage controller - (found %s)\n", disk.StorageController.String())
			continue
		}

		// ignore NVMe devices (not currently supported) TBA
		if disk.StorageController == ghw.STORAGE_CONTROLLER_NVME {
			c.logger.Debugf(" => Ignore: NVMe storage controller - (found %s)\n", disk.StorageController.String())
			continue
		}

		// Skip unknown storage controllers, not usually S.M.A.R.T compatible.
		if disk.StorageController == ghw.STORAGE_CONTROLLER_UNKNOWN {
			c.logger.Debugf(" => Ignore: Unknown storage controller - (found %s)\n", disk.StorageController.String())
			continue
		}

		diskModel := models.Device{
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

func (c *BaseCollector) getJson(url string, target interface{}) error {

	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (c *BaseCollector) postJson(url string, body interface{}, target interface{}) error {
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

func (c *BaseCollector) ExecCmd(cmdName string, cmdArgs []string, workingDir string, environ []string) (string, error) {

	cmd := exec.Command(cmdName, cmdArgs...)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if environ != nil {
		cmd.Env = environ
	}
	if workingDir != "" && path.IsAbs(workingDir) {
		cmd.Dir = workingDir
	} else if workingDir != "" {
		return "", errors.New("Working Directory must be an absolute path")
	}

	err := cmd.Run()
	return stdBuffer.String(), err

}

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
