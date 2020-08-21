//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "fmt"

// BIOSInfo defines BIOS release information
type BIOSInfo struct {
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
	Date    string `json:"date"`
}

func (i *BIOSInfo) String() string {

	vendorStr := ""
	if i.Vendor != "" {
		vendorStr = " vendor=" + i.Vendor
	}
	versionStr := ""
	if i.Version != "" {
		versionStr = " version=" + i.Version
	}
	dateStr := ""
	if i.Date != "" && i.Date != UNKNOWN {
		dateStr = " date=" + i.Date
	}

	res := fmt.Sprintf(
		"bios%s%s%s",
		vendorStr,
		versionStr,
		dateStr,
	)
	return res
}

// BIOS returns a pointer to a BIOSInfo struct containing information
// about the host's BIOS
func BIOS(opts ...*WithOption) (*BIOSInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &BIOSInfo{}
	if err := ctx.biosFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate BIOS information in a top-level
// "bios" YAML/JSON map/object key
type biosPrinter struct {
	Info *BIOSInfo `json:"bios"`
}

// YAMLString returns a string with the BIOS information formatted as YAML
// under a top-level "dmi:" key
func (info *BIOSInfo) YAMLString() string {
	return safeYAML(biosPrinter{info})
}

// JSONString returns a string with the BIOS information formatted as JSON
// under a top-level "bios:" key
func (info *BIOSInfo) JSONString(indent bool) string {
	return safeJSON(biosPrinter{info}, indent)
}
