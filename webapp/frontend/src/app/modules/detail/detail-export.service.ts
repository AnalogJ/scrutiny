import {Injectable} from '@angular/core';
import {jsPDF} from 'jspdf';
import {AppConfig, MetricsStatusThreshold} from 'app/core/config/app.config';
import {DeviceModel} from 'app/core/models/device-model';
import {SmartModel} from 'app/core/models/measurements/smart-model';
import {DeviceHoursPipe} from 'app/shared/device-hours.pipe';
import {DeviceStatusPipe} from 'app/shared/device-status.pipe';
import {FileSizePipe} from 'app/shared/file-size.pipe';
import {TemperaturePipe} from 'app/shared/temperature.pipe';
import {AttributeMetadataModel} from 'app/core/models/thresholds/attribute-metadata-model';
import {SmartAttributeModel} from 'app/core/models/measurements/smart-attribute-model';

interface DriveReport {
    driveName: string;
    driveTitle: string;
    device: DeviceModel;
    latestSmart?: SmartModel;
    metadata: { [key: string]: AttributeMetadataModel } | { [key: number]: AttributeMetadataModel };
    config: AppConfig;
}

@Injectable({providedIn: 'root'})
export class DetailExportService {
    async exportLogoPdf(report: DriveReport, exportedAt: Date = new Date()): Promise<void> {
        const logo = await this.loadLogo();
        const document = new jsPDF({unit: 'mm', format: 'a4'});

        // Mirror the web layout: a primary-color backdrop with a raised, rounded
        // header container and the drive-details heading below the logo bar.
        document.setFillColor(81, 69, 205);
        document.rect(0, 0, 210, 28, 'F');
        document.setFillColor(255, 255, 255);
        document.roundedRect(12, 10, 186, 277, 3, 3, 'F');
        document.addImage(logo, 'PNG', 20, 20, 70, 10.85);
        document.setDrawColor(226, 232, 240);
        document.line(12, 42, 198, 42);
        document.setTextColor(36, 43, 56);
        document.setFont('helvetica', 'bold');
        document.setFontSize(18);
        document.text(`Drive Details - ${report.driveTitle}`, 20, 57);
        document.setTextColor(100, 116, 139);
        document.setFont('helvetica', 'normal');
        document.setFontSize(10);
        document.text('Dive into S.M.A.R.T data', 20, 64);
        document.line(20, 70, 190, 70);

        this.drawDeviceSummary(document, report);
        this.drawSmartAttributes(document, report);
        document.save(this.buildFilename(report.driveName, exportedAt));
    }

    private drawSmartAttributes(document: jsPDF, report: DriveReport): void {
        const attributes = Object.values(report.latestSmart?.attrs ?? {})
            .sort((left, right) => Number(left.attribute_id) - Number(right.attribute_id));

        document.setTextColor(100, 116, 139);
        document.setFont('helvetica', 'bold');
        document.setFontSize(10);
        document.text(`S.M.A.R.T ${report.device.device_protocol} ATTRIBUTES`, 20, 142);

        const isAta = report.device.device_protocol === 'ATA';
        const isScsi = report.device.device_protocol === 'SCSI';
        const columns = isAta
            ? [
                {key: 'status', label: 'Status', x: 20},
                {key: 'id', label: 'ID', x: 42},
                {key: 'name', label: 'Name', x: 58},
                {key: 'value', label: 'Value', x: 120},
                {key: 'threshold', label: 'Threshold', x: 139},
                {key: 'ideal', label: 'Ideal', x: 161},
                {key: 'failure', label: 'Failure Rate', x: 176},
            ]
            : isScsi
                ? [
                    {key: 'status', label: 'Status', x: 20},
                    {key: 'name', label: 'Name', x: 50},
                    {key: 'value', label: 'Value', x: 145},
                    {key: 'threshold', label: 'Threshold', x: 170},
                ]
                : [
                    {key: 'status', label: 'Status', x: 20},
                    {key: 'name', label: 'Name', x: 50},
                    {key: 'value', label: 'Value', x: 130},
                    {key: 'threshold', label: 'Threshold', x: 151},
                    {key: 'ideal', label: 'Ideal', x: 174},
                ];
        document.setFillColor(241, 245, 249);
        document.rect(16, 145, 178, 7, 'F');
        document.setFontSize(7);
        columns.forEach(column => document.text(column.label, column.x, 149.5));

        attributes.forEach((attribute, index) => {
            const y = 156 + index * 7;
            const metadata = report.metadata[attribute.attribute_id];
            const status = this.attributeStatus(attribute.status);

            document.setDrawColor(226, 232, 240);
            document.line(16, y + 2.5, 194, y + 2.5);
            document.setFont('helvetica', 'normal');
            document.setFontSize(7);
            columns.forEach(column => {
                if (column.key === 'status') {
                    this.drawStatusChip(document, status, status.toUpperCase(), column.x, y, true);
                    return;
                }
                document.setTextColor(36, 43, 56);
                document.text(this.attributeCell(column.key, attribute, metadata), column.x, y);
            });
        });
    }

    private attributeStatus(status: number): string {
        // Keep these bit flags aligned with the status constants used by the details UI.
        // tslint:disable-next-line:no-bitwise
        if ((status & 5) !== 0) {
            return 'failed';
        }
        // tslint:disable-next-line:no-bitwise
        if ((status & 2) !== 0) {
            return 'warn';
        }
        return 'passed';
    }

