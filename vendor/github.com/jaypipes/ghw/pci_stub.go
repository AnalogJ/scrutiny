// +build !linux
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"runtime"

	"github.com/pkg/errors"
)

func (ctx *context) pciFillInfo(info *PCIInfo) error {
	return errors.New("pciFillInfo not implemented on " + runtime.GOOS)
}

// GetDevice returns a pointer to a PCIDevice struct that describes the PCI
// device at the requested address. If no such device could be found, returns
// nil
func (info *PCIInfo) GetDevice(address string) *PCIDevice {
	return nil
}

// ListDevices returns a list of pointers to PCIDevice structs present on the
// host system
func (info *PCIInfo) ListDevices() []*PCIDevice {
	return nil
}
