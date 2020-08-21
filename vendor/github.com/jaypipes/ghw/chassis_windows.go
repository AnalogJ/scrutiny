// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

const wqlChassis = "SELECT Caption, Description, Name, Manufacturer, Model, SerialNumber, Tag, TypeDescriptions, Version FROM CIM_Chassis"

type win32Chassis struct {
	Caption          *string
	Description      *string
	Name             *string
	Manufacturer     *string
	Model            *string
	SerialNumber     *string
	Tag              *string
	TypeDescriptions []string
	Version          *string
}

func (ctx *context) chassisFillInfo(info *ChassisInfo) error {
	// Getting data from WMI
	var win32ChassisDescriptions []win32Chassis
	if err := wmi.Query(wqlChassis, &win32ChassisDescriptions); err != nil {
		return err
	}
	if len(win32ChassisDescriptions) > 0 {
		info.AssetTag = *win32ChassisDescriptions[0].Tag
		info.SerialNumber = *win32ChassisDescriptions[0].SerialNumber
		info.Type = UNKNOWN // TODO:
		info.TypeDescription = *win32ChassisDescriptions[0].Model
		info.Vendor = *win32ChassisDescriptions[0].Manufacturer
		info.Version = *win32ChassisDescriptions[0].Version
	}
	return nil
}
