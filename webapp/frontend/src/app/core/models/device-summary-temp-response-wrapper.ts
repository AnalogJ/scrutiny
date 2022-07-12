import {SmartTemperatureModel} from './measurements/smart-temperature-model';

export interface DeviceSummaryTempResponseWrapper {
    success: boolean;
    errors: any[];
    data: {
        temp_history: { [key: string]: SmartTemperatureModel[]; }
    }
}
