package measurements

import (
	"fmt"
	"strconv"
	"strings"
)

const SmartAttributeStatusPassed = "passed"
const SmartAttributeStatusFailed = "failed"
const SmartAttributeStatusWarning = "warn"

type SmartAtaAttribute struct {
	AttributeId int    `json:"attribute_id"`
	Name        string `json:"name"`
	Value       int64  `json:"value"`
	Threshold   int64  `json:"thresh"`
	Worst       int64  `json:"worst"`
	RawValue    int64  `json:"raw_value"`
	RawString   string `json:"raw_string"`
	WhenFailed  string `json:"when_failed"`

	//Generated data
	TransformedValue int64   `json:"transformed_value"`
	Status           string  `json:"status,omitempty"`
	StatusReason     string  `json:"status_reason,omitempty"`
	FailureRate      float64 `json:"failure_rate,omitempty"`
}

func (sa *SmartAtaAttribute) Flatten() map[string]interface{} {

	idString := strconv.Itoa(sa.AttributeId)

	return map[string]interface{}{
		fmt.Sprintf("attr.%s.attribute_id", idString): idString,
		fmt.Sprintf("attr.%s.name", idString):         sa.Name,
		fmt.Sprintf("attr.%s.value", idString):        sa.Value,
		fmt.Sprintf("attr.%s.worst", idString):        sa.Worst,
		fmt.Sprintf("attr.%s.thresh", idString):       sa.Threshold,
		fmt.Sprintf("attr.%s.raw_value", idString):    sa.RawValue,
		fmt.Sprintf("attr.%s.raw_string", idString):   sa.RawString,
		fmt.Sprintf("attr.%s.when_failed", idString):  sa.WhenFailed,
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
	case "name":
		sa.Name = val.(string)
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
	}
}

//
////populate attribute status, using SMART Thresholds & Observed Metadata
//func (sa *SmartAtaAttribute) PopulateAttributeStatus() {
//	if strings.ToUpper(sa.WhenFailed) == SmartWhenFailedFailingNow {
//		//this attribute has previously failed
//		sa.Status = SmartAttributeStatusFailed
//		sa.StatusReason = "Attribute is failing manufacturer SMART threshold"
//
//	} else if strings.ToUpper(sa.WhenFailed) == SmartWhenFailedInThePast {
//		sa.Status = SmartAttributeStatusWarning
//		sa.StatusReason = "Attribute has previously failed manufacturer SMART threshold"
//	}
//
//	if smartMetadata, ok := metadata.AtaMetadata[sa.AttributeId]; ok {
//		sa.MetadataObservedThresholdStatus(smartMetadata)
//	}
//
//	//check if status is blank, set to "passed"
//	if len(sa.Status) == 0 {
//		sa.Status = SmartAttributeStatusPassed
//	}
//}
//
//// compare the attribute (raw, normalized, transformed) value to observed thresholds, and update status if necessary
//func (sa *SmartAtaAttribute) MetadataObservedThresholdStatus(smartMetadata metadata.AtaAttributeMetadata) {
//	//TODO: multiple rules
//	// try to predict the failure rates for observed thresholds that have 0 failure rate and error bars.
//	// - if the attribute is critical
//	//		- the failure rate is over 10 - set to failed
//	//		- the attribute does not match any threshold, set to warn
//	// - if the attribute is not critical
//	//		- if failure rate is above 20 - set to failed
//	// 		- if failure rate is above 10 but below 20 - set to warn
//
//	//update the smart attribute status based on Observed thresholds.
//	var value int64
//	if smartMetadata.DisplayType == metadata.AtaSmartAttributeDisplayTypeNormalized {
//		value = int64(sa.Value)
//	} else if smartMetadata.DisplayType == metadata.AtaSmartAttributeDisplayTypeTransformed {
//		value = sa.TransformedValue
//	} else {
//		value = sa.RawValue
//	}
//
//	for _, obsThresh := range smartMetadata.ObservedThresholds {
//
//		//check if "value" is in this bucket
//		if ((obsThresh.Low == obsThresh.High) && value == obsThresh.Low) ||
//			(obsThresh.Low < value && value <= obsThresh.High) {
//			sa.FailureRate = obsThresh.AnnualFailureRate
//
//			if smartMetadata.Critical {
//				if obsThresh.AnnualFailureRate >= 0.10 {
//					sa.Status = SmartAttributeStatusFailed
//					sa.StatusReason = "Observed Failure Rate for Critical Attribute is greater than 10%"
//				}
//			} else {
//				if obsThresh.AnnualFailureRate >= 0.20 {
//					sa.Status = SmartAttributeStatusFailed
//					sa.StatusReason = "Observed Failure Rate for Attribute is greater than 20%"
//				} else if obsThresh.AnnualFailureRate >= 0.10 {
//					sa.Status = SmartAttributeStatusWarning
//					sa.StatusReason = "Observed Failure Rate for Attribute is greater than 10%"
//				}
//			}
//
//			//we've found the correct bucket, we can drop out of this loop
//			return
//		}
//	}
//	// no bucket found
//	if smartMetadata.Critical {
//		sa.Status = SmartAttributeStatusWarning
//		sa.StatusReason = "Could not determine Observed Failure Rate for Critical Attribute"
//	}
//
//	return
//}
