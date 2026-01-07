import {Pipe, PipeTransform} from '@angular/core';
import {DeviceModel} from 'app/core/models/device-model';

@Pipe({
  name: 'deviceTitle'
})
export class DeviceTitlePipe implements PipeTransform {

    static deviceTitleForType(device: DeviceModel, titleType: string): string {
        const titleParts = []
        switch(titleType){
            case 'name':
                titleParts.push(`/dev/${device.device_name}`)
                if (device.device_type && device.device_type !== 'scsi' && device.device_type !== 'ata'){
                    titleParts.push(device.device_type)
                }
                titleParts.push(device.model_name)

                break;
            case 'serial_id':
                if(!device.device_serial_id) return ''
                titleParts.push(`/by-id/${device.device_serial_id}`)
                break;
            case 'uuid':
                if(!device.device_uuid) return ''
                titleParts.push(`/by-uuid/${device.device_uuid}`)
                break;
            case 'label':
                if(device.label){
                    titleParts.push(device.label)
                } else if(device.device_label){
                    titleParts.push(`/by-label/${device.device_label}`)
                }
                break;
        }
        return titleParts.join(' - ')
    }

    static deviceTitleWithFallback(device: DeviceModel, titleType: string): string {
        const titleParts = []
        if (device.host_id) titleParts.push(device.host_id)

        // add device identifier (fallback to generated device name)
        titleParts.push(DeviceTitlePipe.deviceTitleForType(device, titleType) || DeviceTitlePipe.deviceTitleForType(device, 'name'))

        return titleParts.join(' - ')
    }


    transform(device: DeviceModel, titleType: string = 'name'): string {
        return DeviceTitlePipe.deviceTitleWithFallback(device, titleType)
    }

}
