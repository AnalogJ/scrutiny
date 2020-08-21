// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

func (ctx *context) biosFillInfo(info *BIOSInfo) error {
	info.Vendor = ctx.dmiItem("bios_vendor")
	info.Version = ctx.dmiItem("bios_version")
	info.Date = ctx.dmiItem("bios_date")

	return nil
}