    private drawStatusChip(document: jsPDF, status: string, label: string, x: number, y: number, compact = false): void {
        const palette = status === 'failed'
            ? {background: [251, 213, 213], text: [155, 28, 28]}
            : status === 'warn'
                ? {background: [252, 233, 106], text: [114, 59, 19]}
                : {background: [188, 240, 218], text: [3, 84, 63]};
        document.setFont('helvetica', 'bold');
        document.setFontSize(7);
        const standardWidths = compact
            ? {PASSED: 16, FAILED: 16, WARN: 14}
            : {PASSED: 18, FAILED: 18, WARN: 16};
        const horizontalPadding = compact ? 4 : 6;
        const minimumWidth = compact ? 14 : 16;
        const width = standardWidths[label] ?? Math.max(minimumWidth, document.getTextWidth(label) + horizontalPadding);
        const height = compact ? 4.8 : 6;
        const radius = height / 2;
        const top = compact ? y - 3.7 : y - 4.5;
        const textX = compact ? x + 1.5 : x + 2;
        const textY = compact ? y - 0.4 : y - 0.5;

        document.setFillColor(...palette.background as [number, number, number]);
        document.roundedRect(x - 1, top, width, height, radius, radius, 'F');
        document.setTextColor(...palette.text as [number, number, number]);
        document.text(label, textX, textY);
    }

    private attributeCell(key: string, attribute: SmartAttributeModel, metadata?: AttributeMetadataModel): string {
        switch (key) {
            case 'id':
                return `${attribute.attribute_id} (0x${Number(attribute.attribute_id).toString(16).toUpperCase().padStart(2, '0')})`;
            case 'name':
                return metadata?.display_name ?? 'Unknown Attribute Name';
            case 'value':
                return String(this.attributeValue(attribute, metadata));
            case 'threshold':
                return attribute.thresh === -1 ? '' : String(attribute.thresh ?? '');
            case 'ideal':
                return metadata?.display_type === 'raw' ? metadata.ideal ?? '' : '';
            case 'failure':
                return attribute.failure_rate == null ? '' : `${Math.round(attribute.failure_rate * 100)}%`;
            default:
                return '';
        }
    }

    private attributeValue(attribute: SmartAttributeModel, metadata?: AttributeMetadataModel): number {
        if (metadata?.display_type === 'raw') {
            return attribute.raw_value ?? attribute.value;
        }
        if (metadata?.display_type === 'transformed' && attribute.transformed_value) {
            return attribute.transformed_value;
        }
        return attribute.value;
    }

    private drawDeviceSummary(document: jsPDF, report: DriveReport): void {
        const smart = report.latestSmart;
        const status = DeviceStatusPipe.deviceStatusForModelWithThreshold(
            report.device,
            !!smart,
            report.config.metrics?.status_threshold ?? MetricsStatusThreshold.Both,
            true,
        ).toUpperCase();
        const fields = [
            {label: 'Status', value: status},
            {label: 'Model Family', value: report.device.manufacturer},
            {label: 'Device Model', value: report.device.model_name},
            {label: 'Serial Number', value: report.device.serial_number},
            {label: 'LU WWN Device Id', value: report.device.wwn},
            {label: 'Firmware Version', value: report.device.firmware},
            {label: 'Capacity', value: new FileSizePipe().transform(report.device.capacity, report.config.file_size_si_units)},
            {label: 'Protocol', value: report.device.device_protocol},
            {label: 'Power Cycle Count', value: smart?.power_cycle_count ?? 'Unknown'},
            {
                label: 'Powered On',
                value: smart
                    ? DeviceHoursPipe.format(smart.power_on_hours, report.config.powered_on_hours_unit, {
                        round: true,
                        largest: 1,
                        units: ['y', 'd', 'h'],
                    })
                    : 'Unknown',
            },
            {
                label: 'Temperature',
                value: smart
                    ? new TemperaturePipe().transform(smart.temp, report.config.temperature_unit, true)
                    : 'Unknown',
            },
        ];
        const columnX = [20, 64, 108, 152];

        fields.forEach((field, index) => {
            const x = columnX[index % 4];
            const y = 88 + Math.floor(index / 4) * 16;
            if (field.label === 'Status') {
                this.drawStatusChip(document, status.includes('FAILED') ? 'failed' : 'passed', status, x, y + 0.5);
            } else {
                document.setTextColor(36, 43, 56);
                document.setFont('helvetica', 'bold');
                document.setFontSize(9);
                document.text(String(field.value), x, y);
            }
            document.setTextColor(100, 116, 139);
            document.setFont('helvetica', 'normal');
            document.setFontSize(8);
            document.text(field.label, x, y + 5);
        });
    }

    private buildFilename(driveName: string, exportedAt: Date): string {
        const pad = (value: number): string => value.toString().padStart(2, '0');
        const timestamp = [
            exportedAt.getFullYear(),
            pad(exportedAt.getMonth() + 1),
            pad(exportedAt.getDate()),
            pad(exportedAt.getHours()),
            pad(exportedAt.getMinutes()),
        ].join('-');
        const safeDriveName = driveName.trim().replace(/\s+/g, '_');

        return `${timestamp}-scrutiny-drive-report-${safeDriveName}.pdf`;
    }

    private async loadLogo(): Promise<string> {
        const response = await fetch('assets/images/logo/scrutiny-logo-dark-text.png');
        if (!response.ok) {
            throw new Error(`Unable to load Scrutiny logo: ${response.status}`);
        }

        const blob = await response.blob();
        return new Promise<string>((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(reader.result as string);
            reader.onerror = () => reject(reader.error);
            reader.readAsDataURL(blob);
        });
    }
}
