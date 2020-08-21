// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	"howett.net/plist"
)

type diskOrPartitionPlistNode struct {
	Content          string
	DeviceIdentifier string
	DiskUUID         string
	VolumeName       string
	VolumeUUID       string
	Size             int64
	MountPoint       string
	Partitions       []diskOrPartitionPlistNode
	APFSVolumes      []diskOrPartitionPlistNode
}

type diskUtilListPlist struct {
	AllDisks              []string
	AllDisksAndPartitions []diskOrPartitionPlistNode
	VolumesFromDisks      []string
	WholeDisks            []string
}

type diskUtilInfoPlist struct {
	AESHardware                                 bool   // true
	Bootable                                    bool   // true
	BooterDeviceIdentifier                      string // disk1s2
	BusProtocol                                 string // PCI-Express
	CanBeMadeBootable                           bool   // false
	CanBeMadeBootableRequiresDestroy            bool   // false
	Content                                     string // some-uuid-foo-bar
	DeviceBlockSize                             int64  // 4096
	DeviceIdentifier                            string // disk1s1
	DeviceNode                                  string // /dev/disk1s1
	DeviceTreePath                              string // IODeviceTree:/PCI0@0/RP17@1B/ANS2@0/AppleANS2Controller
	DiskUUID                                    string // some-uuid-foo-bar
	Ejectable                                   bool   // false
	EjectableMediaAutomaticUnderSoftwareControl bool   // false
	EjectableOnly                               bool   // false
	FilesystemName                              string // APFS
	FilesystemType                              string // apfs
	FilesystemUserVisibleName                   string // APFS
	FreeSpace                                   int64  // 343975677952
	GlobalPermissionsEnabled                    bool   // true
	IOKitSize                                   int64  // 499963174912
	IORegistryEntryName                         string // Macintosh HD
	Internal                                    bool   // true
	MediaName                                   string //
	MediaType                                   string // Generic
	MountPoint                                  string // /
	ParentWholeDisk                             string // disk1
	PartitionMapPartition                       bool   // false
	RAIDMaster                                  bool   // false
	RAIDSlice                                   bool   // false
	RecoveryDeviceIdentifier                    string // disk1s3
	Removable                                   bool   // false
	RemovableMedia                              bool   // false
	RemovableMediaOrExternalDevice              bool   // false
	SMARTStatus                                 string // Verified
	Size                                        int64  // 499963174912
	SolidState                                  bool   // true
	SupportsGlobalPermissionsDisable            bool   // true
	SystemImage                                 bool   // false
	TotalSize                                   int64  // 499963174912
	VolumeAllocationBlockSize                   int64  // 4096
	VolumeName                                  string // Macintosh HD
	VolumeSize                                  int64  // 499963174912
	VolumeUUID                                  string // some-uuid-foo-bar
	WholeDisk                                   bool   // false
	Writable                                    bool   // true
	WritableMedia                               bool   // true
	WritableVolume                              bool   // true
	// also has a SMARTDeviceSpecificKeysMayVaryNotGuaranteed dict with various info
	// NOTE: VolumeUUID sometimes == DiskUUID, but not always. So far Content is always a different UUID.
}

type ioregPlist struct {
	// there's a lot more than just this...
	ModelNumber  string `plist:"Model Number"`
	SerialNumber string `plist:"Serial Number"`
	VendorName   string `plist:"Vendor Name"`
}

func (ctx *context) getDiskUtilListPlist() (*diskUtilListPlist, error) {
	out, err := exec.Command("diskutil", "list", "-plist").Output()
	if err != nil {
		return nil, errors.Wrap(err, "diskutil list failed")
	}

	var data diskUtilListPlist
	if _, err := plist.Unmarshal(out, &data); err != nil {
		return nil, errors.Wrap(err, "diskutil list plist unmarshal failed")
	}

	return &data, nil
}

func (ctx *context) getDiskUtilInfoPlist(device string) (*diskUtilInfoPlist, error) {
	out, err := exec.Command("diskutil", "info", "-plist", device).Output()
	if err != nil {
		return nil, errors.Wrapf(err, "diskutil info for %q failed", device)
	}

	var data diskUtilInfoPlist
	if _, err := plist.Unmarshal(out, &data); err != nil {
		return nil, errors.Wrapf(err, "diskutil info plist unmarshal for %q failed", device)
	}

	return &data, nil
}

func (ctx *context) getIoregPlist(ioDeviceTreePath string) (*ioregPlist, error) {
	name := path.Base(ioDeviceTreePath)

	args := []string{
		"ioreg",
		"-a",      // use XML output
		"-d", "1", // limit device tree output depth to root node
		"-r",       // root device tree at matched node
		"-n", name, // match by name
	}
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return nil, errors.Wrapf(err, "ioreg query for %q failed", ioDeviceTreePath)
	}
	if out == nil || len(out) == 0 {
		return nil, nil
	}

	var data []ioregPlist
	if _, err := plist.Unmarshal(out, &data); err != nil {
		return nil, errors.Wrapf(err, "ioreg unmarshal for %q failed", ioDeviceTreePath)
	}
	if len(data) != 1 {
		err := errors.Errorf("ioreg unmarshal resulted in %d I/O device tree nodes (expected 1)", len(data))
		return nil, err
	}

	return &data[0], nil
}

