// maps to webapp/backend/pkg/models/measurements/smart.go
import {SmartAttributeModel} from './smart-attribute-model';

export interface SmartModel {
    date: string;
    device_wwn: string;
    device_protocol: string;

    temp: number;
    power_on_hours: number;
    power_cycle_count: number;
    logical_block_size?: number; //logical block size in bytes (typically 512 or 4096)
    attrs: { [key: string]: SmartAttributeModel }
}
