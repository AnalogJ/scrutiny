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
				"scan_processed", "scan_errors", "scan_issued", "updated_at",
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