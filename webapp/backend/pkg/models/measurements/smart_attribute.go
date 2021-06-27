package measurements

type SmartAttribute interface {
	Flatten() (fields map[string]interface{})
	Inflate(key string, val interface{})
}
