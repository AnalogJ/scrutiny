package database

import (
	mock_config "github.com/analogj/scrutiny/webapp/backend/pkg/config/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_DownsampleScript_Weekly(t *testing.T) {
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

	aggregationType := "weekly"

	//test
	influxDbScript := deviceRepo.DownsampleScript(aggregationType)

	//assert
	require.Equal(t, `
  sourceBucket = "metrics"
  rangeStart = -2w
  rangeEnd = -1w
  aggWindow = 1w
  destBucket = "metrics_weekly"
  destOrg = "scrutiny"

  from(bucket: sourceBucket)
  |> range(start: rangeStart, stop: rangeEnd)
  |> filter(fn: (r) => r["_measurement"] == "smart" )
  |> group(columns: ["device_wwn", "_field"])
  |> aggregateWindow(every: aggWindow, fn: last, createEmpty: false)
  |> to(bucket: destBucket, org: destOrg)

  temp_data = from(bucket: sourceBucket)
  |> range(start: rangeStart, stop: rangeEnd)
  |> filter(fn: (r) => r["_measurement"] == "temp")
  |> group(columns: ["device_wwn"])
  |> toInt()

  temp_data
  |> aggregateWindow(fn: mean, every: aggWindow, createEmpty: false)
  |> to(bucket: destBucket, org: destOrg)
		`, influxDbScript)
}

func Test_DownsampleScript_Monthly(t *testing.T) {
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

	aggregationType := "monthly"

	//test
	influxDbScript := deviceRepo.DownsampleScript(aggregationType)

	//assert
	require.Equal(t, `
  sourceBucket = "metrics_weekly"
  rangeStart = -2mo
  rangeEnd = -1mo
  aggWindow = 1mo
  destBucket = "metrics_monthly"
  destOrg = "scrutiny"

  from(bucket: sourceBucket)
  |> range(start: rangeStart, stop: rangeEnd)
  |> filter(fn: (r) => r["_measurement"] == "smart" )
  |> group(columns: ["device_wwn", "_field"])
  |> aggregateWindow(every: aggWindow, fn: last, createEmpty: false)
  |> to(bucket: destBucket, org: destOrg)

  temp_data = from(bucket: sourceBucket)
  |> range(start: rangeStart, stop: rangeEnd)
  |> filter(fn: (r) => r["_measurement"] == "temp")
  |> group(columns: ["device_wwn"])
  |> toInt()

  temp_data
  |> aggregateWindow(fn: mean, every: aggWindow, createEmpty: false)
  |> to(bucket: destBucket, org: destOrg)
		`, influxDbScript)
}

func Test_DownsampleScript_Yearly(t *testing.T) {
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

	aggregationType := "yearly"

	//test
	influxDbScript := deviceRepo.DownsampleScript(aggregationType)

	//assert
	require.Equal(t, `
  sourceBucket = "metrics_monthly"
  rangeStart = -2y
  rangeEnd = -1y
  aggWindow = 1y
  destBucket = "metrics_yearly"
  destOrg = "scrutiny"

  from(bucket: sourceBucket)
  |> range(start: rangeStart, stop: rangeEnd)
  |> filter(fn: (r) => r["_measurement"] == "smart" )
  |> group(columns: ["device_wwn", "_field"])
  |> aggregateWindow(every: aggWindow, fn: last, createEmpty: false)
  |> to(bucket: destBucket, org: destOrg)

  temp_data = from(bucket: sourceBucket)
  |> range(start: rangeStart, stop: rangeEnd)
  |> filter(fn: (r) => r["_measurement"] == "temp")
  |> group(columns: ["device_wwn"])
  |> toInt()

  temp_data
  |> aggregateWindow(fn: mean, every: aggWindow, createEmpty: false)
  |> to(bucket: destBucket, org: destOrg)
		`, influxDbScript)
}
