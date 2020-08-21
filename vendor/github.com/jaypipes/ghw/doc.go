//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

/*
	package ghw can determine various hardware-related
	information about the host computer:

	* Memory
	* CPU
	* Block storage
	* Topology
	* Network
	* PCI
	* GPU

	Memory

	Information about the host computer's memory can be retrieved using the
	Memory function which returns a pointer to a MemoryInfo struct.

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

	CPU

	The CPU function returns a CPUInfo struct that contains information about
	the CPUs on the host system.

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

	Block storage

	Information about the host computer's local block storage is returned from
	the Block function. This function returns a pointer to a BlockInfo struct.

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

	Topology

	Information about the host computer's architecture (NUMA vs. SMP), the
	host's node layout and processor caches can be retrieved from the Topology
	function. This function returns a pointer to a TopologyInfo struct.

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

	Network

	Information about the host computer's networking hardware is returned from
	the Network function. This function returns a pointer to a NetworkInfo
	struct.

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

	PCI

	ghw contains a PCI database inspection and querying facility that allows
	developers to not only gather information about devices on a local PCI bus
	but also query for information about hardware device classes, vendor and
	product information.

	**NOTE**: Parsing of the PCI-IDS file database is provided by the separate
	http://github.com/jaypipes/pcidb library. You can read that library's
	README for more information about the various structs that are exposed on
	the PCIInfo struct.

	PCIInfo.ListDevices is used to iterate over a host's PCI devices:

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

	The following code snippet shows how to call the PCIInfo.GetDevice method
	and use its returned PCIDevice struct pointer:

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

	GPU

	Information about the host computer's graphics hardware is returned from
	the GPU function. This function returns a pointer to a GPUInfo struct.

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
*/
package ghw
