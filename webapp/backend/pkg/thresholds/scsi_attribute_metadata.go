package thresholds

type ScsiAttributeMetadata struct {
	ID          string `json:"-"`
	DisplayName string `json:"display_name"`
	Ideal       string `json:"ideal"`
	Critical    bool   `json:"critical"`
	Description string `json:"description"`

	Transform          func(int64, int64, string) int64 `json:"-"` // this should be a method to extract/tranform the normalized or raw data to a chartable format. Str
	TransformValueUnit string                           `json:"transform_value_unit,omitempty"`
	DisplayType        string                           `json:"display_type"` //"raw" "normalized" or "transformed"
}

var ScsiMetadata = map[string]ScsiAttributeMetadata{
	"scsi_grown_defect_list": {
		ID:          "scsi_grown_defect_list",
		DisplayName: "Grown Defect List",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "The grown defect count shows the amount of swapped (defective) blocks since the drive was shipped by it's vendor. Each additional defective block increases the count by one.",
	},
	"read_errors_corrected_by_eccfast": {
		ID:          "read_errors_corrected_by_eccfast",
		DisplayName: "Read Errors Corrected by ECC Fast",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "An error correction was applied to get perfect data (a.k.a. ECC on-the-fly). \"Without substantial delay\" means the correction did not postpone reading of later sectors (e.g. a revolution was not lost). The counter is incremented once for each logical block that requires correction. Two different blocks corrected during the same command are counted as two events.",
	},
	"read_errors_corrected_by_eccdelayed": {
		ID:          "read_errors_corrected_by_eccdelayed",
		DisplayName: "Read Errors Corrected by ECC Delayed",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "An error code or algorithm (e.g. ECC, checksum) is applied in order to get perfect data with substantial delay. \"With possible delay\" means the correction took longer than a sector time so that reading/writing of subsequent sectors was delayed (e.g. a lost revolution). The counter is incremented once for each logical block that requires correction. A block with a double error that is correctable counts as one event and two different blocks corrected during the same command count as two events. ",
	},
	"read_errors_corrected_by_rereads_rewrites": {
		ID:          "read_errors_corrected_by_rereads_rewrites",
		DisplayName: "Read Errors Corrected by ReReads/ReWrites",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "This parameter code specifies the counter counting the number of errors that are corrected by applying retries. This counts errors recovered, not the number of retries. If five retries were required to recover one block of data, the counter increments by one, not five. The counter is incremented once for each logical block that is recovered using retries. If an error is not recoverable while applying retries and is recovered by ECC, it isn't counted by this counter; it will be counted by the counter specified by parameter code 01h - Errors Corrected With Possible Delays. ",
	},
	"read_total_errors_corrected": {
		ID:          "read_total_errors_corrected",
		DisplayName: "Read Total Errors Corrected",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "This counter counts the total of parameter code errors 00h, 01h and 02h (i.e. error corrected by ECC: fast and delayed plus errors corrected by rereads and rewrites). There is no \"double counting\" of data errors among these three counters. The sum of all correctable errors can be reached by adding parameter code 01h and 02h errors, not by using this total.",
	},
	"read_correction_algorithm_invocations": {
		ID:          "read_correction_algorithm_invocations",
		DisplayName: "Read Correction Algorithm Invocations",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "This parameter code specifies the counter that counts the total number of retries, or \"times the retry algorithm is invoked\". If after five attempts a counter 02h type error is recovered, then five is added to this counter. If three retries are required to get stable ECC syndrome before a counter 01h type error is corrected, then those three retries are also counted here. The number of retries applied to unsuccessfully recover an error (counter 06h type error) are also counted by this counter. ",
	},
	"read_total_uncorrected_errors": {
		ID:          "read_total_uncorrected_errors",
		DisplayName: "Read Total Uncorrected Errors",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "This parameter code specifies the counter that contains the total number of blocks for which an uncorrected data error has occurred. ",
	},
	"write_errors_corrected_by_eccfast": {
		ID:          "write_errors_corrected_by_eccfast",
		DisplayName: "Write Errors Corrected by ECC Fast",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "An error correction was applied to get perfect data (a.k.a. ECC on-the-fly). \"Without substantial delay\" means the correction did not postpone reading of later sectors (e.g. a revolution was not lost). The counter is incremented once for each logical block that requires correction. Two different blocks corrected during the same command are counted as two events. ",
	},
	"write_errors_corrected_by_eccdelayed": {
		ID:          "write_errors_corrected_by_eccdelayed",
		DisplayName: "Write Errors Corrected by ECC Delayed",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "An error code or algorithm (e.g. ECC, checksum) is applied in order to get perfect data with substantial delay. \"With possible delay\" means the correction took longer than a sector time so that reading/writing of subsequent sectors was delayed (e.g. a lost revolution). The counter is incremented once for each logical block that requires correction. A block with a double error that is correctable counts as one event and two different blocks corrected during the same command count as two events. ",
	},
	"write_errors_corrected_by_rereads_rewrites": {
		ID:          "write_errors_corrected_by_rereads_rewrites",
		DisplayName: "Write Errors Corrected by ReReads/ReWrites",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "This parameter code specifies the counter counting the number of errors that are corrected by applying retries. This counts errors recovered, not the number of retries. If five retries were required to recover one block of data, the counter increments by one, not five. The counter is incremented once for each logical block that is recovered using retries. If an error is not recoverable while applying retries and is recovered by ECC, it isn't counted by this counter; it will be counted by the counter specified by parameter code 01h - Errors Corrected With Possible Delays.",
	},
	"write_total_errors_corrected": {
		ID:          "write_total_errors_corrected",
		DisplayName: "Write Total Errors Corrected",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "This counter counts the total of parameter code errors 00h, 01h and 02h (i.e. error corrected by ECC: fast and delayed plus errors corrected by rereads and rewrites). There is no \"double counting\" of data errors among these three counters. The sum of all correctable errors can be reached by adding parameter code 01h and 02h errors, not by using this total.",
	},
	"write_correction_algorithm_invocations": {
		ID:          "write_correction_algorithm_invocations",
		DisplayName: "Write Correction Algorithm Invocations",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "This parameter code specifies the counter that counts the total number of retries, or \"times the retry algorithm is invoked\". If after five attempts a counter 02h type error is recovered, then five is added to this counter. If three retries are required to get stable ECC syndrome before a counter 01h type error is corrected, then those three retries are also counted here. The number of retries applied to unsuccessfully recover an error (counter 06h type error) are also counted by this counter. ",
	},
	"write_total_uncorrected_errors": {
		ID:          "write_total_uncorrected_errors",
		DisplayName: "Write Total Uncorrected Errors",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: " This parameter code specifies the counter that contains the total number of blocks for which an uncorrected data error has occurred.",
	},
}
