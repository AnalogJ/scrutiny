// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func (ctx *context) cachesForNode(nodeID int) ([]*MemoryCache, error) {
	// The /sys/devices/node/nodeX directory contains a subdirectory called
	// 'cpuX' for each logical processor assigned to the node. Each of those
	// subdirectories containers a 'cache' subdirectory which contains a number
	// of subdirectories beginning with 'index' and ending in the cache's
	// internal 0-based identifier. Those subdirectories contain a number of
	// files, including 'shared_cpu_list', 'size', and 'type' which we use to
	// determine cache characteristics.
	path := filepath.Join(
		ctx.pathSysDevicesSystemNode(),
		fmt.Sprintf("node%d", nodeID),
	)
	caches := make(map[string]*MemoryCache)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "cpu") {
			continue
		}
		if filename == "cpumap" || filename == "cpulist" {
			// There are two files in the node directory that start with 'cpu'
			// but are not subdirectories ('cpulist' and 'cpumap'). Ignore
			// these files.
			continue
		}
		// Grab the logical processor ID by cutting the integer from the
		// /sys/devices/system/node/nodeX/cpuX filename
		cpuPath := filepath.Join(path, filename)
		lpID, _ := strconv.Atoi(filename[3:])

		// Inspect the caches for each logical processor. There will be a
		// /sys/devices/system/node/nodeX/cpuX/cache directory containing a
		// number of directories beginning with the prefix "index" followed by
		// a number. The number indicates the level of the cache, which
		// indicates the "distance" from the processor. Each of these
		// directories contains information about the size of that level of
		// cache and the processors mapped to it.
		cachePath := filepath.Join(cpuPath, "cache")
		if _, err = os.Stat(cachePath); os.IsNotExist(err) {
			continue
		}
		cacheDirFiles, err := ioutil.ReadDir(cachePath)
		if err != nil {
			return nil, err
		}
		for _, cacheDirFile := range cacheDirFiles {
			cacheDirFileName := cacheDirFile.Name()
			if !strings.HasPrefix(cacheDirFileName, "index") {
				continue
			}
			cacheIndex, _ := strconv.Atoi(cacheDirFileName[5:])

			// The cache information is repeated for each node, so here, we
			// just ensure that we only have a one MemoryCache object for each
			// unique combination of level, type and processor map
			level := ctx.memoryCacheLevel(nodeID, lpID, cacheIndex)
			cacheType := ctx.memoryCacheType(nodeID, lpID, cacheIndex)
			sharedCpuMap := ctx.memoryCacheSharedCPUMap(nodeID, lpID, cacheIndex)
			cacheKey := fmt.Sprintf("%d-%d-%s", level, cacheType, sharedCpuMap)

			cache, exists := caches[cacheKey]
			if !exists {
				size := ctx.memoryCacheSize(nodeID, lpID, level)
				cache = &MemoryCache{
					Level:             uint8(level),
					Type:              cacheType,
					SizeBytes:         uint64(size) * uint64(KB),
					LogicalProcessors: make([]uint32, 0),
				}
				caches[cacheKey] = cache
			}
			cache.LogicalProcessors = append(
				cache.LogicalProcessors,
				uint32(lpID),
			)
		}
	}

	cacheVals := make([]*MemoryCache, len(caches))
	x := 0
	for _, c := range caches {
		// ensure the cache's processor set is sorted by logical process ID
		sort.Sort(SortByLogicalProcessorId(c.LogicalProcessors))
		cacheVals[x] = c
		x++
	}

	return cacheVals, nil
}

func (ctx *context) pathNodeCPU(nodeID int, lpID int) string {
	return filepath.Join(
		ctx.pathSysDevicesSystemNode(),
		fmt.Sprintf("node%d", nodeID),
		fmt.Sprintf("cpu%d", lpID),
	)
}

func (ctx *context) pathNodeCPUCache(nodeID int, lpID int) string {
	return filepath.Join(
		ctx.pathNodeCPU(nodeID, lpID),
		"cache",
	)
}

func (ctx *context) pathNodeCPUCacheIndex(nodeID int, lpID int, cacheIndex int) string {
	return filepath.Join(
		ctx.pathNodeCPUCache(nodeID, lpID),
		fmt.Sprintf("index%d", cacheIndex),
	)
}

func (ctx *context) memoryCacheLevel(nodeID int, lpID int, cacheIndex int) int {
	levelPath := filepath.Join(
		ctx.pathNodeCPUCacheIndex(nodeID, lpID, cacheIndex),
		"level",
	)
	levelContents, err := ioutil.ReadFile(levelPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return -1
	}
	// levelContents is now a []byte with the last byte being a newline
	// character. Trim that off and convert the contents to an integer.
	level, err := strconv.Atoi(string(levelContents[:len(levelContents)-1]))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to parse int from %s\n", levelContents)
		return -1
	}
	return level
}

func (ctx *context) memoryCacheSize(nodeID int, lpID int, cacheIndex int) int {
	sizePath := filepath.Join(
		ctx.pathNodeCPUCacheIndex(nodeID, lpID, cacheIndex),
		"size",
	)
	sizeContents, err := ioutil.ReadFile(sizePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return -1
	}
	// size comes as XK\n, so we trim off the K and the newline.
	size, err := strconv.Atoi(string(sizeContents[:len(sizeContents)-2]))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to parse int from %s\n", sizeContents)
		return -1
	}
	return size
}

func (ctx *context) memoryCacheType(nodeID int, lpID int, cacheIndex int) MemoryCacheType {
	typePath := filepath.Join(
		ctx.pathNodeCPUCacheIndex(nodeID, lpID, cacheIndex),
		"type",
	)
	cacheTypeContents, err := ioutil.ReadFile(typePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return MEMORY_CACHE_TYPE_UNIFIED
	}
	switch string(cacheTypeContents[:len(cacheTypeContents)-1]) {
	case "Data":
		return MEMORY_CACHE_TYPE_DATA
	case "Instruction":
		return MEMORY_CACHE_TYPE_INSTRUCTION
	default:
		return MEMORY_CACHE_TYPE_UNIFIED
	}
}

func (ctx *context) memoryCacheSharedCPUMap(nodeID int, lpID int, cacheIndex int) string {
	scpuPath := filepath.Join(
		ctx.pathNodeCPUCacheIndex(nodeID, lpID, cacheIndex),
		"shared_cpu_map",
	)
	sharedCpuMap, err := ioutil.ReadFile(scpuPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return ""
	}
	return string(sharedCpuMap[:len(sharedCpuMap)-1])
}
