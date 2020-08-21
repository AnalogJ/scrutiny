// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

const wqlProduct = "SELECT Caption, Description, IdentifyingNumber, Name, SKUNumber, Vendor, Version, UUID FROM Win32_ComputerSystemProduct"

type win32Product struct {
	Caption           *string
	Description       *string
	IdentifyingNumber *string
	Name              *string
	SKUNumber         *string
	Vendor            *string
	Version           *string
	UUID              *string
}

func (ctx *context) productFillInfo(info *ProductInfo) error {
	// Getting data from WMI
	var win32ProductDescriptions []win32Product
	// Assuming the first product is the host...
	if err := wmi.Query(wqlProduct, &win32ProductDescriptions); err != nil {
		return err
	}
	if len(win32ProductDescriptions) > 0 {
		info.Family = UNKNOWN
		info.Name = *win32ProductDescriptions[0].Name
		info.Vendor = *win32ProductDescriptions[0].Vendor
		info.SerialNumber = *win32ProductDescriptions[0].IdentifyingNumber
		info.UUID = *win32ProductDescriptions[0].UUID
		info.SKU = *win32ProductDescriptions[0].SKUNumber
		info.Version = *win32ProductDescriptions[0].Version
	}

	return nil
}
