package measurements

import (
	"fmt"
	"strings"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
)

type SmartNvmeAttribute struct {
	AttributeId string `json:"attribute_id"` // json string from smartctl
	Value       int64  `json:"value"`
	Threshold   int64  `json:"thresh"`

	TransformedValue int64               `json:"transformed_value"`
	Status           pkg.AttributeStatus `json:"status"`
	StatusReason     string              `json:"status_reason,omitempty"`
	FailureRate      float64             `json:"failure_rate,omitempty"`
}

func (sa *SmartNvmeAttribute) GetTransformedValue() int64 {
	return sa.TransformedValue
}

func (sa *SmartNvmeAttribute) GetStatus() pkg.AttributeStatus {
	return sa.Status
}

func (sa *SmartNvmeAttribute) Flatten() map[string]interface{} {
	return map[string]interface{}{
		fmt.Sprintf("attr.%s.attribute_id", sa.AttributeId): sa.AttributeId,
		fmt.Sprintf("attr.%s.value", sa.AttributeId):        sa.Value,
		fmt.Sprintf("attr.%s.thresh", sa.AttributeId):       sa.Threshold,

		// Generated Data
		fmt.Sprintf("attr.%s.transformed_value", sa.AttributeId): sa.TransformedValue,
		fmt.Sprintf("attr.%s.status", sa.AttributeId):            int64(sa.Status),
		fmt.Sprintf("attr.%s.status_reason", sa.AttributeId):     sa.StatusReason,
		fmt.Sprintf("attr.%s.failure_rate", sa.AttributeId):      sa.FailureRate,
	}
}

func (sa *SmartNvmeAttribute) Inflate(key string, val interface{}) {
	if val == nil {
		return
	}

	keyParts := strings.Split(key, ".")

	switch keyParts[2] {
	case "attribute_id":
		sa.AttributeId = val.(string)
	case "value":
		sa.Value = val.(int64)
	case "thresh":
		sa.Threshold = val.(int64)

	// generated
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
func (sa *SmartNvmeAttribute) PopulateAttributeStatus() *SmartNvmeAttribute {
	//-1 is a special number meaning no threshold.
	if sa.Threshold != -1 {
		if smartMetadata, ok := thresholds.NmveMetadata[sa.AttributeId]; ok {
			// check what the ideal is. Ideal tells us if we our recorded value needs to be above, or below the threshold
			if (smartMetadata.Ideal == "low" && sa.Value > sa.Threshold) ||
				(smartMetadata.Ideal == "high" && sa.Value < sa.Threshold) {
				sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusFailedScrutiny)
				sa.StatusReason += "Attribute is failing recommended SMART threshold"
			}
		}
	}
	// TODO: eventually figure out the critical_warning bits and determine correct error messages here.

	return sa
}
