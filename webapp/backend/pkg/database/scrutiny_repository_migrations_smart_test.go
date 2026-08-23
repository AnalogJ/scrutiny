package database

import (
	"testing"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database/migrations/m20201107210306"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// preInfluxDB spins up an in-memory copy of the pre-InfluxDB schema and returns a gorm handle.
func preInfluxDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_pragma=foreign_keys(0)"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&m20201107210306.Device{},
		&m20201107210306.Smart{},
		&m20201107210306.SmartAtaAttribute{},
		&m20201107210306.SmartNvmeAttribute{},
		&m20201107210306.SmartScsiAttribute{},
	))
	return db
}

// This helper converts pre-InfluxDB sqlite rows into a measurements.Smart, and is the one place
// the migration path re-derives attribute statuses. It is exercised here for all three protocols.
func Test_m20201107210306_FromPreInfluxDBSmartResults_ATA(t *testing.T) {
	db := preInfluxDB(t)

	device := m20201107210306.Device{WWN: "ata-wwn-1", DeviceProtocol: m20201107210306.DeviceProtocolAta}
	require.NoError(t, db.Create(&device).Error)

	smart := m20201107210306.Smart{
		DeviceWWN: device.WWN, TestDate: time.Now(), Temp: 40, PowerOnHours: 1000,
		AtaAttributes: []m20201107210306.SmartAtaAttribute{
			// raw 116 lands in the 70..130 observed bucket (22.3% AFR) -> fails scrutiny analysis
			{AttributeId: 199, Name: "CRC_Error_Count", Value: 200, Worst: 200, Threshold: 0, RawValue: 116, RawString: "116"},
			{AttributeId: 5, Name: "Reallocated_Sector_Ct", Value: 100, Worst: 100, Threshold: 10, RawValue: 0, RawString: "0"},
		},
	}
	require.NoError(t, db.Create(&smart).Error)

	err, post := m20201107210306_FromPreInfluxDBSmartResultsCreatePostInfluxDBSmartResults(db, device, smart)
	require.NoError(t, err)
	require.Equal(t, pkg.DeviceProtocolAta, post.DeviceProtocol)
	require.Len(t, post.Attributes, 2)
	require.Equal(t, pkg.AttributeStatusFailedScrutiny, post.Attributes["199"].GetStatus())
	require.Equal(t, pkg.AttributeStatusPassed, post.Attributes["5"].GetStatus())
	require.Equal(t, pkg.DeviceStatusFailedScrutiny, post.Status)
}

func Test_m20201107210306_FromPreInfluxDBSmartResults_NVMe(t *testing.T) {
	db := preInfluxDB(t)

	device := m20201107210306.Device{WWN: "nvme-wwn-1", DeviceProtocol: m20201107210306.DeviceProtocolNvme}
	require.NoError(t, db.Create(&device).Error)

	smart := m20201107210306.Smart{
		DeviceWWN: device.WWN, TestDate: time.Now(),
		NvmeAttributes: []m20201107210306.SmartNvmeAttribute{
			{AttributeId: "media_errors", Name: "Media Errors", Value: 10},
			{AttributeId: "percentage_used", Name: "Percentage Used", Value: 5},
		},
	}
	require.NoError(t, db.Create(&smart).Error)

	err, post := m20201107210306_FromPreInfluxDBSmartResultsCreatePostInfluxDBSmartResults(db, device, smart)
	require.NoError(t, err)
	require.Equal(t, pkg.DeviceProtocolNvme, post.DeviceProtocol)
	// media_errors has a recommended threshold of 0, so a non-zero count fails
	require.Equal(t, pkg.AttributeStatusFailedScrutiny, post.Attributes["media_errors"].GetStatus())
	require.Equal(t, pkg.DeviceStatusFailedScrutiny, post.Status)
}

func Test_m20201107210306_FromPreInfluxDBSmartResults_SCSI(t *testing.T) {
	db := preInfluxDB(t)

	device := m20201107210306.Device{WWN: "scsi-wwn-1", DeviceProtocol: m20201107210306.DeviceProtocolScsi}
	require.NoError(t, db.Create(&device).Error)

	smart := m20201107210306.Smart{
		DeviceWWN: device.WWN, TestDate: time.Now(),
		ScsiAttributes: []m20201107210306.SmartScsiAttribute{
			{AttributeId: "scsi_grown_defect_list", Name: "Grown Defect List", Value: 3},
			{AttributeId: "read_errors_corrected_by_eccfast", Name: "Read ECC Fast", Value: 1},
		},
	}
	require.NoError(t, db.Create(&smart).Error)

	err, post := m20201107210306_FromPreInfluxDBSmartResultsCreatePostInfluxDBSmartResults(db, device, smart)
	require.NoError(t, err)
	require.Equal(t, pkg.DeviceProtocolScsi, post.DeviceProtocol)
	require.NotEmpty(t, post.Attributes)
}
