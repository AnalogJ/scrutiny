// maps to webapp/backend/pkg/models/measurements/smart.go
import {SmartAttributeModel} from './smart-attribute-model';

export interface SmartModel {
    date: string;
    device_wwn: string;
    device_protocol: string;

    temp: number;
    power_on_hours: number;
    power_cycle_count: number
    attrs: { [key: string]: SmartAttributeModel }
}
