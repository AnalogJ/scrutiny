import {DetailExportService} from './detail-export.service';
import {jsPDF} from 'jspdf';
import {DeviceModel} from 'app/core/models/device-model';
import {SmartModel} from 'app/core/models/measurements/smart-model';

describe('DetailExportService', () => {
    it('downloads a one-page PDF containing the Scrutiny logo', async () => {
        const service = new DetailExportService();
        const logo = 'data:image/png;base64,logo';
        spyOn<any>(service, 'loadLogo').and.resolveTo(logo);
        const addImage = spyOn(jsPDF.prototype, 'addImage').and.callFake(() => jsPDF.prototype);
        const setFillColor = spyOn(jsPDF.prototype, 'setFillColor').and.callFake(() => jsPDF.prototype);
        const roundedRect = spyOn(jsPDF.prototype, 'roundedRect').and.callFake(() => jsPDF.prototype);
        const text = spyOn(jsPDF.prototype, 'text').and.callFake(() => jsPDF.prototype);
        const save = spyOn(jsPDF.prototype, 'save');

        await service.exportLogoPdf({
            driveName: 'WDC WD140EDFZ-11A0VA0',
            driveTitle: '/dev/sdc - WDC WD140EDFZ-11A0VA0',
            device: {
                manufacturer: 'ATA',
                model_name: 'WDC WD140EDFZ-11A0VA0',
                serial_number: '9RK4XXXXX',
                scrutiny_uuid: '42caca8a-9b95-5c75-b059-305771a2a193',
                wwn: '0x5000cca264ec3183',
                firmware: 'MS1OA650',
                capacity: 14000519643136,
                device_protocol: 'ATA',
                device_status: 3,
            } as DeviceModel,
            latestSmart: {
                power_cycle_count: 86,
                power_on_hours: 61320,
                temp: 25,
                attrs: {
                    '1': {
                        attribute_id: 1,
                        value: 100,
                        thresh: 16,
                        transformed_value: 0,
                        status: 0,
                        failure_rate: 0.01,
                    },
                    '5': {
                        attribute_id: 5,
                        value: 1975,
                        thresh: 5,
                        transformed_value: 0,
                        status: 4,
                        failure_rate: 1.5,
                    },
                },
            } as unknown as SmartModel,
            metadata: {
                '1': {
                    display_name: 'Read Error Rate',
                    ideal: 'high',
                    critical: false,
                    description: '',
                    display_type: 'normalized',
                },
                '5': {
                    display_name: 'Reallocated Sectors Count',
                    ideal: 'low',
                    critical: true,
                    description: '',
                    display_type: 'raw',
                },
            },
            config: {
                file_size_si_units: false,
                powered_on_hours_unit: 'humanize',
                temperature_unit: 'celsius',
                metrics: {status_threshold: 3},
            },
        },
            new Date(2026, 7, 12, 3, 45),
        );

        expect(addImage.calls.mostRecent().args as unknown[]).toEqual([logo, 'PNG', 20, 20, 70, 10.85]);
        expect(setFillColor).toHaveBeenCalledWith(81, 69, 205);
        expect(roundedRect).toHaveBeenCalledWith(12, 10, 186, 277, 3, 3, 'F');
        expect(text).toHaveBeenCalledWith('Drive Details - /dev/sdc - WDC WD140EDFZ-11A0VA0', 20, 57);
        expect(text).toHaveBeenCalledWith('Dive into S.M.A.R.T data', 20, 64);
        expect(text).toHaveBeenCalledWith('FAILED: BOTH', 22, 88);
        expect(text).toHaveBeenCalledWith('Status', 20, 93);
        expect(text).toHaveBeenCalledWith('ATA', 64, 88);
        expect(text).toHaveBeenCalledWith('Model Family', 64, 93);
        expect(text).toHaveBeenCalledWith('WDC WD140EDFZ-11A0VA0', 108, 88);
        expect(text).toHaveBeenCalledWith('Device Model', 108, 93);
        expect(text).toHaveBeenCalledWith('9RK4XXXXX', 152, 88);
        expect(text).toHaveBeenCalledWith('Serial Number', 152, 93);
        expect(text).toHaveBeenCalledWith('12.7 TiB', 108, 104);
        expect(text).toHaveBeenCalledWith('Capacity', 108, 109);
        expect(text.calls.allArgs().some(args => args[0] === 'Scrutiny UUID')).toBeFalse();
        expect(text).toHaveBeenCalledWith('S.M.A.R.T ATA ATTRIBUTES', 20, 142);
        expect(text).toHaveBeenCalledWith('Read Error Rate', 58, 156);
        expect(text).toHaveBeenCalledWith('Reallocated Sectors Count', 58, 163);
        expect(setFillColor).toHaveBeenCalledWith(188, 240, 218);
        expect(roundedRect).toHaveBeenCalledWith(19, 152.3, 16, 4.8, 2.4, 2.4, 'F');
        expect(text).toHaveBeenCalledWith('PASSED', 21.5, 155.6);
        expect(text.calls.allArgs().some(args => args[0] === 'History')).toBeFalse();
        expect(text.calls.allArgs().some(args => String(args[0]).includes('visible'))).toBeFalse();
        expect(save).toHaveBeenCalledWith('2026-08-12-03-45-scrutiny-drive-report-WDC_WD140EDFZ-11A0VA0.pdf');
        expect(addImage).toHaveBeenCalledTimes(1);
    });

    it('uses protocol-specific columns for SCSI attributes', async () => {
        const service = new DetailExportService();
        spyOn<any>(service, 'loadLogo').and.resolveTo('data:image/png;base64,logo');
        const text = spyOn(jsPDF.prototype, 'text').and.callFake(() => jsPDF.prototype);
        spyOn(jsPDF.prototype, 'addImage').and.callFake(() => jsPDF.prototype);
        spyOn(jsPDF.prototype, 'save');

        await service.exportLogoPdf({
            driveName: 'SCSI Drive',
            driveTitle: '/dev/sde - SCSI Drive',
            device: {
                manufacturer: 'SEAGATE',
                model_name: 'SCSI Drive',
                serial_number: 'serial',
                scrutiny_uuid: 'uuid',
                wwn: 'wwn',
                firmware: '',
                capacity: 1000,
                device_protocol: 'SCSI',
                device_status: 0,
            } as DeviceModel,
            latestSmart: {
                attrs: {
                    read_correction_algorithm_invocations: {
                        attribute_id: 'read_correction_algorithm_invocations',
                        value: 0,
                        thresh: -1,
                        transformed_value: 0,
                        status: 0,
                    },
                },
            } as unknown as SmartModel,
            metadata: {
                read_correction_algorithm_invocations: {
                    display_name: 'Read Correction Algorithm Invocations',
                    ideal: '',
                    critical: false,
                    description: '',
                    display_type: 'raw',
                },
            },
            config: {metrics: {status_threshold: 3}},
        }, new Date(2026, 7, 12, 3, 45));

        expect(text).toHaveBeenCalledWith('Name', 50, 149.5);
        expect(text).toHaveBeenCalledWith('Read Correction Algorithm Invocations', 50, 156);
        expect(text.calls.allArgs().some(args => args[0] === 'ID')).toBeFalse();
        expect(text.calls.allArgs().some(args => args[0] === 'Ideal')).toBeFalse();
        expect(text.calls.allArgs().some(args => args[0] === 'Failure Rate')).toBeFalse();
    });
});
