import { Pipe, PipeTransform } from '@angular/core';
import {DeviceTitlePipe} from 'app/shared/device-title.pipe';

@Pipe({
    name: 'deviceSort',
    standalone: false
})
export class DeviceSortPipe implements PipeTransform {

    statusCompareFn(a: any, b: any): number {
        function deviceStatus(deviceSummary): number {
            if(!deviceSummary.smart){
                return 0
            } else if (deviceSummary.device.device_status === 0){
                return 1
            } else {
                return deviceSummary.device.device_status * -1 // will return range from -1, -2, -3
            }
        }

        const left = deviceStatus(a)
        const right = deviceStatus(b)

        return left - right;
    }

    titleCompareFn(dashboardDisplay: string) {
        return function (a: any, b: any){
            const _dashboardDisplay = dashboardDisplay
            const left = DeviceTitlePipe.deviceTitleForType(a.device, _dashboardDisplay) || DeviceTitlePipe.deviceTitleForType(a.device, 'name')
            const right = DeviceTitlePipe.deviceTitleForType(b.device, _dashboardDisplay) || DeviceTitlePipe.deviceTitleForType(b.device, 'name')

            if( left < right )
                return -1;

            if( left > right )
                return 1;

            return 0;
        }
    }

    ageCompareFn(a: any, b: any): number {
        const left = a.smart?.power_on_hours ?? 0;
        const right = b.smart?.power_on_hours ?? 0;

        return left - right;
    }

    capacityCompareFn(a: any, b: any): number {
        const left = a.device?.capacity || 0;
        const right = b.device?.capacity || 0;

        return left - right;
    }

    temperatureCompareFn(a: any, b: any): number {
        const left = a.smart?.temp ?? Number.MAX_SAFE_INTEGER;
        const right = b.smart?.temp ?? Number.MAX_SAFE_INTEGER;

        return left - right;
    }


  transform(deviceSummaries: Array<unknown>, sortBy = 'status', dashboardDisplay = 'name'): Array<unknown> {
    // Map legacy values to new format for backward compatibility
    const legacyMap: Record<string, string> = {
        'status': 'status_desc',
        'title': 'title_asc',
        'age': 'age_asc'
    };
    const normalizedSort = legacyMap[sortBy] || sortBy;

    // Parse sort field and direction
    const isDesc = normalizedSort.endsWith('_desc');
    const sortField = normalizedSort.replace(/_(?:asc|desc)$/, '');
    const direction = isDesc ? -1 : 1;

    // Get compare function
    let compareFn: (a: any, b: any) => number;
    switch (sortField) {
        case 'status':
            compareFn = this.statusCompareFn;
            break;
        case 'title':
            compareFn = this.titleCompareFn(dashboardDisplay);
            break;
        case 'age':
            compareFn = this.ageCompareFn;
            break;
        case 'capacity':
            compareFn = this.capacityCompareFn;
            break;
        case 'temperature':
            compareFn = this.temperatureCompareFn;
            break;
        default:
            compareFn = this.statusCompareFn;
    }

    // Apply sort with direction
    deviceSummaries.sort((a, b) => direction * compareFn(a, b));

    return deviceSummaries;
  }

}
