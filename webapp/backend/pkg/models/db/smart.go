package db

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metadata"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"gorm.io/gorm"
	"time"
)

const SmartWhenFailedFailingNow = "FAILING_NOW"
const SmartWhenFailedInThePast = "IN_THE_PAST"

const SmartStatusPassed = "passed"
const SmartStatusFailed = "failed"

type Smart struct {
	gorm.Model

	DeviceWWN string `json:"device_wwn"`
	Device    Device `json:"-" gorm:"foreignkey:DeviceWWN"` // use DeviceWWN as foreign key

	TestDate    time.Time `json:"date"`
	SmartStatus string    `json:"smart_status"` // SmartStatusPassed or SmartStatusFailed

	//Metrics
	Temp            int64 `json:"temp"`
	PowerOnHours    int64 `json:"power_on_hours"`
	PowerCycleCount int64 `json:"power_cycle_count"`

	AtaAttributes  []SmartAtaAttribute  `json:"ata_attributes" gorm:"foreignkey:SmartId"`
	NvmeAttributes []SmartNvmeAttribute `json:"nvme_attributes" gorm:"foreignkey:SmartId"`
	ScsiAttributes []SmartScsiAttribute `json:"scsi_attributes" gorm:"foreignkey:SmartId"`
}

//Parse Collector SMART data results and create Smart object (and associated SmartAtaAttribute entries)
func (sm *Smart) FromCollectorSmartInfo(wwn string, info collector.SmartInfo) error {
	sm.DeviceWWN = wwn
	sm.TestDate = time.Unix(info.LocalTime.TimeT, 0)

	//smart metrics
	sm.Temp = info.Temperature.Current
	sm.PowerCycleCount = info.PowerCycleCount
	sm.PowerOnHours = info.PowerOnTime.Hours

	// process ATA/NVME/SCSI protocol data
	if info.Device.Protocol == DeviceProtocolAta {
		sm.ProcessAtaSmartInfo(info)
	} else if info.Device.Protocol == DeviceProtocolNvme {
		sm.ProcessNvmeSmartInfo(info)
	} else if info.Device.Protocol == DeviceProtocolScsi {
		sm.ProcessScsiSmartInfo(info)
	}

	if info.SmartStatus.Passed {
		sm.SmartStatus = SmartStatusPassed
	} else {
		sm.SmartStatus = SmartStatusFailed
	}
	return nil
}

//generate SmartAtaAttribute entries from Scrutiny Collector Smart data.
func (sm *Smart) ProcessAtaSmartInfo(info collector.SmartInfo) {
	sm.AtaAttributes = []SmartAtaAttribute{}
	for _, collectorAttr := range info.AtaSmartAttributes.Table {
		attrModel := SmartAtaAttribute{
			AttributeId: collectorAttr.ID,
			Name:        collectorAttr.Name,
			Value:       collectorAttr.Value,
			Worst:       collectorAttr.Worst,
			Threshold:   collectorAttr.Thresh,
			RawValue:    collectorAttr.Raw.Value,
			RawString:   collectorAttr.Raw.String,
			WhenFailed:  collectorAttr.WhenFailed,
		}

		//now that we've parsed the data from the smartctl response, lets match it against our metadata rules and add additional Scrutiny specific data.
		if smartMetadata, ok := metadata.AtaMetadata[collectorAttr.ID]; ok {
			attrModel.Name = smartMetadata.DisplayName
			if smartMetadata.Transform != nil {
				attrModel.TransformedValue = smartMetadata.Transform(attrModel.Value, attrModel.RawValue, attrModel.RawString)
			}
		}
		sm.AtaAttributes = append(sm.AtaAttributes, attrModel)
	}
}

