package measurements

import "github.com/analogj/scrutiny/webapp/backend/pkg"

type SmartAttribute interface {
	Flatten() (fields map[string]interface{})
	Inflate(key string, val interface{})
	GetStatus() pkg.AttributeStatus
	GetTransformedValue() int64
}
