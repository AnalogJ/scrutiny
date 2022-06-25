package pkg

const DeviceProtocolAta = "ATA"
const DeviceProtocolScsi = "SCSI"
const DeviceProtocolNvme = "NVMe"

const NotifyFilterAttributesAll = "all"
const NotifyFilterAttributesCritical = "critical"

const NotifyLevelFail = "fail"
const NotifyLevelFailScrutiny = "fail_scrutiny"
const NotifyLevelFailSmart = "fail_smart"

type AttributeStatus uint8

const (
	// AttributeStatusPassed binary, 1,2,4,8,16,32,etc
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

type DeviceStatus uint8

const (
	// DeviceStatusPassed binary, 1,2,4,8,16,32,etc
	DeviceStatusPassed         DeviceStatus = 0
	DeviceStatusFailedSmart    DeviceStatus = 1
	DeviceStatusFailedScrutiny DeviceStatus = 2
)

func DeviceStatusSet(b, flag DeviceStatus) DeviceStatus    { return b | flag }
func DeviceStatusClear(b, flag DeviceStatus) DeviceStatus  { return b &^ flag }
func DeviceStatusToggle(b, flag DeviceStatus) DeviceStatus { return b ^ flag }
func DeviceStatusHas(b, flag DeviceStatus) bool            { return b&flag != 0 }
