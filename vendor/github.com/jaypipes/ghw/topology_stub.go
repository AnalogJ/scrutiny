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

func (ctx *context) topologyFillInfo(info *TopologyInfo) error {
	return errors.New("topologyFillInfo not implemented on " + runtime.GOOS)
}

// TopologyNodes has been deprecated in 0.2. Please use the TopologyInfo.Nodes
// attribute.
// TODO(jaypipes): Remove in 1.0.
func TopologyNodes() ([]*TopologyNode, error) {
	msg := `
The TopologyNodes() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the TopologyInfo.Nodes attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.topologyNodes()
}

func (ctx *context) topologyNodes() ([]*TopologyNode, error) {
	return nil, errors.New("Don't know how to get topology on " + runtime.GOOS)
}
