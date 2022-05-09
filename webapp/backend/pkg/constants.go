package pkg

const DeviceProtocolAta = "ATA"
const DeviceProtocolScsi = "SCSI"
const DeviceProtocolNvme = "NVMe"

const SmartAttributeStatusPassed = 0
const SmartAttributeStatusFailed = 1
const SmartAttributeStatusWarning = 2

const SmartWhenFailedFailingNow = "FAILING_NOW"
const SmartWhenFailedInThePast = "IN_THE_PAST"

//const SmartStatusPassed = "passed"
//const SmartStatusFailed = "failed"

type DeviceStatus int

const (
	DeviceStatusPassed         DeviceStatus = 0
	DeviceStatusFailedSmart    DeviceStatus = iota
	DeviceStatusFailedScrutiny DeviceStatus = iota
)

func Set(b, flag DeviceStatus) DeviceStatus    { return b | flag }
func Clear(b, flag DeviceStatus) DeviceStatus  { return b &^ flag }
func Toggle(b, flag DeviceStatus) DeviceStatus { return b ^ flag }
func Has(b, flag DeviceStatus) bool            { return b&flag != 0 }
