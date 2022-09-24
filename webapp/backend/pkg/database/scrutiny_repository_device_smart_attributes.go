package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// SMART
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (sr *scrutinyRepository) SaveSmartAttributes(ctx context.Context, wwn string, collectorSmartData collector.SmartInfo) (measurements.Smart, error) {
	deviceSmartData := measurements.Smart{}
	err := deviceSmartData.FromCollectorSmartInfo(sr.appConfig, wwn, collectorSmartData)
	if err != nil {
		sr.logger.Errorln("Could not process SMART metrics", err)
		return measurements.Smart{}, err
	}

	tags, fields := deviceSmartData.Flatten()

	// write point immediately
	return deviceSmartData, sr.saveDatapoint(sr.influxWriteApi, "smart", tags, fields, deviceSmartData.Date, ctx)
}

// GetSmartAttributeHistory MUST return in sorted order, where newest entries are at the beginning of the list, and oldest are at the end.
func (sr *scrutinyRepository) GetSmartAttributeHistory(ctx context.Context, wwn string, durationKey string, attributes []string) ([]measurements.Smart, error) {
	// Get SMartResults from InfluxDB

	//TODO: change the filter startrange to a real number.

	// Get parser flux query result
	//appConfig.GetString("web.influxdb.bucket")
	queryStr := sr.aggregateSmartAttributesQuery(wwn, durationKey)
	log.Infoln(queryStr)

	smartResults := []measurements.Smart{}

	result, err := sr.influxQueryApi.Query(ctx, queryStr)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				//fmt.Printf("table: %s\n", result.TableMetadata().String())
			}

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

	//we have to sort the smartResults again, because the `union` command will return multiple 'tables' and only sort the records in each table.
	sortSmartMeasurementsDesc(smartResults)

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
	//sr.logger.Debugf("Storing datapoint in measurement '%s'. tags: %d fields: %d", measurement, len(tags), len(fields))
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
		|> filter(fn: (r) => r["device_wwn"] == "0x5000c5002df89099" )
		|> schema.fieldsAsCols()

		monthData = from(bucket: "metrics_weekly")
		|> range(start: -1mo, stop: -1w)
		|> filter(fn: (r) => r["_measurement"] == "smart" )
		|> filter(fn: (r) => r["device_wwn"] == "0x5000c5002df89099" )
		|> schema.fieldsAsCols()

		yearData = from(bucket: "metrics_monthly")
		|> range(start: -1y, stop: -1mo)
		|> filter(fn: (r) => r["_measurement"] == "smart" )
		|> filter(fn: (r) => r["device_wwn"] == "0x5000c5002df89099" )
		|> schema.fieldsAsCols()

		foreverData = from(bucket: "metrics_yearly")
		|> range(start: -10y, stop: -1y)
		|> filter(fn: (r) => r["_measurement"] == "smart" )
		|> filter(fn: (r) => r["device_wwn"] == "0x5000c5002df89099" )
		|> schema.fieldsAsCols()

		union(tables: [weekData, monthData, yearData, foreverData])
		|> sort(columns: ["_time"], desc: false)
		|> yield(name: "last")

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
