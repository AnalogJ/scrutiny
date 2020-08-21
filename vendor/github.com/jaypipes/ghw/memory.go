//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"math"
)

type MemoryModule struct {
	Label        string `json:"label"`
	Location     string `json:"location"`
	SerialNumber string `json:"serial_number"`
	SizeBytes    int64  `json:"size_bytes"`
	Vendor       string `json:"vendor"`
}

type MemoryInfo struct {
	TotalPhysicalBytes int64 `json:"total_physical_bytes"`
	TotalUsableBytes   int64 `json:"total_usable_bytes"`
	// An array of sizes, in bytes, of memory pages supported by the host
	SupportedPageSizes []uint64        `json:"supported_page_sizes"`
	Modules            []*MemoryModule `json:"modules"`
}

func Memory(opts ...*WithOption) (*MemoryInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &MemoryInfo{}
	if err := ctx.memFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

func (i *MemoryInfo) String() string {
	tpbs := UNKNOWN
	if i.TotalPhysicalBytes > 0 {
		tpb := i.TotalPhysicalBytes
		unit, unitStr := unitWithString(tpb)
		tpb = int64(math.Ceil(float64(i.TotalPhysicalBytes) / float64(unit)))
		tpbs = fmt.Sprintf("%d%s", tpb, unitStr)
	}
	tubs := UNKNOWN
	if i.TotalUsableBytes > 0 {
		tub := i.TotalUsableBytes
		unit, unitStr := unitWithString(tub)
		tub = int64(math.Ceil(float64(i.TotalUsableBytes) / float64(unit)))
		tubs = fmt.Sprintf("%d%s", tub, unitStr)
	}
	return fmt.Sprintf("memory (%s physical, %s usable)", tpbs, tubs)
}

// simple private struct used to encapsulate memory information in a top-level
// "memory" YAML/JSON map/object key
type memoryPrinter struct {
	Info *MemoryInfo `json:"memory"`
}

// YAMLString returns a string with the memory information formatted as YAML
// under a top-level "memory:" key
func (i *MemoryInfo) YAMLString() string {
	return safeYAML(memoryPrinter{i})
}

// JSONString returns a string with the memory information formatted as JSON
// under a top-level "memory:" key
func (i *MemoryInfo) JSONString(indent bool) string {
	return safeJSON(memoryPrinter{i}, indent)
}
