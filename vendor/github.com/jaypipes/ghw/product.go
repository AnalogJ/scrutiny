//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "fmt"

// ProductInfo defines product information
type ProductInfo struct {
	Family       string `json:"family"`
	Name         string `json:"name"`
	Vendor       string `json:"vendor"`
	SerialNumber string `json:"serial_number"`
	UUID         string `json:"uuid"`
	SKU          string `json:"sku"`
	Version      string `json:"version"`
}

func (i *ProductInfo) String() string {
	familyStr := ""
	if i.Family != "" {
		familyStr = " family=" + i.Family
	}
	nameStr := ""
	if i.Name != "" {
		nameStr = " name=" + i.Name
	}
	vendorStr := ""
	if i.Vendor != "" {
		vendorStr = " vendor=" + i.Vendor
	}
	serialStr := ""
	if i.SerialNumber != "" && i.SerialNumber != UNKNOWN {
		serialStr = " serial=" + i.SerialNumber
	}
	uuidStr := ""
	if i.UUID != "" && i.UUID != UNKNOWN {
		uuidStr = " uuid=" + i.UUID
	}
	skuStr := ""
	if i.SKU != "" {
		skuStr = " sku=" + i.SKU
	}
	versionStr := ""
	if i.Version != "" {
		versionStr = " version=" + i.Version
	}

	res := fmt.Sprintf(
		"product%s%s%s%s%s%s%s",
		familyStr,
		nameStr,
		vendorStr,
		serialStr,
		uuidStr,
		skuStr,
		versionStr,
	)
	return res
}

// Product returns a pointer to a ProductInfo struct containing information
// about the host's product
func Product(opts ...*WithOption) (*ProductInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &ProductInfo{}
	if err := ctx.productFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate product information in a top-level
// "product" YAML/JSON map/object key
type productPrinter struct {
	Info *ProductInfo `json:"product"`
}

// YAMLString returns a string with the product information formatted as YAML
// under a top-level "dmi:" key
func (info *ProductInfo) YAMLString() string {
	return safeYAML(productPrinter{info})
}

// JSONString returns a string with the product information formatted as JSON
// under a top-level "product:" key
func (info *ProductInfo) JSONString(indent bool) string {
	return safeJSON(productPrinter{info}, indent)
}
