package database

import (
	"context"
	"fmt"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"gorm.io/gorm/clause"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ZFS Pool
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// RegisterZFSPool inserts or updates a ZFS pool in the database
func (sr *scrutinyRepository) RegisterZFSPool(ctx context.Context, pool models.ZFSPool) error {
	// First, handle the pool itself (upsert)
	if err := sr.gormClient.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "guid"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "host_id", "status", "health",
			"size", "allocated", "free", "fragmentation", "capacity_percent",
			"ashift",
			"scrub_state", "scrub_start_time", "scrub_end_time",
			"scrub_scanned_bytes", "scrub_issued_bytes", "scrub_total_bytes",
			"scrub_errors_count", "scrub_percent_complete",
			"total_read_errors", "total_write_errors", "total_checksum_errors",
		}),
	}).Create(&pool).Error; err != nil {
		return err
	}

	// Handle vdevs - delete existing and recreate
	if len(pool.Vdevs) > 0 {
		// Delete existing vdevs for this pool
		if err := sr.gormClient.WithContext(ctx).Where("pool_guid = ?", pool.GUID).Delete(&models.ZFSVdev{}).Error; err != nil {
			return err
		}

		// Insert new vdevs with hierarchy
		if err := sr.insertVdevsRecursive(ctx, pool.GUID, pool.Vdevs, nil); err != nil {
			return err
		}
	}

	return nil
}

