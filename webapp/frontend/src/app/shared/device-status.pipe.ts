import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'deviceStatus'
})
export class DeviceStatusPipe implements PipeTransform {

  transform(deviceStatusFlag: number): string {
      if(deviceStatusFlag === 0){
          return 'passed'
      } else if(deviceStatusFlag === 3){
          return 'failed: both'
      } else if(deviceStatusFlag === 2) {
          return 'failed: scrutiny'
      } else if(deviceStatusFlag === 1) {
          return 'failed: smart'
      }
    return 'unknown'
  }

}
