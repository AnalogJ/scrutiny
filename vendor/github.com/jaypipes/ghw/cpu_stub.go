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

func (ctx *context) cpuFillInfo(info *CPUInfo) error {
	return errors.New("cpuFillInfo not implemented on " + runtime.GOOS)
}

// Processors has been DEPRECATED in 0.2 and will be REMOVED in 1.0. Please use
// the CPUInfo.Processors attribute.
// TODO(jaypipes): Remove in 1.0
func Processors() []*Processor {
	return nil
}

// TODO: remove
func (ctx *context) coresForNode(nodeID int) ([]*ProcessorCore, error) {
	return nil, errors.New("coresForNode not implemented on " + runtime.GOOS)
}
