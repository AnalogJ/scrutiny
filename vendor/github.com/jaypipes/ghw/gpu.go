//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
)

type GraphicsCard struct {
	// the PCI address where the graphics card can be found
	Address string `json:"address"`
	// The "index" of the card on the bus (generally not useful information,
	// but might as well include it)
	Index int `json:"index"`
	// pointer to a PCIDevice struct that describes the vendor and product
	// model, etc
	// TODO(jaypipes): Rename this field to PCI, instead of DeviceInfo
	DeviceInfo *PCIDevice `json:"pci"`
	// Topology node that the graphics card is affined to. Will be nil if the
	// architecture is not NUMA.
	Node *TopologyNode `json:"node,omitempty"`
}

func (card *GraphicsCard) String() string {
	deviceStr := card.Address
	if card.DeviceInfo != nil {
		deviceStr = card.DeviceInfo.String()
	}
	nodeStr := ""
	if card.Node != nil {
		nodeStr = fmt.Sprintf(" [affined to NUMA node %d]", card.Node.Id)
	}
	return fmt.Sprintf(
		"card #%d %s@%s",
		card.Index,
		nodeStr,
		deviceStr,
	)
}

type GPUInfo struct {
	GraphicsCards []*GraphicsCard `json:"cards"`
}

func GPU(opts ...*WithOption) (*GPUInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &GPUInfo{}
	if err := ctx.gpuFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

func (i *GPUInfo) String() string {
	numCardsStr := "cards"
	if len(i.GraphicsCards) == 1 {
		numCardsStr = "card"
	}
	return fmt.Sprintf(
		"gpu (%d graphics %s)",
		len(i.GraphicsCards),
		numCardsStr,
	)
}

// simple private struct used to encapsulate gpu information in a top-level
// "gpu" YAML/JSON map/object key
type gpuPrinter struct {
	Info *GPUInfo `json:"gpu"`
}

// YAMLString returns a string with the gpu information formatted as YAML
// under a top-level "gpu:" key
func (i *GPUInfo) YAMLString() string {
	return safeYAML(gpuPrinter{i})
}

// JSONString returns a string with the gpu information formatted as JSON
// under a top-level "gpu:" key
func (i *GPUInfo) JSONString(indent bool) string {
	return safeJSON(gpuPrinter{i}, indent)
}
