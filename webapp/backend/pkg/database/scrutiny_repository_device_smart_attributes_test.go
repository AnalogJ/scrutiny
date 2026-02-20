package database

import (
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestApplyAttributeOverrides(t *testing.T) {
	// Setup
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel) // Silence logs during tests
	sr := &scrutinyRepository{
		logger: logger,
	}

	wwn := "0x5000c5002df89099"
	protocol := "ATA"

	// Helper to create a basic smart struct
	createSmart := func(val int64, status pkg.AttributeStatus) *measurements.Smart {
		return &measurements.Smart{
			DeviceProtocol: protocol,
			Status:         pkg.DeviceStatusPassed,
			Attributes: map[string]measurements.SmartAttribute{
				"1": &measurements.SmartAtaAttribute{
					AttributeId: 1,
					Value:       100,
					RawValue:    val,
					Status:      status,
				},
			},
		}
	}

	t.Run("No Overrides", func(t *testing.T) {
		smart := createSmart(0, pkg.AttributeStatusPassed)
		overrides := []models.AttributeOverride{}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.Equal(t, pkg.DeviceStatusPassed, smart.Status)
		assert.Equal(t, pkg.AttributeStatusPassed, smart.Attributes["1"].GetStatus())
	})

	t.Run("Ignore Attribute", func(t *testing.T) {
		smart := createSmart(10, pkg.AttributeStatusFailedScrutiny)
		smart.Status = pkg.DeviceStatusFailedScrutiny // Initially failed

		overrides := []models.AttributeOverride{
			{
				Protocol:    protocol,
				AttributeId: "1",
				WWN:         wwn,
				Action:      "ignore",
			},
		}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.Equal(t, pkg.DeviceStatusPassed, smart.Status)
		assert.Equal(t, pkg.AttributeStatusPassed, smart.Attributes["1"].GetStatus())
		assert.Contains(t, smart.Attributes["1"].(*measurements.SmartAtaAttribute).StatusReason, "Ignored")
	})

	t.Run("Force Status Warn", func(t *testing.T) {
		smart := createSmart(0, pkg.AttributeStatusPassed)

		overrides := []models.AttributeOverride{
			{
				Protocol:    protocol,
				AttributeId: "1",
				Action:      "force_status",
				Status:      "warn",
			},
		}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.False(t, pkg.DeviceStatusHas(smart.Status, pkg.DeviceStatusFailedScrutiny))
		assert.Equal(t, pkg.AttributeStatusWarningScrutiny, smart.Attributes["1"].GetStatus())
	})

	t.Run("Force Status Fail", func(t *testing.T) {
		smart := createSmart(0, pkg.AttributeStatusPassed)

		overrides := []models.AttributeOverride{
			{
				Protocol:    protocol,
				AttributeId: "1",
				Action:      "force_status",
				Status:      "failed",
			},
		}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.True(t, pkg.DeviceStatusHas(smart.Status, pkg.DeviceStatusFailedScrutiny))
		assert.Equal(t, pkg.AttributeStatusFailedScrutiny, smart.Attributes["1"].GetStatus())
	})

	t.Run("Threshold Warn", func(t *testing.T) {
		smart := createSmart(50, pkg.AttributeStatusPassed)

		overrides := []models.AttributeOverride{
			{
				Protocol:     protocol,
				AttributeId:  "1",
				WarnAbove:    40,
				WarnAboveSet: true,
			},
		}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.Equal(t, pkg.AttributeStatusWarningScrutiny, smart.Attributes["1"].GetStatus())
	})

	t.Run("Threshold Fail", func(t *testing.T) {
		smart := createSmart(50, pkg.AttributeStatusPassed)

		overrides := []models.AttributeOverride{
			{
				Protocol:     protocol,
				AttributeId:  "1",
				FailAbove:    40,
				FailAboveSet: true,
			},
		}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.True(t, pkg.DeviceStatusHas(smart.Status, pkg.DeviceStatusFailedScrutiny))
		assert.Equal(t, pkg.AttributeStatusFailedScrutiny, smart.Attributes["1"].GetStatus())
	})

	t.Run("Threshold Precedence (Fail > Warn)", func(t *testing.T) {
		smart := createSmart(50, pkg.AttributeStatusPassed)

		overrides := []models.AttributeOverride{
			{
				Protocol:     protocol,
				AttributeId:  "1",
				WarnAbove:    30,
				WarnAboveSet: true,
				FailAbove:    40,
				FailAboveSet: true,
			},
		}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.True(t, pkg.DeviceStatusHas(smart.Status, pkg.DeviceStatusFailedScrutiny))
		assert.Equal(t, pkg.AttributeStatusFailedScrutiny, smart.Attributes["1"].GetStatus())
	})
	
	t.Run("Threshold Not Exceeded", func(t *testing.T) {
		smart := createSmart(10, pkg.AttributeStatusPassed)

		overrides := []models.AttributeOverride{
			{
				Protocol:     protocol,
				AttributeId:  "1",
				WarnAbove:    20,
				WarnAboveSet: true,
			},
		}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.Equal(t, pkg.AttributeStatusPassed, smart.Attributes["1"].GetStatus())
	})

	t.Run("Protocol Mismatch", func(t *testing.T) {
		smart := createSmart(10, pkg.AttributeStatusPassed)

		overrides := []models.AttributeOverride{
			{
				Protocol:    "NVMe", // Mismatch
				AttributeId: "1",
				Action:      "force_status",
				Status:      "failed",
			},
		}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.Equal(t, pkg.AttributeStatusPassed, smart.Attributes["1"].GetStatus())
	})
	
	t.Run("WWN Mismatch", func(t *testing.T) {
		smart := createSmart(10, pkg.AttributeStatusPassed)

		overrides := []models.AttributeOverride{
			{
				Protocol:    protocol,
				AttributeId: "1",
				WWN:         "0xOtherWWN", // Mismatch
				Action:      "force_status",
				Status:      "failed",
			},
		}

		sr.applyAttributeOverrides(smart, wwn, overrides)

		assert.Equal(t, pkg.AttributeStatusPassed, smart.Attributes["1"].GetStatus())
	})
}
