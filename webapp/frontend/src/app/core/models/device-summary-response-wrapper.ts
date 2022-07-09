import {DeviceSummaryModel} from 'app/core/models/device-summary-model';

// maps to webapp/backend/pkg/models/device_summary.go
export interface DeviceSummaryResponseWrapper {
    success: boolean;
    errors: any[];
    data: {
        summary: { [key: string]: DeviceSummaryModel }
    }
}
