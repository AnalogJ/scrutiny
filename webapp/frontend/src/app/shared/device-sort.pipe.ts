import { Pipe, PipeTransform } from '@angular/core';
import {DeviceTitlePipe} from "./device-title.pipe";

@Pipe({
  name: 'deviceSort'
})
export class DeviceSortPipe implements PipeTransform {

    statusCompareFn(a: any, b: any) {
        function deviceStatus(deviceSummary): number {
            if(!deviceSummary.smart){
                return 0
            } else if (deviceSummary.device.device_status == 0){
                return 1
            } else {
                return deviceSummary.device.device_status * -1 // will return range from -1, -2, -3
            }
        }

        let left = deviceStatus(a)
        let right = deviceStatus(b)

        return left - right;
    }

    titleCompareFn(dashboardDisplay: string) {
        return function (a: any, b: any){
            let _dashboardDisplay = dashboardDisplay
            let left = DeviceTitlePipe.deviceTitleForType(a.device, _dashboardDisplay) || DeviceTitlePipe.deviceTitleForType(a.device, 'name')
            let right = DeviceTitlePipe.deviceTitleForType(b.device, _dashboardDisplay) || DeviceTitlePipe.deviceTitleForType(b.device, 'name')

            if( left < right )
                return -1;

            if( left > right )
                return 1;

            return 0;
        }
    }

    ageCompareFn(a: any, b: any) {
        const left = a.smart?.power_on_hours
        const right = b.smart?.power_on_hours

        return left - right;
    }


  transform(deviceSummaries: Array<unknown>, sortBy = 'status', dashboardDisplay = 'name'): Array<unknown> {
    let compareFn: any
    switch (sortBy) {
        case 'status':
            compareFn = this.statusCompareFn
            break;
        case 'title':
            compareFn = this.titleCompareFn(dashboardDisplay)
            break;
        case 'age':
            compareFn = this.ageCompareFn
            break;
    }

      // failed, unknown/empty, passed
      deviceSummaries.sort(compareFn);

    return deviceSummaries;
  }

}
