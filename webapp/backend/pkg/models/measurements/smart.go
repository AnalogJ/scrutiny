package measurements

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
	"log"
	"strconv"
	"strings"
	"time"
)

type Smart struct {
	Date           time.Time `json:"date"`
	DeviceWWN      string    `json:"device_wwn"` //(tag)
	DeviceProtocol string    `json:"device_protocol"`

	//Metrics (fields)
	Temp            int64 `json:"temp"`
	PowerOnHours    int64 `json:"power_on_hours"`
	PowerCycleCount int64 `json:"power_cycle_count"`

	//Attributes (fields)
	Attributes map[string]SmartAttribute `json:"attrs"`

	//status
	Status pkg.DeviceStatus
}

func (sm *Smart) Flatten() (tags map[string]string, fields map[string]interface{}) {
	tags = map[string]string{
		"device_wwn":      sm.DeviceWWN,
		"device_protocol": sm.DeviceProtocol,
	}

	fields = map[string]interface{}{
		"temp":              sm.Temp,
		"power_on_hours":    sm.PowerOnHours,
		"power_cycle_count": sm.PowerCycleCount,
	}

	for _, attr := range sm.Attributes {
		for attrKey, attrVal := range attr.Flatten() {
			fields[attrKey] = attrVal
		}
	}

	return tags, fields
}

func NewSmartFromInfluxDB(attrs map[string]interface{}) (*Smart, error) {
	//go though the massive map returned from influxdb. If a key is associated with the Smart struct, assign it. If it starts with "attr.*" group it by attributeId, and pass to attribute inflate.

	sm := Smart{
		//required fields
		Date:           attrs["_time"].(time.Time),
		DeviceWWN:      attrs["device_wwn"].(string),
		DeviceProtocol: attrs["device_protocol"].(string),

		Attributes: map[string]SmartAttribute{},
	}

	for key, val := range attrs {
		switch key {
		case "temp":
			sm.Temp = val.(int64)
		case "power_on_hours":
			sm.PowerOnHours = val.(int64)
		case "power_cycle_count":
			sm.PowerCycleCount = val.(int64)
		default:
			// this key is unknown.
			if !strings.HasPrefix(key, "attr.") {
				continue
			}
			//this is a attribute, lets group it with its related "siblings", populating a SmartAttribute object
			keyParts := strings.Split(key, ".")
			attributeId := keyParts[1]
			if _, ok := sm.Attributes[attributeId]; !ok {
				// init the attribute group
				if sm.DeviceProtocol == pkg.DeviceProtocolAta {
					sm.Attributes[attributeId] = &SmartAtaAttribute{}
				} else if sm.DeviceProtocol == pkg.DeviceProtocolNvme {
					sm.Attributes[attributeId] = &SmartNvmeAttribute{}
				} else if sm.DeviceProtocol == pkg.DeviceProtocolScsi {
					sm.Attributes[attributeId] = &SmartScsiAttribute{}
				} else {
					return nil, fmt.Errorf("Unknown Device Protocol: %s", sm.DeviceProtocol)
				}
			}

			sm.Attributes[attributeId].Inflate(key, val)
		}

	}

	log.Printf("Found Smart Device (%s) Attributes (%v)", sm.DeviceWWN, len(sm.Attributes))

	return &sm, nil
}

//Parse Collector SMART data results and create Smart object (and associated SmartAtaAttribute entries)
func (sm *Smart) FromCollectorSmartInfo(wwn string, info collector.SmartInfo) error {
	sm.DeviceWWN = wwn
	sm.Date = time.Unix(info.LocalTime.TimeT, 0)

	//smart metrics
	sm.Temp = info.Temperature.Current
	sm.PowerCycleCount = info.PowerCycleCount
	sm.PowerOnHours = info.PowerOnTime.Hours
	if !info.SmartStatus.Passed {
		sm.Status = pkg.DeviceStatusFailedSmart
	}

	sm.DeviceProtocol = info.Device.Protocol
	// process ATA/NVME/SCSI protocol data
	sm.Attributes = map[string]SmartAttribute{}
	if sm.DeviceProtocol == pkg.DeviceProtocolAta {
		sm.ProcessAtaSmartInfo(info.AtaSmartAttributes.Table)
	} else if sm.DeviceProtocol == pkg.DeviceProtocolNvme {
		sm.ProcessNvmeSmartInfo(info.NvmeSmartHealthInformationLog)
	} else if sm.DeviceProtocol == pkg.DeviceProtocolScsi {
		sm.ProcessScsiSmartInfo(info.ScsiGrownDefectList, info.ScsiErrorCounterLog)
	}

	return nil
}

