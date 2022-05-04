package database

import (
	"context"
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

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

	// write point immediately
	return deviceSmartData, sr.saveDatapoint(sr.influxWriteApi, "smart", tags, fields, deviceSmartData.Date, ctx)
}

func (sr *scrutinyRepository) GetSmartAttributeHistory(ctx context.Context, wwn string, durationKey string, attributes []string) ([]measurements.Smart, error) {
	// Get SMartResults from InfluxDB

	fmt.Println("GetDeviceDetails from INFLUXDB")

	//TODO: change the filter startrange to a real number.

	// Get parser flux query result
	//appConfig.GetString("web.influxdb.bucket")
	queryStr := sr.aggregateSmartAttributesQuery(wwn, durationKey)
	log.Infoln(queryStr)

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
// Helper Methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (sr *scrutinyRepository) saveDatapoint(influxWriteApi api.WriteAPIBlocking, measurement string, tags map[string]string, fields map[string]interface{}, date time.Time, ctx context.Context) error {
	sr.logger.Debugf("Storing datapoint in measurement '%s'. tags: %d fields: %d", measurement, len(tags), len(fields))
	p := influxdb2.NewPoint(measurement,
		tags,
		fields,
		date)

	// write point immediately
	return influxWriteApi.WritePoint(ctx, p)
}

func (sr *scrutinyRepository) aggregateSmartAttributesQuery(wwn string, durationKey string) string {

	/*

		import "influxdata/influxdb/schema"
		weekData = from(bucket: "metrics")
		  |> range(start: -1w, stop: now())
		  |> filter(fn: (r) => r["_measurement"] == "smart" )
		  |> filter(fn: (r) => r["device_wwn"] == "%s" )
		  |> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
		  |> group(columns: ["device_wwn"])
		  |> toInt()

		monthData = from(bucket: "metrics_weekly")
		  |> range(start: -1mo, stop: now())
		  |> filter(fn: (r) => r["_measurement"] == "smart" )
		  |> filter(fn: (r) => r["device_wwn"] == "%s" )
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
			`|> filter(fn: (r) => r["_measurement"] == "smart" )`,
			fmt.Sprintf(`|> filter(fn: (r) => r["device_wwn"] == "%s" )`, wwn),
			"|> schema.fieldsAsCols()",
		}...)
	}

	if len(subQueryNames) == 1 {
		//there's only one bucket being queried, no need to union, just aggregate the dataset and return
		partialQueryStr = append(partialQueryStr, []string{
			subQueryNames[0],
			`|> yield()`,
		}...)
	} else {
		partialQueryStr = append(partialQueryStr, []string{
			fmt.Sprintf("union(tables: [%s])", strings.Join(subQueryNames, ", ")),
			`|> sort(columns: ["_time"], desc: false)`,
			`|> yield(name: "last")`,
		}...)
	}

	return strings.Join(partialQueryStr, "\n")
}
