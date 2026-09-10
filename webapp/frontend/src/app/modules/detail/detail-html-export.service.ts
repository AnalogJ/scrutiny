import {DOCUMENT} from '@angular/common';
import {Inject, Injectable} from '@angular/core';

interface HtmlReportOptions {
    title: string;
    logoDataUrl: string;
}

@Injectable({providedIn: 'root'})
export class DetailHtmlExportService {
    constructor(@Inject(DOCUMENT) private readonly document: Document) {
    }

    async exportReport(
        source: HTMLElement,
        driveName: string,
        title: string,
        exportedAt: Date = new Date(),
    ): Promise<void> {
        const logoDataUrl = await this.loadLogo();
        const html = this.buildReport(source, {title, logoDataUrl});
        this.downloadHtml(html, this.buildFilename(driveName, exportedAt));
    }

    buildReport(source: HTMLElement, options: HtmlReportOptions): string {
        const reportContent = source.cloneNode(true) as HTMLElement;
        this.removeInteractiveAndExternalContent(reportContent);
        this.stripApplicationPresentation(reportContent);

        return `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>${this.escapeHtml(options.title)}</title>
    <style>${this.reportStyles()}</style>
</head>
<body>
    <main class="scrutiny-report">
        <header class="scrutiny-report__header">
            <img src="${options.logoDataUrl}" alt="Scrutiny">
        </header>
        <div class="scrutiny-report__content">${reportContent.outerHTML}</div>
    </main>
</body>
</html>`;
    }

    private removeInteractiveAndExternalContent(reportContent: HTMLElement): void {
        reportContent.querySelectorAll(
            'script, link, iframe, object, embed, mat-icon, [data-report-exclude], .mat-tooltip',
        ).forEach(element => element.remove());

        reportContent.querySelectorAll('*').forEach(element => {
            Array.from(element.attributes).forEach(attribute => {
                if (attribute.name.startsWith('on')) {
                    element.removeAttribute(attribute.name);
                }
            });

            ['src', 'href', 'poster'].forEach(attributeName => {
                const value = element.getAttribute(attributeName);
                if (value && !value.startsWith('data:') && !value.startsWith('#')) {
                    element.removeAttribute(attributeName);
                }
            });
        });
    }

    private stripApplicationPresentation(reportContent: HTMLElement): void {
        [reportContent, ...Array.from(reportContent.querySelectorAll<HTMLElement>('*'))]
            .forEach(element => {
                Array.from(element.attributes).forEach(attribute => {
                    if (attribute.name === 'style' ||
                        attribute.name === 'width' ||
                        attribute.name === 'height' ||
                        attribute.name.startsWith('_ngcontent-') ||
                        attribute.name.startsWith('_nghost-') ||
                        attribute.name.startsWith('ng-reflect-')) {
                        element.removeAttribute(attribute.name);
                    }
                });
            });
    }

    private async loadLogo(): Promise<string> {
        const response = await fetch('assets/images/logo/scrutiny-logo-dark-text.png');
        if (!response.ok) {
            throw new Error(`Unable to load Scrutiny logo: ${response.status}`);
        }
        const blob = await response.blob();
        return new Promise<string>((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(String(reader.result));
            reader.onerror = () => reject(reader.error);
            reader.readAsDataURL(blob);
        });
    }

    private downloadHtml(html: string, filename: string): void {
        const blob = new Blob([html], {type: 'text/html;charset=utf-8'});
        const url = URL.createObjectURL(blob);
        const anchor = this.document.createElement('a');
        anchor.href = url;
        anchor.download = filename;
        anchor.click();
        URL.revokeObjectURL(url);
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
        const sanitizedDriveName = driveName.trim()
            .replace(/\s+/g, '_')
            .replace(/[^A-Za-z0-9._-]/g, '_');
        return `${timestamp}-scrutiny-drive-report-${sanitizedDriveName}.html`;
    }

    private escapeHtml(value: string): string {
        return value
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');
    }

