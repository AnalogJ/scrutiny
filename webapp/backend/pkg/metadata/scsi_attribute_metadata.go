package metadata

type ScsiAttributeMetadata struct {
	ID          string `json:"-"`
	DisplayName string `json:"-"`
	Ideal       string `json:"ideal"`
	Critical    bool   `json:"critical"`
	Description string `json:"description"`

	Transform          func(int64, int64, string) int64 `json:"-"` //this should be a method to extract/tranform the normalized or raw data to a chartable format. Str
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
		Description: "",
	},
	"read_errors_corrected_by_eccfast": {
		ID:          "read_errors_corrected_by_eccfast",
		DisplayName: "Read Errors Corrected by ECC Fast",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"read_errors_corrected_by_eccdelayed": {
		ID:          "read_errors_corrected_by_eccdelayed",
		DisplayName: "Read Errors Corrected by ECC Delayed",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"read_errors_corrected_by_rereads_rewrites": {
		ID:          "read_errors_corrected_by_rereads_rewrites",
		DisplayName: "Read Errors Corrected by ReReads/ReWrites",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "",
	},
	"read_total_errors_corrected": {
		ID:          "read_total_errors_corrected",
		DisplayName: "Read Total Errors Corrected",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"read_correction_algorithm_invocations": {
		ID:          "read_correction_algorithm_invocations",
		DisplayName: "Read Correction Algorithm Invocations",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"read_total_uncorrected_errors": {
		ID:          "read_total_uncorrected_errors",
		DisplayName: "Read Total Uncorrected Errors",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "",
	},
	"write_errors_corrected_by_eccfast": {
		ID:          "write_errors_corrected_by_eccfast",
		DisplayName: "Write Errors Corrected by ECC Fast",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"write_errors_corrected_by_eccdelayed": {
		ID:          "write_errors_corrected_by_eccdelayed",
		DisplayName: "Write Errors Corrected by ECC Delayed",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"write_errors_corrected_by_rereads_rewrites": {
		ID:          "write_errors_corrected_by_rereads_rewrites",
		DisplayName: "Write Errors Corrected by ReReads/ReWrites",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "",
	},
	"write_total_errors_corrected": {
		ID:          "write_total_errors_corrected",
		DisplayName: "Write Total Errors Corrected",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"write_correction_algorithm_invocations": {
		ID:          "write_correction_algorithm_invocations",
		DisplayName: "Write Correction Algorithm Invocations",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"write_total_uncorrected_errors": {
		ID:          "write_total_uncorrected_errors",
		DisplayName: "Write Total Uncorrected Errors",
		DisplayType: "",
		Ideal:       "low",
		Critical:    true,
		Description: "",
	},
}
