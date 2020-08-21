// +build !linux,!windows
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"runtime"

	"github.com/pkg/errors"
)

func (ctx *context) productFillInfo(info *ProductInfo) error {
	return errors.New("productFillInfo not implemented on " + runtime.GOOS)
}
