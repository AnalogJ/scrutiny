package database

import (
	"context"
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
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

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Gorm/SQLite setup
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	fmt.Printf("Trying to connect to database stored: %s\n", appConfig.GetString("web.database.location"))
	database, err := gorm.Open(sqlite.Open(appConfig.GetString("web.database.location")), &gorm.Config{
		//TODO: figure out how to log database queries again.
		//Logger: logger
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database!")
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
		onboardingResponse, err := client.Setup(
			context.Background(),
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

	if writeAPI == nil || queryAPI == nil {
		return nil, fmt.Errorf("Failed to connect to influxdb!")
	}

	deviceRepo := scrutinyRepository{
		appConfig:      appConfig,
		logger:         globalLogger,
		influxClient:   client,
		influxWriteApi: writeAPI,
		influxQueryApi: queryAPI,
		gormClient:     database,
	}

	return &deviceRepo, nil
}

type scrutinyRepository struct {
	appConfig config.Interface
	logger    logrus.FieldLogger

	influxWriteApi api.WriteAPIBlocking
	influxQueryApi api.QueryAPI
	influxClient   influxdb2.Client

	gormClient *gorm.DB
}

func (sr *scrutinyRepository) Close() error {
	sr.influxClient.Close()
	return nil
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
		return nil, fmt.Errorf("Could not get device summary from DB", err)
	}
	return devices, nil
}

// update device (only metadata) from collector
func (sr *scrutinyRepository) UpdateDevice(ctx context.Context, wwn string, collectorSmartData collector.SmartInfo) (models.Device, error) {
	var device models.Device
	if err := sr.gormClient.WithContext(ctx).Where("wwn = ?", wwn).First(&device).Error; err != nil {
		return device, fmt.Errorf("Could not get device from DB", err)
	}

	//TODO catch GormClient err
	err := device.UpdateFromCollectorSmartInfo(collectorSmartData)
	if err != nil {
		return device, err
	}
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

func (sr *scrutinyRepository) GetSmartTemperatureHistory(ctx context.Context) (map[string][]measurements.SmartTemperature, error) {

	deviceTempHistory := map[string][]measurements.SmartTemperature{}

	//TODO: change the query range to a variable.
	queryStr := fmt.Sprintf(`
  import "influxdata/influxdb/schema"
  from(bucket: "%s")
  |> range(start: -3y, stop: now())
  |> filter(fn: (r) => r["_measurement"] == "temp" )
  |> filter(fn: (r) => r["_field"] == "temp")
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

	deviceTempHistory, err := sr.GetSmartTemperatureHistory(ctx)
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
