package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// SMART
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (sr *scrutinyRepository) SaveSmartAttributes(ctx context.Context, wwn string, collectorSmartData collector.SmartInfo) (measurements.Smart, error) {
	sr.logger.Infof("SaveSmartAttributes called for wwn=%s", wwn)
	deviceSmartData := measurements.Smart{}
	err := deviceSmartData.FromCollectorSmartInfo(wwn, collectorSmartData)
	if err != nil {
		sr.logger.Errorln("Could not process SMART metrics", err)
		return measurements.Smart{}, err
	}

	// apply host-level attribute overrides before persisting or notifying
	attributeOverrides := sr.loadAttributeOverrides()
	sr.applyAttributeOverrides(&deviceSmartData, wwn, attributeOverrides)

	tags, fields := deviceSmartData.Flatten()

	// write point immediately
	return deviceSmartData, sr.saveDatapoint(sr.influxWriteApi, "smart", tags, fields, deviceSmartData.Date, ctx)
}

// GetSmartAttributeHistory MUST return in sorted order, where newest entries are at the beginning of the list, and oldest are at the end.
// When selectEntries is > 0, only the most recent selectEntries database entries are returned, starting from the selectEntriesOffset entry.
// For example, with selectEntries = 5, selectEntries = 0, the most recent 5 are returned. With selectEntries = 3, selectEntries = 2, entries
// 2 to 4 are returned (2 being the third newest, since it is zero-indexed)
func (sr *scrutinyRepository) GetSmartAttributeHistory(ctx context.Context, wwn string, durationKey string, selectEntries int, selectEntriesOffset int, attributes []string) ([]measurements.Smart, error) {
	// Get SMartResults from InfluxDB

	//TODO: change the filter startrange to a real number.

	// Get parser flux query result
	//appConfig.GetString("web.influxdb.bucket")
	queryStr := sr.aggregateSmartAttributesQuery(wwn, durationKey, selectEntries, selectEntriesOffset, attributes)
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

func (sr *scrutinyRepository) aggregateSmartAttributesQuery(wwn string, durationKey string, selectEntries int, selectEntriesOffset int, attributes []string) string {

	/*

		import "influxdata/influxdb/schema"
		weekData = from(bucket: "metrics")
		|> range(start: -1w, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "smart" )
		|> filter(fn: (r) => r["device_wwn"] == "0x5000c5002df89099" )
		|> tail(n: 10, offset: 0)
		|> schema.fieldsAsCols()

		monthData = from(bucket: "metrics_weekly")
		|> range(start: -1mo, stop: -1w)
		|> filter(fn: (r) => r["_measurement"] == "smart" )
		|> filter(fn: (r) => r["device_wwn"] == "0x5000c5002df89099" )
		|> tail(n: 10, offset: 0)
		|> schema.fieldsAsCols()

		yearData = from(bucket: "metrics_monthly")
		|> range(start: -1y, stop: -1mo)
		|> filter(fn: (r) => r["_measurement"] == "smart" )
		|> filter(fn: (r) => r["device_wwn"] == "0x5000c5002df89099" )
		|> tail(n: 10, offset: 0)
		|> schema.fieldsAsCols()

		foreverData = from(bucket: "metrics_yearly")
		|> range(start: -10y, stop: -1y)
		|> filter(fn: (r) => r["_measurement"] == "smart" )
		|> filter(fn: (r) => r["device_wwn"] == "0x5000c5002df89099" )
		|> tail(n: 10, offset: 0)
		|> schema.fieldsAsCols()

		union(tables: [weekData, monthData, yearData, foreverData])
		|> group()
		|> sort(columns: ["_time"], desc: true)
		|> tail(n: 6, offset: 4)
		|> yield(name: "last")

	*/

	partialQueryStr := []string{
		`import "influxdata/influxdb/schema"`,
	}

	nestedDurationKeys := sr.lookupNestedDurationKeys(durationKey)

	if len(nestedDurationKeys) == 1 {
		//there's only one bucket being queried, no need to union, just aggregate the dataset and return
		partialQueryStr = append(partialQueryStr, []string{
			sr.generateSmartAttributesSubquery(wwn, nestedDurationKeys[0], selectEntries, selectEntriesOffset, attributes),
			fmt.Sprintf(`%sData`, nestedDurationKeys[0]),
			`|> sort(columns: ["_time"], desc: true)`,
			`|> yield()`,
		}...)
		return strings.Join(partialQueryStr, "\n")
	}

	subQueries := []string{}
	subQueryNames := []string{}
	for _, nestedDurationKey := range nestedDurationKeys {
		subQueryNames = append(subQueryNames, fmt.Sprintf(`%sData`, nestedDurationKey))
		if selectEntries > 0 {
			// We only need the last `n + offset` # of entries from each table to guarantee we can
			// get the last `n` # of entries starting from `offset` of the union
			subQueries = append(subQueries, sr.generateSmartAttributesSubquery(wwn, nestedDurationKey, selectEntries+selectEntriesOffset, 0, attributes))
		} else {
			subQueries = append(subQueries, sr.generateSmartAttributesSubquery(wwn, nestedDurationKey, 0, 0, attributes))
		}
	}
	partialQueryStr = append(partialQueryStr, subQueries...)
	partialQueryStr = append(partialQueryStr, []string{
		fmt.Sprintf("union(tables: [%s])", strings.Join(subQueryNames, ", ")),
		`|> group()`,
		`|> sort(columns: ["_time"], desc: true)`,
	}...)
	if selectEntries > 0 {
		partialQueryStr = append(partialQueryStr, fmt.Sprintf(`|> tail(n: %d, offset: %d)`, selectEntries, selectEntriesOffset))
	}
	partialQueryStr = append(partialQueryStr, `|> yield(name: "last")`)

	return strings.Join(partialQueryStr, "\n")
}

func (sr *scrutinyRepository) generateSmartAttributesSubquery(wwn string, durationKey string, selectEntries int, selectEntriesOffset int, attributes []string) string {
	bucketName := sr.lookupBucketName(durationKey)
	durationRange := sr.lookupDuration(durationKey)

	partialQueryStr := []string{
		fmt.Sprintf(`%sData = from(bucket: "%s")`, durationKey, bucketName),
		fmt.Sprintf(`|> range(start: %s, stop: %s)`, durationRange[0], durationRange[1]),
		`|> filter(fn: (r) => r["_measurement"] == "smart" )`,
		fmt.Sprintf(`|> filter(fn: (r) => r["device_wwn"] == "%s" )`, wwn),
	}

	partialQueryStr = append(partialQueryStr, `|> aggregateWindow(every: 1d, fn: last, createEmpty: false)`)

	if selectEntries > 0 {
		partialQueryStr = append(partialQueryStr, fmt.Sprintf(`|> tail(n: %d, offset: %d)`, selectEntries, selectEntriesOffset))
	}
	partialQueryStr = append(partialQueryStr, "|> schema.fieldsAsCols()")

	return strings.Join(partialQueryStr, "\n")
}

// loadAttributeOverrides retrieves the user-provided overrides from configuration.
func (sr *scrutinyRepository) loadAttributeOverrides() []models.AttributeOverride {
	// Load raw maps so we can detect presence of threshold keys even when value is zero.
	rawOverrides := []map[string]interface{}{}
	if err := sr.appConfig.UnmarshalKey("smart.attribute_overrides", &rawOverrides); err != nil {
		sr.logger.Debugf("failed to parse smart.attribute_overrides: %v", err)
		return []models.AttributeOverride{}
	}

	overrides := make([]models.AttributeOverride, 0, len(rawOverrides))
	for _, raw := range rawOverrides {
		var ao models.AttributeOverride
		if err := mapstructure.Decode(raw, &ao); err != nil {
			sr.logger.Debugf("failed to decode attribute override entry: %v", err)
			continue
		}
		if _, ok := raw["warn_above"]; ok {
			ao.WarnAboveSet = true
		}
		if _, ok := raw["fail_above"]; ok {
			ao.FailAboveSet = true
		}
		overrides = append(overrides, ao)
	}

	sr.logger.Infof("Loaded %d attribute overrides from config", len(overrides))
	for i, o := range overrides {
		sr.logger.Debugf("  Override %d: protocol=%s attributeId=%s wwn=%s action=%s warnAbove=%d warnSet=%t failAbove=%d failSet=%t", i, o.Protocol, o.AttributeId, o.WWN, o.Action, o.WarnAbove, o.WarnAboveSet, o.FailAbove, o.FailAboveSet)
	}
	return overrides
}

// applyAttributeOverrides adjusts attribute statuses according to configured overrides
// and recomputes the device status accordingly (while preserving SMART failure bits).
func (sr *scrutinyRepository) applyAttributeOverrides(smart *measurements.Smart, wwn string, overrides []models.AttributeOverride) {
	if len(overrides) == 0 {
		return
	}

	sr.logger.Debugf("Applying attribute overrides to device wwn=%s protocol=%s with %d attributes", wwn, smart.DeviceProtocol, len(smart.Attributes))

	failedSmart := pkg.DeviceStatusHas(smart.Status, pkg.DeviceStatusFailedSmart)
	// reset and rebuild device status; keep SMART failure flag intact.
	smart.Status = pkg.DeviceStatusPassed
	if failedSmart {
		smart.Status = pkg.DeviceStatusFailedSmart
	}

	for attrKey, attrData := range smart.Attributes {
		override := sr.matchingOverride(smart.DeviceProtocol, wwn, attrKey, overrides)
		if override != nil {
			sr.logger.Infof("Applying override to attribute %s: action=%s", attrKey, override.Action)
			attrData = sr.applyOverrideToAttribute(attrData, *override)
			smart.Attributes[attrKey] = attrData
		}

		// rebuild device status from attribute statuses after overrides are applied
		if pkg.AttributeStatusHas(attrData.GetStatus(), pkg.AttributeStatusFailedScrutiny) {
			smart.Status = pkg.DeviceStatusSet(smart.Status, pkg.DeviceStatusFailedScrutiny)
		}
	}
}

func (sr *scrutinyRepository) matchingOverride(protocol string, wwn string, attributeId string, overrides []models.AttributeOverride) *models.AttributeOverride {
	for ndx := range overrides {
		o := overrides[ndx]
		sr.logger.Debugf("Checking override %d: seeking protocol=%s attributeId=%s wwn=%s against override protocol=%s attributeId=%s wwn=%s", ndx, protocol, attributeId, wwn, o.Protocol, o.AttributeId, o.WWN)
		if !strings.EqualFold(o.Protocol, protocol) {
			sr.logger.Debugf("  Protocol mismatch")
			continue
		}
		if o.AttributeId != "" && o.AttributeId != attributeId {
			sr.logger.Debugf("  AttributeId mismatch: '%s' != '%s'", o.AttributeId, attributeId)
			continue
		}
		if o.WWN != "" && !strings.EqualFold(o.WWN, wwn) {
			sr.logger.Debugf("  WWN mismatch: '%s' != '%s'", o.WWN, wwn)
			continue
		}

		sr.logger.Debugf("  MATCH!")
		return &o
	}
	return nil
}

func (sr *scrutinyRepository) applyOverrideToAttribute(attr measurements.SmartAttribute, override models.AttributeOverride) measurements.SmartAttribute {
	action := strings.ToLower(strings.TrimSpace(override.Action))
	// threshold overrides take precedence over generic force_status/pass logic (unless ignored)
	if action == "ignore" {
		return setAttributeStatus(attr, pkg.AttributeStatusPassed, "Ignored by attribute override")
	}

	if override.FailAboveSet || override.WarnAboveSet {
		return applyThresholdOverride(attr, override)
	}

	switch action {
	case "ignore":
		return setAttributeStatus(attr, pkg.AttributeStatusPassed, "Ignored by attribute override")
	case "force_status":
		status := strings.ToLower(strings.TrimSpace(override.Status))
		switch status {
		case "warn", "warning":
			return setAttributeStatus(attr, pkg.AttributeStatusWarningScrutiny, "Status forced to warning by attribute override")
		case "failed", "fail", "error":
			return setAttributeStatus(attr, pkg.AttributeStatusFailedScrutiny, "Status forced to failed by attribute override")
		default:
			return setAttributeStatus(attr, pkg.AttributeStatusPassed, "Status forced to passed by attribute override")
		}
	default:
		return attr
	}
}

// applyThresholdOverride compares the attribute value to warn/fail thresholds.
// If both warn and fail are set, fail takes precedence when exceeded.
func applyThresholdOverride(attr measurements.SmartAttribute, override models.AttributeOverride) measurements.SmartAttribute {
	val := currentAttributeValue(attr)
	status := pkg.AttributeStatusPassed
	reason := "Status forced to passed by attribute override"

	// fail_above takes priority over warn_above when both are exceeded
	if override.FailAboveSet && val > override.FailAbove {
		status = pkg.AttributeStatusFailedScrutiny
		reason = fmt.Sprintf("Value %d exceeded fail threshold %d", val, override.FailAbove)
	} else if override.WarnAboveSet && val > override.WarnAbove {
		status = pkg.AttributeStatusWarningScrutiny
		reason = fmt.Sprintf("Value %d exceeded warn threshold %d", val, override.WarnAbove)
	}

	return setAttributeStatus(attr, status, reason)
}

// currentAttributeValue returns the most relevant numeric value for comparison.
// ATA uses raw counts; NVMe/SCSI use the current Value field.
func currentAttributeValue(attr measurements.SmartAttribute) int64 {
	switch v := attr.(type) {
	case *measurements.SmartAtaAttribute:
		if v.RawValue != 0 {
			return v.RawValue
		}
		if v.TransformedValue != 0 {
			return v.TransformedValue
		}
		return v.Value
	case *measurements.SmartNvmeAttribute:
		if v.TransformedValue != 0 {
			return v.TransformedValue
		}
		return v.Value
	case *measurements.SmartScsiAttribute:
		if v.TransformedValue != 0 {
			return v.TransformedValue
		}
		return v.Value
	default:
		return 0
	}
}

// setAttributeStatus updates the concrete attribute status/reason fields in-place.
func setAttributeStatus(attr measurements.SmartAttribute, status pkg.AttributeStatus, reason string) measurements.SmartAttribute {
	switch v := attr.(type) {
	case *measurements.SmartAtaAttribute:
		v.Status = status
		v.StatusReason = reason
		v.FailureRate = 0
		return v
	case *measurements.SmartNvmeAttribute:
		v.Status = status
		v.StatusReason = reason
		v.FailureRate = 0
		return v
	case *measurements.SmartScsiAttribute:
		v.Status = status
		v.StatusReason = reason
		v.FailureRate = 0
		return v
	default:
		return attr
	}
}
