package db

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metadata"
	"github.com/jinzhu/gorm"
)

const SmartAttributeStatusPassed = "passed"
const SmartAttributeStatusFailed = "failed"
const SmartAttributeStatusWarning = "warn"

type SmartAtaAttribute struct {
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

	TransformedValue int64               `json:"transformed_value"`
	Status           string              `gorm:"-" json:"status,omitempty"`
	StatusReason     string              `gorm:"-" json:"status_reason,omitempty"`
	FailureRate      float64             `gorm:"-" json:"failure_rate,omitempty"`
	History          []SmartAtaAttribute `gorm:"-" json:"history,omitempty"`
}

// compare the attribute (raw, normalized, transformed) value to observed thresholds, and update status if necessary
func (sa *SmartAtaAttribute) MetadataObservedThresholdStatus(smartMetadata metadata.AtaAttributeMetadata) {
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
