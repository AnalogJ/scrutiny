package database

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"sort"
)

func sortSmartMeasurementsDesc(smartResults []measurements.Smart) {
	sort.SliceStable(smartResults, func(i, j int) bool {
		return smartResults[i].Date.After(smartResults[j].Date)
	})
}
