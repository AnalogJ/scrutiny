import {DeviceModel} from 'app/core/models/device-model';
import {SmartModel} from 'app/core/models/measurements/smart-model';
import {AttributeMetadataModel} from 'app/core/models/thresholds/attribute-metadata-model';

// maps to webapp/backend/pkg/models/device_summary.go
export interface DeviceDetailsResponseWrapper {
    success: boolean;
    errors?: any[];
    data: {
        device: DeviceModel;
        smart_results: SmartModel[];
    },
    metadata: { [key: string]: AttributeMetadataModel } | { [key: number]: AttributeMetadataModel };
}
