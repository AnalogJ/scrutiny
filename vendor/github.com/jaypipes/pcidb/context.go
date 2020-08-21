package pcidb

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

// Concrete merged set of configuration switches that get passed to pcidb
// internal functions
type context struct {
	chroot              string
	cacheOnly           bool
	cachePath           string
	disableNetworkFetch bool
	searchPaths         []string
}

func contextFromOptions(merged *WithOption) *context {
	ctx := &context{
		chroot:              *merged.Chroot,
		cacheOnly:           *merged.CacheOnly,
		cachePath:           getCachePath(),
		disableNetworkFetch: *merged.DisableNetworkFetch,
		searchPaths:         make([]string, 0),
	}
	ctx.setSearchPaths()
	return ctx
}

func getCachePath() string {
	hdir, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed getting homedir.Dir(): %v", err)
		return ""
	}
	fp, err := homedir.Expand(filepath.Join(hdir, ".cache", "pci.ids"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed expanding local cache path: %v", err)
		return ""
	}
	return fp
}

// Depending on the operating system, sets the context's searchPaths to a set
// of local filepaths to search for a pci.ids database file
func (ctx *context) setSearchPaths() {
	// A set of filepaths we will first try to search for the pci-ids DB file
	// on the local machine. If we fail to find one, we'll try pulling the
	// latest pci-ids file from the network
	ctx.searchPaths = append(ctx.searchPaths, ctx.cachePath)
	if ctx.cacheOnly {
		return
	}

	rootPath := ctx.chroot

	if runtime.GOOS != "windows" {
		ctx.searchPaths = append(
			ctx.searchPaths,
			filepath.Join(rootPath, "usr", "share", "hwdata", "pci.ids"),
		)
		ctx.searchPaths = append(
			ctx.searchPaths,
			filepath.Join(rootPath, "usr", "share", "misc", "pci.ids"),
		)
	}
}
