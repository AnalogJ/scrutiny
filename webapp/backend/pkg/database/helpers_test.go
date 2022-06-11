package database

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_sortSmartMeasurementsDesc_LatestFirst(t *testing.T) {
	//setup
	timeNow := time.Now()
	smartResults := []measurements.Smart{
		{
			Date: timeNow.AddDate(0, 0, -2),
		},
		{
			Date: timeNow,
		},
		{
			Date: timeNow.AddDate(0, 0, -1),
		},
	}

	//test
	sortSmartMeasurementsDesc(smartResults)

	//assert
	require.Equal(t, smartResults[0].Date, timeNow)
}
