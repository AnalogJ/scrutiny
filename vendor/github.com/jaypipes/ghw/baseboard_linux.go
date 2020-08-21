// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

func (ctx *context) baseboardFillInfo(info *BaseboardInfo) error {
	info.AssetTag = ctx.dmiItem("board_asset_tag")
	info.SerialNumber = ctx.dmiItem("board_serial")
	info.Vendor = ctx.dmiItem("board_vendor")
	info.Version = ctx.dmiItem("board_version")

	return nil
}
