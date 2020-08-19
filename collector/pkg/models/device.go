package models

type Device struct {
	WWN string `json:"wwn" gorm:"primary_key"`

	DeviceName     string `json:"device_name"`
	Manufacturer   string `json:"manufacturer"`
	ModelName      string `json:"model_name"`
	InterfaceType  string `json:"interface_type"`
	InterfaceSpeed string `json:"interface_speed"`
	SerialNumber   string `json:"serial_name"`
	Capacity       int64  `json:"capacity"`
	Firmware       string `json:"firmware"`
	RotationSpeed  int    `json:"rotational_speed"`
}

type DeviceRespWrapper struct {
	Success bool     `json:"success"`
	Errors  []error  `json:"errors"`
	Data    []Device `json:"data"`
}
