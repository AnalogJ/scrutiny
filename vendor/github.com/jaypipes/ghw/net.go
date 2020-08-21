//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
)

type NICCapability struct {
	Name      string `json:"name"`
	IsEnabled bool   `json:"is_enabled"`
	CanEnable bool   `json:"can_enable"`
}

type NIC struct {
	Name         string           `json:"name"`
	MacAddress   string           `json:"mac_address"`
	IsVirtual    bool             `json:"is_virtual"`
	Capabilities []*NICCapability `json:"capabilities"`
	// TODO(jaypipes): Add PCI field for accessing PCI device information
	// PCI *PCIDevice `json:"pci"`
}

func (n *NIC) String() string {
	isVirtualStr := ""
	if n.IsVirtual {
		isVirtualStr = " (virtual)"
	}
	return fmt.Sprintf(
		"%s%s",
		n.Name,
		isVirtualStr,
	)
}

type NetworkInfo struct {
	NICs []*NIC `json:"nics"`
}

func Network(opts ...*WithOption) (*NetworkInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &NetworkInfo{}
	if err := ctx.netFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

func (i *NetworkInfo) String() string {
	return fmt.Sprintf(
		"net (%d NICs)",
		len(i.NICs),
	)
}

// simple private struct used to encapsulate net information in a
// top-level "net" YAML/JSON map/object key
type netPrinter struct {
	Info *NetworkInfo `json:"network"`
}

// YAMLString returns a string with the net information formatted as YAML
// under a top-level "net:" key
func (i *NetworkInfo) YAMLString() string {
	return safeYAML(netPrinter{i})
}

// JSONString returns a string with the net information formatted as JSON
// under a top-level "net:" key
func (i *NetworkInfo) JSONString(indent bool) string {
	return safeJSON(netPrinter{i}, indent)
}