//generate SmartAtaAttribute entries from Scrutiny Collector Smart data.
func (sm *Smart) ProcessAtaSmartInfo(tableItems []collector.AtaSmartAttributesTableItem) {
	for _, collectorAttr := range tableItems {
		attrModel := SmartAtaAttribute{
			AttributeId: collectorAttr.ID,
			Value:       collectorAttr.Value,
			Worst:       collectorAttr.Worst,
			Threshold:   collectorAttr.Thresh,
			RawValue:    collectorAttr.Raw.Value,
			RawString:   collectorAttr.Raw.String,
			WhenFailed:  collectorAttr.WhenFailed,
		}

		//now that we've parsed the data from the smartctl response, lets match it against our metadata rules and add additional Scrutiny specific data.
		if smartMetadata, ok := thresholds.AtaMetadata[collectorAttr.ID]; ok {
			if smartMetadata.Transform != nil {
				attrModel.TransformedValue = smartMetadata.Transform(attrModel.Value, attrModel.RawValue, attrModel.RawString)
			}
		}
		attrModel.PopulateAttributeStatus()
		sm.Attributes[strconv.Itoa(collectorAttr.ID)] = &attrModel
		if attrModel.Status == pkg.SmartAttributeStatusFailed {
			sm.Status = pkg.Set(sm.Status, pkg.DeviceStatusFailedScrutiny)
		}
	}
}

//generate SmartNvmeAttribute entries from Scrutiny Collector Smart data.
func (sm *Smart) ProcessNvmeSmartInfo(nvmeSmartHealthInformationLog collector.NvmeSmartHealthInformationLog) {

	sm.Attributes = map[string]SmartAttribute{
		"critical_warning":     (&SmartNvmeAttribute{AttributeId: "critical_warning", Value: nvmeSmartHealthInformationLog.CriticalWarning, Threshold: 0}).PopulateAttributeStatus(),
		"temperature":          (&SmartNvmeAttribute{AttributeId: "temperature", Value: nvmeSmartHealthInformationLog.Temperature, Threshold: -1}).PopulateAttributeStatus(),
		"available_spare":      (&SmartNvmeAttribute{AttributeId: "available_spare", Value: nvmeSmartHealthInformationLog.AvailableSpare, Threshold: nvmeSmartHealthInformationLog.AvailableSpareThreshold}).PopulateAttributeStatus(),
		"percentage_used":      (&SmartNvmeAttribute{AttributeId: "percentage_used", Value: nvmeSmartHealthInformationLog.PercentageUsed, Threshold: 100}).PopulateAttributeStatus(),
		"data_units_read":      (&SmartNvmeAttribute{AttributeId: "data_units_read", Value: nvmeSmartHealthInformationLog.DataUnitsRead, Threshold: -1}).PopulateAttributeStatus(),
		"data_units_written":   (&SmartNvmeAttribute{AttributeId: "data_units_written", Value: nvmeSmartHealthInformationLog.DataUnitsWritten, Threshold: -1}).PopulateAttributeStatus(),
		"host_reads":           (&SmartNvmeAttribute{AttributeId: "host_reads", Value: nvmeSmartHealthInformationLog.HostReads, Threshold: -1}).PopulateAttributeStatus(),
		"host_writes":          (&SmartNvmeAttribute{AttributeId: "host_writes", Value: nvmeSmartHealthInformationLog.HostWrites, Threshold: -1}).PopulateAttributeStatus(),
		"controller_busy_time": (&SmartNvmeAttribute{AttributeId: "controller_busy_time", Value: nvmeSmartHealthInformationLog.ControllerBusyTime, Threshold: -1}).PopulateAttributeStatus(),
		"power_cycles":         (&SmartNvmeAttribute{AttributeId: "power_cycles", Value: nvmeSmartHealthInformationLog.PowerCycles, Threshold: -1}).PopulateAttributeStatus(),
		"power_on_hours":       (&SmartNvmeAttribute{AttributeId: "power_on_hours", Value: nvmeSmartHealthInformationLog.PowerOnHours, Threshold: -1}).PopulateAttributeStatus(),
		"unsafe_shutdowns":     (&SmartNvmeAttribute{AttributeId: "unsafe_shutdowns", Value: nvmeSmartHealthInformationLog.UnsafeShutdowns, Threshold: -1}).PopulateAttributeStatus(),
		"media_errors":         (&SmartNvmeAttribute{AttributeId: "media_errors", Value: nvmeSmartHealthInformationLog.MediaErrors, Threshold: 0}).PopulateAttributeStatus(),
		"num_err_log_entries":  (&SmartNvmeAttribute{AttributeId: "num_err_log_entries", Value: nvmeSmartHealthInformationLog.NumErrLogEntries, Threshold: 0}).PopulateAttributeStatus(),
		"warning_temp_time":    (&SmartNvmeAttribute{AttributeId: "warning_temp_time", Value: nvmeSmartHealthInformationLog.WarningTempTime, Threshold: -1}).PopulateAttributeStatus(),
		"critical_comp_time":   (&SmartNvmeAttribute{AttributeId: "critical_comp_time", Value: nvmeSmartHealthInformationLog.CriticalCompTime, Threshold: -1}).PopulateAttributeStatus(),
	}

	//find analyzed attribute status
	for _, val := range sm.Attributes {
		if val.GetStatus() == pkg.SmartAttributeStatusFailed {
			sm.Status = pkg.Set(sm.Status, pkg.DeviceStatusFailedScrutiny)
		}
	}
}

