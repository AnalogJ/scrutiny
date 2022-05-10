package database

import (
	"context"
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database/migrations/m20201107210306"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database/migrations/m20220503120000"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/go-gormigrate/gormigrate/v2"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/gorm"
	"strconv"
	"time"
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
				//get a list of all devices:
				//	get a list of all smart scans in the last 2 weeks:
				//		get a list of associated smart attribute data:
				//			translate to a measurements.Smart{} object
				//			call CUSTOM INFLUXDB SAVE FUNCTION (taking bucket as parameter)
				//	get a list of all smart scans in the last 9 weeks:
				//		do same as above (select 1 scan per week)
				//	get a list of all smart scans in the last 25 months:
				//		do same as above (select 1 scan per month)
				//	get a list of all smart scans:
				//		do same as above (select 1 scan per year)

				preDevices := []m20201107210306.Device{} //pre-migration device information
				if err = tx.Preload("SmartResults", func(db *gorm.DB) *gorm.DB {
					return db.Order("smarts.created_at ASC") //OLD: .Limit(devicesCount)
				}).Find(&preDevices).Error; err != nil {
					sr.logger.Errorln("Could not get device summary from DB", err)
					return err
				}

				//calculate bucket oldest dates
				today := time.Now()
				dailyBucketMax := today.Add(-RETENTION_PERIOD_15_DAYS_IN_SECONDS * time.Second)     //15 days
				weeklyBucketMax := today.Add(-RETENTION_PERIOD_9_WEEKS_IN_SECONDS * time.Second)    //9 weeks
				monthlyBucketMax := today.Add(-RETENTION_PERIOD_25_MONTHS_IN_SECONDS * time.Second) //25 weeks

				for _, preDevice := range preDevices {
					sr.logger.Debugf("====================================")
					sr.logger.Infof("begin processing device: %s", preDevice.WWN)

					//weekly, monthly, yearly lookup storage, so we don't add more data to the buckets than necessary.
					weeklyLookup := map[string]bool{}
					monthlyLookup := map[string]bool{}
					yearlyLookup := map[string]bool{}
					for _, preSmartResult := range preDevice.SmartResults { //pre-migration smart results

						//we're looping in ASC mode, so from oldest entry to most current.

						err, postSmartResults := m20201107210306_FromPreInfluxDBSmartResultsCreatePostInfluxDBSmartResults(tx, preDevice, preSmartResult)
						if err != nil {
							return err
						}
						smartTags, smartFields := postSmartResults.Flatten()

						err, postSmartTemp := m20201107210306_FromPreInfluxDBTempCreatePostInfluxDBTemp(preDevice, preSmartResult)
						if err != nil {
							return err
						}
						tempTags, tempFields := postSmartTemp.Flatten()
						tempTags["device_wwn"] = preDevice.WWN

						year, week := postSmartResults.Date.ISOWeek()
						month := postSmartResults.Date.Month()

						yearStr := strconv.Itoa(year)
						yearMonthStr := fmt.Sprintf("%d-%d", year, month)
						yearWeekStr := fmt.Sprintf("%d-%d", year, week)

						//write data to daily bucket if in the last 15 days
						if postSmartResults.Date.After(dailyBucketMax) {
							sr.logger.Debugf("device (%s) smart data added to bucket: daily", preDevice.WWN)
							// write point immediately
							err = sr.saveDatapoint(
								sr.influxClient.WriteAPIBlocking(sr.appConfig.GetString("web.influxdb.org"), sr.appConfig.GetString("web.influxdb.bucket")),
								"smart",
								smartTags,
								smartFields,
								postSmartResults.Date, ctx)
							if err != nil {
								return err
							}

							err = sr.saveDatapoint(
								sr.influxClient.WriteAPIBlocking(sr.appConfig.GetString("web.influxdb.org"), sr.appConfig.GetString("web.influxdb.bucket")),
								"temp",
								tempTags,
								tempFields,
								postSmartResults.Date, ctx)
							if err != nil {
								return err
							}
						}

						//write data to the weekly bucket if in the last 9 weeks, and week has not been processed yet
						if _, weekExists := weeklyLookup[yearWeekStr]; !weekExists && postSmartResults.Date.After(weeklyBucketMax) {
							sr.logger.Debugf("device (%s) smart data added to bucket: weekly", preDevice.WWN)

							//this week/year pair has not been processed
							weeklyLookup[yearWeekStr] = true
							// write point immediately
							err = sr.saveDatapoint(
								sr.influxClient.WriteAPIBlocking(sr.appConfig.GetString("web.influxdb.org"), fmt.Sprintf("%s_weekly", sr.appConfig.GetString("web.influxdb.bucket"))),
								"smart",
								smartTags,
								smartFields,
								postSmartResults.Date, ctx)

							if err != nil {
								return err
							}

							err = sr.saveDatapoint(
								sr.influxClient.WriteAPIBlocking(sr.appConfig.GetString("web.influxdb.org"), fmt.Sprintf("%s_weekly", sr.appConfig.GetString("web.influxdb.bucket"))),
								"temp",
								tempTags,
								tempFields,
								postSmartResults.Date, ctx)
							if err != nil {
								return err
							}
						}

						//write data to the monthly bucket if in the last 9 weeks, and week has not been processed yet
						if _, monthExists := monthlyLookup[yearMonthStr]; !monthExists && postSmartResults.Date.After(monthlyBucketMax) {
							sr.logger.Debugf("device (%s) smart data added to bucket: monthly", preDevice.WWN)

							//this month/year pair has not been processed
							monthlyLookup[yearMonthStr] = true
							// write point immediately
							err = sr.saveDatapoint(
								sr.influxClient.WriteAPIBlocking(sr.appConfig.GetString("web.influxdb.org"), fmt.Sprintf("%s_monthly", sr.appConfig.GetString("web.influxdb.bucket"))),
								"smart",
								smartTags,
								smartFields,
								postSmartResults.Date, ctx)
							if err != nil {
								return err
							}

							err = sr.saveDatapoint(
								sr.influxClient.WriteAPIBlocking(sr.appConfig.GetString("web.influxdb.org"), fmt.Sprintf("%s_monthly", sr.appConfig.GetString("web.influxdb.bucket"))),
								"temp",
								tempTags,
								tempFields,
								postSmartResults.Date, ctx)
							if err != nil {
								return err
							}
						}

						if _, yearExists := yearlyLookup[yearStr]; !yearExists && year != today.Year() {
							sr.logger.Debugf("device (%s) smart data added to bucket: yearly", preDevice.WWN)

							//this year has not been processed
							yearlyLookup[yearStr] = true
							// write point immediately
							err = sr.saveDatapoint(
								sr.influxClient.WriteAPIBlocking(sr.appConfig.GetString("web.influxdb.org"), fmt.Sprintf("%s_yearly", sr.appConfig.GetString("web.influxdb.bucket"))),
								"smart",
								smartTags,
								smartFields,
								postSmartResults.Date, ctx)
							if err != nil {
								return err
							}

							err = sr.saveDatapoint(
								sr.influxClient.WriteAPIBlocking(sr.appConfig.GetString("web.influxdb.org"), fmt.Sprintf("%s_yearly", sr.appConfig.GetString("web.influxdb.bucket"))),
								"temp",
								tempTags,
								tempFields,
								postSmartResults.Date, ctx)
							if err != nil {
								return err
							}
						}
					}
					sr.logger.Infof("finished processing device %s. weekly: %d, monthly: %d, yearly: %d", preDevice.WWN, len(weeklyLookup), len(monthlyLookup), len(yearlyLookup))

				}

				return nil
			},
		},
		{
			ID: "20220503120000", // cleanup - v0.4.0 - influxdb schema
			Migrate: func(tx *gorm.DB) error {
				// delete unnecessary tables.
				err := tx.Migrator().DropTable(
					&m20201107210306.Smart{},
					&m20201107210306.SmartAtaAttribute{},
					&m20201107210306.SmartNvmeAttribute{},
					&m20201107210306.SmartScsiAttribute{},
				)
				if err != nil {
					return err
				}

				//migrate the device database to the current version
				return tx.AutoMigrate(m20220503120000.Device{})
			},
		},
	})

	if err := m.Migrate(); err != nil {
		sr.logger.Errorf("Database migration failed with error: %w", err)
		return err
	}
	sr.logger.Infoln("Database migration completed successfully")
	return nil
}

