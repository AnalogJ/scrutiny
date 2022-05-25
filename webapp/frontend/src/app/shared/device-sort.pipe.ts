import { Pipe, PipeTransform } from '@angular/core';
import {deviceDisplayTitle} from "app/layout/common/dashboard-device/dashboard-device.component";

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
            let left = deviceDisplayTitle(a.device, _dashboardDisplay) || deviceDisplayTitle(a.device, 'name')
            let right = deviceDisplayTitle(b.device, _dashboardDisplay) || deviceDisplayTitle(b.device, 'name')

            if( left < right )
                return -1;

            if( left > right )
                return 1;

            return 0;
        }
    }


  transform(deviceSummaries: Array<unknown>, sortBy = 'status', dashboardDisplay = "name"): Array<unknown> {
    let compareFn = undefined
    switch (sortBy) {
        case 'status':
            compareFn = this.statusCompareFn
            break;
        case 'title':
            compareFn = this.titleCompareFn(dashboardDisplay)
            break;
    }

      //failed, unknown/empty, passed
      deviceSummaries.sort(compareFn);

    return deviceSummaries;
  }

}
