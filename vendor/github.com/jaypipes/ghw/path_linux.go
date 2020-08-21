// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"path/filepath"
)

func (ctx *context) pathVarLog() string {
	return filepath.Join(ctx.chroot, "var", "log")
}

func (ctx *context) pathProcMeminfo() string {
	return filepath.Join(ctx.chroot, "proc", "meminfo")
}

func (ctx *context) pathSysKernelMMHugepages() string {
	return filepath.Join(ctx.chroot, "sys", "kernel", "mm", "hugepages")
}

func (ctx *context) pathProcCpuinfo() string {
	return filepath.Join(ctx.chroot, "proc", "cpuinfo")
}

func (ctx *context) pathEtcMtab() string {
	return filepath.Join(ctx.chroot, "etc", "mtab")
}

func (ctx *context) pathSysBlock() string {
	return filepath.Join(ctx.chroot, "sys", "block")
}

func (ctx *context) pathSysDevicesSystemNode() string {
	return filepath.Join(ctx.chroot, "sys", "devices", "system", "node")
}

func (ctx *context) pathSysBusPciDevices() string {
	return filepath.Join(ctx.chroot, "sys", "bus", "pci", "devices")
}

func (ctx *context) pathSysClassDrm() string {
	return filepath.Join(ctx.chroot, "sys", "class", "drm")
}

func (ctx *context) pathSysClassDMI() string {
	return filepath.Join(ctx.chroot, "sys", "class", "dmi")
}

func (ctx *context) pathSysClassNet() string {
	return filepath.Join(ctx.chroot, "sys", "class", "net")
}

func (ctx *context) pathRunUdevData() string {
	return filepath.Join(ctx.chroot, "run", "udev", "data")
}
