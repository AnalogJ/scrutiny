import {Pipe, PipeTransform} from '@angular/core';
import {MetricsStatusThreshold} from '../core/config/app.config';
import {DeviceModel} from '../core/models/device-model';

const DEVICE_STATUS_NAMES: { [key: number]: string } = {
    0: 'passed',
    1: 'failed',
    2: 'failed',
    3: 'failed'
};

const DEVICE_STATUS_NAMES_WITH_REASON: { [key: number]: string } = {
    0: 'passed',
    1: 'failed: smart',
    2: 'failed: scrutiny',
    3: 'failed: both'
};


@Pipe({
    name: 'deviceStatus'
})
export class DeviceStatusPipe implements PipeTransform {


    static deviceStatusForModelWithThreshold(
        deviceModel: DeviceModel,
        hasSmartResults: boolean = true,
        threshold: MetricsStatusThreshold = MetricsStatusThreshold.Both,
        includeReason: boolean = false
    ): string {
        // no smart data, so treat the device status as unknown
        if (!hasSmartResults) {
            return 'unknown'
        }

        let statusNameLookup = DEVICE_STATUS_NAMES
        if (includeReason) {
            statusNameLookup = DEVICE_STATUS_NAMES_WITH_REASON
        }
        // determine the device status, by comparing it against the allowed threshold
        // tslint:disable-next-line:no-bitwise
        const deviceStatus = deviceModel.device_status & threshold
        return statusNameLookup[deviceStatus]
    }

    // static deviceStatusForModelWithThreshold(deviceModel: DeviceModel | any, threshold: MetricsStatusThreshold): string {
    //     // tslint:disable-next-line:no-bitwise
    //     const deviceStatus = deviceModel?.device_status & threshold
    //     if(deviceStatus === 0){
    //         return 'passed'
    //     } else if(deviceStatus === 3){
    //         return 'failed: both'
    //     } else if(deviceStatus === 2) {
    //         return 'failed: scrutiny'
    //     } else if(deviceStatus === 1) {
    //         return 'failed: smart'
    //     }
    //     return 'unknown'
    // }

    transform(
        deviceModel: DeviceModel,
        hasSmartResults: boolean = true,
        threshold: MetricsStatusThreshold = MetricsStatusThreshold.Both,
        includeReason: boolean = false
    ): string {
        return DeviceStatusPipe.deviceStatusForModelWithThreshold(deviceModel, hasSmartResults, threshold, includeReason)
    }

}
