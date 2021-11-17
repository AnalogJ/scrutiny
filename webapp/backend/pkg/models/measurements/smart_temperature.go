package measurements

import (
	"time"
)

type SmartTemperature struct {
	Date time.Time `json:"date"`
	Temp int64     `json:"temp"`
}

func (st *SmartTemperature) Flatten() (tags map[string]string, fields map[string]interface{}) {
	fields = map[string]interface{}{
		"temp": st.Temp,
	}
	tags = map[string]string{}

	return tags, fields
}

func (st *SmartTemperature) Inflate(key string, val interface{}) {
	if val == nil {
		return
	}

	if key == "temp" {
		switch t := val.(type) {
		case int64:
			st.Temp = t
		case float64:
			st.Temp = int64(t)
		}
	}
}
