package m20201107210306

import "gorm.io/gorm"

type SmartScsiAttribute struct {
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
	History          []SmartScsiAttribute `gorm:"-" json:"history,omitempty"`
}
