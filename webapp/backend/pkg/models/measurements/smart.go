package measurements

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
	"log"
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

	log.Printf("Prefetched Smart: %v\n", sm)

	//two steps (because we dont know the
	for key, val := range attrs {
		log.Printf("Found Attribute (%s = %v)\n", key, val)

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

	log.Printf("########NUMBER OF ATTRIBUTES: %v", len(sm.Attributes))
	log.Printf("########SMART: %v", sm)

	//panic("ERROR HERE.")

	//log.Printf("Sm.Attributes: %v", sm.Attributes)
	//log.Printf("sm.Attributes[attributeId]: %v", sm.Attributes[attributeId])

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

	sm.DeviceProtocol = info.Device.Protocol
	// process ATA/NVME/SCSI protocol data
	sm.Attributes = map[string]SmartAttribute{}
	if sm.DeviceProtocol == pkg.DeviceProtocolAta {
		sm.ProcessAtaSmartInfo(info)
	} else if sm.DeviceProtocol == pkg.DeviceProtocolNvme {
		sm.ProcessNvmeSmartInfo(info)
	} else if sm.DeviceProtocol == pkg.DeviceProtocolScsi {
		sm.ProcessScsiSmartInfo(info)
	}

	return nil
}

//generate SmartAtaAttribute entries from Scrutiny Collector Smart data.
func (sm *Smart) ProcessAtaSmartInfo(info collector.SmartInfo) {
	sm.Status = pkg.DeviceStatusPassed
	for _, collectorAttr := range info.AtaSmartAttributes.Table {
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
		sm.Attributes[string(collectorAttr.ID)] = &attrModel
		if attrModel.Status == pkg.SmartAttributeStatusFailed {
			sm.Status = pkg.DeviceStatusFailedScrutiny
		}
	}
}

//generate SmartNvmeAttribute entries from Scrutiny Collector Smart data.
func (sm *Smart) ProcessNvmeSmartInfo(info collector.SmartInfo) {

	sm.Attributes = map[string]SmartAttribute{
		"critical_warning":     (&SmartNvmeAttribute{AttributeId: "critical_warning", Value: info.NvmeSmartHealthInformationLog.CriticalWarning, Threshold: 0}).PopulateAttributeStatus(),
		"temperature":          (&SmartNvmeAttribute{AttributeId: "temperature", Value: info.NvmeSmartHealthInformationLog.Temperature, Threshold: -1}).PopulateAttributeStatus(),
		"available_spare":      (&SmartNvmeAttribute{AttributeId: "available_spare", Value: info.NvmeSmartHealthInformationLog.AvailableSpare, Threshold: info.NvmeSmartHealthInformationLog.AvailableSpareThreshold}).PopulateAttributeStatus(),
		"percentage_used":      (&SmartNvmeAttribute{AttributeId: "percentage_used", Value: info.NvmeSmartHealthInformationLog.PercentageUsed, Threshold: 100}).PopulateAttributeStatus(),
		"data_units_read":      (&SmartNvmeAttribute{AttributeId: "data_units_read", Value: info.NvmeSmartHealthInformationLog.DataUnitsRead, Threshold: -1}).PopulateAttributeStatus(),
		"data_units_written":   (&SmartNvmeAttribute{AttributeId: "data_units_written", Value: info.NvmeSmartHealthInformationLog.DataUnitsWritten, Threshold: -1}).PopulateAttributeStatus(),
		"host_reads":           (&SmartNvmeAttribute{AttributeId: "host_reads", Value: info.NvmeSmartHealthInformationLog.HostReads, Threshold: -1}).PopulateAttributeStatus(),
		"host_writes":          (&SmartNvmeAttribute{AttributeId: "host_writes", Value: info.NvmeSmartHealthInformationLog.HostWrites, Threshold: -1}).PopulateAttributeStatus(),
		"controller_busy_time": (&SmartNvmeAttribute{AttributeId: "controller_busy_time", Value: info.NvmeSmartHealthInformationLog.ControllerBusyTime, Threshold: -1}).PopulateAttributeStatus(),
		"power_cycles":         (&SmartNvmeAttribute{AttributeId: "power_cycles", Value: info.NvmeSmartHealthInformationLog.PowerCycles, Threshold: -1}).PopulateAttributeStatus(),
		"power_on_hours":       (&SmartNvmeAttribute{AttributeId: "power_on_hours", Value: info.NvmeSmartHealthInformationLog.PowerOnHours, Threshold: -1}).PopulateAttributeStatus(),
		"unsafe_shutdowns":     (&SmartNvmeAttribute{AttributeId: "unsafe_shutdowns", Value: info.NvmeSmartHealthInformationLog.UnsafeShutdowns, Threshold: -1}).PopulateAttributeStatus(),
		"media_errors":         (&SmartNvmeAttribute{AttributeId: "media_errors", Value: info.NvmeSmartHealthInformationLog.MediaErrors, Threshold: 0}).PopulateAttributeStatus(),
		"num_err_log_entries":  (&SmartNvmeAttribute{AttributeId: "num_err_log_entries", Value: info.NvmeSmartHealthInformationLog.NumErrLogEntries, Threshold: 0}).PopulateAttributeStatus(),
		"warning_temp_time":    (&SmartNvmeAttribute{AttributeId: "warning_temp_time", Value: info.NvmeSmartHealthInformationLog.WarningTempTime, Threshold: -1}).PopulateAttributeStatus(),
		"critical_comp_time":   (&SmartNvmeAttribute{AttributeId: "critical_comp_time", Value: info.NvmeSmartHealthInformationLog.CriticalCompTime, Threshold: -1}).PopulateAttributeStatus(),
	}

	//find analyzed attribute status
	for _, val := range sm.Attributes {
		if val.GetStatus() == pkg.SmartAttributeStatusFailed {
			sm.Status = pkg.DeviceStatusFailedScrutiny
		}
	}
}

