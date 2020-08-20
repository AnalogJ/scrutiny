package models

type Device struct {
	WWN string `json:"wwn"`

	DeviceName     string `json:"device_name"`
	Manufacturer   string `json:"manufacturer"`
	ModelName      string `json:"model_name"`
	InterfaceType  string `json:"interface_type"`
	InterfaceSpeed string `json:"interface_speed"`
	SerialNumber   string `json:"serial_number"`
	Firmware       string `json:"firmware"`
	RotationSpeed  int    `json:"rotational_speed"`
	Capacity       int64  `json:"capacity"`
	FormFactor     string `json:"form_factor"`
	SmartSupport   bool   `json:"smart_support"`
}

type DeviceWrapper struct {
	Success bool     `json:"success,omitempty"`
	Errors  []error  `json:"errors,omitempty"`
	Data    []Device `json:"data"`
}