func (ctx *context) makePartition(disk, s diskOrPartitionPlistNode, isAPFS bool) (*Partition, error) {
	if s.Size < 0 {
		return nil, errors.Errorf("invalid size %q of partition %q", s.Size, s.DeviceIdentifier)
	}

	var partType string
	if isAPFS {
		partType = "APFS Volume"
	} else {
		partType = s.Content
	}

	info, err := ctx.getDiskUtilInfoPlist(s.DeviceIdentifier)
	if err != nil {
		return nil, err
	}

	return &Partition{
		Disk:       nil, // filled in later
		Name:       s.DeviceIdentifier,
		Label:      s.VolumeName,
		MountPoint: s.MountPoint,
		SizeBytes:  uint64(s.Size),
		Type:       partType,
		IsReadOnly: !info.WritableVolume,
	}, nil
}

// driveTypeFromPlist looks at the supplied property list struct and attempts to
// determine the disk type
func (ctx *context) driveTypeFromPlist(
	infoPlist *diskUtilInfoPlist,
) DriveType {
	dt := DRIVE_TYPE_HDD
	if infoPlist.SolidState {
		dt = DRIVE_TYPE_SSD
	}
	// TODO(jaypipes): Figure out how to determine floppy and/or CD/optical
	// drive type on Mac
	return dt
}

// storageControllerFromPlist looks at the supplied property list struct and
// attempts to determine the storage controller in use for the device
func (ctx *context) storageControllerFromPlist(
	infoPlist *diskUtilInfoPlist,
) StorageController {
	sc := STORAGE_CONTROLLER_SCSI
	if strings.HasSuffix(infoPlist.DeviceTreePath, "IONVMeController") {
		sc = STORAGE_CONTROLLER_NVME
	}
	// TODO(jaypipes): I don't know if Mac even supports IDE controllers and
	// the "virtio" controller is libvirt-specific
	return sc
}

// busTypeFromPlist looks at the supplied property list struct and attempts to
// determine the bus type in use for the device
func (ctx *context) busTypeFromPlist(
	infoPlist *diskUtilInfoPlist,
) BusType {
	// TODO(jaypipes): Find out if Macs support any bus other than
	// PCIe... it doesn't seem like they do
	return BUS_TYPE_PCI
}

func (ctx *context) blockFillInfo(info *BlockInfo) error {
	listPlist, err := ctx.getDiskUtilListPlist()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}

	info.TotalPhysicalBytes = 0
	info.Disks = make([]*Disk, 0, len(listPlist.AllDisksAndPartitions))
	info.Partitions = []*Partition{}

	for _, disk := range listPlist.AllDisksAndPartitions {
		if disk.Size < 0 {
			return errors.Errorf("invalid size %q of disk %q", disk.Size, disk.DeviceIdentifier)
		}

		infoPlist, err := ctx.getDiskUtilInfoPlist(disk.DeviceIdentifier)
		if err != nil {
			return err
		}
		if infoPlist.DeviceBlockSize < 0 {
			return errors.Errorf("invalid block size %q of disk %q", infoPlist.DeviceBlockSize, disk.DeviceIdentifier)
		}

		busPath := strings.TrimPrefix(infoPlist.DeviceTreePath, "IODeviceTree:")

		ioregPlist, err := ctx.getIoregPlist(infoPlist.DeviceTreePath)
		if err != nil {
			return err
		}
		if ioregPlist == nil {
			continue
		}

		// The NUMA node & WWN don't seem to be reported by any tools available by default in macOS.
		diskReport := &Disk{
			Name:                   disk.DeviceIdentifier,
			SizeBytes:              uint64(disk.Size),
			PhysicalBlockSizeBytes: uint64(infoPlist.DeviceBlockSize),
			DriveType:              ctx.driveTypeFromPlist(infoPlist),
			IsRemovable:            infoPlist.Removable,
			StorageController:      ctx.storageControllerFromPlist(infoPlist),
			BusType:                ctx.busTypeFromPlist(infoPlist),
			BusPath:                busPath,
			NUMANodeID:             -1,
			Vendor:                 ioregPlist.VendorName,
			Model:                  ioregPlist.ModelNumber,
			SerialNumber:           ioregPlist.SerialNumber,
			WWN:                    "",
			Partitions:             make([]*Partition, 0, len(disk.Partitions)+len(disk.APFSVolumes)),
		}

		for _, partition := range disk.Partitions {
			part, err := ctx.makePartition(disk, partition, false)
			if err != nil {
				return err
			}
			part.Disk = diskReport
			diskReport.Partitions = append(diskReport.Partitions, part)
		}
		for _, volume := range disk.APFSVolumes {
			part, err := ctx.makePartition(disk, volume, true)
			if err != nil {
				return err
			}
			part.Disk = diskReport
			diskReport.Partitions = append(diskReport.Partitions, part)
		}

		info.TotalPhysicalBytes += uint64(disk.Size)
		info.Disks = append(info.Disks, diskReport)
		info.Partitions = append(info.Partitions, diskReport.Partitions...)
	}

	return nil
}
