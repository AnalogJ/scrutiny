package thresholds

// https://media.kingston.com/support/downloads/MKP_521.6_SMART-DCP1000_attribute.pdf
// https://www.percona.com/blog/2017/02/09/using-nvme-command-line-tools-to-check-nvme-flash-health/
// https://nvmexpress.org/resources/nvm-express-technology-features/nvme-features-for-error-reporting-smart-log-pages-failures-and-management-capabilities-in-nvme-architectures/
// https://www.micromat.com/product_manuals/drive_scope_manual_01.pdf
type NvmeAttributeMetadata struct {
	ID          string `json:"-"`
	DisplayName string `json:"-"`
	Ideal       string `json:"ideal"`
	Critical    bool   `json:"critical"`
	Description string `json:"description"`

	Transform          func(int64, int64, string) int64 `json:"-"` //this should be a method to extract/tranform the normalized or raw data to a chartable format. Str
	TransformValueUnit string                           `json:"transform_value_unit,omitempty"`
	DisplayType        string                           `json:"display_type"` //"raw" "normalized" or "transformed"
}

var NmveMetadata = map[string]NvmeAttributeMetadata{
	"critical_warning": {
		ID:          "critical_warning",
		DisplayName: "Critical Warning",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "This field indicates critical warnings for the state of the controller. Each bit corresponds to a critical warning type; multiple bits may be set. If a bit is cleared to ‘0’, then that critical warning does not apply. Critical warnings may result in an asynchronous event notification to the host. Bits in this field represent the current associated state and are not persistent.",
	},
	"temperature": {
		ID:          "temperature",
		DisplayName: "Temperature",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"available_spare": {
		ID:          "available_spare",
		DisplayName: "Available Spare",
		DisplayType: "",
		Ideal:       "high",
		Critical:    true,
		Description: "Contains a normalized percentage (0 to 100%) of the remaining spare capacity available.",
	},
	"percentage_used": {
		ID:          "percentage_used",
		DisplayName: "Percentage Used",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "Contains a vendor specific estimate of the percentage of NVM subsystem life used based on the actual usage and the manufacturer’s prediction of NVM life. A value of 100 indicates that the estimated endurance of the NVM in the NVM subsystem has been consumed, but may not indicate an NVM subsystem failure. The value is allowed to exceed 100. Percentages greater than 254 shall be represented as 255. This value shall be updated once per power-on hour (when the controller is not in a sleep state).",
	},
	"data_units_read": {
		ID:          "data_units_read",
		DisplayName: "Data Units Read",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the number of 512 byte data units the host has read from the controller; this value does not include metadata. This value is reported in thousands (i.e., a value of 1 corresponds to 1000 units of 512 bytes read) and is rounded up. When the LBA size is a value other than 512 bytes, the controller shall convert the amount of data read to 512 byte units.",
	},
	"data_units_written": {
		ID:          "data_units_written",
		DisplayName: "Data Units Written",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the number of 512 byte data units the host has written to the controller; this value does not include metadata. This value is reported in thousands (i.e., a value of 1 corresponds to 1000 units of 512 bytes written) and is rounded up. When the LBA size is a value other than 512 bytes, the controller shall convert the amount of data written to 512 byte units.",
	},
	"host_reads": {
		ID:          "host_reads",
		DisplayName: "Host Reads",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the number of read commands completed by the controller",
	},
	"host_writes": {
		ID:          "host_writes",
		DisplayName: "Host Writes",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the number of write commands completed by the controller",
	},
	"controller_busy_time": {
		ID:          "controller_busy_time",
		DisplayName: "Controller Busy Time",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the amount of time the controller is busy with I/O commands. The controller is busy when there is a command outstanding to an I/O Queue (specifically, a command was issued via an I/O Submission Queue Tail doorbell write and the corresponding completion queue entry has not been posted yet to the associated I/O Completion Queue). This value is reported in minutes.",
	},
	"power_cycles": {
		ID:          "power_cycles",
		DisplayName: "Power Cycles",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the number of power cycles.",
	},
	"power_on_hours": {
		ID:          "power_on_hours",
		DisplayName: "Power on Hours",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the number of power-on hours. Power on hours is always logging, even when in low power mode.",
	},
	"unsafe_shutdowns": {
		ID:          "unsafe_shutdowns",
		DisplayName: "Unsafe Shutdowns",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the number of unsafe shutdowns. This count is incremented when a shutdown notification (CC.SHN) is not received prior to loss of power.",
	},
	"media_errors": {
		ID:          "media_errors",
		DisplayName: "Media Errors",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "Contains the number of occurrences where the controller detected an unrecovered data integrity error. Errors such as uncorrectable ECC, CRC checksum failure, or LBA tag mismatch are included in this field.",
	},
	"num_err_log_entries": {
		ID:          "num_err_log_entries",
		DisplayName: "Numb Err Log Entries",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "Contains the number of Error Information log entries over the life of the controller.",
	},
	"warning_temp_time": {
		ID:          "warning_temp_time",
		DisplayName: "Warning Temp Time",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the amount of time in minutes that the controller is operational and the Composite Temperature is greater than or equal to the Warning Composite Temperature Threshold (WCTEMP) field and less than the Critical Composite Temperature Threshold (CCTEMP) field in the Identify Controller data structure.",
	},
	"critical_comp_time": {
		ID:          "critical_comp_time",
		DisplayName: "Critical CompTime",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "Contains the amount of time in minutes that the controller is operational and the Composite Temperature is greater the Critical Composite Temperature Threshold (CCTEMP) field in the Identify Controller data structure.",
	},
}
