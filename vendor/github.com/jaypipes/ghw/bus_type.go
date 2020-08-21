//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "strings"

type BusType int

const (
	BUS_TYPE_UNKNOWN BusType = iota
	BUS_TYPE_IDE
	BUS_TYPE_PCI
	BUS_TYPE_SCSI
	BUS_TYPE_NVME
	BUS_TYPE_VIRTIO
)

var (
	busTypeString = map[BusType]string{
		BUS_TYPE_UNKNOWN: "Unknown",
		BUS_TYPE_IDE:     "IDE",
		BUS_TYPE_PCI:     "PCI",
		BUS_TYPE_SCSI:    "SCSI",
		BUS_TYPE_NVME:    "NVMe",
		BUS_TYPE_VIRTIO:  "Virtio",
	}
)

func (bt BusType) String() string {
	return busTypeString[bt]
}

// NOTE(jaypipes): since serialized output is as "official" as we're going to
// get, let's lowercase the string output when serializing, in order to
// "normalize" the expected serialized output
func (bt BusType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strings.ToLower(bt.String()) + "\""), nil
}
