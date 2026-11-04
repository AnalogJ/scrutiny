package database

import (
	"context"
	"encoding/json"
	"testing"

	mock_config "github.com/analogj/scrutiny/webapp/backend/pkg/config/mock"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/gofrs/uuid/v5"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type countingWriteAPI struct {
	api.WriteAPIBlocking
	points int
}

func (c *countingWriteAPI) WritePoint(_ context.Context, points ...*write.Point) error {
	c.points += len(points)
	return nil
}

func Test_SaveSmartTemperature_CurrentTemperature(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name     string
		input    string
		expected int
	}{
		{"unknown temperature writes no point", `{}`, 0},
		{"reported 0 still writes a point", `{"temperature": {"current": 0}}`, 1},
		{"reported temperature writes a point", `{"temperature": {"current": 28}}`, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var smartInfo collector.SmartInfo
			require.NoError(t, json.Unmarshal([]byte(tt.input), &smartInfo))

			writeAPI := &countingWriteAPI{}
			deviceRepo := scrutinyRepository{influxWriteApi: writeAPI}

			require.NoError(t, deviceRepo.SaveSmartTemperature(context.Background(), uuid.Must(uuid.NewV4()), "SCSI", smartInfo, false))
			require.Equal(t, tt.expected, writeAPI.points)
		})
	}
}

func Test_aggregateTempQuery_Week(t *testing.T) {
	t.Parallel()

	//setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()

	deviceRepo := scrutinyRepository{
		appConfig: fakeConfig,
	}

	aggregationType := DURATION_KEY_WEEK

	//test
	influxDbScript := deviceRepo.aggregateTempQuery(aggregationType)

	//assert
	require.Equal(t, `import "influxdata/influxdb/schema"
weekData = from(bucket: "metrics")
|> range(start: -1w, stop: now())
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

weekData
|> schema.fieldsAsCols()
|> yield()`, influxDbScript)
}

func Test_aggregateTempQuery_UnknownDurationUsesWeek(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()

	deviceRepo := scrutinyRepository{
		appConfig: fakeConfig,
	}

	week := deviceRepo.aggregateTempQuery(DURATION_KEY_WEEK)
	require.Equal(t, week, deviceRepo.aggregateTempQuery("bogus"))
	require.Equal(t, week, deviceRepo.aggregateTempQuery(""))
}

func Test_aggregateTempQuery_Month(t *testing.T) {
	t.Parallel()

	//setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()

	deviceRepo := scrutinyRepository{
		appConfig: fakeConfig,
	}

	aggregationType := DURATION_KEY_MONTH

	//test
	influxDbScript := deviceRepo.aggregateTempQuery(aggregationType)

	//assert
	require.Equal(t, `import "influxdata/influxdb/schema"
weekData = from(bucket: "metrics")
|> range(start: -1w, stop: now())
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

monthData = from(bucket: "metrics_weekly")
|> range(start: -1mo, stop: -1w)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

union(tables: [weekData, monthData])
|> group(columns: ["scrutiny_uuid"])
|> sort(columns: ["_time"], desc: false)
|> schema.fieldsAsCols()`, influxDbScript)
}

func Test_aggregateTempQuery_Year(t *testing.T) {
	t.Parallel()

	//setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()

	deviceRepo := scrutinyRepository{
		appConfig: fakeConfig,
	}

	aggregationType := DURATION_KEY_YEAR

	//test
	influxDbScript := deviceRepo.aggregateTempQuery(aggregationType)

	//assert
	require.Equal(t, `import "influxdata/influxdb/schema"
weekData = from(bucket: "metrics")
|> range(start: -1w, stop: now())
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

monthData = from(bucket: "metrics_weekly")
|> range(start: -1mo, stop: -1w)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

yearData = from(bucket: "metrics_monthly")
|> range(start: -1y, stop: -1mo)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

union(tables: [weekData, monthData, yearData])
|> group(columns: ["scrutiny_uuid"])
|> sort(columns: ["_time"], desc: false)
|> schema.fieldsAsCols()`, influxDbScript)
}

func Test_aggregateTempQuery_Forever(t *testing.T) {
	t.Parallel()

	//setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()

	deviceRepo := scrutinyRepository{
		appConfig: fakeConfig,
	}

	aggregationType := DURATION_KEY_FOREVER

	//test
	influxDbScript := deviceRepo.aggregateTempQuery(aggregationType)

	//assert
	require.Equal(t, `import "influxdata/influxdb/schema"
weekData = from(bucket: "metrics")
|> range(start: -1w, stop: now())
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

monthData = from(bucket: "metrics_weekly")
|> range(start: -1mo, stop: -1w)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

yearData = from(bucket: "metrics_monthly")
|> range(start: -1y, stop: -1mo)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

foreverData = from(bucket: "metrics_yearly")
|> range(start: -10y, stop: -1y)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["scrutiny_uuid"])
|> toInt()

union(tables: [weekData, monthData, yearData, foreverData])
|> group(columns: ["scrutiny_uuid"])
|> sort(columns: ["_time"], desc: false)
|> schema.fieldsAsCols()`, influxDbScript)
}
