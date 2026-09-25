package collector

type CollectorError struct {
	Error        string `json:"error"`
	HostId       string `json:"host_id,omitempty"`
	DeviceName   string `json:"device_name,omitempty"`
	DeviceType   string `json:"device_type,omitempty"`
	ScrutinyUUID string `json:"scrutiny_uuid,omitempty"`
}
