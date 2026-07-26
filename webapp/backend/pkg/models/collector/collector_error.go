package collector

// CollectorError is reported by the collector when it could not gather data. Everything except
// Error is best-effort: a scan failure has no device, and a device whose `smartctl --info` failed
// never got a ScrutinyUUID.
type CollectorError struct {
	Error        string `json:"error"`
	HostId       string `json:"host_id,omitempty"`
	DeviceName   string `json:"device_name,omitempty"`
	DeviceType   string `json:"device_type,omitempty"`
	ScrutinyUUID string `json:"scrutiny_uuid,omitempty"`
}
