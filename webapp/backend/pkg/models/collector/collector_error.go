package collector

type CollectorError struct {
	Error  string `json:"error"`
	HostId string `json:"host_id,omitempty"`
}
