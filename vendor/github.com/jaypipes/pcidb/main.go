//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pcidb

import (
	"fmt"
	"os"
	"strconv"
)

var (
	ERR_NO_DB = fmt.Errorf("No pci-ids DB files found (and network fetch disabled)")
	trueVar   = true
)

// ProgrammingInterface is the PCI programming interface for a class of PCI
// devices
type ProgrammingInterface struct {
	// hex-encoded PCI_ID of the programming interface
	ID string `json:"id"`
	// common string name for the programming interface
	Name string `json:"name"`
}

// Subclass is a subdivision of a PCI class
type Subclass struct {
	// hex-encoded PCI_ID for the device subclass
	ID string `json:"id"`
	// common string name for the subclass
	Name string `json:"name"`
	// any programming interfaces this subclass might have
	ProgrammingInterfaces []*ProgrammingInterface `json:"programming_interfaces"`
}

// Class is the PCI class
type Class struct {
	// hex-encoded PCI_ID for the device class
	ID string `json:"id"`
	// common string name for the class
	Name string `json:"name"`
	// any subclasses belonging to this class
	Subclasses []*Subclass `json:"subclasses"`
}

// Product provides information about a PCI device model
// NOTE(jaypipes): In the hardware world, the PCI "device_id" is the identifier
// for the product/model
type Product struct {
	// vendor ID for the product
	VendorID string `json:"vendor_id"`
	// hex-encoded PCI_ID for the product/model
	ID string `json:"id"`
	// common string name of the vendor
	Name string `json:"name"`
	// "subdevices" or "subsystems" for the product
	Subsystems []*Product `json:"subsystems"`
}

// Vendor provides information about a device vendor
type Vendor struct {
	// hex-encoded PCI_ID for the vendor
	ID string `json:"id"`
	// common string name of the vendor
	Name string `json:"name"`
	// all top-level devices for the vendor
	Products []*Product `json:"products"`
}

type PCIDB struct {
	// hash of class ID -> class information
	Classes map[string]*Class `json:"classes"`
	// hash of vendor ID -> vendor information
	Vendors map[string]*Vendor `json:"vendors"`
	// hash of vendor ID + product/device ID -> product information
	Products map[string]*Product `json:"products"`
}

// WithOption is used to represent optionally-configured settings
type WithOption struct {
	// Chroot is the directory that pcidb uses when attempting to discover
	// pciids DB files
	Chroot *string
	// CacheOnly is mostly just useful for testing. It essentially disables
	// looking for any non ~/.cache/pci.ids filepaths (which is useful when we
	// want to test the fetch-from-network code paths
	CacheOnly *bool
	// Disables the default behaviour of fetching a pci-ids from a known
	// location on the network if no local pci-ids DB files can be found.
	// Useful for secure environments or environments with no network
	// connectivity.
	DisableNetworkFetch *bool
}

func WithChroot(dir string) *WithOption {
	return &WithOption{Chroot: &dir}
}

func WithCacheOnly() *WithOption {
	return &WithOption{CacheOnly: &trueVar}
}

func WithDisableNetworkFetch() *WithOption {
	return &WithOption{DisableNetworkFetch: &trueVar}
}

func mergeOptions(opts ...*WithOption) *WithOption {
	// Grab options from the environs by default
	defaultChroot := "/"
	if val, exists := os.LookupEnv("PCIDB_CHROOT"); exists {
		defaultChroot = val
	}
	defaultCacheOnly := false
	if val, exists := os.LookupEnv("PCIDB_CACHE_ONLY"); exists {
		if parsed, err := strconv.ParseBool(val); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Failed parsing a bool from PCIDB_CACHE_ONLY "+
					"environ value of %s",
				val,
			)
		} else if parsed {
			defaultCacheOnly = parsed
		}
	}
	defaultDisableNetworkFetch := false
	if val, exists := os.LookupEnv("PCIDB_DISABLE_NETWORK_FETCH"); exists {
		if parsed, err := strconv.ParseBool(val); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Failed parsing a bool from PCIDB_DISABLE_NETWORK_FETCH "+
					"environ value of %s",
				val,
			)
		} else if parsed {
			defaultDisableNetworkFetch = parsed
		}
	}

	merged := &WithOption{}
	for _, opt := range opts {
		if opt.Chroot != nil {
			merged.Chroot = opt.Chroot
		}
		if opt.CacheOnly != nil {
			merged.CacheOnly = opt.CacheOnly
		}
		if opt.DisableNetworkFetch != nil {
			merged.DisableNetworkFetch = opt.DisableNetworkFetch
		}
	}
	// Set the default value if missing from merged
	if merged.Chroot == nil {
		merged.Chroot = &defaultChroot
	}
	if merged.CacheOnly == nil {
		merged.CacheOnly = &defaultCacheOnly
	}
	if merged.DisableNetworkFetch == nil {
		merged.DisableNetworkFetch = &defaultDisableNetworkFetch
	}
	return merged
}

// New returns a pointer to a PCIDB struct which contains information you can
// use to query PCI vendor, product and class information. It accepts zero or
// more pointers to WithOption structs. If you want to modify the behaviour of
// pcidb, use one of the option modifiers when calling New. For example, to
// change the root directory that pcidb uses when discovering pciids DB files,
// call New(WithChroot("/my/root/override"))
func New(opts ...*WithOption) (*PCIDB, error) {
	ctx := contextFromOptions(mergeOptions(opts...))
	db := &PCIDB{}
	if err := db.load(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
