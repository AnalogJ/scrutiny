package ghw

import "os"

const (
	defaultChroot = "/"
	envKeyChroot  = "GHW_CHROOT"
)

// optDefaultChroot returns the value of the GHW_CHROOT environs variable or
// the default value of "/" if not set
func envOrDefaultChroot() string {
	// Grab options from the environs by default
	if val, exists := os.LookupEnv(envKeyChroot); exists {
		return val
	}
	return defaultChroot
}

// WithOption is used to represent optionally-configured settings. Each field
// is a pointer to some concrete value so that we can tell when something has
// been set or left unset.
type WithOption struct {
	// To facilitate querying of sysfs filesystems that are bind-mounted to a
	// non-default root mountpoint, we allow users to set the GHW_CHROOT environ
	// vairable to an alternate mountpoint. For instance, assume that the user of
	// ghw is a Golang binary being executed from an application container that has
	// certain host filesystems bind-mounted into the container at /host. The user
	// would ensure the GHW_CHROOT environ variable is set to "/host" and ghw will
	// build its paths from that location instead of /
	Chroot *string
}

func WithChroot(dir string) *WithOption {
	return &WithOption{Chroot: &dir}
}

func mergeOptions(opts ...*WithOption) *WithOption {
	merged := &WithOption{}
	for _, opt := range opts {
		if opt.Chroot != nil {
			merged.Chroot = opt.Chroot
		}
	}
	// Set the default value if missing from mergeOpts
	if merged.Chroot == nil {
		chroot := envOrDefaultChroot()
		merged.Chroot = &chroot
	}
	return merged
}

// contextFromEnv returns a pointer to a context struct that has been populated
// from the environs or default options values
func contextFromEnv() *context {
	return &context{
		chroot: envOrDefaultChroot(),
	}
}
