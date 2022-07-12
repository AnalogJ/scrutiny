package database

import (
	mock_config "github.com/analogj/scrutiny/webapp/backend/pkg/config/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_aggregateTempQuery_Week(t *testing.T) {
	t.Parallel()

	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
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
|> group(columns: ["device_wwn"])
|> toInt()

weekData
|> schema.fieldsAsCols()
|> yield()`, influxDbScript)
}

func Test_aggregateTempQuery_Month(t *testing.T) {
	t.Parallel()

	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
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
|> group(columns: ["device_wwn"])
|> toInt()

monthData = from(bucket: "metrics_weekly")
|> range(start: -1mo, stop: -1w)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["device_wwn"])
|> toInt()

union(tables: [weekData, monthData])
|> group(columns: ["device_wwn"])
|> sort(columns: ["_time"], desc: false)
|> schema.fieldsAsCols()`, influxDbScript)
}

func Test_aggregateTempQuery_Year(t *testing.T) {
	t.Parallel()

	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
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
|> group(columns: ["device_wwn"])
|> toInt()

monthData = from(bucket: "metrics_weekly")
|> range(start: -1mo, stop: -1w)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["device_wwn"])
|> toInt()

yearData = from(bucket: "metrics_monthly")
|> range(start: -1y, stop: -1mo)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["device_wwn"])
|> toInt()

union(tables: [weekData, monthData, yearData])
|> group(columns: ["device_wwn"])
|> sort(columns: ["_time"], desc: false)
|> schema.fieldsAsCols()`, influxDbScript)
}

func Test_aggregateTempQuery_Forever(t *testing.T) {
	t.Parallel()

	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
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
|> group(columns: ["device_wwn"])
|> toInt()

monthData = from(bucket: "metrics_weekly")
|> range(start: -1mo, stop: -1w)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["device_wwn"])
|> toInt()

yearData = from(bucket: "metrics_monthly")
|> range(start: -1y, stop: -1mo)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["device_wwn"])
|> toInt()

foreverData = from(bucket: "metrics_yearly")
|> range(start: -10y, stop: -1y)
|> filter(fn: (r) => r["_measurement"] == "temp" )
|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
|> group(columns: ["device_wwn"])
|> toInt()

union(tables: [weekData, monthData, yearData, foreverData])
|> group(columns: ["device_wwn"])
|> sort(columns: ["_time"], desc: false)
|> schema.fieldsAsCols()`, influxDbScript)
}
