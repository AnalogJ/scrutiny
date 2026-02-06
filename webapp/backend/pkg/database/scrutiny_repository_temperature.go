package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Temperature Data
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (sr *scrutinyRepository) SaveSmartTemperature(ctx context.Context, wwn string, deviceProtocol string, collectorSmartData collector.SmartInfo, discardSCTTempHistory bool) error {
	if len(collectorSmartData.AtaSctTemperatureHistory.Table) > 0 && !discardSCTTempHistory {

		for ndx, temp := range collectorSmartData.AtaSctTemperatureHistory.Table {
			//temp value may be null, we must skip/ignore them. See #393
			if temp == 0 {
				continue
			}

			intervalSec := collectorSmartData.AtaSctTemperatureHistory.LoggingIntervalMinutes * 60
			datapointTime := collectorSmartData.LocalTime.TimeT - int64(ndx) * intervalSec
			alignedDatapointTime := datapointTime - datapointTime % intervalSec
			smartTemp := measurements.SmartTemperature{
				Date: time.Unix(alignedDatapointTime, 0),
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
	}


	// Even if ata_sct_temperature_history is present, also add current temperature. See #824
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
// Helper Methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
