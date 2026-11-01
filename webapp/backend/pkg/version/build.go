package version

// BuildVersion is injected via -ldflags from git.version.sh: the release tag
// (e.g. v0.9.5) or branch#commit, matching what the web GUI shows. It lives
// apart from version.go, which the release tooling rewrites.
var BuildVersion string

// String returns the version the web API reports and the collector checks it against.
func String() string {
	if BuildVersion != "" {
		return BuildVersion
	}
	return "v" + VERSION
}
