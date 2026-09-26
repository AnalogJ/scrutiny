package measurements

import "github.com/analogj/scrutiny/webapp/backend/pkg"

type SmartAttribute interface {
	Flatten() (fields map[string]interface{})
	Inflate(key string, val interface{})
	GetStatus() pkg.AttributeStatus
	// GetComparableValue returns the value the attribute's status is derived from, so changes between submissions can be detected
	GetComparableValue() int64
}
