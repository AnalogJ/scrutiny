# `ghw` - Golang HardWare discovery/inspection library [![Build Status](https://travis-ci.org/jaypipes/ghw.svg?branch=master)](https://travis-ci.org/jaypipes/ghw)
![ghw mascot](images/ghw-gopher.png)
<br /><br />
`ghw` is a small Golang library providing hardware inspection and discovery
for Linux and Windows. There currently exists partial support for MacOSX.

## Design Principles

* No root privileges needed for discovery

  `ghw` goes the extra mile to be useful without root priveleges. We query for
  host hardware information as directly as possible without relying on shellouts
  to programs like `dmidecode` that require root privileges to execute.

  Elevated privileges are indeed required to query for some information, but
  `ghw` will never error out if blocked from reading that information. Instead,
  `ghw` will print a warning message about the information that could not be
  retrieved. You may disable these warning messages with `GHW_DISABLE_WARNINGS`
  environment variable.

* Well-documented code and plenty of example code

  The code itself should be well-documented with lots of usage
  examples.

* Interfaces should be consistent across modules

  Each module in the library should be structured in a consistent fashion, and
  the structs returned by various library functions should have consistent
  attribute and method names.

## Inspecting != Monitoring

`ghw` is a tool for gathering information about your hardware's **capacity**
and **capabilities**.

