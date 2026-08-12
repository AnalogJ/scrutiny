package collector

import "strings"

// SmartctlExitStatus is the bitmask smartctl uses for its process exit code, and which it also
// reports as `smartctl.exit_status` in its JSON output.
// https://github.com/smartmontools/smartmontools/blob/RELEASE_7_5/smartmontools/smartctl.h#L16-L47
type SmartctlExitStatus int

const (
	// SmartctlExitStatusFailCmd - the command line did not parse, or smartctl hit an internal error.
	SmartctlExitStatusFailCmd SmartctlExitStatus = 1 << iota
	// SmartctlExitStatusFailDev - the device could not be opened, did not return an IDENTIFY DEVICE
	// structure, or is in a low power mode and `-n` asked smartctl to exit.
	SmartctlExitStatusFailDev
	// SmartctlExitStatusFailSmart - a SMART command failed, or a SMART data structure had a bad checksum.
	SmartctlExitStatusFailSmart
	// SmartctlExitStatusFailStatus - the SMART status check returned "DISK FAILING".
	SmartctlExitStatusFailStatus
	// SmartctlExitStatusFailAttr - prefail attributes are less than or equal to their threshold.
	SmartctlExitStatusFailAttr
	// SmartctlExitStatusFailAge - attributes were less than or equal to their threshold at some point in the past.
	SmartctlExitStatusFailAge
	// SmartctlExitStatusFailErr - the device error log contains errors.
	SmartctlExitStatusFailErr
	// SmartctlExitStatusFailLog - the device self-test log contains errors.
	SmartctlExitStatusFailLog
)

// SmartctlExitStatusFatal covers the bits that mean smartctl could not produce usable output. The
// remaining bits describe problems with the disk itself, which is exactly the data worth keeping.
const SmartctlExitStatusFatal = SmartctlExitStatusFailCmd | SmartctlExitStatusFailDev | SmartctlExitStatusFailSmart

var smartctlExitStatusDescriptions = []struct {
	flag        SmartctlExitStatus
	description string
}{
	{SmartctlExitStatusFailCmd, "smartctl could not parse the commandline"},
	{SmartctlExitStatusFailDev, "smartctl could not open the device, or the device is in a low power mode"},
	{SmartctlExitStatusFailSmart, "a SMART command failed, or a SMART data structure had a checksum error"},
	{SmartctlExitStatusFailStatus, "smartctl detected a failing disk"},
	{SmartctlExitStatusFailAttr, "smartctl detected a disk in pre-fail"},
	{SmartctlExitStatusFailAge, "smartctl detected a disk that was close to failure in the past"},
	{SmartctlExitStatusFailErr, "smartctl detected an error log with errors"},
	{SmartctlExitStatusFailLog, "smartctl detected a self test log with errors"},
}

func (s SmartctlExitStatus) Has(flags SmartctlExitStatus) bool {
	return s&flags != 0
}

// IsFatal reports whether smartctl failed in a way that leaves its output unusable.
func (s SmartctlExitStatus) IsFatal() bool {
	return s.Has(SmartctlExitStatusFatal)
}

// Descriptions returns a description for every bit that is set.
func (s SmartctlExitStatus) Descriptions() []string {
	descriptions := []string{}
	for _, d := range smartctlExitStatusDescriptions {
		if s.Has(d.flag) {
			descriptions = append(descriptions, d.description)
		}
	}
	return descriptions
}

func (s SmartctlExitStatus) String() string {
	return strings.Join(s.Descriptions(), "; ")
}
