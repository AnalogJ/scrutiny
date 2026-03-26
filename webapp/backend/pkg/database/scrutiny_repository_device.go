package database

import (
	"context"
	"fmt"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm/clause"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Device
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// insert device into DB (and update specified columns if device is already registered)
// update device fields that may change: (DeviceType, HostID)
func (sr *scrutinyRepository) RegisterDevice(ctx context.Context, dev models.Device) error {
	if err := sr.gormClient.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "scrutiny_uuid"}},
		DoUpdates: clause.AssignmentColumns([]string{"host_id", "device_name", "device_type", "device_uuid", "device_serial_id", "device_label"}),
	}).Create(&dev).Error; err != nil {
		return err
	}
	return nil
}

// get a list of all devices (only device metadata, no SMART data)
func (sr *scrutinyRepository) GetDevices(ctx context.Context) ([]models.Device, error) {
	//Get a list of all the active devices.
	devices := []models.Device{}
	if err := sr.gormClient.WithContext(ctx).Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("could not get device summary from DB: %v", err)
	}
	return devices, nil
}

// update device (only metadata) from collector
func (sr *scrutinyRepository) UpdateDevice(ctx context.Context, scrutiny_uuid uuid.UUID, collectorSmartData collector.SmartInfo) (models.Device, error) {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("scrutiny_uuid = ?", scrutiny_uuid.String()).First(&device).Error; err != nil {
		return device, fmt.Errorf("could not get device from DB: %v", err)
	}

	//TODO catch GormClient err
	err := device.UpdateFromCollectorSmartInfo(collectorSmartData)
	if err != nil {
		return device, err
	}
	return device, sr.gormClient.Model(&device).Updates(device).Error
}

// Update Device Status
func (sr *scrutinyRepository) UpdateDeviceStatus(ctx context.Context, scrutiny_uuid uuid.UUID, status pkg.DeviceStatus) (models.Device, error) {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("scrutiny_uuid = ?", scrutiny_uuid.String()).First(&device).Error; err != nil {
		return device, fmt.Errorf("could not get device from DB: %v", err)
	}

	device.DeviceStatus = pkg.DeviceStatusSet(device.DeviceStatus, status)
	return device, sr.gormClient.Model(&device).Updates(device).Error
}

func (sr *scrutinyRepository) GetDeviceDetails(ctx context.Context, scrutiny_uuid uuid.UUID) (models.Device, error) {
	var device models.Device

	fmt.Println("GetDeviceDetails from GORM")

	if err := sr.gormClient.WithContext(ctx).Where("scrutiny_uuid = ?", scrutiny_uuid.String()).First(&device).Error; err != nil {
		return models.Device{}, err
	}

	return device, nil
}

// Update Device Archived State
func (sr *scrutinyRepository) UpdateDeviceArchived(ctx context.Context, scrutiny_uuid uuid.UUID, archived bool) error {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("scrutiny_uuid = ?", scrutiny_uuid.String()).First(&device).Error; err != nil {
		return fmt.Errorf("could not get device from DB: %v", err)
	}

	return sr.gormClient.Model(&device).Where("scrutiny_uuid = ?", scrutiny_uuid.String()).Update("archived", archived).Error
}

func (sr *scrutinyRepository) DeleteDevice(ctx context.Context, scrutiny_uuid uuid.UUID) error {
	if err := sr.gormClient.WithContext(ctx).Where("scrutiny_uuid = ?", scrutiny_uuid.String()).Delete(&models.Device{}).Error; err != nil {
		return err
	}

	//delete data from influxdb.
	buckets := []string{
		sr.appConfig.GetString("web.influxdb.bucket"),
		fmt.Sprintf("%s_weekly", sr.appConfig.GetString("web.influxdb.bucket")),
		fmt.Sprintf("%s_monthly", sr.appConfig.GetString("web.influxdb.bucket")),
		fmt.Sprintf("%s_yearly", sr.appConfig.GetString("web.influxdb.bucket")),
	}

	for _, bucket := range buckets {
		sr.logger.Infof("Deleting data for %s in bucket: %s", scrutiny_uuid.String(), bucket)
		if err := sr.influxClient.DeleteAPI().DeleteWithName(
			ctx,
			sr.appConfig.GetString("web.influxdb.org"),
			bucket,
			time.Now().AddDate(-10, 0, 0),
			time.Now(),
			fmt.Sprintf(`scrutiny_uuid="%s"`, scrutiny_uuid.String()),
		); err != nil {
			return err
		}
	}

	return nil
}
