package measurements

// TransformDataUnits converts data units to bytes and determines the appropriate unit (MB, GB, TB, PB)
// Returns the transformed value and the unit as a string
func TransformDataUnits(value int64) (int64, string) {
	// Convert to bytes: value * 1000 * 512 (1000 units of 512 byte data units)
	// According to NVMe spec: "This value is reported in thousands (i.e., a value of 1 corresponds to 1000 units of 512 bytes read)"
	bytes := float64(value) * 1000 * 512

	if bytes < 1000*1000*1000 {
		// Less than 1 GB, show in MB
		return int64(bytes / 1000 / 1000), "MB"
	} else if bytes < 1000*1000*1000*1000 {
		// Less than 1 TB, show in GB
		return int64(bytes / 1000 / 1000 / 1000), "GB"
	} else if bytes < 1000*1000*1000*1000*1000 {
		// Less than 1 PB, show in TB
		return int64(bytes / 1000 / 1000 / 1000 / 1000), "TB"
	} else {
		// Show in PB
		return int64(bytes / 1000 / 1000 / 1000 / 1000 / 1000), "PB"
	}
} 