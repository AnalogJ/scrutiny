// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

func (ctx *context) productFillInfo(info *ProductInfo) error {

	info.Family = ctx.dmiItem("product_family")
	info.Name = ctx.dmiItem("product_name")
	info.Vendor = ctx.dmiItem("sys_vendor")
	info.SerialNumber = ctx.dmiItem("product_serial")
	info.UUID = ctx.dmiItem("product_uuid")
	info.SKU = ctx.dmiItem("product_sku")
	info.Version = ctx.dmiItem("product_version")

	return nil
}
