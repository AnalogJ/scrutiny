package database

import (
	"context"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database/migrations/m20201107210306"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/go-gormigrate/gormigrate/v2"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// SQLite migrations
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//database.AutoMigrate(&models.Device{})

func (sr *scrutinyRepository) Migrate(ctx context.Context) error {

	sr.logger.Infoln("Database migration starting")

	m := gormigrate.New(sr.gormClient, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20201107210306", // v0.3.13 (pre-influxdb schema). 9fac3c6308dc6cb6cd5bbc43a68cd93e8fb20b87
			Migrate: func(tx *gorm.DB) error {
				// it's a good practice to copy the struct inside the function,

				return tx.AutoMigrate(
					&m20201107210306.Device{},
					&m20201107210306.Smart{},
					&m20201107210306.SmartAtaAttribute{},
					&m20201107210306.SmartNvmeAttribute{},
					&m20201107210306.SmartNvmeAttribute{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					&m20201107210306.Device{},
					&m20201107210306.Smart{},
					&m20201107210306.SmartAtaAttribute{},
					&m20201107210306.SmartNvmeAttribute{},
					&m20201107210306.SmartNvmeAttribute{},
					"self_tests",
				)
			},
		},
		{
			ID: "20220503113100", // backwards compatible - influxdb schema
			Migrate: func(tx *gorm.DB) error {
				// delete unnecessary table.
				err := tx.Migrator().DropTable("self_tests")
				if err != nil {
					return err
				}

				//add columns to the Device schema, so we can start adding data to the database & influxdb
				err = tx.Migrator().AddColumn(&models.Device{}, "Label") //Label  string `json:"label"`
				if err != nil {
					return err
				}
				err = tx.Migrator().AddColumn(&models.Device{}, "DeviceStatus") //DeviceStatus pkg.DeviceStatus `json:"device_status"`
				if err != nil {
					return err
				}

				//TODO: migrate the data from GORM to influxdb.

				return nil
			},
		},
		//{
		//	ID: "20220503120000", // v0.4.0 - influxdb schema
		//	Migrate: func(tx *gorm.DB) error {
		//		// delete unnecessary tables.
		//		err := tx.Migrator().DropTable(
		//			&m20201107210306.Smart{},
		//			&m20201107210306.SmartAtaAttribute{},
		//			&m20201107210306.SmartNvmeAttribute{},
		//			&m20201107210306.SmartNvmeAttribute{},
		//		)
		//		if err != nil {
		//			return err
		//		}
		//
		//		//migrate the device database to the final version
		//		return tx.AutoMigrate(models.Device{})
		//	},
		//},
	})

	if err := m.Migrate(); err != nil {
		sr.logger.Errorf("Database migration failed with error: %w", err)
		return err
	}
	sr.logger.Infoln("Database migration completed successfully")
	return nil
}
