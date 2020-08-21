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

func (ctx *context) netFillInfo(info *NetworkInfo) error {
	return errors.New("netFillInfo not implemented on " + runtime.GOOS)
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
	return nil
}
