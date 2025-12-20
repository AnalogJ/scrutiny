package database

import (
	"context"
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"gorm.io/gorm/clause"
	"time"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Device
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// insert device into DB (and update specified columns if device is already registered)
// update device fields that may change: (DeviceType, HostID)
func (sr *scrutinyRepository) RegisterDevice(ctx context.Context, dev models.Device) error {
	if err := sr.gormClient.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "wwn"}},
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
		return nil, fmt.Errorf("Could not get device summary from DB: %v", err)
	}
	return devices, nil
}

// update device (only metadata) from collector
func (sr *scrutinyRepository) UpdateDevice(ctx context.Context, wwn string, collectorSmartData collector.SmartInfo) (models.Device, error) {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("wwn = ?", wwn).First(&device).Error; err != nil {
		return device, fmt.Errorf("Could not get device from DB: %v", err)
	}

	//TODO catch GormClient err
	err := device.UpdateFromCollectorSmartInfo(collectorSmartData)
	if err != nil {
		return device, err
	}
	return device, sr.gormClient.Model(&device).Updates(device).Error
}

// Update Device Status
func (sr *scrutinyRepository) UpdateDeviceStatus(ctx context.Context, wwn string, status pkg.DeviceStatus) (models.Device, error) {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("wwn = ?", wwn).First(&device).Error; err != nil {
		return device, fmt.Errorf("Could not get device from DB: %v", err)
	}

	device.DeviceStatus = pkg.DeviceStatusSet(device.DeviceStatus, status)
	return device, sr.gormClient.Model(&device).Updates(device).Error
}

func (sr *scrutinyRepository) GetDeviceDetails(ctx context.Context, wwn string) (models.Device, error) {
	var device models.Device

	fmt.Println("GetDeviceDetails from GORM")

	if err := sr.gormClient.WithContext(ctx).Where("wwn = ?", wwn).First(&device).Error; err != nil {
		return models.Device{}, err
	}

	return device, nil
}

// Update Device Archived State
func (sr *scrutinyRepository) UpdateDeviceArchived(ctx context.Context, wwn string, archived bool) error {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("wwn = ?", wwn).First(&device).Error; err != nil {
		return fmt.Errorf("Could not get device from DB: %v", err)
	}

	return sr.gormClient.Model(&device).Where("wwn = ?", wwn).Update("archived", archived).Error
}

// Update Device Muted State
func (sr *scrutinyRepository) UpdateDeviceMuted(ctx context.Context, wwn string, muted bool) error {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("wwn = ?", wwn).First(&device).Error; err != nil {
		return fmt.Errorf("Could not get device from DB: %v", err)
	}

	return sr.gormClient.Model(&device).Where("wwn = ?", wwn).Update("muted", muted).Error
}

// Update Device Label (custom user-provided name)
func (sr *scrutinyRepository) UpdateDeviceLabel(ctx context.Context, wwn string, label string) error {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("wwn = ?", wwn).First(&device).Error; err != nil {
		return fmt.Errorf("Could not get device from DB: %v", err)
	}

	return sr.gormClient.Model(&device).Where("wwn = ?", wwn).Update("label", label).Error
}

func (sr *scrutinyRepository) DeleteDevice(ctx context.Context, wwn string) error {
	if err := sr.gormClient.WithContext(ctx).Where("wwn = ?", wwn).Delete(&models.Device{}).Error; err != nil {
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
		sr.logger.Infof("Deleting data for %s in bucket: %s", wwn, bucket)
		if err := sr.influxClient.DeleteAPI().DeleteWithName(
			ctx,
			sr.appConfig.GetString("web.influxdb.org"),
			bucket,
			time.Now().AddDate(-10, 0, 0),
			time.Now(),
			fmt.Sprintf(`device_wwn="%s"`, wwn),
		); err != nil {
			return err
		}
	}

	return nil
}
