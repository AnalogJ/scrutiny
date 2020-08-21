// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaypipes/pcidb"
)

func (ctx *context) pciFillInfo(info *PCIInfo) error {
	db, err := pcidb.New(pcidb.WithChroot(ctx.chroot))
	if err != nil {
		return err
	}
	info.Classes = db.Classes
	info.Vendors = db.Vendors
	info.Products = db.Products
	info.ctx = ctx
	return nil
}

func (ctx *context) getPCIDeviceModaliasPath(address string) string {
	pciAddr := PCIAddressFromString(address)
	if pciAddr == nil {
		return ""
	}
	return filepath.Join(
		ctx.pathSysBusPciDevices(),
		pciAddr.Domain+":"+pciAddr.Bus+":"+pciAddr.Slot+"."+pciAddr.Function,
		"modalias",
	)
}

type deviceModaliasInfo struct {
	vendorID     string
	productID    string
	subproductID string
	subvendorID  string
	classID      string
	subclassID   string
	progIfaceID  string
}

func parseModaliasFile(fp string) *deviceModaliasInfo {
	if _, err := os.Stat(fp); err != nil {
		return nil
	}
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil
	}
	// The modalias file is an encoded file that looks like this:
	//
	// $ cat /sys/devices/pci0000\:00/0000\:00\:03.0/0000\:03\:00.0/modalias
	// pci:v000010DEd00001C82sv00001043sd00008613bc03sc00i00
	//
	// It is interpreted like so:
	//
	// pci: -- ignore
	// v000010DE -- PCI vendor ID
	// d00001C82 -- PCI device ID (the product/model ID)
	// sv00001043 -- PCI subsystem vendor ID
	// sd00008613 -- PCI subsystem device ID (subdevice product/model ID)
	// bc03 -- PCI base class
	// sc00 -- PCI subclass
	// i00 -- programming interface
	vendorID := strings.ToLower(string(data[9:13]))
	productID := strings.ToLower(string(data[18:22]))
	subvendorID := strings.ToLower(string(data[28:32]))
	subproductID := strings.ToLower(string(data[38:42]))
	classID := string(data[44:46])
	subclassID := string(data[48:50])
	progIfaceID := string(data[51:53])
	return &deviceModaliasInfo{
		vendorID:     vendorID,
		productID:    productID,
		subproductID: subproductID,
		subvendorID:  subvendorID,
		classID:      classID,
		subclassID:   subclassID,
		progIfaceID:  progIfaceID,
	}
}

// Returns a pointer to a pcidb.Vendor struct matching the supplied vendor
// ID string. If no such vendor ID string could be found, returns the
// pcidb.Vendor struct populated with "unknown" vendor Name attribute and
// empty Products attribute.
func findPCIVendor(info *PCIInfo, vendorID string) *pcidb.Vendor {
	vendor := info.Vendors[vendorID]
	if vendor == nil {
		return &pcidb.Vendor{
			ID:       vendorID,
			Name:     UNKNOWN,
			Products: []*pcidb.Product{},
		}
	}
	return vendor
}

// Returns a pointer to a pcidb.Product struct matching the supplied vendor
// and product ID strings. If no such product could be found, returns the
// pcidb.Product struct populated with "unknown" product Name attribute and
// empty Subsystems attribute.
func findPCIProduct(
	info *PCIInfo,
	vendorID string,
	productID string,
) *pcidb.Product {
	product := info.Products[vendorID+productID]
	if product == nil {
		return &pcidb.Product{
			ID:         productID,
			Name:       UNKNOWN,
			Subsystems: []*pcidb.Product{},
		}
	}
	return product
}

// Returns a pointer to a pcidb.Product struct matching the supplied vendor,
// product, subvendor and subproduct ID strings. If no such product could be
// found, returns the pcidb.Product struct populated with "unknown" product
// Name attribute and empty Subsystems attribute.
func findPCISubsystem(
	info *PCIInfo,
	vendorID string,
	productID string,
	subvendorID string,
	subproductID string,
) *pcidb.Product {
	product := info.Products[vendorID+productID]
	subvendor := info.Vendors[subvendorID]
	if subvendor != nil && product != nil {
		for _, p := range product.Subsystems {
			if p.ID == subproductID {
				return p
			}
		}
	}
	return &pcidb.Product{
		VendorID: subvendorID,
		ID:       subproductID,
		Name:     UNKNOWN,
	}
}

