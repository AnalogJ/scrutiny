package models

// SmartctlExitStatus is a bitmask representing the exit status of smartctl.
// https://github.com/smartmontools/smartmontools/blob/b1bb7d73c37c16ddddf73906ac9e9e9ad673d481/src/smartctl.h#L16-L47
type SmartctlExitStatus int

// Return codes (bitmask)
const (
	// command line did not parse, or internal error occurred in smartctl
	SmartctlExitStatusFailCmd SmartctlExitStatus = 0x01 << iota

	// device open failed
	// device is in low power mode and -n option requests to exit
	// read device identity (ATA only) failed
	SmartctlExitStatusFailDev

	// smart command failed, or ATA identify device structure missing information
	SmartctlExitStatusFailSmart

	// SMART STATUS returned FAILURE
	SmartctlExitStatusFailStatus

	// Attributes found <= threshold with prefail=1
	SmartctlExitStatusFailAttr

	// SMART STATUS returned GOOD but age attributes failed or prefail
	// attributes have failed in the past
	SmartctlExitStatusFailAge

	// Device had Errors in the error log
	SmartctlExitStatusFailErr

	// Device had Errors in the self-test log
	SmartctlExitStatusFailLog
)

func NewSmartctlExitStatus(exitCode int) SmartctlExitStatus {
	return SmartctlExitStatus(exitCode)
}

func (s SmartctlExitStatus) HasFailCmd() bool {
	return s&SmartctlExitStatusFailCmd != 0
}

func (s SmartctlExitStatus) HasFailDev() bool {
	return s&SmartctlExitStatusFailDev != 0
}

func (s SmartctlExitStatus) HasFailSmart() bool {
	return s&SmartctlExitStatusFailSmart != 0
}

func (s SmartctlExitStatus) HasFailStatus() bool {
	return s&SmartctlExitStatusFailStatus != 0
}

func (s SmartctlExitStatus) HasFailAttr() bool {
	return s&SmartctlExitStatusFailAttr != 0
}

func (s SmartctlExitStatus) HasFailAge() bool {
	return s&SmartctlExitStatusFailAge != 0
}

func (s SmartctlExitStatus) HasFailErr() bool {
	return s&SmartctlExitStatusFailErr != 0
}

func (s SmartctlExitStatus) HasFailLog() bool {
	return s&SmartctlExitStatusFailLog != 0
}
