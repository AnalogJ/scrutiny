package db

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metadata"
	"gorm.io/gorm"
)

type SmartNvmeAttribute struct {
	gorm.Model

	SmartId int    `json:"smart_id"`
	Smart   Device `json:"-" gorm:"foreignkey:SmartId"` // use SmartId as foreign key

	AttributeId string `json:"attribute_id"` //json string from smartctl
	Name        string `json:"name"`
	Value       int    `json:"value"`
	Threshold   int    `json:"thresh"`

	TransformedValue int64                `json:"transformed_value"`
	Status           string               `gorm:"-" json:"status,omitempty"`
	StatusReason     string               `gorm:"-" json:"status_reason,omitempty"`
	FailureRate      float64              `gorm:"-" json:"failure_rate,omitempty"`
	History          []SmartNvmeAttribute `gorm:"-" json:"history,omitempty"`
}

//populate attribute status, using SMART Thresholds & Observed Metadata
func (sa *SmartNvmeAttribute) PopulateAttributeStatus() {

	//-1 is a special number meaning no threshold.
	if sa.Threshold != -1 {
		if smartMetadata, ok := metadata.NmveMetadata[sa.AttributeId]; ok {
			//check what the ideal is. Ideal tells us if we our recorded value needs to be above, or below the threshold
			if (smartMetadata.Ideal == "low" && sa.Value > sa.Threshold) ||
				(smartMetadata.Ideal == "high" && sa.Value < sa.Threshold) {
				sa.Status = SmartAttributeStatusFailed
				sa.StatusReason = "Attribute is failing recommended SMART threshold"
			}
		}
	}
	//TODO: eventually figure out the critical_warning bits and determine correct error messages here.

	//check if status is blank, set to "passed"
	if len(sa.Status) == 0 {
		sa.Status = SmartAttributeStatusPassed
	}
}
