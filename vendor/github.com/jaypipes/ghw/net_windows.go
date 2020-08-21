// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"strings"

	"github.com/StackExchange/wmi"
)

const wqlNetworkAdapter = "SELECT Description, DeviceID, Index, InterfaceIndex, MACAddress, Manufacturer, Name, NetConnectionID, ProductName, ServiceName  FROM Win32_NetworkAdapter"

type win32NetworkAdapter struct {
	Description     *string
	DeviceID        *string
	Index           *uint32
	InterfaceIndex  *uint32
	MACAddress      *string
	Manufacturer    *string
	Name            *string
	NetConnectionID *string
	ProductName     *string
	ServiceName     *string
}

func (ctx *context) netFillInfo(info *NetworkInfo) error {
	// Getting info from WMI
	var win32NetDescriptions []win32NetworkAdapter
	if err := wmi.Query(wqlNetworkAdapter, &win32NetDescriptions); err != nil {
		return err
	}

	info.NICs = ctx.nics(win32NetDescriptions)
	return nil
}

func (ctx *context) nics(win32NetDescriptions []win32NetworkAdapter) []*NIC {
	// Converting into standard structures
	nics := make([]*NIC, 0)
	for _, nicDescription := range win32NetDescriptions {
		nic := &NIC{
			Name:         ctx.netDeviceName(nicDescription),
			MacAddress:   *nicDescription.MACAddress,
			IsVirtual:    false,
			Capabilities: []*NICCapability{},
		}
		// Appenging NIC to NICs
		nics = append(nics, nic)
	}

	return nics
}

func (ctx *context) netDeviceName(description win32NetworkAdapter) string {
	var name string
	if strings.TrimSpace(*description.NetConnectionID) != "" {
		name = *description.NetConnectionID + " - " + *description.Description
	} else {
		name = *description.Description
	}
	return name
}
