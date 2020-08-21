// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

func (ctx *context) chassisFillInfo(info *ChassisInfo) error {
	info.AssetTag = ctx.dmiItem("chassis_asset_tag")
	info.SerialNumber = ctx.dmiItem("chassis_serial")
	info.Type = ctx.dmiItem("chassis_type")
	typeDesc, found := chassisTypeDescriptions[info.Type]
	if !found {
		typeDesc = UNKNOWN
	}
	info.TypeDescription = typeDesc
	info.Vendor = ctx.dmiItem("chassis_vendor")
	info.Version = ctx.dmiItem("chassis_version")

	return nil
}
