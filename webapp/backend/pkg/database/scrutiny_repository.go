package database

import (
	"context"
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

const (
	// 60seconds * 60minutes * 24hours * 15 days
	RETENTION_PERIOD_15_DAYS_IN_SECONDS = 1_296_000

	// 60seconds * 60minutes * 24hours * 7 days * 9 weeks
	RETENTION_PERIOD_9_WEEKS_IN_SECONDS = 5_443_200

	// 60seconds * 60minutes * 24hours * 7 days * (52 + 52 + 4)weeks
	RETENTION_PERIOD_25_MONTHS_IN_SECONDS = 65_318_400

	DURATION_KEY_WEEK    = "week"
	DURATION_KEY_MONTH   = "month"
	DURATION_KEY_YEAR    = "year"
	DURATION_KEY_FOREVER = "forever"
)

//// GormLogger is a custom logger for Gorm, making it use logrus.
//type GormLogger struct{ Logger logrus.FieldLogger }
//
//// Print handles log events from Gorm for the custom logger.
//func (gl *GormLogger) Print(v ...interface{}) {
//	switch v[0] {
//	case "sql":
//		gl.Logger.WithFields(
//			logrus.Fields{
//				"module":  "gorm",
//				"type":    "sql",
//				"rows":    v[5],
//				"src_ref": v[1],
//				"values":  v[4],
//			},
//		).Debug(v[3])
//	case "log":
//		gl.Logger.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Print(v[2])
//	}
//}

func NewScrutinyRepository(appConfig config.Interface, globalLogger logrus.FieldLogger) (DeviceRepo, error) {
	backgroundContext := context.Background()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Gorm/SQLite setup
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	fmt.Printf("Trying to connect to database stored: %s\n", appConfig.GetString("web.database.location"))
	database, err := gorm.Open(sqlite.Open(appConfig.GetString("web.database.location")), &gorm.Config{
		//TODO: figure out how to log database queries again.
		//Logger: logger
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database! - %v", err)
	}

	//database.SetLogger()
	database.AutoMigrate(&models.Device{})

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// InfluxDB setup
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// Create a new client using an InfluxDB server base URL and an authentication token
	influxdbUrl := fmt.Sprintf("http://%s:%s", appConfig.GetString("web.influxdb.host"), appConfig.GetString("web.influxdb.port"))
	globalLogger.Debugf("InfluxDB url: %s", influxdbUrl)

	client := influxdb2.NewClient(influxdbUrl, appConfig.GetString("web.influxdb.token"))

	if !appConfig.IsSet("web.influxdb.token") {
		globalLogger.Debugf("No influxdb token found, running first-time setup...")

		// if no token is provided, but we have a valid server, we're going to assume this is the first setup of our server.
		// we will initialize with a predetermined username & password, that you should change.

		// metrics bucket will have a retention period of 8 days (since it will be down-sampled once a week)
		// in seconds (60seconds * 60minutes * 24hours * 15 days) = 1_296_000 (see EnsureBucket() function)
		onboardingResponse, err := client.Setup(
			backgroundContext,
			appConfig.GetString("web.influxdb.init_username"),
			appConfig.GetString("web.influxdb.init_password"),
			appConfig.GetString("web.influxdb.org"),
			appConfig.GetString("web.influxdb.bucket"),
			0)
		if err != nil {
			return nil, err
		}

		appConfig.Set("web.influxdb.token", *onboardingResponse.Auth.Token)
		//todo: determine if we should write the config file out here.
	}

	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(appConfig.GetString("web.influxdb.org"), appConfig.GetString("web.influxdb.bucket"))

	// Get query client
	queryAPI := client.QueryAPI(appConfig.GetString("web.influxdb.org"))

	// Get task client
	taskAPI := client.TasksAPI()

	if writeAPI == nil || queryAPI == nil || taskAPI == nil {
		return nil, fmt.Errorf("Failed to connect to influxdb!")
	}

	deviceRepo := scrutinyRepository{
		appConfig:      appConfig,
		logger:         globalLogger,
		influxClient:   client,
		influxWriteApi: writeAPI,
		influxQueryApi: queryAPI,
		influxTaskApi:  taskAPI,
		gormClient:     database,
	}

	orgInfo, err := client.OrganizationsAPI().FindOrganizationByName(backgroundContext, appConfig.GetString("web.influxdb.org"))
	if err != nil {
		return nil, err
	}

	// Initialize Buckets (if necessary)
	err = deviceRepo.EnsureBuckets(backgroundContext, orgInfo)
	if err != nil {
		return nil, err
	}

	// Initialize Background Tasks
	err = deviceRepo.EnsureTasks(backgroundContext, *orgInfo.Id)
	if err != nil {
		return nil, err
	}
	return &deviceRepo, nil
}

type scrutinyRepository struct {
	appConfig config.Interface
	logger    logrus.FieldLogger

	influxWriteApi api.WriteAPIBlocking
	influxQueryApi api.QueryAPI
	influxTaskApi  api.TasksAPI
	influxClient   influxdb2.Client

	gormClient *gorm.DB
}

func (sr *scrutinyRepository) Close() error {
	sr.influxClient.Close()
	return nil
}

func (sr *scrutinyRepository) EnsureBuckets(ctx context.Context, org *domain.Organization) error {

	var mainBucketRetentionRule domain.RetentionRule
	var weeklyBucketRetentionRule domain.RetentionRule
	var monthlyBucketRetentionRule domain.RetentionRule
	if sr.appConfig.GetBool("web.influxdb.retention_policy") {

		// for tetsting purposes, we may not want to set a retention policy, this will allow to set data with old timestamps,
		//then manually run the downsampling scripts
		mainBucketRetentionRule = domain.RetentionRule{EverySeconds: RETENTION_PERIOD_15_DAYS_IN_SECONDS}
		weeklyBucketRetentionRule = domain.RetentionRule{EverySeconds: RETENTION_PERIOD_9_WEEKS_IN_SECONDS}
		monthlyBucketRetentionRule = domain.RetentionRule{EverySeconds: RETENTION_PERIOD_25_MONTHS_IN_SECONDS}
	}

	mainBucket := sr.appConfig.GetString("web.influxdb.bucket")
	if foundMainBucket, foundErr := sr.influxClient.BucketsAPI().FindBucketByName(ctx, mainBucket); foundErr != nil {
		// metrics bucket will have a retention period of 15 days (since it will be down-sampled once a week)
		_, err := sr.influxClient.BucketsAPI().CreateBucketWithName(ctx, org, mainBucket, mainBucketRetentionRule)
		if err != nil {
			return err
		}
	} else if sr.appConfig.GetBool("web.influxdb.retention_policy") {
		//correctly set the retention period for the main bucket (cant do it during setup/creation)
		foundMainBucket.RetentionRules = domain.RetentionRules{mainBucketRetentionRule}
		sr.influxClient.BucketsAPI().UpdateBucket(ctx, foundMainBucket)
	}

	//create buckets (used for downsampling)
	weeklyBucket := fmt.Sprintf("%s_weekly", sr.appConfig.GetString("web.influxdb.bucket"))
	if _, foundErr := sr.influxClient.BucketsAPI().FindBucketByName(ctx, weeklyBucket); foundErr != nil {
		// metrics_weekly bucket will have a retention period of 8+1 weeks (since it will be down-sampled once a month)
		_, err := sr.influxClient.BucketsAPI().CreateBucketWithName(ctx, org, weeklyBucket, weeklyBucketRetentionRule)
		if err != nil {
			return err
		}
	}

	monthlyBucket := fmt.Sprintf("%s_monthly", sr.appConfig.GetString("web.influxdb.bucket"))
	if _, foundErr := sr.influxClient.BucketsAPI().FindBucketByName(ctx, monthlyBucket); foundErr != nil {
		// metrics_monthly bucket will have a retention period of 24+1 months (since it will be down-sampled once a year)
		_, err := sr.influxClient.BucketsAPI().CreateBucketWithName(ctx, org, monthlyBucket, monthlyBucketRetentionRule)
		if err != nil {
			return err
		}
	}

	yearlyBucket := fmt.Sprintf("%s_yearly", sr.appConfig.GetString("web.influxdb.bucket"))
	if _, foundErr := sr.influxClient.BucketsAPI().FindBucketByName(ctx, yearlyBucket); foundErr != nil {
		// metrics_yearly bucket will have an infinite retention period
		_, err := sr.influxClient.BucketsAPI().CreateBucketWithName(ctx, org, yearlyBucket)
		if err != nil {
			return err
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tasks
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (sr *scrutinyRepository) EnsureTasks(ctx context.Context, orgID string) error {
	weeklyTaskName := "tsk-weekly-aggr"
	if found, findErr := sr.influxTaskApi.FindTasks(ctx, &api.TaskFilter{Name: weeklyTaskName}); findErr == nil && len(found) == 0 {
		//weekly on Sunday at 1:00am
		_, err := sr.influxTaskApi.CreateTaskWithCron(ctx, weeklyTaskName, sr.DownsampleScript("weekly"), "0 1 * * 0", orgID)
		if err != nil {
			return err
		}
	}

	monthlyTaskName := "tsk-monthly-aggr"
	if found, findErr := sr.influxTaskApi.FindTasks(ctx, &api.TaskFilter{Name: monthlyTaskName}); findErr == nil && len(found) == 0 {
		//monthly on first day of the month at 1:30am
		_, err := sr.influxTaskApi.CreateTaskWithCron(ctx, monthlyTaskName, sr.DownsampleScript("monthly"), "30 1 1 * *", orgID)
		if err != nil {
			return err
		}
	}

	yearlyTaskName := "tsk-yearly-aggr"
	if found, findErr := sr.influxTaskApi.FindTasks(ctx, &api.TaskFilter{Name: yearlyTaskName}); findErr == nil && len(found) == 0 {
		//yearly on the first day of the year at 2:00am
		_, err := sr.influxTaskApi.CreateTaskWithCron(ctx, yearlyTaskName, sr.DownsampleScript("yearly"), "0 2 1 1 *", orgID)
		if err != nil {
			return err
		}
	}
	return nil
}

/*


sourceBucket = "metrics"
rangeStart = -2w
rangeEnd = -1w
aggWindow = 1w
destBucket = "metrics_weekly"
destOrg = "scrutiny"

smart_data = from(bucket: sourceBucket)
|> range(start: rangeStart, stop: rangeEnd)
|> filter(fn: (r) => r["_measurement"] == "smart" )
|> filter(fn: (r) => r["_field"] !~ /(_measurement|device_protocol|device_wwn|attribute_id|raw_string|status_reason|when_failed)/)
|> yield(name: "last")

smart_data
|> aggregateWindow(fn: mean, every: aggWindow)
|> to(bucket: destBucket, org: destOrg)

temp_data = from(bucket: sourceBucket)
|> range(start: rangeStart, stop: rangeEnd)
|> filter(fn: (r) => r["_measurement"] == "temp")
|> toInt()
|> yield(name: "mean")

temp_data
|> aggregateWindow(fn: mean, every: aggWindow)
|> to(bucket: destBucket, org: destOrg)

*/
func (sr *scrutinyRepository) DownsampleScript(aggregationType string) string {
	var sourceBucket string // the source of the data
	var destBucket string   // the destination for the aggregated data
	var rangeStart string
	var rangeEnd string
	var aggWindow string
	switch aggregationType {
	case "weekly":
		sourceBucket = sr.appConfig.GetString("web.influxdb.bucket")
		destBucket = fmt.Sprintf("%s_weekly", sr.appConfig.GetString("web.influxdb.bucket"))
		rangeStart = "-2w"
		rangeEnd = "-1w"
		aggWindow = "1w"
	case "monthly":
		sourceBucket = fmt.Sprintf("%s_weekly", sr.appConfig.GetString("web.influxdb.bucket"))
		destBucket = fmt.Sprintf("%s_monthly", sr.appConfig.GetString("web.influxdb.bucket"))
		rangeStart = "-2mo"
		rangeEnd = "-1mo"
		aggWindow = "1mo"
	case "yearly":
		sourceBucket = fmt.Sprintf("%s_monthly", sr.appConfig.GetString("web.influxdb.bucket"))
		destBucket = fmt.Sprintf("%s_yearly", sr.appConfig.GetString("web.influxdb.bucket"))
		rangeStart = "-2y"
		rangeEnd = "-1y"
		aggWindow = "1y"
	}

	return fmt.Sprintf(`import "types"
  sourceBucket = "%s"
  rangeStart = %s
  rangeEnd = %s
  aggWindow = %s
  destBucket = "%s"
  destOrg = "%s"

  smart_data = from(bucket: sourceBucket)
  |> range(start: rangeStart, stop: rangeEnd)
  |> filter(fn: (r) => r["_measurement"] == "smart" )
  |> group(columns: ["device_wwn", "_field"])

  non_numeric_smart_data = smart_data
    |> filter(fn: (r) => types.isType(v: r._value, type: "string") or types.isType(v: r._value, type: "bool"))
    |> aggregateWindow(every: aggWindow, fn: last, createEmpty: false)

  numeric_smart_data = smart_data
    |> filter(fn: (r) => types.isType(v: r._value, type: "int") or types.isType(v: r._value, type: "float"))
    |> aggregateWindow(every: aggWindow, fn: mean, createEmpty: false)

  union(tables: [non_numeric_smart_data, numeric_smart_data])
  |> to(bucket: destBucket, org: destOrg)

  temp_data = from(bucket: sourceBucket)
  |> range(start: rangeStart, stop: rangeEnd)
  |> filter(fn: (r) => r["_measurement"] == "temp")
  |> group(columns: ["device_wwn"])
  |> toInt()

  temp_data
  |> aggregateWindow(fn: mean, every: aggWindow)
  |> to(bucket: destBucket, org: destOrg)
		`,
		sourceBucket,
		rangeStart,
		rangeEnd,
		aggWindow,
		destBucket,
		sr.appConfig.GetString("web.influxdb.org"),
	)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Device
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//insert device into DB (and update specified columns if device is already registered)
// update device fields that may change: (DeviceType, HostID)
func (sr *scrutinyRepository) RegisterDevice(ctx context.Context, dev models.Device) error {
	if err := sr.gormClient.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "wwn"}},
		DoUpdates: clause.AssignmentColumns([]string{"host_id", "device_name", "device_type"}),
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

//Update Device Status
func (sr *scrutinyRepository) UpdateDeviceStatus(ctx context.Context, wwn string, status pkg.DeviceStatus) (models.Device, error) {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("wwn = ?", wwn).First(&device).Error; err != nil {
		return device, fmt.Errorf("Could not get device from DB: %v", err)
	}

	device.DeviceStatus = pkg.Set(device.DeviceStatus, status)
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// SMART
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (sr *scrutinyRepository) SaveSmartAttributes(ctx context.Context, wwn string, collectorSmartData collector.SmartInfo) (measurements.Smart, error) {
	deviceSmartData := measurements.Smart{}
	err := deviceSmartData.FromCollectorSmartInfo(wwn, collectorSmartData)
	if err != nil {
		sr.logger.Errorln("Could not process SMART metrics", err)
		return measurements.Smart{}, err
	}

	tags, fields := deviceSmartData.Flatten()
	p := influxdb2.NewPoint("smart",
		tags,
		fields,
		deviceSmartData.Date)

	// write point immediately
	return deviceSmartData, sr.influxWriteApi.WritePoint(ctx, p)
}

func (sr *scrutinyRepository) GetSmartAttributeHistory(ctx context.Context, wwn string, startAt string, attributes []string) ([]measurements.Smart, error) {
	// Get SMartResults from InfluxDB

	fmt.Println("GetDeviceDetails from INFLUXDB")

	//TODO: change the filter startrange to a real number.

	// Get parser flux query result
	//appConfig.GetString("web.influxdb.bucket")
	queryStr := fmt.Sprintf(`
  import "influxdata/influxdb/schema"
  from(bucket: "%s")
  |> range(start: -2y, stop: now())
  |> filter(fn: (r) => r["_measurement"] == "smart" )
  |> filter(fn: (r) => r["device_wwn"] == "%s" )
  |> schema.fieldsAsCols()
  |> group(columns: ["device_wwn"])
  |> yield(name: "last")
		`,
		sr.appConfig.GetString("web.influxdb.bucket"),
		wwn,
	)

	smartResults := []measurements.Smart{}

	result, err := sr.influxQueryApi.Query(ctx, queryStr)
	if err == nil {
		fmt.Println("GetDeviceDetails NO EROR")

		// Use Next() to iterate over query result lines
		for result.Next() {
			fmt.Println("GetDeviceDetails NEXT")

			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				//fmt.Printf("table: %s\n", result.TableMetadata().String())
			}

			fmt.Printf("DECODINIG TABLE VALUES: %v", result.Record().Values())
			smartData, err := measurements.NewSmartFromInfluxDB(result.Record().Values())
			if err != nil {
				return nil, err
			}
			smartResults = append(smartResults, *smartData)

		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	} else {
		return nil, err
	}

	return smartResults, nil

	//if err := device.SquashHistory(); err != nil {
	//	logger.Errorln("An error occurred while squashing device history", err)
	//	c.JSON(http.StatusInternalServerError, gin.H{"success": false})
	//	return
	//}
	//
	//if err := device.ApplyMetadataRules(); err != nil {
	//	logger.Errorln("An error occurred while applying scrutiny thresholds & rules", err)
	//	c.JSON(http.StatusInternalServerError, gin.H{"success": false})
	//	return
	//}

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Temperature Data
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (sr *scrutinyRepository) SaveSmartTemperature(ctx context.Context, wwn string, deviceProtocol string, collectorSmartData collector.SmartInfo) error {
	if len(collectorSmartData.AtaSctTemperatureHistory.Table) > 0 {

		for ndx, temp := range collectorSmartData.AtaSctTemperatureHistory.Table {

			minutesOffset := collectorSmartData.AtaSctTemperatureHistory.LoggingIntervalMinutes * int64(ndx) * 60
			smartTemp := measurements.SmartTemperature{
				Date: time.Unix(collectorSmartData.LocalTime.TimeT-minutesOffset, 0),
				Temp: temp,
			}

			tags, fields := smartTemp.Flatten()
			tags["device_wwn"] = wwn
			p := influxdb2.NewPoint("temp",
				tags,
				fields,
				smartTemp.Date)
			err := sr.influxWriteApi.WritePoint(ctx, p)
			if err != nil {
				return err
			}
		}
		// also add the current temperature.
	} else {

		smartTemp := measurements.SmartTemperature{
			Date: time.Unix(collectorSmartData.LocalTime.TimeT, 0),
			Temp: collectorSmartData.Temperature.Current,
		}

		tags, fields := smartTemp.Flatten()
		tags["device_wwn"] = wwn
		p := influxdb2.NewPoint("temp",
			tags,
			fields,
			smartTemp.Date)
		return sr.influxWriteApi.WritePoint(ctx, p)
	}
	return nil
}

func (sr *scrutinyRepository) GetSmartTemperatureHistory(ctx context.Context, durationKey string) (map[string][]measurements.SmartTemperature, error) {
	//we can get temp history for "week", "month", DURATION_KEY_YEAR, "forever"

	deviceTempHistory := map[string][]measurements.SmartTemperature{}

	//TODO: change the query range to a variable.
	queryStr := sr.aggregateTempQuery(durationKey)

	result, err := sr.influxQueryApi.Query(ctx, queryStr)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {

			if deviceWWN, ok := result.Record().Values()["device_wwn"]; ok {

				//check if deviceWWN has been seen and initialized already
				if _, ok := deviceTempHistory[deviceWWN.(string)]; !ok {
					deviceTempHistory[deviceWWN.(string)] = []measurements.SmartTemperature{}
				}

				currentTempHistory := deviceTempHistory[deviceWWN.(string)]
				smartTemp := measurements.SmartTemperature{}

				for key, val := range result.Record().Values() {
					smartTemp.Inflate(key, val)
				}
				smartTemp.Date = result.Record().Values()["_time"].(time.Time)
				currentTempHistory = append(currentTempHistory, smartTemp)
				deviceTempHistory[deviceWWN.(string)] = currentTempHistory
			}
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	} else {
		return nil, err
	}
	return deviceTempHistory, nil

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// DeviceSummary
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// get a map of all devices and associated SMART data
func (sr *scrutinyRepository) GetSummary(ctx context.Context) (map[string]*models.DeviceSummary, error) {
	devices, err := sr.GetDevices(ctx)
	if err != nil {
		return nil, err
	}

	summaries := map[string]*models.DeviceSummary{}

	for _, device := range devices {
		summaries[device.WWN] = &models.DeviceSummary{Device: device}
	}

	// Get parser flux query result
	//appConfig.GetString("web.influxdb.bucket")
	queryStr := fmt.Sprintf(`
  import "influxdata/influxdb/schema"
  from(bucket: "%s")
  |> range(start: -1y, stop: now())
  |> filter(fn: (r) => r["_measurement"] == "smart" )
  |> filter(fn: (r) => r["_field"] == "temp" or r["_field"] == "power_on_hours" or r["_field"] == "date")
  |> last()
  |> schema.fieldsAsCols()
  |> group(columns: ["device_wwn"])
  |> yield(name: "last")
		`,
		sr.appConfig.GetString("web.influxdb.bucket"),
	)

	result, err := sr.influxQueryApi.Query(ctx, queryStr)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				//fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result

			//get summary data from Influxdb.
			//result.Record().Values()
			if deviceWWN, ok := result.Record().Values()["device_wwn"]; ok {

				//ensure summaries is intialized for this wwn
				if _, exists := summaries[deviceWWN.(string)]; !exists {
					summaries[deviceWWN.(string)] = &models.DeviceSummary{}
				}

				summaries[deviceWWN.(string)].SmartResults = &models.SmartSummary{
					Temp:          result.Record().Values()["temp"].(int64),
					PowerOnHours:  result.Record().Values()["power_on_hours"].(int64),
					CollectorDate: result.Record().Values()["_time"].(time.Time),
				}
			}
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	} else {
		return nil, err
	}

	deviceTempHistory, err := sr.GetSmartTemperatureHistory(ctx, DURATION_KEY_WEEK)
	if err != nil {
		sr.logger.Printf("========================>>>>>>>>======================")
		sr.logger.Printf("========================>>>>>>>>======================")
		sr.logger.Printf("========================>>>>>>>>======================")
		sr.logger.Printf("========================>>>>>>>>======================")
		sr.logger.Printf("========================>>>>>>>>======================")
		sr.logger.Printf("Error: %v", err)
	}
	for wwn, tempHistory := range deviceTempHistory {
		summaries[wwn].TempHistory = tempHistory
	}

	return summaries, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper Methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (sr *scrutinyRepository) lookupBucketName(durationKey string) string {
	switch durationKey {
	case DURATION_KEY_WEEK:
		//data stored in the last week
		return sr.appConfig.GetString("web.influxdb.bucket")
	case DURATION_KEY_MONTH:
		// data stored in the last month (after the first week)
		return fmt.Sprintf("%s_weekly", sr.appConfig.GetString("web.influxdb.bucket"))
	case DURATION_KEY_YEAR:
		// data stored in the last year (after the first month)
		return fmt.Sprintf("%s_monthly", sr.appConfig.GetString("web.influxdb.bucket"))
	case DURATION_KEY_FOREVER:
		//data stored before the last year
		return fmt.Sprintf("%s_yearly", sr.appConfig.GetString("web.influxdb.bucket"))
	}
	return sr.appConfig.GetString("web.influxdb.bucket")
}

func (sr *scrutinyRepository) lookupDuration(durationKey string) []string {

	switch durationKey {
	case DURATION_KEY_WEEK:
		//data stored in the last week
		return []string{"-1w", "now()"}
	case DURATION_KEY_MONTH:
		// data stored in the last month (after the first week)
		return []string{"-1mo", "-1w"}
	case DURATION_KEY_YEAR:
		// data stored in the last year (after the first month)
		return []string{"-1y", "-1mo"}
	case DURATION_KEY_FOREVER:
		//data stored before the last year
		return []string{"-10y", "-1y"}
	}
	return []string{"-1w", "now()"}
}

func (sr *scrutinyRepository) lookupNestedDurationKeys(durationKey string) []string {
	switch durationKey {
	case DURATION_KEY_WEEK:
		//all data is stored in a single bucket
		return []string{DURATION_KEY_WEEK}
	case DURATION_KEY_MONTH:
		//data is stored in the week bucket and the month bucket
		return []string{DURATION_KEY_WEEK, DURATION_KEY_MONTH}
	case DURATION_KEY_YEAR:
		// data stored in the last year (after the first month)
		return []string{DURATION_KEY_WEEK, DURATION_KEY_MONTH, DURATION_KEY_YEAR}
	case DURATION_KEY_FOREVER:
		//data stored before the last year
		return []string{DURATION_KEY_WEEK, DURATION_KEY_MONTH, DURATION_KEY_YEAR}
	}
	return []string{DURATION_KEY_WEEK}
}

func (sr *scrutinyRepository) aggregateTempQuery(durationKey string) string {

	/*
		import "influxdata/influxdb/schema"
		weekData = from(bucket: "metrics")
		  |> range(start: -1w, stop: now())
		  |> filter(fn: (r) => r["_measurement"] == "temp" )
		  |> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
		  |> group(columns: ["device_wwn"])
		  |> toInt()

		monthData = from(bucket: "metrics_weekly")
		  |> range(start: -1mo, stop: now())
		  |> filter(fn: (r) => r["_measurement"] == "temp" )
		  |> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
		  |> group(columns: ["device_wwn"])
		  |> toInt()

		union(tables: [weekData, monthData])
		  |> group(columns: ["device_wwn"])
		  |> sort(columns: ["_time"], desc: false)
		  |> schema.fieldsAsCols()

	*/

	partialQueryStr := []string{
		`import "influxdata/influxdb/schema"`,
	}

	nestedDurationKeys := sr.lookupNestedDurationKeys(durationKey)

	subQueryNames := []string{}
	for _, nestedDurationKey := range nestedDurationKeys {
		bucketName := sr.lookupBucketName(nestedDurationKey)
		durationRange := sr.lookupDuration(nestedDurationKey)

		subQueryNames = append(subQueryNames, fmt.Sprintf(`%sData`, nestedDurationKey))
		partialQueryStr = append(partialQueryStr, []string{
			fmt.Sprintf(`%sData = from(bucket: "%s")`, nestedDurationKey, bucketName),
			fmt.Sprintf(`|> range(start: %s, stop: %s)`, durationRange[0], durationRange[1]),
			`|> filter(fn: (r) => r["_measurement"] == "temp" )`,
			`|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)`,
			`|> group(columns: ["device_wwn"])`,
			`|> toInt()`,
			"",
		}...)
	}

	if len(subQueryNames) == 1 {
		//there's only one bucket being queried, no need to union, just aggregate the dataset and return
		partialQueryStr = append(partialQueryStr, []string{
			subQueryNames[0],
			"|> schema.fieldsAsCols()",
			"|> yield()",
		}...)
	} else {
		partialQueryStr = append(partialQueryStr, []string{
			fmt.Sprintf("union(tables: [%s])", strings.Join(subQueryNames, ", ")),
			`|> group(columns: ["device_wwn"])`,
			`|> sort(columns: ["_time"], desc: false)`,
			"|> schema.fieldsAsCols()",
		}...)
	}

	return strings.Join(partialQueryStr, "\n")
}
