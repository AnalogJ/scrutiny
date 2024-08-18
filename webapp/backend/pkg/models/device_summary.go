package models

import (
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
)

type DeviceSummaryWrapper struct {
	Success bool    `json:"success"`
	Errors  []error `json:"errors"`
	Data    struct {
		Summary map[string]*DeviceSummary `json:"summary"`
	} `json:"data"`
}

type DeviceSummary struct {
	Device Device `json:"device"`

	SmartResults *SmartSummary                   `json:"smart,omitempty"`
	TempHistory  []measurements.SmartTemperature `json:"temp_history,omitempty"`
}
type SmartSummary struct {
	// Collector Summary Data
	CollectorDate time.Time `json:"collector_date,omitempty"`
	Temp          int64     `json:"temp,omitempty"`
	PowerOnHours  int64     `json:"power_on_hours,omitempty"`
}
