package pkg

const DeviceProtocolAta = "ATA"
const DeviceProtocolScsi = "SCSI"
const DeviceProtocolNvme = "NVMe"

//go:generate stringer -type=AttributeStatus
// AttributeStatus bitwise flag, 1,2,4,8,16,32,etc
type AttributeStatus uint8

const (
	AttributeStatusPassed          AttributeStatus = 0
	AttributeStatusFailedSmart     AttributeStatus = 1
	AttributeStatusWarningScrutiny AttributeStatus = 2
	AttributeStatusFailedScrutiny  AttributeStatus = 4
)

const AttributeWhenFailedFailingNow = "FAILING_NOW"
const AttributeWhenFailedInThePast = "IN_THE_PAST"

func AttributeStatusSet(b, flag AttributeStatus) AttributeStatus    { return b | flag }
func AttributeStatusClear(b, flag AttributeStatus) AttributeStatus  { return b &^ flag }
func AttributeStatusToggle(b, flag AttributeStatus) AttributeStatus { return b ^ flag }
func AttributeStatusHas(b, flag AttributeStatus) bool               { return b&flag != 0 }

//go:generate stringer -type=DeviceStatus
// DeviceStatus bitwise flag, 1,2,4,8,16,32,etc
type DeviceStatus uint8

const (
	DeviceStatusPassed         DeviceStatus = 0
	DeviceStatusFailedSmart    DeviceStatus = 1
	DeviceStatusFailedScrutiny DeviceStatus = 2
)

func DeviceStatusSet(b, flag DeviceStatus) DeviceStatus    { return b | flag }
func DeviceStatusClear(b, flag DeviceStatus) DeviceStatus  { return b &^ flag }
func DeviceStatusToggle(b, flag DeviceStatus) DeviceStatus { return b ^ flag }
func DeviceStatusHas(b, flag DeviceStatus) bool            { return b&flag != 0 }

// Metrics Specific Filtering & Threshold Constants
type MetricsNotifyLevel int64

const (
	MetricsNotifyLevelWarn MetricsNotifyLevel = 1
	MetricsNotifyLevelFail MetricsNotifyLevel = 2
)

type MetricsStatusFilterAttributes int64

const (
	MetricsStatusFilterAttributesAll      MetricsStatusFilterAttributes = 0
	MetricsStatusFilterAttributesCritical MetricsStatusFilterAttributes = 1
)

// MetricsStatusThreshold bitwise flag, 1,2,4,8,16,32,etc
type MetricsStatusThreshold int64

const (
	MetricsStatusThresholdSmart    MetricsStatusThreshold = 1
	MetricsStatusThresholdScrutiny MetricsStatusThreshold = 2

	//shortcut
	MetricsStatusThresholdBoth MetricsStatusThreshold = 3
)

// Deprecated
const NotifyFilterAttributesAll = "all"

// Deprecated
const NotifyFilterAttributesCritical = "critical"

// Deprecated
const NotifyLevelFail = "fail"

// Deprecated
const NotifyLevelFailScrutiny = "fail_scrutiny"

// Deprecated
const NotifyLevelFailSmart = "fail_smart"
