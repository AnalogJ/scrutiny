// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func (ctx *context) topologyFillInfo(info *TopologyInfo) error {
	info.Nodes = ctx.topologyNodes()
	if len(info.Nodes) == 1 {
		info.Architecture = ARCHITECTURE_SMP
	} else {
		info.Architecture = ARCHITECTURE_NUMA
	}
	return nil
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
	return ctx.topologyNodes(), nil
}

func (ctx *context) topologyNodes() []*TopologyNode {
	nodes := make([]*TopologyNode, 0)

	files, err := ioutil.ReadDir(ctx.pathSysDevicesSystemNode())
	if err != nil {
		warn("failed to determine nodes: %s\n", err)
		return nodes
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "node") {
			continue
		}
		node := &TopologyNode{}
		nodeID, err := strconv.Atoi(filename[4:])
		if err != nil {
			warn("failed to determine node ID: %s\n", err)
			return nodes
		}
		node.ID = nodeID
		cores, err := ctx.coresForNode(nodeID)
		if err != nil {
			warn("failed to determine cores for node: %s\n", err)
			return nodes
		}
		node.Cores = cores
		caches, err := ctx.cachesForNode(nodeID)
		if err != nil {
			warn("failed to determine caches for node: %s\n", err)
			return nodes
		}
		node.Caches = caches
		nodes = append(nodes, node)
	}
	return nodes
}
