// map to webapp/backend/pkg/thresholds/ata_attribute_metadata.go
// map to webapp/backend/pkg/thresholds/nvme_attribute_metadata.go
// map to webapp/backend/pkg/thresholds/scsi_attribute_metadata.go
export interface AttributeMetadataModel {
    display_name: string
    ideal: string
    critical: boolean
    description: string

    transform_value_unit?: string
    observed_thresholds?: any[]
    display_type: string
}
