/// <reference types="jasmine" />

import {ChangeDetectorRef, ElementRef} from '@angular/core';
import {DetailComponent} from './detail.component';
import {DetailHtmlExportService} from './detail-html-export.service';

describe('DetailComponent HTML export', () => {
    it('exports all attributes and restores the critical-only view', async () => {
        const exportService = jasmine.createSpyObj<DetailHtmlExportService>('DetailHtmlExportService', ['exportReport']);
        let exportedAttributeCount = 0;
        exportService.exportReport.and.callFake(() => {
            exportedAttributeCount = component.smartAttributeDataSource.data.length;
            return Promise.resolve();
        });
        const changeDetector = jasmine.createSpyObj<ChangeDetectorRef>('ChangeDetectorRef', ['detectChanges']);
        const component = new DetailComponent(
            {} as any,
            {} as any,
            {} as any,
            'en-US',
            changeDetector,
            exportService,
        );
        const reportElement = document.createElement('div');
        component.reportContent = new ElementRef(reportElement);
        component.device = {
            device_name: 'sda',
            model_name: 'Test Drive',
            device_protocol: 'ATA',
            device_status: 0,
        } as any;
        component.config = {dashboard_display: 'name'} as any;
        component.metadata = {
            1: {display_name: 'Critical', critical: true},
            2: {display_name: 'Informational', critical: false},
        } as any;
        component.smart_results = [{
            attrs: {
                1: {attribute_id: 1, value: 100, thresh: 1, status: 0, chartData: []},
                2: {attribute_id: 2, value: 100, thresh: 1, status: 0, chartData: []},
            },
        }] as any;
        component.smartAttributeDataSource.data = [component.smart_results[0].attrs[1]];

        await component.exportHtmlReport();

        expect(exportedAttributeCount).toBe(2);
        expect(exportService.exportReport).toHaveBeenCalledWith(
            reportElement,
            'Test Drive',
            'Drive Details - /dev/sda - Test Drive',
        );
        expect(component.onlyCritical).toBeTrue();
        expect(component.smartAttributeDataSource.data.length).toBe(1);
        expect(changeDetector.detectChanges).toHaveBeenCalledTimes(2);
    });
});
