package models

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"time"
)

type DeviceWrapper struct {
	Success bool     `json:"success"`
	Errors  []error  `json:"errors"`
	Data    []Device `json:"data"`
}

type Device struct {
	//GORM attributes, see: http://gorm.io/docs/conventions.html
	Archived  bool `json:"archived"`
	Muted     bool `json:"muted"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	WWN string `json:"wwn" gorm:"primary_key"`

	DeviceName     string `json:"device_name"`
	DeviceUUID     string `json:"device_uuid"`
	DeviceSerialID string `json:"device_serial_id"`
	DeviceLabel    string `json:"device_label"`

	Manufacturer   string `json:"manufacturer"`
	ModelName      string `json:"model_name"`
	InterfaceType  string `json:"interface_type"`
	InterfaceSpeed string `json:"interface_speed"`
	SerialNumber   string `json:"serial_number"`
	Firmware       string `json:"firmware"`
	RotationSpeed  int    `json:"rotational_speed"`
	Capacity       int64  `json:"capacity"`
	FormFactor     string `json:"form_factor"`
	SmartSupport   bool   `json:"smart_support"`
	DeviceProtocol string `json:"device_protocol"` //protocol determines which smart attribute types are available (ATA, NVMe, SCSI)
	DeviceType     string `json:"device_type"`     //device type is used for querying with -d/t flag, should only be used by collector.

	// User provided metadata
	Label  string `json:"label"`
	HostId string `json:"host_id"`

	// Data set by Scrutiny
	DeviceStatus pkg.DeviceStatus `json:"device_status"`
}

func (dv *Device) IsAta() bool {
	return dv.DeviceProtocol == pkg.DeviceProtocolAta
}

func (dv *Device) IsScsi() bool {
	return dv.DeviceProtocol == pkg.DeviceProtocolScsi
}

func (dv *Device) IsNvme() bool {
	return dv.DeviceProtocol == pkg.DeviceProtocolNvme
}

//
////This method requires a device with an array of SmartResults.
////It will remove all SmartResults other than the first (the latest one)
////All removed SmartResults, will be processed, grouping SmartAtaAttribute by attribute_id
//// and adding theme to an array called History.
//func (dv *Device) SquashHistory() error {
//	if len(dv.SmartResults) <= 1 {
//		return nil //no ataHistory found. ignore
//	}
//
//	latestSmartResultSlice := dv.SmartResults[0:1]
//	historicalSmartResultSlice := dv.SmartResults[1:]
//
//	//re-assign the latest slice to the SmartResults field
//	dv.SmartResults = latestSmartResultSlice
//
//	//process the historical slice for ATA data
//	if len(dv.SmartResults[0].AtaAttributes) > 0 {
//		ataHistory := map[int][]SmartAtaAttribute{}
//		for _, smartResult := range historicalSmartResultSlice {
//			for _, smartAttribute := range smartResult.AtaAttributes {
//				if _, ok := ataHistory[smartAttribute.AttributeId]; !ok {
//					ataHistory[smartAttribute.AttributeId] = []SmartAtaAttribute{}
//				}
//				ataHistory[smartAttribute.AttributeId] = append(ataHistory[smartAttribute.AttributeId], smartAttribute)
//			}
//		}
//
//		//now assign the historical slices to the AtaAttributes in the latest SmartResults
//		for sandx, smartAttribute := range dv.SmartResults[0].AtaAttributes {
//			if attributeHistory, ok := ataHistory[smartAttribute.AttributeId]; ok {
//				dv.SmartResults[0].AtaAttributes[sandx].History = attributeHistory
//			}
//		}
//	}
//
//	//process the historical slice for Nvme data
//	if len(dv.SmartResults[0].NvmeAttributes) > 0 {
//		nvmeHistory := map[string][]SmartNvmeAttribute{}
//		for _, smartResult := range historicalSmartResultSlice {
//			for _, smartAttribute := range smartResult.NvmeAttributes {
//				if _, ok := nvmeHistory[smartAttribute.AttributeId]; !ok {
//					nvmeHistory[smartAttribute.AttributeId] = []SmartNvmeAttribute{}
//				}
//				nvmeHistory[smartAttribute.AttributeId] = append(nvmeHistory[smartAttribute.AttributeId], smartAttribute)
//			}
//		}
//
//		//now assign the historical slices to the AtaAttributes in the latest SmartResults
//		for sandx, smartAttribute := range dv.SmartResults[0].NvmeAttributes {
//			if attributeHistory, ok := nvmeHistory[smartAttribute.AttributeId]; ok {
//				dv.SmartResults[0].NvmeAttributes[sandx].History = attributeHistory
//			}
//		}
//	}
//	//process the historical slice for Scsi data
//	if len(dv.SmartResults[0].ScsiAttributes) > 0 {
//		scsiHistory := map[string][]SmartScsiAttribute{}
//		for _, smartResult := range historicalSmartResultSlice {
//			for _, smartAttribute := range smartResult.ScsiAttributes {
//				if _, ok := scsiHistory[smartAttribute.AttributeId]; !ok {
//					scsiHistory[smartAttribute.AttributeId] = []SmartScsiAttribute{}
//				}
//				scsiHistory[smartAttribute.AttributeId] = append(scsiHistory[smartAttribute.AttributeId], smartAttribute)
//			}
//		}
//
//		//now assign the historical slices to the AtaAttributes in the latest SmartResults
//		for sandx, smartAttribute := range dv.SmartResults[0].ScsiAttributes {
//			if attributeHistory, ok := scsiHistory[smartAttribute.AttributeId]; ok {
//				dv.SmartResults[0].ScsiAttributes[sandx].History = attributeHistory
//			}
//		}
//	}
//	return nil
//}
//
//func (dv *Device) ApplyMetadataRules() error {
//
//	//embed metadata in the latest smart attributes object
//	if len(dv.SmartResults) > 0 {
//		for ndx, attr := range dv.SmartResults[0].AtaAttributes {
//			attr.PopulateAttributeStatus()
//			dv.SmartResults[0].AtaAttributes[ndx] = attr
//		}
//
//		for ndx, attr := range dv.SmartResults[0].NvmeAttributes {
//			attr.PopulateAttributeStatus()
//			dv.SmartResults[0].NvmeAttributes[ndx] = attr
//
//		}
//
//		for ndx, attr := range dv.SmartResults[0].ScsiAttributes {
//			attr.PopulateAttributeStatus()
//			dv.SmartResults[0].ScsiAttributes[ndx] = attr
//
//		}
//	}
//	return nil
//}

// This function is called every time the collector sends SMART data to the API.
// It can be used to update device data that can change over time.
func (dv *Device) UpdateFromCollectorSmartInfo(info collector.SmartInfo) error {
	dv.Firmware = info.FirmwareVersion
	dv.DeviceProtocol = info.Device.Protocol

	if !info.SmartStatus.Passed {
		dv.DeviceStatus = pkg.DeviceStatusSet(dv.DeviceStatus, pkg.DeviceStatusFailedSmart)
	}

	return nil
}