    private reportStyles(): string {
        return `
            :root { color-scheme: light; }
            * { box-sizing: border-box; }
            body { margin: 0; background: #5145cd; color: #242b38; font-family: Arial, Helvetica, sans-serif; }
            .scrutiny-report { width: min(1200px, calc(100% - 48px)); margin: 24px auto; background: #fff; border-radius: 12px; overflow: hidden; }
            .scrutiny-report__header { min-height: 96px; display: flex; align-items: center; padding: 24px 40px; border-bottom: 1px solid #e2e8f0; }
            .scrutiny-report__header img { display: block; width: 180px; height: auto; }
            .scrutiny-report__content { padding: 24px 40px 40px; }
            .scrutiny-report [data-report-exclude] { display: none !important; }
            [data-report-root], [data-report-layout] { display: block; width: 100%; }
            [data-report-heading] { width: 100%; margin: 4px 0 24px; }
            [data-report-heading] h2 { margin: 0; font-size: 28px; line-height: 1.2; }
            [data-report-heading] .text-secondary { margin-top: 4px; color: #64748b; }
            [data-report-summary] { width: 100%; margin-bottom: 28px; }
            [data-report-summary] treo-card { display: block; padding: 18px 20px; border: 1px solid #e2e8f0; border-radius: 8px; }
            [data-report-summary] treo-card > div { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 16px 24px; }
            [data-report-summary] treo-card > div > div { min-width: 0; overflow-wrap: anywhere; }
            [data-report-summary] treo-card > div > div > div:first-child { font-weight: 600; }
            [data-report-summary] .text-secondary { margin-top: 2px; color: #64748b; font-size: 12px; font-weight: 400; }
            [data-report-table] { width: 100%; }
            [data-report-table] > div { width: 100%; border: 1px solid #e2e8f0; border-radius: 8px; overflow: hidden; }
            [data-report-table] .p-6 { padding: 16px 20px; }
            [data-report-table] .font-bold { font-weight: 700; }
            [data-report-table] .uppercase { text-transform: uppercase; }
            [data-report-table] .tracking-wider { letter-spacing: .04em; }
            table { width: 100%; border-collapse: collapse; table-layout: auto; }
            thead { display: table-header-group; }
            tbody { display: table-row-group; }
            tfoot { display: table-footer-group; }
            tr { display: table-row; }
            th, td { padding: 8px 10px; border-bottom: 1px solid #e2e8f0; text-align: left; vertical-align: middle; line-height: 1.25; }
            th { background: #f8fafc; color: #64748b; font-size: 12px; font-weight: 700; white-space: nowrap; }
            td { color: #242b38; font-size: 13px; }
            .mat-column-status { width: 92px; }
            .mat-column-id { width: 105px; white-space: nowrap; }
            .mat-column-value, .mat-column-thresh, .mat-column-ideal, .mat-column-failure { width: 86px; white-space: nowrap; }
            .mat-column-status > span,
            [data-report-summary] .inline-flex { display: inline-flex; align-items: center; padding: 2px 8px; border-radius: 999px; font-size: 11px; font-weight: 700; line-height: 1.3; text-transform: uppercase; white-space: nowrap; }
            .mat-column-status > span > span:first-child,
            [data-report-summary] .inline-flex > span:first-child { width: 7px; height: 7px; margin-right: 6px; border-radius: 50%; flex: 0 0 auto; }
            .green-200 { background: #bcf0da; color: #03543f; }
            .red-200 { background: #fbd5d5; color: #9b1c1c; }
            .yellow-200 { background: #fce96a; color: #713f12; }
            .bg-green { background: #0e9f6e; }
            .bg-red { background: #f05252; }
            .bg-yellow { background: #d69e2e; }
            .yellow-50 { background: #fdfdea; }
            @media (max-width: 800px) {
                .scrutiny-report { width: min(100% - 24px, 1200px); margin: 12px auto; }
                .scrutiny-report__content { padding: 20px; }
                [data-report-summary] treo-card > div { grid-template-columns: repeat(2, minmax(0, 1fr)); }
            }
            @media (max-width: 520px) {
                [data-report-summary] treo-card > div { grid-template-columns: 1fr; }
            }
            @media print {
                @page { size: A4 landscape; margin: 10mm; }
                body { background: #fff; print-color-adjust: exact; -webkit-print-color-adjust: exact; }
                .scrutiny-report { width: 100%; margin: 0; border-radius: 0; overflow: visible; }
                .scrutiny-report__header { min-height: 0; padding: 4mm 0; }
                .scrutiny-report__header img { width: 42mm; }
                .scrutiny-report__content { padding: 4mm 0 0; }
                [data-report-heading] { margin: 0 0 5mm; }
                [data-report-heading] h2 { font-size: 20px; }
                [data-report-summary] { margin-bottom: 5mm; break-inside: avoid; page-break-inside: avoid; }
                [data-report-summary] treo-card { padding: 3mm 4mm; }
                [data-report-summary] treo-card > div { grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 2.5mm 5mm; }
                [data-report-table], [data-report-table] > div { overflow: visible; border-radius: 0; }
                [data-report-table] .p-6 { padding: 2mm 0; }
                thead { display: table-header-group; }
                tfoot { display: table-footer-group; }
                tr { break-inside: avoid; page-break-inside: avoid; }
                th, td { padding: 1.2mm 1.5mm; font-size: 9px; }
            }
        `;
    }
}
