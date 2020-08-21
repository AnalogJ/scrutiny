// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/jaypipes/pcidb"
)

const wqlVideoController = "SELECT Caption, CreationClassName, Description, DeviceID, Name, PNPDeviceID, SystemCreationClassName, SystemName, VideoArchitecture, VideoMemoryType, VideoModeDescription, VideoProcessor FROM Win32_VideoController"

type win32VideoController struct {
	Caption                 string
	CreationClassName       string
	Description             string
	DeviceID                string
	Name                    string
	PNPDeviceID             string
	SystemCreationClassName string
	SystemName              string
	VideoArchitecture       uint16
	VideoMemoryType         uint16
	VideoModeDescription    string
	VideoProcessor          string
}

const wqlPnPEntity = "SELECT Caption, CreationClassName, Description, DeviceID, Manufacturer, Name, PNPClass, PNPDeviceID FROM Win32_PnPEntity"

type win32PnPEntity struct {
	Caption           string
	CreationClassName string
	Description       string
	DeviceID          string
	Manufacturer      string
	Name              string
	PNPClass          string
	PNPDeviceID       string
}

func (ctx *context) gpuFillInfo(info *GPUInfo) error {
	// Getting data from WMI
	var win32VideoControllerDescriptions []win32VideoController
	if err := wmi.Query(wqlVideoController, &win32VideoControllerDescriptions); err != nil {
		return err
	}

	// Building dynamic WHERE clause with addresses to create a single query collecting all desired data
	queryAddresses := []string{}
	for _, description := range win32VideoControllerDescriptions {
		var queryAddres = strings.Replace(description.PNPDeviceID, "\\", `\\`, -1)
		queryAddresses = append(queryAddresses, "PNPDeviceID='"+queryAddres+"'")
	}
	whereClause := strings.Join(queryAddresses[:], " OR ")

	// Getting data from WMI
	var win32PnPDescriptions []win32PnPEntity
	var wqlPnPDevice = wqlPnPEntity + " WHERE " + whereClause
	if err := wmi.Query(wqlPnPDevice, &win32PnPDescriptions); err != nil {
		return err
	}

	// Converting into standard structures
	cards := make([]*GraphicsCard, 0)
	for _, description := range win32VideoControllerDescriptions {
		card := &GraphicsCard{
			Address:    description.DeviceID, // https://stackoverflow.com/questions/32073667/how-do-i-discover-the-pcie-bus-topology-and-slot-numbers-on-the-board
			Index:      0,
			DeviceInfo: ctx.GetDevice(description.PNPDeviceID, win32PnPDescriptions),
		}
		cards = append(cards, card)
	}
	info.GraphicsCards = cards
	return nil
}

func (ctx *context) GetDevice(id string, entities []win32PnPEntity) *PCIDevice {
	// Backslashing PnP address ID as requested by JSON and VMI query: https://docs.microsoft.com/en-us/windows/win32/wmisdk/where-clause
	var queryAddress = strings.Replace(id, "\\", `\\`, -1)
	// Preparing default structure
	var device = &PCIDevice{
		Address: queryAddress,
		Vendor: &pcidb.Vendor{
			ID:       UNKNOWN,
			Name:     UNKNOWN,
			Products: []*pcidb.Product{},
		},
		Subsystem: &pcidb.Product{
			ID:         UNKNOWN,
			Name:       UNKNOWN,
			Subsystems: []*pcidb.Product{},
		},
		Product: &pcidb.Product{
			ID:         UNKNOWN,
			Name:       UNKNOWN,
			Subsystems: []*pcidb.Product{},
		},
		Class: &pcidb.Class{
			ID:         UNKNOWN,
			Name:       UNKNOWN,
			Subclasses: []*pcidb.Subclass{},
		},
		Subclass: &pcidb.Subclass{
			ID:                    UNKNOWN,
			Name:                  UNKNOWN,
			ProgrammingInterfaces: []*pcidb.ProgrammingInterface{},
		},
		ProgrammingInterface: &pcidb.ProgrammingInterface{
			ID:   UNKNOWN,
			Name: UNKNOWN,
		},
	}
	// If an entity is found we get its data inside the standard structure
	for _, description := range entities {
		if id == description.PNPDeviceID {
			device.Vendor.ID = description.Manufacturer
			device.Vendor.Name = description.Manufacturer
			device.Product.ID = description.Name
			device.Product.Name = description.Description
			break
		}
	}
	return device
}
