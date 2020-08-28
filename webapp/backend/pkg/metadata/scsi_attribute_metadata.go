package metadata

type ScsiAttributeMetadata struct {
	ID          string `json:"-"`
	DisplayName string `json:"-"`
	Ideal       string `json:"ideal"`
	Critical    bool   `json:"critical"`
	Description string `json:"description"`

	Transform          func(int, int64, string) int64 `json:"-"` //this should be a method to extract/tranform the normalized or raw data to a chartable format. Str
	TransformValueUnit string                         `json:"transform_value_unit,omitempty"`
	DisplayType        string                         `json:"display_type"` //"raw" "normalized" or "transformed"
}

var ScsiMetadata = map[string]ScsiAttributeMetadata{
	"scsi_grown_defect_list": {
		ID:          "scsi_grown_defect_list",
		DisplayName: "Grown Defect List",
		DisplayType: "",
		Ideal:       "",
		Critical:    true,
		Description: "",
	},
	"read.errors_corrected_by_eccfast": {
		ID:          "read.errors_corrected_by_eccfast",
		DisplayName: "Read Errors Corrected by ECC Fast",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"read.errors_corrected_by_eccdelayed": {
		ID:          "read.errors_corrected_by_eccdelayed",
		DisplayName: "Read Errors Corrected by ECC Delayed",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"read.errors_corrected_by_rereads_rewrites": {
		ID:          "read.errors_corrected_by_rereads_rewrites",
		DisplayName: "Read Errors Corrected by ReReads/ReWrites",
		DisplayType: "",
		Ideal:       "",
		Critical:    true,
		Description: "",
	},
	"read.total_errors_corrected": {
		ID:          "read.total_errors_corrected",
		DisplayName: "Read Total Errors Corrected",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"read.correction_algorithm_invocations": {
		ID:          "read.correction_algorithm_invocations",
		DisplayName: "Read Correction Algorithm Invocations",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"read.total_uncorrected_errors": {
		ID:          "read.total_uncorrected_errors",
		DisplayName: "Read Total Uncorrected Errors",
		DisplayType: "",
		Ideal:       "",
		Critical:    true,
		Description: "",
	},
	"write.errors_corrected_by_eccfast": {
		ID:          "write.errors_corrected_by_eccfast",
		DisplayName: "Write Errors Corrected by ECC Fast",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"write.errors_corrected_by_eccdelayed": {
		ID:          "write.errors_corrected_by_eccdelayed",
		DisplayName: "Write Errors Corrected by ECC Delayed",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"write.errors_corrected_by_rereads_rewrites": {
		ID:          "write.errors_corrected_by_rereads_rewrites",
		DisplayName: "Write Errors Corrected by ReReads/ReWrites",
		DisplayType: "",
		Ideal:       "",
		Critical:    true,
		Description: "",
	},
	"write.total_errors_corrected": {
		ID:          "write.total_errors_corrected",
		DisplayName: "Write Total Errors Corrected",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"write.correction_algorithm_invocations": {
		ID:          "write.correction_algorithm_invocations",
		DisplayName: "Write Correction Algorithm Invocations",
		DisplayType: "",
		Ideal:       "",
		Critical:    false,
		Description: "",
	},
	"write.total_uncorrected_errors": {
		ID:          "write.total_uncorrected_errors",
		DisplayName: "Write Total Uncorrected Errors",
		DisplayType: "",
		Ideal:       "",
		Critical:    true,
		Description: "",
	},
}