//generate SmartScsiAttribute entries from Scrutiny Collector Smart data.
func (sm *Smart) ProcessScsiSmartInfo(info collector.SmartInfo) {
	sm.Attributes = map[string]SmartAttribute{
		"scsi_grown_defect_list":                     (&SmartScsiAttribute{AttributeId: "scsi_grown_defect_list", Value: info.ScsiGrownDefectList, Threshold: 0}).PopulateAttributeStatus(),
		"read_errors_corrected_by_eccfast":           (&SmartScsiAttribute{AttributeId: "read_errors_corrected_by_eccfast", Value: info.ScsiErrorCounterLog.Read.ErrorsCorrectedByEccfast, Threshold: -1}).PopulateAttributeStatus(),
		"read_errors_corrected_by_eccdelayed":        (&SmartScsiAttribute{AttributeId: "read_errors_corrected_by_eccdelayed", Value: info.ScsiErrorCounterLog.Read.ErrorsCorrectedByEccdelayed, Threshold: -1}).PopulateAttributeStatus(),
		"read_errors_corrected_by_rereads_rewrites":  (&SmartScsiAttribute{AttributeId: "read_errors_corrected_by_rereads_rewrites", Value: info.ScsiErrorCounterLog.Read.ErrorsCorrectedByRereadsRewrites, Threshold: 0}).PopulateAttributeStatus(),
		"read_total_errors_corrected":                (&SmartScsiAttribute{AttributeId: "read_total_errors_corrected", Value: info.ScsiErrorCounterLog.Read.TotalErrorsCorrected, Threshold: -1}).PopulateAttributeStatus(),
		"read_correction_algorithm_invocations":      (&SmartScsiAttribute{AttributeId: "read_correction_algorithm_invocations", Value: info.ScsiErrorCounterLog.Read.CorrectionAlgorithmInvocations, Threshold: -1}).PopulateAttributeStatus(),
		"read_total_uncorrected_errors":              (&SmartScsiAttribute{AttributeId: "read_total_uncorrected_errors", Value: info.ScsiErrorCounterLog.Read.TotalUncorrectedErrors, Threshold: 0}).PopulateAttributeStatus(),
		"write_errors_corrected_by_eccfast":          (&SmartScsiAttribute{AttributeId: "write_errors_corrected_by_eccfast", Value: info.ScsiErrorCounterLog.Write.ErrorsCorrectedByEccfast, Threshold: -1}).PopulateAttributeStatus(),
		"write_errors_corrected_by_eccdelayed":       (&SmartScsiAttribute{AttributeId: "write_errors_corrected_by_eccdelayed", Value: info.ScsiErrorCounterLog.Write.ErrorsCorrectedByEccdelayed, Threshold: -1}).PopulateAttributeStatus(),
		"write_errors_corrected_by_rereads_rewrites": (&SmartScsiAttribute{AttributeId: "write_errors_corrected_by_rereads_rewrites", Value: info.ScsiErrorCounterLog.Write.ErrorsCorrectedByRereadsRewrites, Threshold: 0}).PopulateAttributeStatus(),
		"write_total_errors_corrected":               (&SmartScsiAttribute{AttributeId: "write_total_errors_corrected", Value: info.ScsiErrorCounterLog.Write.TotalErrorsCorrected, Threshold: -1}).PopulateAttributeStatus(),
		"write_correction_algorithm_invocations":     (&SmartScsiAttribute{AttributeId: "write_correction_algorithm_invocations", Value: info.ScsiErrorCounterLog.Write.CorrectionAlgorithmInvocations, Threshold: -1}).PopulateAttributeStatus(),
		"write_total_uncorrected_errors":             (&SmartScsiAttribute{AttributeId: "write_total_uncorrected_errors", Value: info.ScsiErrorCounterLog.Write.TotalUncorrectedErrors, Threshold: 0}).PopulateAttributeStatus(),
	}

	//find analyzed attribute status
	for _, val := range sm.Attributes {
		if val.GetStatus() == pkg.SmartAttributeStatusFailed {
			sm.Status = pkg.DeviceStatusFailedScrutiny
		}
	}
}
