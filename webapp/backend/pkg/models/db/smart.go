package db

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metadata"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/jinzhu/gorm"
	"time"
)

const SmartWhenFailedFailingNow = "FAILING_NOW"
const SmartWhenFailedInThePast = "IN_THE_PAST"

type Smart struct {
	gorm.Model

	DeviceWWN string `json:"device_wwn"`
	Device    Device `json:"-" gorm:"foreignkey:DeviceWWN"` // use DeviceWWN as foreign key

	TestDate    time.Time `json:"date"`
	SmartStatus string    `json:"smart_status"`

	//Metrics
	Temp            int64 `json:"temp"`
	PowerOnHours    int64 `json:"power_on_hours"`
	PowerCycleCount int64 `json:"power_cycle_count"`

	SmartAttributes []SmartAttribute `json:"smart_attributes" gorm:"foreignkey:SmartId"`
}

func (sm *Smart) FromCollectorSmartInfo(wwn string, info collector.SmartInfo) error {
	sm.DeviceWWN = wwn
	sm.TestDate = time.Unix(info.LocalTime.TimeT, 0)

	//smart metrics
	sm.Temp = info.Temperature.Current
	sm.PowerCycleCount = info.PowerCycleCount
	sm.PowerOnHours = info.PowerOnTime.Hours

	sm.SmartAttributes = []SmartAttribute{}
	for _, collectorAttr := range info.AtaSmartAttributes.Table {
		attrModel := SmartAttribute{
			AttributeId: collectorAttr.ID,
			Name:        collectorAttr.Name,
			Value:       collectorAttr.Value,
			Worst:       collectorAttr.Worst,
			Threshold:   collectorAttr.Thresh,
			RawValue:    collectorAttr.Raw.Value,
			RawString:   collectorAttr.Raw.String,
			WhenFailed:  collectorAttr.WhenFailed,
		}

		//now that we've parsed the data from the smartctl response, lets match it against our metadata rules and add additional Scrutiny specific data.
		if smartMetadata, ok := metadata.AtaSmartAttributes[collectorAttr.ID]; ok {
			attrModel.Name = smartMetadata.DisplayName
			if smartMetadata.Transform != nil {
				attrModel.TransformedValue = smartMetadata.Transform(attrModel.Value, attrModel.RawValue, attrModel.RawString)
			}
		}
		sm.SmartAttributes = append(sm.SmartAttributes, attrModel)
	}

	if info.SmartStatus.Passed {
		sm.SmartStatus = "passed"
	} else {
		sm.SmartStatus = "failed"
	}
	return nil
}

const SmartAttributeStatusPassed = "passed"
const SmartAttributeStatusFailed = "failed"
const SmartAttributeStatusWarning = "warn"

type SmartAttribute struct {
	gorm.Model

	SmartId int    `json:"smart_id"`
	Smart   Device `json:"-" gorm:"foreignkey:SmartId"` // use SmartId as foreign key

	AttributeId int    `json:"attribute_id"`
	Name        string `json:"name"`
	Value       int    `json:"value"`
	Worst       int    `json:"worst"`
	Threshold   int    `json:"thresh"`
	RawValue    int64  `json:"raw_value"`
	RawString   string `json:"raw_string"`
	WhenFailed  string `json:"when_failed"`

	TransformedValue int64            `json:"transformed_value"`
	Status           string           `gorm:"-" json:"status,omitempty"`
	StatusReason     string           `gorm:"-" json:"status_reason,omitempty"`
	FailureRate      float64          `gorm:"-" json:"failure_rate,omitempty"`
	History          []SmartAttribute `gorm:"-" json:"history,omitempty"`
}

// compare the attribute (raw, normalized, transformed) value to observed thresholds, and update status if necessary
func (sa *SmartAttribute) MetadataObservedThresholdStatus(smartMetadata metadata.AtaSmartAttribute) {
	//TODO: multiple rules
	// try to predict the failure rates for observed thresholds that have 0 failure rate and error bars.
	// - if the attribute is critical
	//		- the failure rate is over 10 - set to failed
	//		- the attribute does not match any threshold, set to warn
	// - if the attribute is not critical
	//		- if failure rate is above 20 - set to failed
	// 		- if failure rate is above 10 but below 20 - set to warn

	//update the smart attribute status based on Observed thresholds.
	var value int64
	if smartMetadata.DisplayType == metadata.AtaSmartAttributeDisplayTypeNormalized {
		value = int64(sa.Value)
	} else if smartMetadata.DisplayType == metadata.AtaSmartAttributeDisplayTypeTransformed {
		value = sa.TransformedValue
	} else {
		value = sa.RawValue
	}

	for _, obsThresh := range smartMetadata.ObservedThresholds {

		//check if "value" is in this bucket
		if ((obsThresh.Low == obsThresh.High) && value == obsThresh.Low) ||
			(obsThresh.Low < value && value <= obsThresh.High) {
			sa.FailureRate = obsThresh.AnnualFailureRate

			if smartMetadata.Critical {
				if obsThresh.AnnualFailureRate >= 0.10 {
					sa.Status = SmartAttributeStatusFailed
					sa.StatusReason = "Observed Failure Rate for Critical Attribute is greater than 10%"
				}
			} else {
				if obsThresh.AnnualFailureRate >= 0.20 {
					sa.Status = SmartAttributeStatusFailed
					sa.StatusReason = "Observed Failure Rate for Attribute is greater than 20%"
				} else if obsThresh.AnnualFailureRate >= 0.10 {
					sa.Status = SmartAttributeStatusWarning
					sa.StatusReason = "Observed Failure Rate for Attribute is greater than 10%"
				}
			}

			//we've found the correct bucket, we can drop out of this loop
			return
		}
	}
	// no bucket found
	if smartMetadata.Critical {
		sa.Status = SmartAttributeStatusWarning
		sa.StatusReason = "Could not determine Observed Failure Rate for Critical Attribute"
	}

	return
}
