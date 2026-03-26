package detect

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/require"
)

func TestGenerateScrutinyUUID(t *testing.T) {
	t.Run("NVMe device from test data", func(t *testing.T) {
		testData, err := os.ReadFile("testdata/smartctl_info_nvme.json")
		require.NoError(t, err)

		var smartInfo collector.SmartInfo
		err = json.Unmarshal(testData, &smartInfo)
		require.NoError(t, err)

		device := &models.Device{
			ModelName:    smartInfo.ModelName,
			SerialNumber: smartInfo.SerialNumber,
		}
		// NVMe drives don't have a WWN
		// so scrutiny falls back to serial number
		device.WWN = device.SerialNumber

		uuid := GenerateScrutinyUUID(device.ModelName, device.SerialNumber, device.WWN)

		require.NotEmpty(t, uuid.String(), "Generated UUID should not be empty")
		require.Equal(t, uint8(5), uuid.Version(), "Expected UUID version 5")

		uuid2 := GenerateScrutinyUUID(device.ModelName, device.SerialNumber, device.WWN)
		require.True(t, bytes.Equal(uuid.Bytes(), uuid2.Bytes()), "UUID generation should be deterministic for the same input")
	})

	// Test with different device data to ensure uniqueness
	t.Run("different devices produce different UUIDs", func(t *testing.T) {
		device1 := models.Device{
			ModelName:    "Samsung SSD 860 EVO 1TB",
			SerialNumber: "S3ZANX0K123456A",
			WWN:          "5002538e40a22954",
		}

		device2 := device1
		device2.SerialNumber = "S3ZANX0K123456B"

		uuid1 := GenerateScrutinyUUID(device1.ModelName, device1.SerialNumber, device1.WWN)
		uuid2 := GenerateScrutinyUUID(device2.ModelName, device2.SerialNumber, device2.WWN)

		require.False(t, bytes.Equal(uuid1.Bytes(), uuid2.Bytes()), "Different devices should produce different UUIDs")
	})
}

func TestScrutinyNamespaceUUID(t *testing.T) {
	// Make sure no one changes the namespace
	expectedNamespace, err := uuid.FromString("3ea22b35-682b-49fb-a655-abffed108e48")
	if err != nil {
		t.Fatalf("Failed to parse expected namespace UUID: %v", err)
	}

	require.True(t, bytes.Equal(ScrutinyNamespaceUUID.Bytes(), expectedNamespace.Bytes()), "Scrutiny Namespace UUID should never change")
}
