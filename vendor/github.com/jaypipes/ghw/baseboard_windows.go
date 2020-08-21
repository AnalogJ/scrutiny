// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "github.com/StackExchange/wmi"

const wqlBaseboard = "SELECT Manufacturer, SerialNumber, Tag, Version FROM Win32_BaseBoard"

type win32Baseboard struct {
	Manufacturer *string
	SerialNumber *string
	Tag          *string
	Version      *string
}

func (ctx *context) baseboardFillInfo(info *BaseboardInfo) error {
	// Getting data from WMI
	var win32BaseboardDescriptions []win32Baseboard
	if err := wmi.Query(wqlBaseboard, &win32BaseboardDescriptions); err != nil {
		return err
	}
	if len(win32BaseboardDescriptions) > 0 {
		info.AssetTag = *win32BaseboardDescriptions[0].Tag
		info.SerialNumber = *win32BaseboardDescriptions[0].SerialNumber
		info.Vendor = *win32BaseboardDescriptions[0].Manufacturer
		info.Version = *win32BaseboardDescriptions[0].Version
	}

	return nil
}
