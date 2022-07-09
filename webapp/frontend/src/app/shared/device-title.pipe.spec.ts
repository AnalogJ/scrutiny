import {DeviceTitlePipe} from './device-title.pipe';
import {DeviceModel} from 'app/core/models/device-model';

describe('DeviceTitlePipe', () => {
    it('create an instance', () => {
        const pipe = new DeviceTitlePipe();
        expect(pipe).toBeTruthy();
    });

    describe('#deviceTitleForType', () => {
        const testCases = [
            {
                'device': {
                    'device_name': 'sda',
                    'device_type': 'ata',
                    'model_name': 'Samsung',
                },
                'titleType': 'name',
                'result': '/dev/sda - Samsung'
            },{
                'device': {
                    'device_name': 'nvme0',
                    'device_type': 'nvme',
                    'model_name': 'Samsung',
                },
                'titleType': 'name',
                'result': '/dev/nvme0 - nvme - Samsung'
            },{
                'device': {},
                'titleType': 'serial_id',
                'result': ''
            },{
                'device': {
                    'device_serial_id': 'ata-WDC_WD140EDFZ-11AXXXXX_9RXXXXXX',
                },
                'titleType': 'serial_id',
                'result': '/by-id/ata-WDC_WD140EDFZ-11AXXXXX_9RXXXXXX'
            },{
                'device': {},
                'titleType': 'uuid',
                'result': ''
            },{
                'device': {
                    'device_uuid': 'abcdef-1234-4567-8901'
                },
                'titleType': 'uuid',
                'result': '/by-uuid/abcdef-1234-4567-8901'
            },{
                'device': {},
                'titleType': 'label',
                'result': ''
            },{
                'device': {
                    'label': 'custom-device-label'
                },
                'titleType': 'label',
                'result': 'custom-device-label'
            },{
                'device': {
                    'device_label': 'drive-volume-label'
                },
                'titleType': 'label',
                'result': '/by-label/drive-volume-label'
            },
        ]
        testCases.forEach((test, index) => {
            it(`should correctly format device title ${JSON.stringify(test.device)}. (testcase: ${index + 1})`, () => {
                // test
                const formatted = DeviceTitlePipe.deviceTitleForType(test.device as DeviceModel, test.titleType)
                expect(formatted).toEqual(test.result);
            });
        })
    })

    describe('#deviceTitleWithFallback',() => {
        const testCases = [
            {
                'device': {
                    'device_name': 'sda',
                    'device_type': 'ata',
                    'model_name': 'Samsung',
                },
                'titleType': 'name',
                'result': '/dev/sda - Samsung'
            },{
                'device': {
                    'device_name': 'nvme0',
                    'device_type': 'nvme',
                    'model_name': 'Samsung',
                },
                'titleType': 'name',
                'result': '/dev/nvme0 - nvme - Samsung'
            },{
                'device': {
                    'device_name': 'fallback',
                    'device_type': 'ata',
                    'model_name': 'fallback',
                },
                'titleType': 'serial_id',
                'result': '/dev/fallback - fallback'
            },{
                'device': {
                    'device_serial_id': 'ata-WDC_WD140EDFZ-11AXXXXX_9RXXXXXX',
                },
                'titleType': 'serial_id',
                'result': '/by-id/ata-WDC_WD140EDFZ-11AXXXXX_9RXXXXXX'
            },{
                'device': {
                    'device_name': 'fallback',
                    'device_type': 'ata',
                    'model_name': 'fallback',
                },
                'titleType': 'uuid',
                'result': '/dev/fallback - fallback'
            },{
                'device': {
                    'device_uuid': 'abcdef-1234-4567-8901'
                },
                'titleType': 'uuid',
                'result': '/by-uuid/abcdef-1234-4567-8901'
            },{
                'device': {
                    'device_name': 'fallback',
                    'device_type': 'ata',
                    'model_name': 'fallback',
                },
                'titleType': 'label',
                'result': '/dev/fallback - fallback'
            },{
                'device': {
                    'label': 'custom-device-label'
                },
                'titleType': 'label',
                'result': 'custom-device-label'
            },{
                'device': {
                    'device_label': 'drive-volume-label'
                },
                'titleType': 'label',
                'result': '/by-label/drive-volume-label'
            },
        ]
        testCases.forEach((test, index) => {
            it(`should correctly format device title ${JSON.stringify(test.device)}. (testcase: ${index + 1})`, () => {
                // test
                const formatted = DeviceTitlePipe.deviceTitleWithFallback(test.device as DeviceModel, test.titleType)
                expect(formatted).toEqual(test.result);
            });
        })
    })
});
