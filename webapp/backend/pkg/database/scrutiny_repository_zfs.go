package database

import (
	"context"
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ZFS Pool
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// RegisterZfsPools inserts or updates ZFS pools in the database
func (sr *scrutinyRepository) RegisterZfsPools(ctx context.Context, pools []models.ZfsPool) error {
	for _, pool := range pools {
		// First, register the pool
		if err := sr.gormClient.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "pool_guid"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"host_id", "name", "state", "txg", "spa_version", "zpl_version",
				"status", "action", "error_count", "alloc_space", "total_space", "def_space",
				"read_errors", "write_errors", "checksum_errors", "scan_function", "scan_state",
				"scan_start_time", "scan_end_time", "scan_to_examine", "scan_examined",
				"scan_processed", "scan_errors", "scan_issued", "size", "allocated", "free",
				"fragmentation", "capacity_percent", "dedupratio", "updated_at",
			}),
		}).Create(&pool).Error; err != nil {
			return fmt.Errorf("failed to register ZFS pool %s: %v", pool.Name, err)
		}

		// Delete existing vdevs for this pool to avoid duplicates
		if err := sr.gormClient.WithContext(ctx).Where("pool_guid = ?", pool.PoolGuid).Delete(&models.ZfsVdev{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing vdevs for pool %s: %v", pool.Name, err)
		}

		// Register vdevs with proper hierarchy
		if err := sr.registerVdevsHierarchy(ctx, pool.Vdevs, pool.PoolGuid, nil); err != nil {
			return fmt.Errorf("failed to register vdevs for pool %s: %v", pool.Name, err)
		}
	}
	return nil
}

// registerVdevsHierarchy recursively registers vdevs maintaining parent-child relationships
func (sr *scrutinyRepository) registerVdevsHierarchy(ctx context.Context, vdevs []models.ZfsVdev, poolGuid string, parentId *uint) error {
	for _, vdev := range vdevs {
		vdev.PoolGuid = poolGuid
		vdev.ParentId = parentId

		// Create the vdev
		if err := sr.gormClient.WithContext(ctx).Create(&vdev).Error; err != nil {
			return fmt.Errorf("failed to create vdev %s: %v", vdev.Name, err)
		}

		// Register children recursively
		if len(vdev.Children) > 0 {
			if err := sr.registerVdevsHierarchy(ctx, vdev.Children, poolGuid, &vdev.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetZfsPools retrieves all ZFS pools for a specific host
func (sr *scrutinyRepository) GetZfsPools(ctx context.Context) ([]models.ZfsPool, error) {
	var pools []models.ZfsPool
	if err := sr.gormClient.WithContext(ctx).Preload("Vdevs", func(db *gorm.DB) *gorm.DB {
		return db.Order("parent_id ASC, name ASC")
	}).Find(&pools).Error; err != nil {
		return nil, fmt.Errorf("could not get ZFS pools from DB: %v", err)
	}
	return pools, nil
}

// GetZfsPoolByGuid retrieves a specific ZFS pool by its GUID
func (sr *scrutinyRepository) GetZfsPoolByGuid(ctx context.Context, poolGuid string) (models.ZfsPool, error) {
	var pool models.ZfsPool
	if err := sr.gormClient.WithContext(ctx).Preload("Vdevs", func(db *gorm.DB) *gorm.DB {
		return db.Order("parent_id ASC, name ASC")
	}).Where("pool_guid = ?", poolGuid).First(&pool).Error; err != nil {
		return pool, fmt.Errorf("could not get ZFS pool %s from DB: %v", poolGuid, err)
	}
	return pool, nil
}

// GetZfsPoolsByHost retrieves all ZFS pools for a specific host
func (sr *scrutinyRepository) GetZfsPoolsByHost(ctx context.Context, hostId string) ([]models.ZfsPool, error) {
	var pools []models.ZfsPool
	if err := sr.gormClient.WithContext(ctx).Preload("Vdevs", func(db *gorm.DB) *gorm.DB {
		return db.Order("parent_id ASC, name ASC")
	}).Where("host_id = ?", hostId).Find(&pools).Error; err != nil {
		return nil, fmt.Errorf("could not get ZFS pools for host %s from DB: %v", hostId, err)
	}
	return pools, nil
}

// DeleteZfsPool deletes a ZFS pool and all its vdevs
func (sr *scrutinyRepository) DeleteZfsPool(ctx context.Context, poolGuid string) error {
	// Delete vdevs first (due to foreign key constraint)
	if err := sr.gormClient.WithContext(ctx).Where("pool_guid = ?", poolGuid).Delete(&models.ZfsVdev{}).Error; err != nil {
		return fmt.Errorf("failed to delete vdevs for pool %s: %v", poolGuid, err)
	}

	// Delete the pool
	if err := sr.gormClient.WithContext(ctx).Where("pool_guid = ?", poolGuid).Delete(&models.ZfsPool{}).Error; err != nil {
		return fmt.Errorf("failed to delete ZFS pool %s: %v", poolGuid, err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ZFS Dataset
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// RegisterZfsDatasets inserts or updates ZFS datasets in the database
func (sr *scrutinyRepository) RegisterZfsDatasets(ctx context.Context, datasets []models.ZfsDataset) error {
	for _, dataset := range datasets {
		if err := sr.gormClient.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "host_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"type", "pool", "create_txg", "used", "available", "referenced", "mountpoint", "updated_at",
			}),
		}).Create(&dataset).Error; err != nil {
			return fmt.Errorf("failed to register ZFS dataset %s: %v", dataset.Name, err)
		}
	}
	return nil
}

// GetZfsDatasets retrieves all ZFS datasets
func (sr *scrutinyRepository) GetZfsDatasets(ctx context.Context) ([]models.ZfsDataset, error) {
	var datasets []models.ZfsDataset
	if err := sr.gormClient.WithContext(ctx).Order("name ASC").Find(&datasets).Error; err != nil {
		return nil, fmt.Errorf("could not get ZFS datasets from DB: %v", err)
	}
	return datasets, nil
}

// GetZfsDatasetsByPool retrieves ZFS datasets for a specific pool
func (sr *scrutinyRepository) GetZfsDatasetsByPool(ctx context.Context, poolName string) ([]models.ZfsDataset, error) {
	var datasets []models.ZfsDataset
	if err := sr.gormClient.WithContext(ctx).Where("pool = ?", poolName).Order("name ASC").Find(&datasets).Error; err != nil {
		return nil, fmt.Errorf("could not get ZFS datasets for pool %s from DB: %v", poolName, err)
	}
	return datasets, nil
}

// GetZfsDatasetsByHost retrieves ZFS datasets for a specific host
func (sr *scrutinyRepository) GetZfsDatasetsByHost(ctx context.Context, hostId string) ([]models.ZfsDataset, error) {
	var datasets []models.ZfsDataset
	if err := sr.gormClient.WithContext(ctx).Where("host_id = ?", hostId).Order("name ASC").Find(&datasets).Error; err != nil {
		return nil, fmt.Errorf("could not get ZFS datasets for host %s from DB: %v", hostId, err)
	}
	return datasets, nil
}