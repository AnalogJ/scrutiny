# `pcidb` - the Golang PCI DB library [![Build Status](https://travis-ci.org/jaypipes/pcidb.svg?branch=master)](https://travis-ci.org/jaypipes/pcidb)

`pcidb` is a small Golang library for programmatic querying of PCI vendor,
product and class information.

We currently [test](https://travis-ci.org/jaypipes/pcidb/) `pcidb` on Linux, Windows and MacOSX.

## Usage

`pcidb` contains a PCI database inspection and querying facility that allows
developers to query for information about hardware device classes, vendor and
product information.

The `pcidb.New()` function returns a `pcidb.PCIDB` struct or an error if the
PCI database could not be loaded.

> `pcidb`'s default behaviour is to first search for pci-ids DB files on the
> local host system in well-known filesystem paths. If `pcidb` cannot find a
> pci-ids DB file on the local host system, it will then fetch a current
> pci-ids DB file from the network. You can disable this network-fetching
> behaviour with the `pcidb.WithDisableNetworkFetch()` function or set the
> `PCIDB_DISABLE_NETWORK_FETCH` to a non-0 value.

The `pcidb.PCIDB` struct contains a number of fields that may be queried for
PCI information:

* `pcidb.PCIDB.Classes` is a map, keyed by the PCI class ID (a hex-encoded
  string) of pointers to `pcidb.Class` structs, one for each class of PCI
  device known to `pcidb`
* `pcidb.PCIDB.Vendors` is a map, keyed by the PCI vendor ID (a hex-encoded
  string) of pointers to `pcidb.Vendor` structs, one for each PCI vendor
  known to `pcidb`
* `pcidb.PCIDB.Products` is a map, keyed by the PCI product ID* (a hex-encoded
  string) of pointers to `pcidb.Product` structs, one for each PCI product
  known to `pcidb`

**NOTE**: PCI products are often referred to by their "device ID". We use
the term "product ID" in `pcidb` because it more accurately reflects what the
identifier is for: a specific product line produced by the vendor.

### Overriding the root mountpoint `pcidb` uses

The default root mountpoint that `pcidb` uses when looking for information
about the host system is `/`. So, for example, when looking up known PCI IDS DB
files on Linux, `pcidb` will attempt to discover a pciids DB file at
`/usr/share/misc/pci.ids`. If you are calling `pcidb` from a system that has an
alternate root mountpoint, you can either set the `PCIDB_CHROOT` environment
variable to that alternate path, or call the `pcidb.New()` function with the
`pcidb.WithChroot()` modifier.

For example, if you are executing from within an application container that has
bind-mounted the root host filesystem to the mount point `/host`, you would set
`PCIDB_CHROOT` to `/host` so that pcidb can find files like
`/usr/share/misc/pci.ids` at `/host/usr/share/misc/pci.ids`.

Alternately, you can use the `pcidb.WithChroot()` function like so:

```go
pci := pcidb.New(pcidb.WithChroot("/host"))
```

### PCI device classes

Let's take a look at the PCI device class information and how to query the PCI
database for class, subclass, and programming interface information.

Each `pcidb.Class` struct contains the following fields:

* `pcidb.Class.ID` is the hex-encoded string identifier for the device
  class
* `pcidb.Class.Name` is the common name/description of the class
* `pcidb.Class.Subclasses` is an array of pointers to
  `pcidb.Subclass` structs, one for each subclass in the device class

Each `pcidb.Subclass` struct contains the following fields:

* `pcidb.Subclass.ID` is the hex-encoded string identifier for the device
  subclass
* `pcidb.Subclass.Name` is the common name/description of the subclass
* `pcidb.Subclass.ProgrammingInterfaces` is an array of pointers to
  `pcidb.ProgrammingInterface` structs, one for each programming interface
   for the device subclass

Each `pcidb.ProgrammingInterface` struct contains the following fields:

* `pcidb.ProgrammingInterface.ID` is the hex-encoded string identifier for
  the programming interface
* `pcidb.ProgrammingInterface.Name` is the common name/description for the
  programming interface

```go
package main

import (
    "fmt"

    "github.com/jaypipes/pcidb"
)

func main() {
    pci, err := pcidb.New()
    if err != nil {
        fmt.Printf("Error getting PCI info: %v", err)
    }

    for _, devClass := range pci.Classes {
        fmt.Printf(" Device class: %v ('%v')\n", devClass.Name, devClass.ID)
        for _, devSubclass := range devClass.Subclasses {
            fmt.Printf("    Device subclass: %v ('%v')\n", devSubclass.Name, devSubclass.ID)
            for _, progIface := range devSubclass.ProgrammingInterfaces {
                fmt.Printf("        Programming interface: %v ('%v')\n", progIface.Name, progIface.ID)
            }
        }
    }
}
```

Example output from my personal workstation, snipped for brevity:

```
...
 Device class: Serial bus controller ('0c')
    Device subclass: FireWire (IEEE 1394) ('00')
        Programming interface: Generic ('00')
        Programming interface: OHCI ('10')
    Device subclass: ACCESS Bus ('01')
    Device subclass: SSA ('02')
    Device subclass: USB controller ('03')
        Programming interface: UHCI ('00')
        Programming interface: OHCI ('10')
        Programming interface: EHCI ('20')
        Programming interface: XHCI ('30')
        Programming interface: Unspecified ('80')
        Programming interface: USB Device ('fe')
    Device subclass: Fibre Channel ('04')
    Device subclass: SMBus ('05')
    Device subclass: InfiniBand ('06')
    Device subclass: IPMI SMIC interface ('07')
    Device subclass: SERCOS interface ('08')
    Device subclass: CANBUS ('09')
...
```

### PCI vendors and products

Let's take a look at the PCI vendor information and how to query the PCI
database for vendor information and the products a vendor supplies.

Each `pcidb.Vendor` struct contains the following fields:

* `pcidb.Vendor.ID` is the hex-encoded string identifier for the vendor
* `pcidb.Vendor.Name` is the common name/description of the vendor
* `pcidb.Vendor.Products` is an array of pointers to `pcidb.Product`
  structs, one for each product supplied by the vendor

Each `pcidb.Product` struct contains the following fields:

* `pcidb.Product.VendorID` is the hex-encoded string identifier for the
  product's vendor
* `pcidb.Product.ID` is the hex-encoded string identifier for the product
* `pcidb.Product.Name` is the common name/description of the subclass
* `pcidb.Product.Subsystems` is an array of pointers to
  `pcidb.Product` structs, one for each "subsystem" (sometimes called
  "sub-device" in PCI literature) for the product

**NOTE**: A subsystem product may have a different vendor than its "parent" PCI
product. This is sometimes referred to as the "sub-vendor".

Here's some example code that demonstrates listing the PCI vendors with the
most known products:

```go
package main

import (
    "fmt"
    "sort"

    "github.com/jaypipes/pcidb"
)

type ByCountProducts []*pcidb.Vendor

func (v ByCountProducts) Len() int {
    return len(v)
}

func (v ByCountProducts) Swap(i, j int) {
    v[i], v[j] = v[j], v[i]
}

func (v ByCountProducts) Less(i, j int) bool {
    return len(v[i].Products) > len(v[j].Products)
}

func main() {
    pci, err := pcidb.New()
    if err != nil {
        fmt.Printf("Error getting PCI info: %v", err)
    }

    vendors := make([]*pcidb.Vendor, len(pci.Vendors))
    x := 0
    for _, vendor := range pci.Vendors {
        vendors[x] = vendor
        x++
    }

    sort.Sort(ByCountProducts(vendors))

    fmt.Println("Top 5 vendors by product")
    fmt.Println("====================================================")
    for _, vendor := range vendors[0:5] {
        fmt.Printf("%v ('%v') has %d products\n", vendor.Name, vendor.ID, len(vendor.Products))
    }
}
```

which yields (on my local workstation as of July 7th, 2018):

```
Top 5 vendors by product
====================================================
Intel Corporation ('8086') has 3389 products
NVIDIA Corporation ('10de') has 1358 products
Advanced Micro Devices, Inc. [AMD/ATI] ('1002') has 886 products
National Instruments ('1093') has 601 products
Chelsio Communications Inc ('1425') has 525 products
```

The following is an example of querying the PCI product and subsystem
information to find the products which have the most number of subsystems that
have a different vendor than the top-level product. In other words, the two
products which have been re-sold or re-manufactured with the most number of
different companies.

```go
package main

import (
    "fmt"
    "sort"

    "github.com/jaypipes/pcidb"
)

type ByCountSeparateSubvendors []*pcidb.Product

func (v ByCountSeparateSubvendors) Len() int {
    return len(v)
}

func (v ByCountSeparateSubvendors) Swap(i, j int) {
    v[i], v[j] = v[j], v[i]
}

func (v ByCountSeparateSubvendors) Less(i, j int) bool {
    iVendor := v[i].VendorID
    iSetSubvendors := make(map[string]bool, 0)
    iNumDiffSubvendors := 0
    jVendor := v[j].VendorID
    jSetSubvendors := make(map[string]bool, 0)
    jNumDiffSubvendors := 0

    for _, sub := range v[i].Subsystems {
        if sub.VendorID != iVendor {
            iSetSubvendors[sub.VendorID] = true
        }
    }
    iNumDiffSubvendors = len(iSetSubvendors)

    for _, sub := range v[j].Subsystems {
        if sub.VendorID != jVendor {
            jSetSubvendors[sub.VendorID] = true
        }
    }
    jNumDiffSubvendors = len(jSetSubvendors)

    return iNumDiffSubvendors > jNumDiffSubvendors
}

func main() {
    pci, err := pcidb.New()
    if err != nil {
        fmt.Printf("Error getting PCI info: %v", err)
    }

    products := make([]*pcidb.Product, len(pci.Products))
    x := 0
    for _, product := range pci.Products {
        products[x] = product
        x++
    }

    sort.Sort(ByCountSeparateSubvendors(products))

    fmt.Println("Top 2 products by # different subvendors")
    fmt.Println("====================================================")
    for _, product := range products[0:2] {
        vendorID := product.VendorID
        vendor := pci.Vendors[vendorID]
        setSubvendors := make(map[string]bool, 0)

        for _, sub := range product.Subsystems {
            if sub.VendorID != vendorID {
                setSubvendors[sub.VendorID] = true
            }
        }
        fmt.Printf("%v ('%v') from %v\n", product.Name, product.ID, vendor.Name)
        fmt.Printf(" -> %d subsystems under the following different vendors:\n", len(setSubvendors))
        for subvendorID, _ := range setSubvendors {
            subvendor, exists := pci.Vendors[subvendorID]
            subvendorName := "Unknown subvendor"
            if exists {
                subvendorName = subvendor.Name
            }
            fmt.Printf("      - %v ('%v')\n", subvendorName, subvendorID)
        }
    }
}
```

which yields (on my local workstation as of July 7th, 2018):

```
Top 2 products by # different subvendors
====================================================
RTL-8100/8101L/8139 PCI Fast Ethernet Adapter ('8139') from Realtek Semiconductor Co., Ltd.
 -> 34 subsystems under the following different vendors:
      - OVISLINK Corp. ('149c')
      - EPoX Computer Co., Ltd. ('1695')
      - Red Hat, Inc ('1af4')
      - Mitac ('1071')
      - Netgear ('1385')
      - Micro-Star International Co., Ltd. [MSI] ('1462')
      - Hangzhou Silan Microelectronics Co., Ltd. ('1904')
      - Compex ('11f6')
      - Edimax Computer Co. ('1432')
      - KYE Systems Corporation ('1489')
      - ZyXEL Communications Corporation ('187e')
      - Acer Incorporated [ALI] ('1025')
      - Matsushita Electric Industrial Co., Ltd. ('10f7')
      - Ruby Tech Corp. ('146c')
      - Belkin ('1799')
      - Allied Telesis ('1259')
      - Unex Technology Corp. ('1429')
      - CIS Technology Inc ('1436')
      - D-Link System Inc ('1186')
      - Ambicom Inc ('1395')
      - AOPEN Inc. ('a0a0')
      - TTTech Computertechnik AG (Wrong ID) ('0357')
      - Gigabyte Technology Co., Ltd ('1458')
      - Packard Bell B.V. ('1631')
      - Billionton Systems Inc ('14cb')
      - Kingston Technologies ('2646')
      - Accton Technology Corporation ('1113')
      - Samsung Electronics Co Ltd ('144d')
      - Biostar Microtech Int'l Corp ('1565')
      - U.S. Robotics ('16ec')
      - KTI ('8e2e')
      - Hewlett-Packard Company ('103c')
      - ASUSTeK Computer Inc. ('1043')
      - Surecom Technology ('10bd')
Bt878 Video Capture ('036e') from Brooktree Corporation
 -> 30 subsystems under the following different vendors:
      - iTuner ('aa00')
      - Nebula Electronics Ltd. ('0071')
      - DViCO Corporation ('18ac')
      - iTuner ('aa05')
      - iTuner ('aa0d')
      - LeadTek Research Inc. ('107d')
      - Avermedia Technologies Inc ('1461')
      - Chaintech Computer Co. Ltd ('270f')
      - iTuner ('aa07')
      - iTuner ('aa0a')
      - Microtune, Inc. ('1851')
      - iTuner ('aa01')
      - iTuner ('aa04')
      - iTuner ('aa06')
      - iTuner ('aa0f')
      - iTuner ('aa02')
      - iTuner ('aa0b')
      - Pinnacle Systems, Inc. (Wrong ID) ('bd11')
      - Rockwell International ('127a')
      - Askey Computer Corp. ('144f')
      - Twinhan Technology Co. Ltd ('1822')
      - Anritsu Corp. ('1852')
      - iTuner ('aa08')
      - Hauppauge computer works Inc. ('0070')
      - Pinnacle Systems Inc. ('11bd')
      - Conexant Systems, Inc. ('14f1')
      - iTuner ('aa09')
      - iTuner ('aa03')
      - iTuner ('aa0c')
      - iTuner ('aa0e')
```

## Developers

Contributions to `pcidb` are welcomed! Fork the repo on GitHub and submit a pull
request with your proposed changes. Or, feel free to log an issue for a feature
request or bug report.

### Running tests

You can run unit tests easily using the `make test` command, like so:


```
[jaypipes@uberbox pcidb]$ make test
go test github.com/jaypipes/pcidb
ok      github.com/jaypipes/pcidb    0.045s
```
