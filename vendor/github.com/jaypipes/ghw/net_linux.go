// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	_WARN_ETHTOOL_NOT_INSTALLED = `ethtool not installed. Cannot grab NIC capabilities`
)

func (ctx *context) netFillInfo(info *NetworkInfo) error {
	info.NICs = ctx.nics()
	return nil
}

// NICS has been deprecated in 0.2. Please use the NetworkInfo.NICs attribute.
// TODO(jaypipes): Remove in 1.0.
func NICs() []*NIC {
	msg := `
The NICs() function has been DEPRECATED and will be removed in the 1.0 release
of ghw. Please use the NetworkInfo.NICs attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.nics()
}

func (ctx *context) nics() []*NIC {
	nics := make([]*NIC, 0)

	files, err := ioutil.ReadDir(ctx.pathSysClassNet())
	if err != nil {
		return nics
	}

	etInstalled := ethtoolInstalled()
	if !etInstalled {
		warn(_WARN_ETHTOOL_NOT_INSTALLED)
	}
	for _, file := range files {
		filename := file.Name()
		// Ignore loopback...
		if filename == "lo" {
			continue
		}

		netPath := filepath.Join(ctx.pathSysClassNet(), filename)
		dest, _ := os.Readlink(netPath)
		isVirtual := false
		if strings.Contains(dest, "devices/virtual/net") {
			isVirtual = true
		}

		nic := &NIC{
			Name:      filename,
			IsVirtual: isVirtual,
		}

		mac := ctx.netDeviceMacAddress(filename)
		nic.MacAddress = mac
		if etInstalled {
			nic.Capabilities = ctx.netDeviceCapabilities(filename)
		} else {
			nic.Capabilities = []*NICCapability{}
		}
		nics = append(nics, nic)
	}
	return nics
}

func (ctx *context) netDeviceMacAddress(dev string) string {
	// Instead of use udevadm, we can get the device's MAC address by examing
	// the /sys/class/net/$DEVICE/address file in sysfs. However, for devices
	// that have addr_assign_type != 0, return None since the MAC address is
	// random.
	aatPath := filepath.Join(ctx.pathSysClassNet(), dev, "addr_assign_type")
	contents, err := ioutil.ReadFile(aatPath)
	if err != nil {
		return ""
	}
	if strings.TrimSpace(string(contents)) != "0" {
		return ""
	}
	addrPath := filepath.Join(ctx.pathSysClassNet(), dev, "address")
	contents, err = ioutil.ReadFile(addrPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(contents))
}

func ethtoolInstalled() bool {
	_, err := exec.LookPath("ethtool")
	return err == nil
}

func (ctx *context) netDeviceCapabilities(dev string) []*NICCapability {
	caps := make([]*NICCapability, 0)
	path, _ := exec.LookPath("ethtool")
	cmd := exec.Command(path, "-k", dev)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		msg := fmt.Sprintf("could not grab NIC capabilities for %s: %s", dev, err)
		warn(msg)
		return caps
	}

	// The out variable will now contain something that looks like the
	// following.
	//
	// Features for enp58s0f1:
	// rx-checksumming: on
	// tx-checksumming: off
	//     tx-checksum-ipv4: off
	//     tx-checksum-ip-generic: off [fixed]
	//     tx-checksum-ipv6: off
	//     tx-checksum-fcoe-crc: off [fixed]
	//     tx-checksum-sctp: off [fixed]
	// scatter-gather: off
	//     tx-scatter-gather: off
	//     tx-scatter-gather-fraglist: off [fixed]
	// tcp-segmentation-offload: off
	//     tx-tcp-segmentation: off
	//     tx-tcp-ecn-segmentation: off [fixed]
	//     tx-tcp-mangleid-segmentation: off
	//     tx-tcp6-segmentation: off
	// < snipped >
	scanner := bufio.NewScanner(&out)
	// Skip the first line...
	scanner.Scan()
	for scanner.Scan() {
		line := strings.TrimPrefix(scanner.Text(), "\t")
		caps = append(caps, netParseEthtoolFeature(line))
	}
	return caps
}

// netParseEthtoolFeature parses a line from the ethtool -k output and returns
// a NICCapability.
//
// The supplied line will look like the following:
//
// tx-checksum-ip-generic: off [fixed]
//
// [fixed] indicates that the feature may not be turned on/off. Note: it makes
// no difference whether a privileged user runs `ethtool -k` when determining
// whether [fixed] appears for a feature.
func netParseEthtoolFeature(line string) *NICCapability {
	parts := strings.Fields(line)
	cap := strings.TrimSuffix(parts[0], ":")
	enabled := parts[1] == "on"
	fixed := len(parts) == 3 && parts[2] == "[fixed]"
	return &NICCapability{
		Name:      cap,
		IsEnabled: enabled,
		CanEnable: !fixed,
	}
}
