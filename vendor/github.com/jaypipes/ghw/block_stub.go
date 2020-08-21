// +build !linux,!darwin,!windows
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"runtime"

	"github.com/pkg/errors"
)

func (ctx *context) blockFillInfo(info *BlockInfo) error {
	return errors.New("blockFillInfo not implemented on " + runtime.GOOS)
}

// DiskPhysicalBlockSizeBytes has been deprecated in 0.2. Please use the
// Disk.PhysicalBlockSizeBytes attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskPhysicalBlockSizeBytes(disk string) uint64 {
	msg := `
The DiskPhysicalBlockSizeBytes() function has been DEPRECATED and will be
removed in the 1.0 release of ghw. Please use the Disk.PhysicalBlockSizeBytes
attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskPhysicalBlockSizeBytes(disk)
}

func (ctx *context) diskPhysicalBlockSizeBytes(disk string) uint64 {
	return 0
}

// DiskSizeBytes has been deprecated in 0.2. Please use the Disk.SizeBytes
// attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskSizeBytes(disk string) uint64 {
	msg := `
The DiskSizeBytes() function has been DEPRECATED and will be
removed in the 1.0 release of ghw. Please use the Disk.SizeBytes attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskSizeBytes(disk)
}

func (ctx *context) diskSizeBytes(disk string) uint64 {
	return 0
}

// DiskNUMANodeID has been deprecated in 0.2. Please use the Disk.NUMANodeID
// attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskNUMANodeID(disk string) int {
	msg := `
The DiskNUMANodeID() function has been DEPRECATED and will be
removed in the 1.0 release of ghw. Please use the Disk.NUMANodeID attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskNUMANodeID(disk)
}

func (ctx *context) diskNUMANodeID(disk string) int {
	return -1
}

// DiskVendor has been deprecated in 0.2. Please use the Disk.Vendor attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskVendor(disk string) string {
	msg := `
The DiskVendor() function has been DEPRECATED and will be
removed in the 1.0 release of ghw. Please use the Disk.Vendor attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskVendor(disk)
}

func (ctx *context) diskVendor(disk string) string {
	return UNKNOWN
}

// DiskModel has been deprecated in 0.2. Please use the Disk.Model attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskModel(disk string) string {
	msg := `
The DiskModel() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the Disk.Model attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskModel(disk)
}

func (ctx *context) diskModel(disk string) string {
	return UNKNOWN
}

// DiskSerialNumber has been deprecated in 0.2. Please use the Disk.SerialNumber attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskSerialNumber(disk string) string {
	msg := `
The DiskSerialNumber() function has been DEPRECATED and will be removed in the
1.0 release of ghw. Please use the Disk.SerialNumber attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskSerialNumber(disk)
}

func (ctx *context) diskSerialNumber(disk string) string {
	return UNKNOWN
}

// DiskBusPath has been deprecated in 0.2. Please use the Disk.BusPath attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskBusPath(disk string) string {
	msg := `
The DiskBusPath() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the Disk.BusPath attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskBusPath(disk)
}

func (ctx *context) diskBusPath(disk string) string {
	return UNKNOWN
}

// DiskWWN has been deprecated in 0.2. Please use the Disk.WWN attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskWWN(disk string) string {
	msg := `
The DiskWWN() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the Disk.WWN attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskWWN(disk)
}

func (ctx *context) diskWWN(disk string) string {
	return UNKNOWN
}

// DiskPartitions has been deprecated in 0.2. Please use the Disk.Partitions attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskPartitions(disk string) []*Partition {
	msg := `
The DiskPartitions() function has been DEPRECATED and will be removed in the
1.0 release of ghw. Please use the Disk.Partitions attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskPartitions(disk)
}

func (ctx *context) diskPartitions(disk string) []*Partition {
	return nil
}

// Disks has been deprecated in 0.2. Please use the BlockInfo.Disks attribute.
// TODO(jaypipes): Remove in 1.0.
func Disks() []*Disk {
	msg := `
The Disks() function has been DEPRECATED and will be removed in the
1.0 release of ghw. Please use the BlockInfo.Disks attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.disks()
}

func (ctx *context) disks() []*Disk {
	return nil
}

// PartitionSizeBytes has been deprecated in 0.2. Please use the
// Partition.SizeBytes attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionSizeBytes(part string) uint64 {
	msg := `
The PartitionSizeBytes() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.SizeBytes attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.partitionSizeBytes(part)
}

func (ctx *context) partitionSizeBytes(part string) uint64 {
	return 0
}

// PartitionInfo has been deprecated in 0.2. Please use the Partition struct.
// TODO(jaypipes): Remove in 1.0.
func PartitionInfo(part string) (string, string, bool) {
	msg := `
The PartitionInfo() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition struct.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.partitionInfo(part)
}

// Given a full or short partition name, returns the mount point, the type of
// the partition and whether it's readonly
func (ctx *context) partitionInfo(part string) (string, string, bool) {
	// full name, short name, read-only
	return "", "", true
}

// PartitionMountPoint has been deprecated in 0.2. Please use the
// Partition.MountPoint attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionMountPoint(part string) string {
	msg := `
The PartitionMountPoint() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.MountPoint attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.partitionMountPoint(part)
}

func (ctx *context) partitionMountPoint(part string) string {
	mp, _, _ := ctx.partitionInfo(part)
	return mp
}

// PartitionType has been deprecated in 0.2. Please use the
// Partition.Type attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionType(part string) string {
	msg := `
The PartitionType() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.Type attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.partitionType(part)
}

func (ctx *context) partitionType(part string) string {
	_, pt, _ := ctx.partitionInfo(part)
	return pt
}

// PartitionIsReadOnly has been deprecated in 0.2. Please use the
// Partition.IsReadOnly attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionIsReadOnly(part string) bool {
	msg := `
The PartitionIsReadOnly() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.IsReadOnly attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.partitionIsReadOnly(part)
}

func (ctx *context) partitionIsReadOnly(part string) bool {
	_, _, ro := ctx.partitionInfo(part)
	return ro
}
