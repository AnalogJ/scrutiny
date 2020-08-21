//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "fmt"

// BaseboardInfo defines baseboard release information
type BaseboardInfo struct {
	AssetTag     string `json:"asset_tag"`
	SerialNumber string `json:"serial_number"`
	Vendor       string `json:"vendor"`
	Version      string `json:"version"`
}

func (i *BaseboardInfo) String() string {

	vendorStr := ""
	if i.Vendor != "" {
		vendorStr = " vendor=" + i.Vendor
	}
	serialStr := ""
	if i.SerialNumber != "" && i.SerialNumber != UNKNOWN {
		serialStr = " serial=" + i.SerialNumber
	}
	versionStr := ""
	if i.Version != "" {
		versionStr = " version=" + i.Version
	}

	res := fmt.Sprintf(
		"baseboard%s%s%s",
		vendorStr,
		serialStr,
		versionStr,
	)
	return res
}

// Baseboard returns a pointer to a BaseboardInfo struct containing information
// about the host's baseboard
func Baseboard(opts ...*WithOption) (*BaseboardInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &BaseboardInfo{}
	if err := ctx.baseboardFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate baseboard information in a top-level
// "baseboard" YAML/JSON map/object key
type baseboardPrinter struct {
	Info *BaseboardInfo `json:"baseboard"`
}

// YAMLString returns a string with the baseboard information formatted as YAML
// under a top-level "dmi:" key
func (info *BaseboardInfo) YAMLString() string {
	return safeYAML(baseboardPrinter{info})
}

// JSONString returns a string with the baseboard information formatted as JSON
// under a top-level "baseboard:" key
func (info *BaseboardInfo) JSONString(indent bool) string {
	return safeJSON(baseboardPrinter{info}, indent)
}
