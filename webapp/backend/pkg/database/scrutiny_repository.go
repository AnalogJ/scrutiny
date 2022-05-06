package database

import (
	"context"
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	globalLogger.Infof("Trying to connect to scrutiny sqlite db: %s\n", appConfig.GetString("web.database.location"))
	database, err := gorm.Open(sqlite.Open(appConfig.GetString("web.database.location")), &gorm.Config{
		//TODO: figure out how to log database queries again.
		//Logger: logger
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database! - %v", err)
	}
	globalLogger.Infof("Successfully connected to scrutiny sqlite db: %s\n", appConfig.GetString("web.database.location"))

	//database.SetLogger()

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
		// we should write the config file out here. Ignore failures.
		err = appConfig.WriteConfig()
		if err != nil {
			globalLogger.Infof("ignoring error while writing influxdb info to config: %v", err)
		}
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

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// InfluxDB & SQLite migrations
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//database.AutoMigrate(&models.Device{})
	err = deviceRepo.Migrate(backgroundContext)
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

		// in tests, we may not want to set a retention policy. If "false", we can set data with old timestamps,
		// then manually run the down sampling scripts. This should be true for production environments.
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
  	bucketBaseName = "%s"

	dailyData = from(bucket: bucketBaseName)
	|> range(start: -10y, stop: now())
	|> filter(fn: (r) => r["_measurement"] == "smart" )
	|> filter(fn: (r) => r["_field"] == "temp" or r["_field"] == "power_on_hours" or r["_field"] == "date")
	|> last()
	|> schema.fieldsAsCols()
	|> group(columns: ["device_wwn"])
	
	weeklyData = from(bucket: bucketBaseName + "_weekly")
	|> range(start: -10y, stop: now())
	|> filter(fn: (r) => r["_measurement"] == "smart" )
	|> filter(fn: (r) => r["_field"] == "temp" or r["_field"] == "power_on_hours" or r["_field"] == "date")
	|> last()
	|> schema.fieldsAsCols()
	|> group(columns: ["device_wwn"])
	
	monthlyData = from(bucket: bucketBaseName + "_monthly")
	|> range(start: -10y, stop: now())
	|> filter(fn: (r) => r["_measurement"] == "smart" )
	|> filter(fn: (r) => r["_field"] == "temp" or r["_field"] == "power_on_hours" or r["_field"] == "date")
	|> last()
	|> schema.fieldsAsCols()
	|> group(columns: ["device_wwn"])
	
	yearlyData = from(bucket: bucketBaseName + "_yearly")
	|> range(start: -10y, stop: now())
	|> filter(fn: (r) => r["_measurement"] == "smart" )
	|> filter(fn: (r) => r["_field"] == "temp" or r["_field"] == "power_on_hours" or r["_field"] == "date")
	|> last()
	|> schema.fieldsAsCols()
	|> group(columns: ["device_wwn"])
	
	union(tables: [dailyData, weeklyData, monthlyData, yearlyData])
	|> sort(columns: ["_time"], desc: false)
	|> group(columns: ["device_wwn"])
	|> last(column: "device_wwn")
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
		return []string{DURATION_KEY_WEEK, DURATION_KEY_MONTH, DURATION_KEY_YEAR, DURATION_KEY_FOREVER}
	}
	return []string{DURATION_KEY_WEEK}
}
