// maps to webapp/backend/pkg/models/measurements/smart_ata_attribute.go
// maps to webapp/backend/pkg/models/measurements/smart_nvme_attribute.go
// maps to webapp/backend/pkg/models/measurements/smart_scsi_attribute.go
export interface SmartAttributeModel {
    attribute_id: number | string
    value: number
    thresh: number
    worst?: number
    raw_value?: number
    raw_string?: string
    when_failed?: string

    transformed_value: number
    status: number
    status_reason?: string
    failure_rate?: number

    chartData?: any[]
}