// Returns a pointer to a pcidb.Class struct matching the supplied class ID
// string. If no such class ID string could be found, returns the
// pcidb.Class struct populated with "unknown" class Name attribute and
// empty Subclasses attribute.
func findPCIClass(info *PCIInfo, classID string) *pcidb.Class {
	class := info.Classes[classID]
	if class == nil {
		return &pcidb.Class{
			ID:         classID,
			Name:       UNKNOWN,
			Subclasses: []*pcidb.Subclass{},
		}
	}
	return class
}

// Returns a pointer to a pcidb.Subclass struct matching the supplied class
// and subclass ID strings.  If no such subclass could be found, returns the
// pcidb.Subclass struct populated with "unknown" subclass Name attribute
// and empty ProgrammingInterfaces attribute.
func findPCISubclass(
	info *PCIInfo,
	classID string,
	subclassID string,
) *pcidb.Subclass {
	class := info.Classes[classID]
	if class != nil {
		for _, sc := range class.Subclasses {
			if sc.ID == subclassID {
				return sc
			}
		}
	}
	return &pcidb.Subclass{
		ID:   subclassID,
		Name: UNKNOWN,
		ProgrammingInterfaces: []*pcidb.ProgrammingInterface{},
	}
}

// Returns a pointer to a pcidb.ProgrammingInterface struct matching the
// supplied class, subclass and programming interface ID strings.  If no such
// programming interface could be found, returns the
// pcidb.ProgrammingInterface struct populated with "unknown" Name attribute
func findPCIProgrammingInterface(
	info *PCIInfo,
	classID string,
	subclassID string,
	progIfaceID string,
) *pcidb.ProgrammingInterface {
	subclass := findPCISubclass(info, classID, subclassID)
	for _, pi := range subclass.ProgrammingInterfaces {
		if pi.ID == progIfaceID {
			return pi
		}
	}
	return &pcidb.ProgrammingInterface{
		ID:   progIfaceID,
		Name: UNKNOWN,
	}
}

// GetDevice returns a pointer to a PCIDevice struct that describes the PCI
// device at the requested address. If no such device could be found, returns
// nil
func (info *PCIInfo) GetDevice(address string) *PCIDevice {
	fp := info.ctx.getPCIDeviceModaliasPath(address)
	if fp == "" {
		return nil
	}

	modaliasInfo := parseModaliasFile(fp)
	if modaliasInfo == nil {
		return nil
	}

	vendor := findPCIVendor(info, modaliasInfo.vendorID)
	product := findPCIProduct(
		info,
		modaliasInfo.vendorID,
		modaliasInfo.productID,
	)
	subsystem := findPCISubsystem(
		info,
		modaliasInfo.vendorID,
		modaliasInfo.productID,
		modaliasInfo.subvendorID,
		modaliasInfo.subproductID,
	)
	class := findPCIClass(info, modaliasInfo.classID)
	subclass := findPCISubclass(
		info,
		modaliasInfo.classID,
		modaliasInfo.subclassID,
	)
	progIface := findPCIProgrammingInterface(
		info,
		modaliasInfo.classID,
		modaliasInfo.subclassID,
		modaliasInfo.progIfaceID,
	)

	return &PCIDevice{
		Address:              address,
		Vendor:               vendor,
		Subsystem:            subsystem,
		Product:              product,
		Class:                class,
		Subclass:             subclass,
		ProgrammingInterface: progIface,
	}
}

// ListDevices returns a list of pointers to PCIDevice structs present on the
// host system
func (info *PCIInfo) ListDevices() []*PCIDevice {
	devs := make([]*PCIDevice, 0)
	// We scan the /sys/bus/pci/devices directory which contains a collection
	// of symlinks. The names of the symlinks are all the known PCI addresses
	// for the host. For each address, we grab a *PCIDevice matching the
	// address and append to the returned array.
	links, err := ioutil.ReadDir(info.ctx.pathSysBusPciDevices())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: failed to read /sys/bus/pci/devices")
		return nil
	}
	var dev *PCIDevice
	for _, link := range links {
		addr := link.Name()
		dev = info.GetDevice(addr)
		if dev == nil {
			_, _ = fmt.Fprintf(os.Stderr, "error: failed to get device information for PCI address %s\n", addr)
		} else {
			devs = append(devs, dev)
		}
	}
	return devs
}