// insertVdevsRecursive inserts vdevs and their children recursively
func (sr *scrutinyRepository) insertVdevsRecursive(ctx context.Context, poolGUID string, vdevs []models.ZFSVdev, parentID *uint) error {
	for _, vdev := range vdevs {
		vdev.PoolGUID = poolGUID
		vdev.ParentID = parentID
		vdev.ID = 0 // Reset ID to let GORM auto-generate

		children := vdev.Children
		vdev.Children = nil // Don't try to insert children via association

		if err := sr.gormClient.WithContext(ctx).Create(&vdev).Error; err != nil {
			return err
		}

		// Recursively insert children
		if len(children) > 0 {
			if err := sr.insertVdevsRecursive(ctx, poolGUID, children, &vdev.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetZFSPools returns all non-archived ZFS pools
func (sr *scrutinyRepository) GetZFSPools(ctx context.Context) ([]models.ZFSPool, error) {
	pools := []models.ZFSPool{}
	if err := sr.gormClient.WithContext(ctx).Where("archived = ?", false).Find(&pools).Error; err != nil {
		return nil, fmt.Errorf("could not get ZFS pools from DB: %v", err)
	}
	return pools, nil
}

// GetZFSPoolDetails returns a single ZFS pool with its vdev hierarchy
func (sr *scrutinyRepository) GetZFSPoolDetails(ctx context.Context, guid string) (models.ZFSPool, error) {
	var pool models.ZFSPool

	if err := sr.gormClient.WithContext(ctx).Where("guid = ?", guid).First(&pool).Error; err != nil {
		return models.ZFSPool{}, err
	}

	// Load top-level vdevs (those without a parent)
	var vdevs []models.ZFSVdev
	if err := sr.gormClient.WithContext(ctx).Where("pool_guid = ? AND parent_id IS NULL", guid).Find(&vdevs).Error; err != nil {
		return pool, err
	}

	// Load children recursively for each vdev
	for i := range vdevs {
		if err := sr.loadVdevChildren(ctx, &vdevs[i]); err != nil {
			return pool, err
		}
	}

	pool.Vdevs = vdevs
	return pool, nil
}

// loadVdevChildren recursively loads children for a vdev
func (sr *scrutinyRepository) loadVdevChildren(ctx context.Context, vdev *models.ZFSVdev) error {
	var children []models.ZFSVdev
	if err := sr.gormClient.WithContext(ctx).Where("parent_id = ?", vdev.ID).Find(&children).Error; err != nil {
		return err
	}

	for i := range children {
		if err := sr.loadVdevChildren(ctx, &children[i]); err != nil {
			return err
		}
	}

	vdev.Children = children
	return nil
}

// UpdateZFSPoolArchived updates the archived state of a ZFS pool
func (sr *scrutinyRepository) UpdateZFSPoolArchived(ctx context.Context, guid string, archived bool) error {
	var pool models.ZFSPool
	if err := sr.gormClient.WithContext(ctx).Where("guid = ?", guid).First(&pool).Error; err != nil {
		return fmt.Errorf("could not get ZFS pool from DB: %v", err)
	}

	return sr.gormClient.Model(&pool).Where("guid = ?", guid).Update("archived", archived).Error
}

// UpdateZFSPoolMuted updates the muted state of a ZFS pool
func (sr *scrutinyRepository) UpdateZFSPoolMuted(ctx context.Context, guid string, muted bool) error {
	var pool models.ZFSPool
	if err := sr.gormClient.WithContext(ctx).Where("guid = ?", guid).First(&pool).Error; err != nil {
		return fmt.Errorf("could not get ZFS pool from DB: %v", err)
	}

	return sr.gormClient.Model(&pool).Where("guid = ?", guid).Update("muted", muted).Error
}

// UpdateZFSPoolLabel updates the label of a ZFS pool
func (sr *scrutinyRepository) UpdateZFSPoolLabel(ctx context.Context, guid string, label string) error {
	var pool models.ZFSPool
	if err := sr.gormClient.WithContext(ctx).Where("guid = ?", guid).First(&pool).Error; err != nil {
		return fmt.Errorf("could not get ZFS pool from DB: %v", err)
	}

	return sr.gormClient.Model(&pool).Where("guid = ?", guid).Update("label", label).Error
}

// DeleteZFSPool deletes a ZFS pool and its associated data
func (sr *scrutinyRepository) DeleteZFSPool(ctx context.Context, guid string) error {
	// Delete vdevs first (foreign key constraint)
	if err := sr.gormClient.WithContext(ctx).Where("pool_guid = ?", guid).Delete(&models.ZFSVdev{}).Error; err != nil {
		return err
	}

	// Delete the pool
	if err := sr.gormClient.WithContext(ctx).Where("guid = ?", guid).Delete(&models.ZFSPool{}).Error; err != nil {
		return err
	}

	// Delete data from InfluxDB
	buckets := []string{
		sr.appConfig.GetString("web.influxdb.bucket"),
		fmt.Sprintf("%s_weekly", sr.appConfig.GetString("web.influxdb.bucket")),
		fmt.Sprintf("%s_monthly", sr.appConfig.GetString("web.influxdb.bucket")),
		fmt.Sprintf("%s_yearly", sr.appConfig.GetString("web.influxdb.bucket")),
	}

	for _, bucket := range buckets {
		sr.logger.Infof("Deleting ZFS pool data for %s in bucket: %s", guid, bucket)
		if err := sr.influxClient.DeleteAPI().DeleteWithName(
			ctx,
			sr.appConfig.GetString("web.influxdb.org"),
			bucket,
			time.Now().AddDate(-10, 0, 0),
			time.Now(),
			fmt.Sprintf(`pool_guid="%s"`, guid),
		); err != nil {
			return err
		}
	}

	return nil
}

// GetZFSPoolsSummary returns a summary of all non-archived ZFS pools
func (sr *scrutinyRepository) GetZFSPoolsSummary(ctx context.Context) (map[string]*models.ZFSPool, error) {
	pools, err := sr.GetZFSPools(ctx)
	if err != nil {
		return nil, err
	}

	summary := make(map[string]*models.ZFSPool)
	for i := range pools {
		summary[pools[i].GUID] = &pools[i]
	}

	return summary, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ZFS Pool Metrics (InfluxDB)
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// SaveZFSPoolMetrics saves ZFS pool metrics to InfluxDB
func (sr *scrutinyRepository) SaveZFSPoolMetrics(ctx context.Context, pool models.ZFSPool) error {
	// Create metrics from pool data
	metrics := measurements.ZFSPoolMetrics{
		Date:            time.Now(),
		PoolGUID:        pool.GUID,
		PoolName:        pool.Name,
		Size:            pool.Size,
		Allocated:       pool.Allocated,
		Free:            pool.Free,
		CapacityPercent: pool.CapacityPercent,
		Fragmentation:   pool.Fragmentation,
		Status:          string(pool.Status),
		ReadErrors:      pool.TotalReadErrors,
		WriteErrors:     pool.TotalWriteErrors,
		ChecksumErrors:  pool.TotalChecksumErrors,
		ScrubState:      string(pool.ScrubState),
		ScrubPercent:    pool.ScrubPercentComplete,
		ScrubErrors:     pool.ScrubErrorsCount,
	}

	tags, fields := metrics.Flatten()

	// Save to daily bucket
	return sr.saveDatapoint(
		sr.influxWriteApi,
		"zfs_pool",
		tags,
		fields,
		metrics.Date,
		ctx,
	)
}

// GetZFSPoolMetricsHistory retrieves historical metrics for a ZFS pool
func (sr *scrutinyRepository) GetZFSPoolMetricsHistory(ctx context.Context, guid string, durationKey string) ([]measurements.ZFSPoolMetrics, error) {
	// Map duration key to actual duration and bucket
	bucketName := sr.lookupBucketName(durationKey)
	duration := sr.lookupDuration(durationKey)

	queryStr := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "zfs_pool")
		|> filter(fn: (r) => r["pool_guid"] == "%s")
		|> aggregateWindow(every: 1h, fn: last, createEmpty: false)
		|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
		|> sort(columns: ["_time"], desc: false)
	`, bucketName, duration[0], duration[1], guid)

	result, err := sr.influxQueryApi.Query(ctx, queryStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query ZFS pool metrics: %v", err)
	}

	var metricsHistory []measurements.ZFSPoolMetrics
	for result.Next() {
		record := result.Record()
		values := record.Values()

		metrics, err := measurements.NewZFSPoolMetricsFromInfluxDB(values)
		if err != nil {
			sr.logger.Warnf("Failed to parse ZFS pool metrics: %v", err)
			continue
		}

		metricsHistory = append(metricsHistory, *metrics)
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("query error: %v", result.Err())
	}

	return metricsHistory, nil
}