//generate SmartNvmeAttribute entries from Scrutiny Collector Smart data.
func (sm *Smart) ProcessNvmeSmartInfo(info collector.SmartInfo) {
	sm.NvmeAttributes = []SmartNvmeAttribute{
		{AttributeId: "critical_warning", Name: "Critical Warning", Value: info.NvmeSmartHealthInformationLog.CriticalWarning, Threshold: 0},
		{AttributeId: "temperature", Name: "Temperature", Value: info.NvmeSmartHealthInformationLog.Temperature, Threshold: -1},
		{AttributeId: "available_spare", Name: "Available Spare", Value: info.NvmeSmartHealthInformationLog.AvailableSpare, Threshold: info.NvmeSmartHealthInformationLog.AvailableSpareThreshold},
		{AttributeId: "percentage_used", Name: "Percentage Used", Value: info.NvmeSmartHealthInformationLog.PercentageUsed, Threshold: 100},
		{AttributeId: "data_units_read", Name: "Data Units Read", Value: info.NvmeSmartHealthInformationLog.DataUnitsRead, Threshold: -1},
		{AttributeId: "data_units_written", Name: "Data Units Written", Value: info.NvmeSmartHealthInformationLog.DataUnitsWritten, Threshold: -1},
		{AttributeId: "host_reads", Name: "Host Reads", Value: info.NvmeSmartHealthInformationLog.HostReads, Threshold: -1},
		{AttributeId: "host_writes", Name: "Host Writes", Value: info.NvmeSmartHealthInformationLog.HostWrites, Threshold: -1},
		{AttributeId: "controller_busy_time", Name: "Controller Busy Time", Value: info.NvmeSmartHealthInformationLog.ControllerBusyTime, Threshold: -1},
		{AttributeId: "power_cycles", Name: "Power Cycles", Value: info.NvmeSmartHealthInformationLog.PowerCycles, Threshold: -1},
		{AttributeId: "power_on_hours", Name: "Power on Hours", Value: info.NvmeSmartHealthInformationLog.PowerOnHours, Threshold: -1},
		{AttributeId: "unsafe_shutdowns", Name: "Unsafe Shutdowns", Value: info.NvmeSmartHealthInformationLog.UnsafeShutdowns, Threshold: -1},
		{AttributeId: "media_errors", Name: "Media Errors", Value: info.NvmeSmartHealthInformationLog.MediaErrors, Threshold: 0},
		{AttributeId: "num_err_log_entries", Name: "Numb Err Log Entries", Value: info.NvmeSmartHealthInformationLog.NumErrLogEntries, Threshold: 0},
		{AttributeId: "warning_temp_time", Name: "Warning Temp Time", Value: info.NvmeSmartHealthInformationLog.WarningTempTime, Threshold: -1},
		{AttributeId: "critical_comp_time", Name: "Critical CompTime", Value: info.NvmeSmartHealthInformationLog.CriticalCompTime, Threshold: -1},
	}
}

//generate SmartScsiAttribute entries from Scrutiny Collector Smart data.
func (sm *Smart) ProcessScsiSmartInfo(info collector.SmartInfo) {
	sm.ScsiAttributes = []SmartScsiAttribute{
		{AttributeId: "scsi_grown_defect_list", Name: "Grown Defect List", Value: info.ScsiGrownDefectList, Threshold: 0},
		{AttributeId: "read.errors_corrected_by_eccfast", Name: "Read Errors Corrected by ECC Fast", Value: info.ScsiErrorCounterLog.Read.ErrorsCorrectedByEccfast, Threshold: -1},
		{AttributeId: "read.errors_corrected_by_eccdelayed", Name: "Read Errors Corrected by ECC Delayed", Value: info.ScsiErrorCounterLog.Read.ErrorsCorrectedByEccdelayed, Threshold: -1},
		{AttributeId: "read.errors_corrected_by_rereads_rewrites", Name: "Read Errors Corrected by ReReads/ReWrites", Value: info.ScsiErrorCounterLog.Read.ErrorsCorrectedByRereadsRewrites, Threshold: 0},
		{AttributeId: "read.total_errors_corrected", Name: "Read Total Errors Corrected", Value: info.ScsiErrorCounterLog.Read.TotalErrorsCorrected, Threshold: -1},
		{AttributeId: "read.correction_algorithm_invocations", Name: "Read Correction Algorithm Invocations", Value: info.ScsiErrorCounterLog.Read.CorrectionAlgorithmInvocations, Threshold: -1},
		{AttributeId: "read.total_uncorrected_errors", Name: "Read Total Uncorrected Errors", Value: info.ScsiErrorCounterLog.Read.TotalUncorrectedErrors, Threshold: 0},
		{AttributeId: "write.errors_corrected_by_eccfast", Name: "Write Errors Corrected by ECC Fast", Value: info.ScsiErrorCounterLog.Write.ErrorsCorrectedByEccfast, Threshold: -1},
		{AttributeId: "write.errors_corrected_by_eccdelayed", Name: "Write Errors Corrected by ECC Delayed", Value: info.ScsiErrorCounterLog.Write.ErrorsCorrectedByEccdelayed, Threshold: -1},
		{AttributeId: "write.errors_corrected_by_rereads_rewrites", Name: "Write Errors Corrected by ReReads/ReWrites", Value: info.ScsiErrorCounterLog.Write.ErrorsCorrectedByRereadsRewrites, Threshold: 0},
		{AttributeId: "write.total_errors_corrected", Name: "Write Total Errors Corrected", Value: info.ScsiErrorCounterLog.Write.TotalErrorsCorrected, Threshold: -1},
		{AttributeId: "write.correction_algorithm_invocations", Name: "Write Correction Algorithm Invocations", Value: info.ScsiErrorCounterLog.Write.CorrectionAlgorithmInvocations, Threshold: -1},
		{AttributeId: "write.total_uncorrected_errors", Name: "Write Total Uncorrected Errors", Value: info.ScsiErrorCounterLog.Write.TotalUncorrectedErrors, Threshold: 0},
	}
}
