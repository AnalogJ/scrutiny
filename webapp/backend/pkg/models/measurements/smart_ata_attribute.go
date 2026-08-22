package measurements

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
)

type SmartAtaAttribute struct {
	AttributeId int    `json:"attribute_id"`
	Value       int64  `json:"value"`
	Threshold   int64  `json:"thresh"`
	Worst       int64  `json:"worst"`
	RawValue    int64  `json:"raw_value"`
	RawString   string `json:"raw_string"`
	WhenFailed  string `json:"when_failed"`

	//Generated data
	TransformedValue int64               `json:"transformed_value"`
	Status           pkg.AttributeStatus `json:"status"`
	StatusReason     string              `json:"status_reason,omitempty"`
	FailureRate      float64             `json:"failure_rate,omitempty"`
}

func (sa *SmartAtaAttribute) GetTransformedValue() int64 {
	return sa.TransformedValue
}

func (sa *SmartAtaAttribute) GetStatus() pkg.AttributeStatus {
	return sa.Status
}

func (sa *SmartAtaAttribute) Flatten() map[string]interface{} {

	idString := strconv.Itoa(sa.AttributeId)

	return map[string]interface{}{
		fmt.Sprintf("attr.%s.attribute_id", idString): idString,
		fmt.Sprintf("attr.%s.value", idString):        sa.Value,
		fmt.Sprintf("attr.%s.worst", idString):        sa.Worst,
		fmt.Sprintf("attr.%s.thresh", idString):       sa.Threshold,
		fmt.Sprintf("attr.%s.raw_value", idString):    sa.RawValue,
		fmt.Sprintf("attr.%s.raw_string", idString):   sa.RawString,
		fmt.Sprintf("attr.%s.when_failed", idString):  sa.WhenFailed,

		//Generated Data
		fmt.Sprintf("attr.%s.transformed_value", idString): sa.TransformedValue,
		fmt.Sprintf("attr.%s.status", idString):            int64(sa.Status),
		fmt.Sprintf("attr.%s.status_reason", idString):     sa.StatusReason,
		fmt.Sprintf("attr.%s.failure_rate", idString):      sa.FailureRate,
	}
}
func (sa *SmartAtaAttribute) Inflate(key string, val interface{}) {
	if val == nil {
		return
	}
	keyParts := strings.Split(key, ".")

	switch keyParts[2] {
	case "attribute_id":
		attrId, err := strconv.Atoi(val.(string))
		if err == nil {
			sa.AttributeId = attrId
		}
	case "value":
		sa.Value = val.(int64)
	case "worst":
		sa.Worst = val.(int64)
	case "thresh":
		sa.Threshold = val.(int64)
	case "raw_value":
		sa.RawValue = val.(int64)
	case "raw_string":
		sa.RawString = val.(string)
	case "when_failed":
		sa.WhenFailed = val.(string)

	//generated
	case "transformed_value":
		sa.TransformedValue = val.(int64)
	case "status":
		sa.Status = pkg.AttributeStatus(val.(int64))
	case "status_reason":
		sa.StatusReason = val.(string)
	case "failure_rate":
		sa.FailureRate = val.(float64)

	}
}

// populate attribute status, using SMART Thresholds & Observed Metadata
// Chainable
func (sa *SmartAtaAttribute) PopulateAttributeStatus(override *pkg.AttributeOverride) *SmartAtaAttribute {
	if strings.ToUpper(sa.WhenFailed) == pkg.AttributeWhenFailedFailingNow {
		//this attribute has previously failed
		sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusFailedSmart)
		sa.StatusReason += "Attribute is failing manufacturer SMART threshold"
		//if the Smart Status is failed, we should exit early, no need to look at thresholds.
		return sa

	} else if strings.ToUpper(sa.WhenFailed) == pkg.AttributeWhenFailedInThePast {
		sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusWarningScrutiny)
		sa.StatusReason += "Attribute has previously failed manufacturer SMART threshold"
	}

	smartMetadata, hasMetadata := thresholds.AtaMetadata[sa.AttributeId]

	// A user override replaces the observed-threshold analysis entirely. The SMART-derived
	// flags set above are left alone: an override tunes Scrutiny's heuristics, it does not
	// suppress the drive reporting itself as failing.
	if override != nil {
		value := sa.RawValue
		idealLow := true
		if hasMetadata {
			value = sa.observedValue(smartMetadata)
			idealLow = idealIsLow(smartMetadata.Ideal)
		}
		status, reason := evaluateAttributeOverride(value, idealLow, override)
		sa.Status = pkg.AttributeStatusSet(sa.Status, status)
		sa.StatusReason += reason
		return sa
	}

	if hasMetadata {
		sa.ValidateThreshold(smartMetadata)
	}

	return sa
}

// observedValue returns the value this attribute is analyzed against: normalized, transformed
// or raw, depending on the attribute's DisplayType.
func (sa *SmartAtaAttribute) observedValue(smartMetadata thresholds.AtaAttributeMetadata) int64 {
	switch smartMetadata.DisplayType {
	case thresholds.AtaSmartAttributeDisplayTypeNormalized:
		return sa.Value
	case thresholds.AtaSmartAttributeDisplayTypeTransformed:
		return sa.TransformedValue
	default:
		return sa.RawValue
	}
}

// compare the attribute (raw, normalized, transformed) value to observed thresholds, and update status if necessary
func (sa *SmartAtaAttribute) ValidateThreshold(smartMetadata thresholds.AtaAttributeMetadata) {
	//TODO: multiple rules
	// try to predict the failure rates for observed thresholds that have 0 failure rate and error bars.
	// - if the attribute is critical
	//		- the failure rate is over 10 - set to failed
	//		- the attribute does not match any threshold, set to warn
	// - if the attribute is not critical
	//		- if failure rate is above 20 - set to failed
	// 		- if failure rate is above 10 but below 20 - set to warn

	//update the smart attribute status based on Observed thresholds.
	value := sa.observedValue(smartMetadata)

	for _, obsThresh := range smartMetadata.ObservedThresholds {

		//check if "value" is in this bucket
		if ((obsThresh.Low == obsThresh.High) && value == obsThresh.Low) ||
			(obsThresh.Low < value && value <= obsThresh.High) {
			sa.FailureRate = obsThresh.AnnualFailureRate

			if smartMetadata.Critical {
				if obsThresh.AnnualFailureRate >= 0.10 {
					sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusFailedScrutiny)
					sa.StatusReason += "Observed Failure Rate for Critical Attribute is greater than 10%"
				}
			} else {
				if obsThresh.AnnualFailureRate >= 0.20 {
					sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusFailedScrutiny)
					sa.StatusReason += "Observed Failure Rate for Non-Critical Attribute is greater than 20%"
				} else if obsThresh.AnnualFailureRate >= 0.10 {
					sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusWarningScrutiny)
					sa.StatusReason += "Observed Failure Rate for Non-Critical Attribute is greater than 10%"
				}
			}

			//we've found the correct bucket, we can drop out of this loop
			return
		}
	}
	// no bucket found
	if smartMetadata.Critical {
		sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusWarningScrutiny)
		sa.StatusReason = "Could not determine Observed Failure Rate for Critical Attribute"
	}
}
