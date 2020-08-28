import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'deviceSort'
})
export class DeviceSortPipe implements PipeTransform {

    numericalStatus(device): number {
        if(!device.smart_results[0]){
            return 0
        } else if (device.smart_results[0].smart_status == 'passed'){
            return 1
        } else {
            return -1
        }
    }


  transform(devices: Array<unknown>, ...args: unknown[]): Array<unknown> {
      //failed, unknown/empty, passed
      devices.sort((a: any, b: any) => {

          let left = this.numericalStatus(a)
          let right = this.numericalStatus(b)

          return left - right;
      });


    return devices;
  }

}
