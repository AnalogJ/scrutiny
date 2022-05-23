import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'deviceSort'
})
export class DeviceSortPipe implements PipeTransform {

    numericalStatus(deviceSummary): number {
        if(!deviceSummary.smart){
            return 0
        } else if (deviceSummary.device.device_status == 0){
            return 1
        } else {
            return deviceSummary.device.device_status * -1 // will return range from -1, -2, -3
        }
    }


  transform(deviceSummaries: Array<unknown>, sortBy = ''): Array<unknown> {
      //failed, unknown/empty, passed
      deviceSummaries.sort((a: any, b: any) => {

          let left = this.numericalStatus(a)
          let right = this.numericalStatus(b)

          return left - right;
      });


    return deviceSummaries;
  }

}