// Deprecated
func m20201107210306_FromPreInfluxDBTempCreatePostInfluxDBTemp(preDevice m20201107210306.Device, preSmartResult m20201107210306.Smart) (error, measurements.SmartTemperature) {
	//extract temperature data for every datapoint
	postSmartTemp := measurements.SmartTemperature{
		Date: preSmartResult.TestDate,
		Temp: preSmartResult.Temp,
	}

	return nil, postSmartTemp
}

// Deprecated
func m20201107210306_FromPreInfluxDBSmartResultsCreatePostInfluxDBSmartResults(database *gorm.DB, preDevice m20201107210306.Device, preSmartResult m20201107210306.Smart) (error, measurements.Smart) {
	//create a measurements.Smart object (which we will then push to the InfluxDB)
	postDeviceSmartData := measurements.Smart{
		Date:            preSmartResult.TestDate,
		DeviceWWN:       preDevice.WWN,
		DeviceProtocol:  preDevice.DeviceProtocol,
		Temp:            preSmartResult.Temp,
		PowerOnHours:    preSmartResult.PowerOnHours,
		PowerCycleCount: preSmartResult.PowerCycleCount,

		// this needs to be populated using measurements.Smart.ProcessAtaSmartInfo, ProcessScsiSmartInfo or ProcessNvmeSmartInfo
		// because those functions will take into account thresholds (which we didn't consider correctly previously)
		Attributes: map[string]measurements.SmartAttribute{},
	}

	result := database.Preload("AtaAttributes").Preload("NvmeAttributes").Preload("ScsiAttributes").Find(&preSmartResult)
	if result.Error != nil {
		return result.Error, postDeviceSmartData
	}

	if preDevice.IsAta() {
		preAtaSmartAttributesTable := []collector.AtaSmartAttributesTableItem{}
		for _, preAtaAttribute := range preSmartResult.AtaAttributes {
			preAtaSmartAttributesTable = append(preAtaSmartAttributesTable, collector.AtaSmartAttributesTableItem{
				ID:         preAtaAttribute.AttributeId,
				Name:       preAtaAttribute.Name,
				Value:      int64(preAtaAttribute.Value),
				Worst:      int64(preAtaAttribute.Worst),
				Thresh:     int64(preAtaAttribute.Threshold),
				WhenFailed: preAtaAttribute.WhenFailed,
				Flags: struct {
					Value         int    `json:"value"`
					String        string `json:"string"`
					Prefailure    bool   `json:"prefailure"`
					UpdatedOnline bool   `json:"updated_online"`
					Performance   bool   `json:"performance"`
					ErrorRate     bool   `json:"error_rate"`
					EventCount    bool   `json:"event_count"`
					AutoKeep      bool   `json:"auto_keep"`
				}{
					Value:         0,
					String:        "",
					Prefailure:    false,
					UpdatedOnline: false,
					Performance:   false,
					ErrorRate:     false,
					EventCount:    false,
					AutoKeep:      false,
				},
				Raw: struct {
					Value  int64  `json:"value"`
					String string `json:"string"`
				}{
					Value:  preAtaAttribute.RawValue,
					String: preAtaAttribute.RawString,
				},
			})
		}

		postDeviceSmartData.ProcessAtaSmartInfo(preAtaSmartAttributesTable)

	} else if preDevice.IsNvme() {
		//info collector.SmartInfo
		postNvmeSmartHealthInformation := collector.NvmeSmartHealthInformationLog{}

		for _, preNvmeAttribute := range preSmartResult.NvmeAttributes {
			switch preNvmeAttribute.AttributeId {
			case "critical_warning":
				postNvmeSmartHealthInformation.CriticalWarning = int64(preNvmeAttribute.Value)
			case "temperature":
				postNvmeSmartHealthInformation.Temperature = int64(preNvmeAttribute.Value)
			case "available_spare":
				postNvmeSmartHealthInformation.AvailableSpare = int64(preNvmeAttribute.Value)
			case "available_spare_threshold":
				postNvmeSmartHealthInformation.AvailableSpareThreshold = int64(preNvmeAttribute.Value)
			case "percentage_used":
				postNvmeSmartHealthInformation.PercentageUsed = int64(preNvmeAttribute.Value)
			case "data_units_read":
				postNvmeSmartHealthInformation.DataUnitsWritten = int64(preNvmeAttribute.Value)
			case "data_units_written":
				postNvmeSmartHealthInformation.DataUnitsWritten = int64(preNvmeAttribute.Value)
			case "host_reads":
				postNvmeSmartHealthInformation.HostReads = int64(preNvmeAttribute.Value)
			case "host_writes":
				postNvmeSmartHealthInformation.HostWrites = int64(preNvmeAttribute.Value)
			case "controller_busy_time":
				postNvmeSmartHealthInformation.ControllerBusyTime = int64(preNvmeAttribute.Value)
			case "power_cycles":
				postNvmeSmartHealthInformation.PowerCycles = int64(preNvmeAttribute.Value)
			case "power_on_hours":
				postNvmeSmartHealthInformation.PowerOnHours = int64(preNvmeAttribute.Value)
			case "unsafe_shutdowns":
				postNvmeSmartHealthInformation.UnsafeShutdowns = int64(preNvmeAttribute.Value)
			case "media_errors":
				postNvmeSmartHealthInformation.MediaErrors = int64(preNvmeAttribute.Value)
			case "num_err_log_entries":
				postNvmeSmartHealthInformation.NumErrLogEntries = int64(preNvmeAttribute.Value)
			case "warning_temp_time":
				postNvmeSmartHealthInformation.WarningTempTime = int64(preNvmeAttribute.Value)
			case "critical_comp_time":
				postNvmeSmartHealthInformation.CriticalCompTime = int64(preNvmeAttribute.Value)
			}
		}

		postDeviceSmartData.ProcessNvmeSmartInfo(postNvmeSmartHealthInformation)

	} else if preDevice.IsScsi() {
		//info collector.SmartInfo
		var postScsiGrownDefectList int64
		postScsiErrorCounterLog := collector.ScsiErrorCounterLog{
			Read: struct {
				ErrorsCorrectedByEccfast         int64  `json:"errors_corrected_by_eccfast"`
				ErrorsCorrectedByEccdelayed      int64  `json:"errors_corrected_by_eccdelayed"`
				ErrorsCorrectedByRereadsRewrites int64  `json:"errors_corrected_by_rereads_rewrites"`
				TotalErrorsCorrected             int64  `json:"total_errors_corrected"`
				CorrectionAlgorithmInvocations   int64  `json:"correction_algorithm_invocations"`
				GigabytesProcessed               string `json:"gigabytes_processed"`
				TotalUncorrectedErrors           int64  `json:"total_uncorrected_errors"`
			}{},
			Write: struct {
				ErrorsCorrectedByEccfast         int64  `json:"errors_corrected_by_eccfast"`
				ErrorsCorrectedByEccdelayed      int64  `json:"errors_corrected_by_eccdelayed"`
				ErrorsCorrectedByRereadsRewrites int64  `json:"errors_corrected_by_rereads_rewrites"`
				TotalErrorsCorrected             int64  `json:"total_errors_corrected"`
				CorrectionAlgorithmInvocations   int64  `json:"correction_algorithm_invocations"`
				GigabytesProcessed               string `json:"gigabytes_processed"`
				TotalUncorrectedErrors           int64  `json:"total_uncorrected_errors"`
			}{},
		}

		for _, preScsiAttribute := range preSmartResult.ScsiAttributes {
			switch preScsiAttribute.AttributeId {
			case "scsi_grown_defect_list":
				postScsiGrownDefectList = int64(preScsiAttribute.Value)
			case "read.errors_corrected_by_eccfast":
				postScsiErrorCounterLog.Read.ErrorsCorrectedByEccfast = int64(preScsiAttribute.Value)
			case "read.errors_corrected_by_eccdelayed":
				postScsiErrorCounterLog.Read.ErrorsCorrectedByEccdelayed = int64(preScsiAttribute.Value)
			case "read.errors_corrected_by_rereads_rewrites":
				postScsiErrorCounterLog.Read.ErrorsCorrectedByRereadsRewrites = int64(preScsiAttribute.Value)
			case "read.total_errors_corrected":
				postScsiErrorCounterLog.Read.TotalErrorsCorrected = int64(preScsiAttribute.Value)
			case "read.correction_algorithm_invocations":
				postScsiErrorCounterLog.Read.CorrectionAlgorithmInvocations = int64(preScsiAttribute.Value)
			case "read.total_uncorrected_errors":
				postScsiErrorCounterLog.Read.TotalUncorrectedErrors = int64(preScsiAttribute.Value)
			case "write.errors_corrected_by_eccfast":
				postScsiErrorCounterLog.Write.ErrorsCorrectedByEccfast = int64(preScsiAttribute.Value)
			case "write.errors_corrected_by_eccdelayed":
				postScsiErrorCounterLog.Write.ErrorsCorrectedByEccdelayed = int64(preScsiAttribute.Value)
			case "write.errors_corrected_by_rereads_rewrites":
				postScsiErrorCounterLog.Write.ErrorsCorrectedByRereadsRewrites = int64(preScsiAttribute.Value)
			case "write.total_errors_corrected":
				postScsiErrorCounterLog.Write.TotalErrorsCorrected = int64(preScsiAttribute.Value)
			case "write.correction_algorithm_invocations":
				postScsiErrorCounterLog.Write.CorrectionAlgorithmInvocations = int64(preScsiAttribute.Value)
			case "write.total_uncorrected_errors":
				postScsiErrorCounterLog.Write.TotalUncorrectedErrors = int64(preScsiAttribute.Value)
			}
		}
		postDeviceSmartData.ProcessScsiSmartInfo(postScsiGrownDefectList, postScsiErrorCounterLog)
	} else {
		return fmt.Errorf("Unknown device protocol: %s", preDevice.DeviceProtocol), postDeviceSmartData
	}

	return nil, postDeviceSmartData
}