It is important to point out that `ghw` does **NOT** report information that is
temporary or variable. It is **NOT** a system monitor nor is it an appropriate
tool for gathering data points for metrics that change over time.  If you are
looking for a system that tracks usage of CPU, memory, network I/O or disk I/O,
there are plenty of great open source tools that do this! Check out the
[Prometheus project](https://prometheus.io/) for a great example.

## Usage

You can use the functions in `ghw` to determine various hardware-related
information about the host computer:

* [Memory](#memory)
* [CPU](#cpu)
* [Block storage](#block-storage)
* [Topology](#topology)
* [Network](#network)
* [PCI](#pci)
* [GPU](#gpu)
* [Chassis](#chassis)
* [BIOS](#bios)
* [Baseboard](#baseboard)
* [Product](#product)
* [YAML and JSON serialization](#serialization)

### Overriding the root mountpoint `ghw` uses

The default root mountpoint that `ghw` uses when looking for information about
the host system is `/`. So, for example, when looking up CPU information on a
Linux system, `ghw.CPU()` will use the path `/proc/cpuinfo`.

If you are calling `ghw` from a system that has an alternate root mountpoint,
you can either set the `GHW_CHROOT` environment variable to that alternate
path, or call the module constructor function with the `ghw.WithChroot()`
modifier.

For example, if you are executing from within an application container that has
bind-mounted the root host filesystem to the mount point `/host`, you would set
`GHW_CHROOT` to `/host` so that `ghw` can find `/proc/cpuinfo` at
`/host/proc/cpuinfo`.

Alternately, you can use the `ghw.WithChroot()` function like so:

```go
cpu, err := ghw.CPU(ghw.WithChroot("/host"))
```

### Disabling warning messages

When `ghw` isn't able to retrieve some information, it may print certain
warning messages to `stderr`. To disable these warnings, simply set the
`GHW_DISABLE_WARNINGS` environs variable:

```
$ ghwc memory
WARNING:
Could not determine total physical bytes of memory. This may
be due to the host being a virtual machine or container with no
/var/log/syslog file, or the current user may not have necessary
privileges to read the syslog. We are falling back to setting the
total physical amount of memory to the total usable amount of memory
memory (24GB physical, 24GB usable)
```

```
$ GHW_DISABLE_WARNINGS=1 ghwc memory
memory (24GB physical, 24GB usable)
```

### Memory

Information about the host computer's memory can be retrieved using the
`ghw.Memory()` function which returns a pointer to a `ghw.MemoryInfo` struct.

The `ghw.MemoryInfo` struct contains three fields:

* `ghw.MemoryInfo.TotalPhysicalBytes` contains the amount of physical memory on
  the host
* `ghw.MemoryInfo.TotalUsableBytes` contains the amount of memory the
  system can actually use. Usable memory accounts for things like the kernel's
  resident memory size and some reserved system bits
* `ghw.MemoryInfo.SupportedPageSizes` is an array of integers representing the
  size, in bytes, of memory pages the system supports
* `ghw.MemoryInfo.Modules` is an array of pointers to `ghw.MemoryModule`
  structs, one for each physical [DIMM](https://en.wikipedia.org/wiki/DIMM).
  Currently, this information is only included on Windows, with Linux support
  [planned](https://github.com/jaypipes/ghw/pull/171#issuecomment-597082409).

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	memory, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}

	fmt.Println(memory.String())
}
```

Example output from my personal workstation:

```
memory (24GB physical, 24GB usable)
```

#### Physical versus Usable Memory

There has been [some](https://github.com/jaypipes/ghw/pull/171)
[confusion](https://github.com/jaypipes/ghw/issues/183) regarding the
difference between the total physical bytes versus total usable bytes of
memory.

Some of this confusion has been due to a misunderstanding of the term "usable".
As mentioned [above](#inspection!=monitoring), `ghw` does inspection of the
system's capacity.

A host computer has two capacities when it comes to RAM. The first capacity is
the amount of RAM that is contained in all memory banks (DIMMs) that are
attached to the motherboard. `ghw.MemoryInfo.TotalPhysicalBytes` refers to this
first capacity.

There is a (usually small) amount of RAM that is consumed by the bootloader
before the operating system is started (booted). Once the bootloader has booted
the operating system, the amount of RAM that may be used by the operating
system and its applications is fixed. `ghw.MemoryInfo.TotalUsableBytes` refers
to this second capacity.

You can determine the amount of RAM that the bootloader used (that is not made
available to the operating system) by subtracting
`ghw.MemoryInfo.TotalUsableBytes` from `ghw.MemoryInfo.TotalPhysicalBytes`:

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	memory, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}

        phys := memory.TotalPhysicalBytes
        usable := memory.TotalUsableBytes

	fmt.Printf("The bootloader consumes %d bytes of RAM\n", phys - usable)
}
```

Example output from my personal workstation booted into a Windows10 operating
system with a Linux GRUB bootloader:

```
The bootloader consumes 3832720 bytes of RAM
```

### CPU

The `ghw.CPU()` function returns a `ghw.CPUInfo` struct that contains
information about the CPUs on the host system.

`ghw.CPUInfo` contains the following fields:

* `ghw.CPUInfo.TotalCores` has the total number of physical cores the host
  system contains
* `ghw.CPUInfo.TotalThreads` has the total number of hardware threads the
  host system contains
* `ghw.CPUInfo.Processors` is an array of `ghw.Processor` structs, one for each
  physical processor package contained in the host

Each `ghw.Processor` struct contains a number of fields:

* `ghw.Processor.ID` is the physical processor `uint32` ID according to the
  system
* `ghw.Processor.NumCores` is the number of physical cores in the processor
  package
* `ghw.Processor.NumThreads` is the number of hardware threads in the processor
  package
* `ghw.Processor.Vendor` is a string containing the vendor name
* `ghw.Processor.Model` is a string containing the vendor's model name
* `ghw.Processor.Capabilities` is an array of strings indicating the features
  the processor has enabled
* `ghw.Processor.Cores` is an array of `ghw.ProcessorCore` structs that are
  packed onto this physical processor

A `ghw.ProcessorCore` has the following fields:

* `ghw.ProcessorCore.ID` is the `uint32` identifier that the host gave this
  core. Note that this does *not* necessarily equate to a zero-based index of
  the core within a physical package. For example, the core IDs for an Intel Core
  i7 are 0, 1, 2, 8, 9, and 10
* `ghw.ProcessorCore.Index` is the zero-based index of the core on the physical
  processor package
* `ghw.ProcessorCore.NumThreads` is the number of hardware threads associated
  with the core
* `ghw.ProcessorCore.LogicalProcessors` is an array of logical processor IDs
  assigned to any processing unit for the core

```go
package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/jaypipes/ghw"
)

func main() {
	cpu, err := ghw.CPU()
	if err != nil {
		fmt.Printf("Error getting CPU info: %v", err)
	}

	fmt.Printf("%v\n", cpu)

	for _, proc := range cpu.Processors {
		fmt.Printf(" %v\n", proc)
		for _, core := range proc.Cores {
			fmt.Printf("  %v\n", core)
		}
		if len(proc.Capabilities) > 0 {
			// pretty-print the (large) block of capability strings into rows
			// of 6 capability strings
			rows := int(math.Ceil(float64(len(proc.Capabilities)) / float64(6)))
			for row := 1; row < rows; row = row + 1 {
				rowStart := (row * 6) - 1
				rowEnd := int(math.Min(float64(rowStart+6), float64(len(proc.Capabilities))))
				rowElems := proc.Capabilities[rowStart:rowEnd]
				capStr := strings.Join(rowElems, " ")
				if row == 1 {
					fmt.Printf("  capabilities: [%s\n", capStr)
				} else if rowEnd < len(proc.Capabilities) {
					fmt.Printf("                 %s\n", capStr)
				} else {
					fmt.Printf("                 %s]\n", capStr)
				}
			}
		}
	}
}
```

Example output from my personal workstation:

```
cpu (1 physical package, 6 cores, 12 hardware threads)
 physical package #0 (6 cores, 12 hardware threads)
  processor core #0 (2 threads), logical processors [0 6]
  processor core #1 (2 threads), logical processors [1 7]
  processor core #2 (2 threads), logical processors [2 8]
  processor core #3 (2 threads), logical processors [3 9]
  processor core #4 (2 threads), logical processors [4 10]
  processor core #5 (2 threads), logical processors [5 11]
  capabilities: [msr pae mce cx8 apic sep
                 mtrr pge mca cmov pat pse36
                 clflush dts acpi mmx fxsr sse
                 sse2 ss ht tm pbe syscall
                 nx pdpe1gb rdtscp lm constant_tsc arch_perfmon
                 pebs bts rep_good nopl xtopology nonstop_tsc
                 cpuid aperfmperf pni pclmulqdq dtes64 monitor
                 ds_cpl vmx est tm2 ssse3 cx16
                 xtpr pdcm pcid sse4_1 sse4_2 popcnt
                 aes lahf_lm pti retpoline tpr_shadow vnmi
                 flexpriority ept vpid dtherm ida arat]
```

### Block storage

Information about the host computer's local block storage is returned from the
`ghw.Block()` function. This function returns a pointer to a `ghw.BlockInfo`
struct.

The `ghw.BlockInfo` struct contains two fields:

* `ghw.BlockInfo.TotalPhysicalBytes` contains the amount of physical block
  storage on the host
* `ghw.BlockInfo.Disks` is an array of pointers to `ghw.Disk` structs, one for
  each disk drive found by the system

Each `ghw.Disk` struct contains the following fields:

* `ghw.Disk.Name` contains a string with the short name of the disk, e.g. "sda"
* `ghw.Disk.SizeBytes` contains the amount of storage the disk provides
* `ghw.Disk.PhysicalBlockSizeBytes` contains the size of the physical blocks
  used on the disk, in bytes
* `ghw.Disk.IsRemovable` contains a boolean indicating if the disk drive is
  removable
* `ghw.Disk.DriveType` is the type of drive. It is of type `ghw.DriveType`
  which has a `ghw.DriveType.String()` method that can be called to return a
  string representation of the bus. This string will be "HDD", "FDD", "ODD",
  or "SSD", which correspond to a hard disk drive (rotational), floppy drive,
  optical (CD/DVD) drive and solid-state drive.
* `ghw.Disk.StorageController` is the type of storage controller/drive. It is
  of type `ghw.StorageController` which has a `ghw.StorageController.String()`
  method that can be called to return a string representation of the bus. This
  string will be "SCSI", "IDE", "virtio", "MMC", or "NVMe"
* `ghw.Disk.NUMANodeID` is the numeric index of the NUMA node this disk is
  local to, or -1
* `ghw.Disk.Vendor` contains a string with the name of the hardware vendor for
  the disk drive
* `ghw.Disk.Model` contains a string with the vendor-assigned disk model name
* `ghw.Disk.SerialNumber` contains a string with the disk's serial number
* `ghw.Disk.WWN` contains a string with the disk's
  [World Wide Name](https://en.wikipedia.org/wiki/World_Wide_Name)
* `ghw.Disk.Partitions` contains an array of pointers to `ghw.Partition`
  structs, one for each partition on the disk

Each `ghw.Partition` struct contains these fields:

* `ghw.Partition.Name` contains a string with the short name of the partition,
  e.g. "sda1"
* `ghw.Partition.SizeBytes` contains the amount of storage the partition
  provides
* `ghw.Partition.MountPoint` contains a string with the partition's mount
  point, or "" if no mount point was discovered
* `ghw.Partition.Type` contains a string indicated the filesystem type for the
  partition, or "" if the system could not determine the type
* `ghw.Partition.IsReadOnly` is a bool indicating the partition is read-only
* `ghw.Partition.Disk` is a pointer to the `ghw.Disk` object associated with
  the partition. This will be `nil` if the `ghw.Partition` struct was returned
  by the `ghw.DiskPartitions()` library function.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
	}

	fmt.Printf("%v\n", block)

	for _, disk := range block.Disks {
		fmt.Printf(" %v\n", disk)
		for _, part := range disk.Partitions {
			fmt.Printf("  %v\n", part)
		}
	}
}
```

Example output from my personal workstation:

```
block storage (1 disk, 2TB physical storage)
 sda HDD (2TB) SCSI [@pci-0000:04:00.0-scsi-0:1:0:0 (node #0)] vendor=LSI model=Logical_Volume serial=600508e000000000f8253aac9a1abd0c WWN=0x600508e000000000f8253aac9a1abd0c
  /dev/sda1 (100MB)
  /dev/sda2 (187GB)
  /dev/sda3 (449MB)
  /dev/sda4 (1KB)
  /dev/sda5 (15GB)
  /dev/sda6 (2TB) [ext4] mounted@/
```

> Note that `ghw` looks in the udev runtime database for some information. If
> you are using `ghw` in a container, remember to bind mount `/dev/disk` and
> `/run` into your container, otherwise `ghw` won't be able to query the udev
> DB or sysfs paths for information.

### Topology

> **NOTE**: Topology support is currently Linux-only. Windows support is
> [planned](https://github.com/jaypipes/ghw/issues/166).

Information about the host computer's architecture (NUMA vs. SMP), the host's
node layout and processor caches can be retrieved from the `ghw.Topology()`
function. This function returns a pointer to a `ghw.TopologyInfo` struct.

The `ghw.TopologyInfo` struct contains two fields:

* `ghw.TopologyInfo.Architecture` contains an enum with the value `ghw.NUMA` or
  `ghw.SMP` depending on what the topology of the system is
* `ghw.TopologyInfo.Nodes` is an array of pointers to `ghw.TopologyNode`
  structs, one for each topology node (typically physical processor package)
  found by the system

Each `ghw.TopologyNode` struct contains the following fields:

* `ghw.TopologyNode.ID` is the system's `uint32` identifier for the node
* `ghw.TopologyNode.Cores` is an array of pointers to `ghw.ProcessorCore` structs that
  are contained in this node
* `ghw.TopologyNode.Caches` is an array of pointers to `ghw.MemoryCache` structs that
  represent the low-level caches associated with processors and cores on the
  system

See above in the [CPU](#cpu) section for information about the
`ghw.ProcessorCore` struct and how to use and query it.

Each `ghw.MemoryCache` struct contains the following fields:

* `ghw.MemoryCache.Type` is an enum that contains one of `ghw.DATA`,
  `ghw.INSTRUCTION` or `ghw.UNIFIED` depending on whether the cache stores CPU
  instructions, program data, or both
* `ghw.MemoryCache.Level` is a positive integer indicating how close the cache
  is to the processor
* `ghw.MemoryCache.SizeBytes` is an integer containing the number of bytes the
  cache can contain
* `ghw.MemoryCache.LogicalProcessors` is an array of integers representing the
  logical processors that use the cache

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	topology, err := ghw.Topology()
	if err != nil {
		fmt.Printf("Error getting topology info: %v", err)
	}

	fmt.Printf("%v\n", topology)

	for _, node := range topology.Nodes {
		fmt.Printf(" %v\n", node)
		for _, cache := range node.Caches {
			fmt.Printf("  %v\n", cache)
		}
	}
}
```

Example output from my personal workstation:

```
topology SMP (1 nodes)
 node #0 (6 cores)
  L1i cache (32 KB) shared with logical processors: 3,9
  L1i cache (32 KB) shared with logical processors: 2,8
  L1i cache (32 KB) shared with logical processors: 11,5
  L1i cache (32 KB) shared with logical processors: 10,4
  L1i cache (32 KB) shared with logical processors: 0,6
  L1i cache (32 KB) shared with logical processors: 1,7
  L1d cache (32 KB) shared with logical processors: 11,5
  L1d cache (32 KB) shared with logical processors: 10,4
  L1d cache (32 KB) shared with logical processors: 3,9
  L1d cache (32 KB) shared with logical processors: 1,7
  L1d cache (32 KB) shared with logical processors: 0,6
  L1d cache (32 KB) shared with logical processors: 2,8
  L2 cache (256 KB) shared with logical processors: 2,8
  L2 cache (256 KB) shared with logical processors: 3,9
  L2 cache (256 KB) shared with logical processors: 0,6
  L2 cache (256 KB) shared with logical processors: 10,4
  L2 cache (256 KB) shared with logical processors: 1,7
  L2 cache (256 KB) shared with logical processors: 11,5
  L3 cache (12288 KB) shared with logical processors: 0,1,10,11,2,3,4,5,6,7,8,9
```

### Network

Information about the host computer's networking hardware is returned from the
`ghw.Network()` function. This function returns a pointer to a
`ghw.NetworkInfo` struct.

The `ghw.NetworkInfo` struct contains one field:

* `ghw.NetworkInfo.NICs` is an array of pointers to `ghw.NIC` structs, one
  for each network interface controller found for the systen

Each `ghw.NIC` struct contains the following fields:

* `ghw.NIC.Name` is the system's identifier for the NIC
* `ghw.NIC.MacAddress` is the MAC address for the NIC, if any
* `ghw.NIC.IsVirtual` is a boolean indicating if the NIC is a virtualized
  device
* `ghw.NIC.Capabilities` is an array of pointers to `ghw.NICCapability` structs
  that can describe the things the NIC supports. These capabilities match the
  returned values from the `ethtool -k <DEVICE>` call on Linux

The `ghw.NICCapability` struct contains the following fields:

* `ghw.NICCapability.Name` is the string name of the capability (e.g.
  "tcp-segmentation-offload")
* `ghw.NICCapability.IsEnabled` is a boolean indicating whether the capability
  is currently enabled/active on the NIC
* `ghw.NICCapability.CanEnable` is a boolean indicating whether the capability
  may be enabled

```go
package main

import (
    "fmt"

    "github.com/jaypipes/ghw"
)

func main() {
    net, err := ghw.Network()
    if err != nil {
        fmt.Printf("Error getting network info: %v", err)
    }

    fmt.Printf("%v\n", net)

    for _, nic := range net.NICs {
        fmt.Printf(" %v\n", nic)

        enabledCaps := make([]int, 0)
        for x, cap := range nic.Capabilities {
            if cap.IsEnabled {
                enabledCaps = append(enabledCaps, x)
            }
        }
        if len(enabledCaps) > 0 {
            fmt.Printf("  enabled capabilities:\n")
            for _, x := range enabledCaps {
                fmt.Printf("   - %s\n", nic.Capabilities[x].Name)
            }
        }
    }
}
```

Example output from my personal laptop:

```
net (3 NICs)
 docker0
  enabled capabilities:
   - tx-checksumming
   - tx-checksum-ip-generic
   - scatter-gather
   - tx-scatter-gather
   - tx-scatter-gather-fraglist
   - tcp-segmentation-offload
   - tx-tcp-segmentation
   - tx-tcp-ecn-segmentation
   - tx-tcp-mangleid-segmentation
   - tx-tcp6-segmentation
   - udp-fragmentation-offload
   - generic-segmentation-offload
   - generic-receive-offload
   - tx-vlan-offload
   - highdma
   - tx-lockless
   - netns-local
   - tx-gso-robust
   - tx-fcoe-segmentation
   - tx-gre-segmentation
   - tx-gre-csum-segmentation
   - tx-ipxip4-segmentation
   - tx-ipxip6-segmentation
   - tx-udp_tnl-segmentation
   - tx-udp_tnl-csum-segmentation
   - tx-gso-partial
   - tx-sctp-segmentation
   - tx-esp-segmentation
   - tx-vlan-stag-hw-insert
 enp58s0f1
  enabled capabilities:
   - rx-checksumming
   - generic-receive-offload
   - rx-vlan-offload
   - tx-vlan-offload
   - highdma
 wlp59s0
  enabled capabilities:
   - scatter-gather
   - tx-scatter-gather
   - generic-segmentation-offload
   - generic-receive-offload
   - highdma
   - netns-local
```

### PCI

`ghw` contains a PCI database inspection and querying facility that allows
developers to not only gather information about devices on a local PCI bus but
also query for information about hardware device classes, vendor and product
information.

**NOTE**: Parsing of the PCI-IDS file database is provided by the separate
[github.com/jaypipes/pcidb library](http://github.com/jaypipes/pcidb). You can
read that library's README for more information about the various structs that
are exposed on the `ghw.PCIInfo` struct.

The `ghw.PCI()` function returns a `ghw.PCIInfo` struct. The `ghw.PCIInfo`
struct contains a number of fields that may be queried for PCI information:

* `ghw.PCIInfo.Classes` is a map, keyed by the PCI class ID (a hex-encoded
  string) of pointers to `pcidb.Class` structs, one for each class of PCI
  device known to `ghw`
* `ghw.PCIInfo.Vendors` is a map, keyed by the PCI vendor ID (a hex-encoded
  string) of pointers to `pcidb.Vendor` structs, one for each PCI vendor
  known to `ghw`
* `ghw.PCIInfo.Products` is a map, keyed by the PCI product ID* (a hex-encoded
  string) of pointers to `pcidb.Product` structs, one for each PCI product
  known to `ghw`

**NOTE**: PCI products are often referred to by their "device ID". We use
the term "product ID" in `ghw` because it more accurately reflects what the
identifier is for: a specific product line produced by the vendor.

#### Listing and accessing host PCI device information

In addition to the above information, the `ghw.PCIInfo` struct has the
following methods:

* `ghw.PCIInfo.ListDevices() []*PCIDevice`
* `ghw.PCIInfo.GetDevice(address string) *PCIDevice`

This methods return either an array of or a single pointer to a `ghw.PCIDevice`
struct, which has the following fields:


* `ghw.PCIDevice.Vendor` is a pointer to a `pcidb.Vendor` struct that
  describes the device's primary vendor. This will always be non-nil.
* `ghw.PCIDevice.Product` is a pointer to a `pcidb.Product` struct that
  describes the device's primary product. This will always be non-nil.
* `ghw.PCIDevice.Subsystem` is a pointer to a `pcidb.Product` struct that
  describes the device's secondary/sub-product. This will always be non-nil.
* `ghw.PCIDevice.Class` is a pointer to a `pcidb.Class` struct that
  describes the device's class. This will always be non-nil.
* `ghw.PCIDevice.Subclass` is a pointer to a `pcidb.Subclass` struct
  that describes the device's subclass. This will always be non-nil.
* `ghw.PCIDevice.ProgrammingInterface` is a pointer to a
  `pcidb.ProgrammingInterface` struct that describes the device subclass'
  programming interface. This will always be non-nil.

The following code snippet shows how to call the `ghw.PCIInfo.ListDevices()`
method and output a simple list of PCI address and vendor/product information:

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}
	fmt.Printf("host PCI devices:\n")
	fmt.Println("====================================================")
	devices := pci.ListDevices()
	if len(devices) == 0 {
		fmt.Printf("error: could not retrieve PCI devices\n")
		return
	}

	for _, device := range devices {
		vendor := device.Vendor
		vendorName := vendor.Name
		if len(vendor.Name) > 20 {
			vendorName = string([]byte(vendorName)[0:17]) + "..."
		}
		product := device.Product
		productName := product.Name
		if len(product.Name) > 40 {
			productName = string([]byte(productName)[0:37]) + "..."
		}
		fmt.Printf("%-12s\t%-20s\t%-40s\n", device.Address, vendorName, productName)
	}
}
```

on my local workstation the output of the above looks like the following:

```
host PCI devices:
====================================================
0000:00:00.0	Intel Corporation   	5520/5500/X58 I/O Hub to ESI Port
0000:00:01.0	Intel Corporation   	5520/5500/X58 I/O Hub PCI Express Roo...
0000:00:02.0	Intel Corporation   	5520/5500/X58 I/O Hub PCI Express Roo...
0000:00:03.0	Intel Corporation   	5520/5500/X58 I/O Hub PCI Express Roo...
0000:00:07.0	Intel Corporation   	5520/5500/X58 I/O Hub PCI Express Roo...
0000:00:10.0	Intel Corporation   	7500/5520/5500/X58 Physical and Link ...
0000:00:10.1	Intel Corporation   	7500/5520/5500/X58 Routing and Protoc...
0000:00:14.0	Intel Corporation   	7500/5520/5500/X58 I/O Hub System Man...
0000:00:14.1	Intel Corporation   	7500/5520/5500/X58 I/O Hub GPIO and S...
0000:00:14.2	Intel Corporation   	7500/5520/5500/X58 I/O Hub Control St...
0000:00:14.3	Intel Corporation   	7500/5520/5500/X58 I/O Hub Throttle R...
0000:00:19.0	Intel Corporation   	82567LF-2 Gigabit Network Connection
0000:00:1a.0	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1a.1	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1a.2	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1a.7	Intel Corporation   	82801JI (ICH10 Family) USB2 EHCI Cont...
0000:00:1b.0	Intel Corporation   	82801JI (ICH10 Family) HD Audio Contr...
0000:00:1c.0	Intel Corporation   	82801JI (ICH10 Family) PCI Express Ro...
0000:00:1c.1	Intel Corporation   	82801JI (ICH10 Family) PCI Express Po...
0000:00:1c.4	Intel Corporation   	82801JI (ICH10 Family) PCI Express Ro...
0000:00:1d.0	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1d.1	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1d.2	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1d.7	Intel Corporation   	82801JI (ICH10 Family) USB2 EHCI Cont...
0000:00:1e.0	Intel Corporation   	82801 PCI Bridge
0000:00:1f.0	Intel Corporation   	82801JIR (ICH10R) LPC Interface Contr...
0000:00:1f.2	Intel Corporation   	82801JI (ICH10 Family) SATA AHCI Cont...
0000:00:1f.3	Intel Corporation   	82801JI (ICH10 Family) SMBus Controller
0000:01:00.0	NEC Corporation     	uPD720200 USB 3.0 Host Controller
0000:02:00.0	Marvell Technolog...	88SE9123 PCIe SATA 6.0 Gb/s controller
0000:02:00.1	Marvell Technolog...	88SE912x IDE Controller
0000:03:00.0	NVIDIA Corporation  	GP107 [GeForce GTX 1050 Ti]
0000:03:00.1	NVIDIA Corporation  	UNKNOWN
0000:04:00.0	LSI Logic / Symbi...	SAS2004 PCI-Express Fusion-MPT SAS-2 ...
0000:06:00.0	Qualcomm Atheros    	AR5418 Wireless Network Adapter [AR50...
0000:08:03.0	LSI Corporation     	FW322/323 [TrueFire] 1394a Controller
0000:3f:00.0	Intel Corporation   	UNKNOWN
0000:3f:00.1	Intel Corporation   	Xeon 5600 Series QuickPath Architectu...
0000:3f:02.0	Intel Corporation   	Xeon 5600 Series QPI Link 0
0000:3f:02.1	Intel Corporation   	Xeon 5600 Series QPI Physical 0
0000:3f:02.2	Intel Corporation   	Xeon 5600 Series Mirror Port Link 0
0000:3f:02.3	Intel Corporation   	Xeon 5600 Series Mirror Port Link 1
0000:3f:03.0	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:03.1	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:03.4	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:04.0	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:04.1	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:04.2	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:04.3	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:05.0	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:05.1	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:05.2	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:05.3	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:06.0	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:06.1	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:06.2	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:06.3	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
```

The following code snippet shows how to call the `ghw.PCIInfo.GetDevice()`
method and use its returned `ghw.PCIDevice` struct pointer:

```go
package main

import (
	"fmt"
	"os"

	"github.com/jaypipes/ghw"
)

func main() {
	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}

	addr := "0000:00:00.0"
	if len(os.Args) == 2 {
		addr = os.Args[1]
	}
	fmt.Printf("PCI device information for %s\n", addr)
	fmt.Println("====================================================")
	deviceInfo := pci.GetDevice(addr)
	if deviceInfo == nil {
		fmt.Printf("could not retrieve PCI device information for %s\n", addr)
		return
	}

	vendor := deviceInfo.Vendor
	fmt.Printf("Vendor: %s [%s]\n", vendor.Name, vendor.ID)
	product := deviceInfo.Product
	fmt.Printf("Product: %s [%s]\n", product.Name, product.ID)
	subsystem := deviceInfo.Subsystem
	subvendor := pci.Vendors[subsystem.VendorID]
	subvendorName := "UNKNOWN"
	if subvendor != nil {
		subvendorName = subvendor.Name
	}
	fmt.Printf("Subsystem: %s [%s] (Subvendor: %s)\n", subsystem.Name, subsystem.ID, subvendorName)
	class := deviceInfo.Class
	fmt.Printf("Class: %s [%s]\n", class.Name, class.ID)
	subclass := deviceInfo.Subclass
	fmt.Printf("Subclass: %s [%s]\n", subclass.Name, subclass.ID)
	progIface := deviceInfo.ProgrammingInterface
	fmt.Printf("Programming Interface: %s [%s]\n", progIface.Name, progIface.ID)
}
```

Here's a sample output from my local workstation:

```
PCI device information for 0000:03:00.0
====================================================
Vendor: NVIDIA Corporation [10de]
Product: GP107 [GeForce GTX 1050 Ti] [1c82]
Subsystem: UNKNOWN [8613] (Subvendor: ASUSTeK Computer Inc.)
Class: Display controller [03]
Subclass: VGA compatible controller [00]
Programming Interface: VGA controller [00]
```

### GPU

Information about the host computer's graphics hardware is returned from the
`ghw.GPU()` function. This function returns a pointer to a `ghw.GPUInfo`
struct.

The `ghw.GPUInfo` struct contains one field:

* `ghw.GPUInfo.GraphicCards` is an array of pointers to `ghw.GraphicsCard`
  structs, one for each graphics card found for the systen

Each `ghw.GraphicsCard` struct contains the following fields:

* `ghw.GraphicsCard.Index` is the system's numeric zero-based index for the
  card on the bus
* `ghw.GraphicsCard.Address` is the PCI address for the graphics card
* `ghw.GraphicsCard.DeviceInfo` is a pointer to a `ghw.PCIDevice` struct
  describing the graphics card. This may be `nil` if no PCI device information
  could be determined for the card.
* `ghw.GraphicsCard.Node` is an pointer to a `ghw.TopologyNode` struct that the
  GPU/graphics card is affined to. On non-NUMA systems, this will always be
  `nil`.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	gpu, err := ghw.GPU()
	if err != nil {
		fmt.Printf("Error getting GPU info: %v", err)
	}

	fmt.Printf("%v\n", gpu)

	for _, card := range gpu.GraphicsCards {
		fmt.Printf(" %v\n", card)
	}
}
```

Example output from my personal workstation:

```
gpu (1 graphics card)
 card #0 @0000:03:00.0 -> class: 'Display controller' vendor: 'NVIDIA Corporation' product: 'GP107 [GeForce GTX 1050 Ti]'
```

**NOTE**: You can [read more](#pci) about the fields of the `ghw.PCIDevice`
struct if you'd like to dig deeper into PCI subsystem and programming interface
information

**NOTE**: You can [read more](#topology) about the fields of the
`ghw.TopologyNode` struct if you'd like to dig deeper into the NUMA/topology
subsystem

### Chassis

The host's chassis information is accessible with the `ghw.Chassis()` function.  This
function returns a pointer to a `ghw.ChassisInfo` struct.

The `ghw.ChassisInfo` struct contains multiple fields:

* `ghw.ChassisInfo.AssetTag` is a string with the chassis asset tag
* `ghw.ChassisInfo.SerialNumber` is a string with the chassis serial number
* `ghw.ChassisInfo.Type` is a string with the chassis type *code*
* `ghw.ChassisInfo.TypeDescription` is a string with a description of the chassis type
* `ghw.ChassisInfo.Vendor` is a string with the chassis vendor
* `ghw.ChassisInfo.Version` is a string with the chassis version

**NOTE**: These fields are often missing for non-server hardware. Don't be
surprised to see empty string or "None" values.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	chassis, err := ghw.Chassis()
	if err != nil {
		fmt.Printf("Error getting chassis info: %v", err)
	}

	fmt.Printf("%v\n", chassis)
}
```

Example output from my personal workstation:

```
chassis type=Desktop vendor=System76 version=thelio-r1
```

**NOTE**: Some of the values such as serial numbers are shown as unknown because
the Linux kernel by default disallows access to those fields if you're not running
as root.  They will be populated if it runs as root or otherwise you may see warnings
like the following:

```
WARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied
```

You can ignore them or use the [Disabling warning messages](#disabling-warning-messages)
feature to quiet things down.

### BIOS

The host's basis input/output system (BIOS) information is accessible with the `ghw.BIOS()` function.  This
function returns a pointer to a `ghw.BIOSInfo` struct.

The `ghw.BIOSInfo` struct contains multiple fields:

* `ghw.BIOSInfo.Vendor` is a string with the BIOS vendor
* `ghw.BIOSInfo.Version` is a string with the BIOS version
* `ghw.BIOSInfo.Date` is a string with the date the BIOS was flashed/created

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	bios, err := ghw.BIOS()
	if err != nil {
		fmt.Printf("Error getting BIOS info: %v", err)
	}

	fmt.Printf("%v\n", bios)
}
```

Example output from my personal workstation:

```
bios vendor=System76 version=F2 Z5 date=11/14/2018
```

### Baseboard

The host's baseboard information is accessible with the `ghw.Baseboard()` function.  This
function returns a pointer to a `ghw.BaseboardInfo` struct.

The `ghw.BaseboardInfo` struct contains multiple fields:

* `ghw.BaseboardInfo.AssetTag` is a string with the baseboard asset tag
* `ghw.BaseboardInfo.SerialNumber` is a string with the baseboard serial number
* `ghw.BaseboardInfo.Vendor` is a string with the baseboard vendor
* `ghw.BaseboardInfo.Version` is a string with the baseboard version

**NOTE**: These fields are often missing for non-server hardware. Don't be
surprised to see empty string or "None" values.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	baseboard, err := ghw.Baseboard()
	if err != nil {
		fmt.Printf("Error getting baseboard info: %v", err)
	}

	fmt.Printf("%v\n", baseboard)
}
```

Example output from my personal workstation:

```
baseboard vendor=System76 version=thelio-r1
```

**NOTE**: Some of the values such as serial numbers are shown as unknown because
the Linux kernel by default disallows access to those fields if you're not running
as root.  They will be populated if it runs as root or otherwise you may see warnings
like the following:

```
WARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied
```

You can ignore them or use the [Disabling warning messages](#disabling-warning-messages)
feature to quiet things down.

### Product

The host's product information is accessible with the `ghw.Product()` function.  This
function returns a pointer to a `ghw.ProductInfo` struct.

The `ghw.ProductInfo` struct contains multiple fields:

* `ghw.ProductInfo.Family` is a string describing the product family
* `ghw.ProductInfo.Name` is a string with the product name
* `ghw.ProductInfo.SerialNumber` is a string with the product serial number
* `ghw.ProductInfo.UUID` is a string with the product UUID
* `ghw.ProductInfo.SKU` is a string with the product stock unit identifier (SKU)
* `ghw.ProductInfo.Vendor` is a string with the product vendor
* `ghw.ProductInfo.Version` is a string with the product version

**NOTE**: These fields are often missing for non-server hardware. Don't be
surprised to see empty string, "Default string" or "None" values.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	product, err := ghw.Product()
	if err != nil {
		fmt.Printf("Error getting product info: %v", err)
	}

	fmt.Printf("%v\n", product)
}
```

Example output from my personal workstation:

```
product family=Default string name=Thelio vendor=System76 sku=Default string version=thelio-r1
```

**NOTE**: Some of the values such as serial numbers are shown as unknown because
the Linux kernel by default disallows access to those fields if you're not running
as root.  They will be populated if it runs as root or otherwise you may see warnings
like the following:

```
WARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied
```

You can ignore them or use the [Disabling warning messages](#disabling-warning-messages)
feature to quiet things down.

## Serialization

All of the `ghw` `XXXInfo` structs -- e.g. `ghw.CPUInfo` -- have two methods
for producing a serialized JSON or YAML string representation of the contained
information:

* `JSONString()` returns a string containing the information serialized into
  JSON. It accepts a single boolean parameter indicating whether to use
  indentation when outputting the string
* `YAMLString()` returns a string containing the information serialized into
  YAML

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	mem, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}

	fmt.Printf("%s", mem.YAMLString())
}
```

the above example code prints the following out on my local workstation:

```
memory:
  supported_page_sizes:
  - 1073741824
  - 2097152
  total_physical_bytes: 25263415296
  total_usable_bytes: 25263415296
```

## Developers

Contributions to `ghw` are welcomed! Fork the repo on GitHub and submit a pull
request with your proposed changes. Or, feel free to log an issue for a feature
request or bug report.

### Running tests

You can run unit tests easily using the `make test` command, like so:


```
[jaypipes@uberbox ghw]$ make test
go test github.com/jaypipes/ghw github.com/jaypipes/ghw/cmd/ghwc
ok      github.com/jaypipes/ghw 0.084s
?       github.com/jaypipes/ghw/cmd/ghwc    [no test files]
```
