package db

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metadata"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"strings"
	"time"
)

type DeviceWrapper struct {
	Success bool     `json:"success"`
	Errors  []error  `json:"errors"`
	Data    []Device `json:"data"`
}

type Device struct {
	//GORM attributes, see: http://gorm.io/docs/conventions.html
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	WWN string `json:"wwn" gorm:"primary_key"`

	DeviceName     string `json:"device_name"`
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

	SmartResults []Smart `gorm:"foreignkey:DeviceWWN" json:"smart_results"`
}

//This method requires a device with an array of SmartResults.
//It will remove all SmartResults other than the first (the latest one)
//All removed SmartResults, will be processed, grouping SmartAttribute by attribute_id
// and adding theme to an array called History.
func (dv *Device) SquashHistory() error {
	if len(dv.SmartResults) <= 1 {
		return nil //no history found. ignore
	}

	latestSmartResultSlice := dv.SmartResults[0:1]
	historicalSmartResultSlice := dv.SmartResults[1:]

	//re-assign the latest slice to the SmartResults field
	dv.SmartResults = latestSmartResultSlice

	//process the historical slice
	history := map[int][]SmartAttribute{}
	for _, smartResult := range historicalSmartResultSlice {
		for _, smartAttribute := range smartResult.SmartAttributes {
			if _, ok := history[smartAttribute.AttributeId]; !ok {
				history[smartAttribute.AttributeId] = []SmartAttribute{}
			}
			history[smartAttribute.AttributeId] = append(history[smartAttribute.AttributeId], smartAttribute)
		}
	}

	//now assign the historical slices to the SmartAttributes in the latest SmartResults
	for sandx, smartAttribute := range dv.SmartResults[0].SmartAttributes {
		if attributeHistory, ok := history[smartAttribute.AttributeId]; ok {
			dv.SmartResults[0].SmartAttributes[sandx].History = attributeHistory
		}
	}

	return nil
}

func (dv *Device) ApplyMetadataRules() error {
	//embed metadata in the latest smart attributes object
	if len(dv.SmartResults) > 0 {
		for ndx, attr := range dv.SmartResults[0].SmartAttributes {
			if strings.ToUpper(attr.WhenFailed) == SmartWhenFailedFailingNow {
				//this attribute has previously failed
				dv.SmartResults[0].SmartAttributes[ndx].Status = SmartAttributeStatusFailed
				dv.SmartResults[0].SmartAttributes[ndx].StatusReason = "Attribute is failing manufacturer SMART threshold"

			} else if strings.ToUpper(attr.WhenFailed) == SmartWhenFailedInThePast {
				dv.SmartResults[0].SmartAttributes[ndx].Status = SmartAttributeStatusWarning
				dv.SmartResults[0].SmartAttributes[ndx].StatusReason = "Attribute has previously failed manufacturer SMART threshold"
			}

			if smartMetadata, ok := metadata.AtaSmartAttributes[attr.AttributeId]; ok {
				dv.SmartResults[0].SmartAttributes[ndx].MetadataObservedThresholdStatus(smartMetadata)
			}

			//check if status is blank, set to "passed"
			if len(dv.SmartResults[0].SmartAttributes[ndx].Status) == 0 {
				dv.SmartResults[0].SmartAttributes[ndx].Status = SmartAttributeStatusPassed
			}
		}
	}
	return nil
}

func (dv *Device) UpdateFromCollectorSmartInfo(info collector.SmartInfo) error {
	dv.InterfaceSpeed = info.InterfaceSpeed.Current.String
	dv.Firmware = info.FirmwareVersion
	dv.RotationSpeed = info.RotationRate
	dv.Capacity = info.UserCapacity.Bytes
	dv.FormFactor = info.FormFactor.Name
	//dv.SmartSupport =
	return nil
}
