/// <reference types="jasmine" />

import {DetailHtmlExportService} from './detail-html-export.service';

describe('DetailHtmlExportService', () => {
    let service: DetailHtmlExportService;

    beforeEach(() => {
        service = new DetailHtmlExportService(document);
    });

    it('builds a self-contained, non-interactive HTML report from rendered content', () => {
        const source = document.createElement('section');
        source.innerHTML = `
            <h1>Drive Details - /dev/sda</h1>
            <img src="https://example.com/external.png" alt="external">
            <button data-report-exclude>Export</button>
            <script>window.reportScriptRan = true;</script>
            <link rel="stylesheet" href="https://example.com/report.css">
            <div data-report-table style="width: 1000px; min-width: 1000px; height: 800px; min-height: 800px">
                <div style="width: 1000px; height: 800px">
                    <table style="width: 1000px; height: 800px" height="800">
                        <tbody style="width: 1000px; height: 800px">
                            <tr style="width: 1000px; height: 50px"><td style="width: 200px; height: 50px">PASSED</td></tr>
                        </tbody>
                    </table>
                </div>
            </div>
        `;
        source.style.color = 'rgb(36, 43, 56)';
        document.body.appendChild(source);

        const html = service.buildReport(source, {
            title: 'Drive Details - /dev/sda',
            logoDataUrl: 'data:image/png;base64,c2NydXRpbnk=',
        });

        expect(html).toContain('<!doctype html>');
        expect(html).toContain('<style>');
        expect(html).toContain('Drive Details - /dev/sda');
        expect(html).toContain('data:image/png;base64,c2NydXRpbnk=');
        expect(html).toContain('PASSED');
        expect(html).not.toContain('<script');
        expect(html).not.toContain('<link');
        expect(html).not.toContain('<button');
        expect(html).not.toContain('https://example.com');
        expect(html).not.toContain('window.reportScriptRan');
        expect(html).toContain('[data-report-root], [data-report-layout] { display: block; width: 100%; }');

        const exportedDocument = new DOMParser().parseFromString(html, 'text/html');
        expect(exportedDocument.querySelector('.scrutiny-report__content [style]')).toBeNull();
        exportedDocument.querySelectorAll<HTMLElement>('[data-report-table], [data-report-table] *')
            .forEach(element => {
                expect(element.style.width).toBe('');
                expect(element.style.minWidth).toBe('');
                expect(element.style.maxWidth).toBe('');
                expect(element.style.height).toBe('');
                expect(element.style.minHeight).toBe('');
                expect(element.style.maxHeight).toBe('');
                expect(element.hasAttribute('height')).toBeFalse();
            });
        expect(html).toContain('th, td { padding: 8px 10px;');
        expect(html).toContain('[data-report-summary] treo-card > div { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr));');
        expect(html).toContain('[data-report-table] { width: 100%; }');
        expect(html).toContain('break-inside: avoid; page-break-inside: avoid;');

        source.remove();
    });

    it('downloads the report with a sortable timestamp and sanitized drive name', () => {
        const source = document.createElement('section');
        source.textContent = 'report';
        const download = spyOn<any>(service, 'downloadHtml');
        spyOn<any>(service, 'loadLogo').and.resolveTo('data:image/png;base64,c2NydXRpbnk=');

        service.exportReport(source, 'WDC WD140EDFZ-11A0VA0', 'Drive Details', new Date(2026, 7, 18, 3, 45));

        return Promise.resolve().then(() => {
            expect(download).toHaveBeenCalled();
            expect(download.calls.mostRecent().args[1]).toBe(
                '2026-08-18-03-45-scrutiny-drive-report-WDC_WD140EDFZ-11A0VA0.html',
            );
        });
    });
});