//generate SmartScsiAttribute entries from Scrutiny Collector Smart data.
func (sm *Smart) ProcessScsiSmartInfo(defectGrownList int64, scsiErrorCounterLog collector.ScsiErrorCounterLog) {
	sm.Attributes = map[string]SmartAttribute{
		"scsi_grown_defect_list":                     (&SmartScsiAttribute{AttributeId: "scsi_grown_defect_list", Value: defectGrownList, Threshold: 0}).PopulateAttributeStatus(),
		"read_errors_corrected_by_eccfast":           (&SmartScsiAttribute{AttributeId: "read_errors_corrected_by_eccfast", Value: scsiErrorCounterLog.Read.ErrorsCorrectedByEccfast, Threshold: -1}).PopulateAttributeStatus(),
		"read_errors_corrected_by_eccdelayed":        (&SmartScsiAttribute{AttributeId: "read_errors_corrected_by_eccdelayed", Value: scsiErrorCounterLog.Read.ErrorsCorrectedByEccdelayed, Threshold: -1}).PopulateAttributeStatus(),
		"read_errors_corrected_by_rereads_rewrites":  (&SmartScsiAttribute{AttributeId: "read_errors_corrected_by_rereads_rewrites", Value: scsiErrorCounterLog.Read.ErrorsCorrectedByRereadsRewrites, Threshold: 0}).PopulateAttributeStatus(),
		"read_total_errors_corrected":                (&SmartScsiAttribute{AttributeId: "read_total_errors_corrected", Value: scsiErrorCounterLog.Read.TotalErrorsCorrected, Threshold: -1}).PopulateAttributeStatus(),
		"read_correction_algorithm_invocations":      (&SmartScsiAttribute{AttributeId: "read_correction_algorithm_invocations", Value: scsiErrorCounterLog.Read.CorrectionAlgorithmInvocations, Threshold: -1}).PopulateAttributeStatus(),
		"read_total_uncorrected_errors":              (&SmartScsiAttribute{AttributeId: "read_total_uncorrected_errors", Value: scsiErrorCounterLog.Read.TotalUncorrectedErrors, Threshold: 0}).PopulateAttributeStatus(),
		"write_errors_corrected_by_eccfast":          (&SmartScsiAttribute{AttributeId: "write_errors_corrected_by_eccfast", Value: scsiErrorCounterLog.Write.ErrorsCorrectedByEccfast, Threshold: -1}).PopulateAttributeStatus(),
		"write_errors_corrected_by_eccdelayed":       (&SmartScsiAttribute{AttributeId: "write_errors_corrected_by_eccdelayed", Value: scsiErrorCounterLog.Write.ErrorsCorrectedByEccdelayed, Threshold: -1}).PopulateAttributeStatus(),
		"write_errors_corrected_by_rereads_rewrites": (&SmartScsiAttribute{AttributeId: "write_errors_corrected_by_rereads_rewrites", Value: scsiErrorCounterLog.Write.ErrorsCorrectedByRereadsRewrites, Threshold: 0}).PopulateAttributeStatus(),
		"write_total_errors_corrected":               (&SmartScsiAttribute{AttributeId: "write_total_errors_corrected", Value: scsiErrorCounterLog.Write.TotalErrorsCorrected, Threshold: -1}).PopulateAttributeStatus(),
		"write_correction_algorithm_invocations":     (&SmartScsiAttribute{AttributeId: "write_correction_algorithm_invocations", Value: scsiErrorCounterLog.Write.CorrectionAlgorithmInvocations, Threshold: -1}).PopulateAttributeStatus(),
		"write_total_uncorrected_errors":             (&SmartScsiAttribute{AttributeId: "write_total_uncorrected_errors", Value: scsiErrorCounterLog.Write.TotalUncorrectedErrors, Threshold: 0}).PopulateAttributeStatus(),
	}

	//find analyzed attribute status
	for _, val := range sm.Attributes {
		if val.GetStatus() == pkg.SmartAttributeStatusFailed {
			sm.Status = pkg.Set(sm.Status, pkg.DeviceStatusFailedScrutiny)
		}
	}
}
