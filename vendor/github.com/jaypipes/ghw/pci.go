//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/jaypipes/pcidb"
)

var (
	regexPCIAddress *regexp.Regexp = regexp.MustCompile(
		`^(([0-9a-f]{0,4}):)?([0-9a-f]{2}):([0-9a-f]{2})\.([0-9a-f]{1})$`,
	)
)

type PCIDevice struct {
	// The PCI address of the device
	Address   string         `json:"address"`
	Vendor    *pcidb.Vendor  `json:"vendor"`
	Product   *pcidb.Product `json:"product"`
	Subsystem *pcidb.Product `json:"subsystem"`
	// optional subvendor/sub-device information
	Class *pcidb.Class `json:"class"`
	// optional sub-class for the device
	Subclass *pcidb.Subclass `json:"subclass"`
	// optional programming interface
	ProgrammingInterface *pcidb.ProgrammingInterface `json:"programming_interface"`
}

// NOTE(jaypipes) PCIDevice has a custom JSON marshaller because we don't want
// to serialize the entire PCIDB information for the Vendor (which includes all
// of the vendor's products, etc). Instead, we simply serialize the ID and
// human-readable name of the vendor, product, class, etc.
func (pd *PCIDevice) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString("{")
	b.WriteString(fmt.Sprintf("\"address\":\"%s\"", pd.Address))
	b.WriteString(",\"vendor\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			pd.Vendor.ID,
			pd.Vendor.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"product\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			pd.Product.ID,
			pd.Product.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"subsystem\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			pd.Subsystem.ID,
			pd.Subsystem.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"class\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			pd.Class.ID,
			pd.Class.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"subclass\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			pd.Subclass.ID,
			pd.Subclass.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"programming_interface\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			pd.ProgrammingInterface.ID,
			pd.ProgrammingInterface.Name,
		),
	)
	b.WriteString("}")
	b.WriteString("}")
	return b.Bytes(), nil
}

func (di *PCIDevice) String() string {
	vendorName := UNKNOWN
	if di.Vendor != nil {
		vendorName = di.Vendor.Name
	}
	productName := UNKNOWN
	if di.Product != nil {
		productName = di.Product.Name
	}
	className := UNKNOWN
	if di.Class != nil {
		className = di.Class.Name
	}
	return fmt.Sprintf(
		"%s -> class: '%s' vendor: '%s' product: '%s'",
		di.Address,
		className,
		vendorName,
		productName,
	)
}

type PCIInfo struct {
	ctx *context
	// hash of class ID -> class information
	Classes map[string]*pcidb.Class
	// hash of vendor ID -> vendor information
	Vendors map[string]*pcidb.Vendor
	// hash of vendor ID + product/device ID -> product information
	Products map[string]*pcidb.Product
}

type PCIAddress struct {
	Domain   string
	Bus      string
	Slot     string
	Function string
}

// Given a string address, returns a complete PCIAddress struct, filled in with
// domain, bus, slot and function components. The address string may either
// be in $BUS:$SLOT.$FUNCTION (BSF) format or it can be a full PCI address
// that includes the 4-digit $DOMAIN information as well:
// $DOMAIN:$BUS:$SLOT.$FUNCTION.
//
// Returns "" if the address string wasn't a valid PCI address.
func PCIAddressFromString(address string) *PCIAddress {
	addrLowered := strings.ToLower(address)
	matches := regexPCIAddress.FindStringSubmatch(addrLowered)
	if len(matches) == 6 {
		dom := "0000"
		if matches[1] != "" {
			dom = matches[2]
		}
		return &PCIAddress{
			Domain:   dom,
			Bus:      matches[3],
			Slot:     matches[4],
			Function: matches[5],
		}
	}
	return nil
}

func PCI(opts ...*WithOption) (*PCIInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &PCIInfo{}
	if err := ctx.pciFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}
