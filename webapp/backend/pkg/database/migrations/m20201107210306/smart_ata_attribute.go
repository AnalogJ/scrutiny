package m20201107210306

import "gorm.io/gorm"

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
