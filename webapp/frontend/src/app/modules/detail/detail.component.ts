import humanizeDuration from 'humanize-duration';
import {AfterViewInit, Component, ElementRef, Inject, LOCALE_ID, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {ApexOptions} from 'ng-apexcharts';
import {AppConfig} from 'app/core/config/app.config';
import {DetailService} from './detail.service';
import {DetailSettingsComponent} from 'app/layout/common/detail-settings/detail-settings.component';
import {MatDialog as MatDialog} from '@angular/material/dialog';
import {MatSort} from '@angular/material/sort';
import {MatTableDataSource as MatTableDataSource} from '@angular/material/table';
import {Subject} from 'rxjs';
import {ScrutinyConfigService} from 'app/core/config/scrutiny-config.service';
import {animate, state, style, transition, trigger} from '@angular/animations';
import {formatDate} from '@angular/common';
import {takeUntil} from 'rxjs/operators';
import {DeviceModel} from 'app/core/models/device-model';
import {SmartModel} from 'app/core/models/measurements/smart-model';
import {SmartAttributeModel} from 'app/core/models/measurements/smart-attribute-model';
import {AttributeMetadataModel} from 'app/core/models/thresholds/attribute-metadata-model';
import {DeviceStatusPipe} from 'app/shared/device-status.pipe';

// from Constants.go - these must match
const AttributeStatusPassed = 0
const AttributeStatusFailedSmart = 1
const AttributeStatusWarningScrutiny = 2
const AttributeStatusFailedScrutiny = 4


@Component({
    selector: 'detail',
    templateUrl: './detail.component.html',
    styleUrls: ['./detail.component.scss'],
    animations: [
        trigger('detailExpand', [
            state('collapsed', style({ height: '0px', minHeight: '0' })),
            state('expanded', style({ height: '*' })),
            transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
        ]),
    ],
    standalone: false
})

export class DetailComponent implements OnInit, AfterViewInit, OnDestroy {

    /**
     * Constructor
     *
     * @param {DetailService} _detailService
     * @param {MatDialog} dialog
     * @param {ScrutinyConfigService} _configService
     * @param {string} locale
     */
    constructor(
        private _detailService: DetailService,
        public dialog: MatDialog,
        private _configService: ScrutinyConfigService,
        @Inject(LOCALE_ID) public locale: string,
        private readonly elementRef: ElementRef
    ) {
        // Set the private defaults
        this._unsubscribeAll = new Subject();

        // Set the defaults
        this.smartAttributeDataSource = new MatTableDataSource();
        // this.recentTransactionsTableColumns = ['status', 'id', 'name', 'value', 'worst', 'thresh'];
        this.smartAttributeTableColumns = ['status', 'id', 'name', 'value', 'worst', 'thresh', 'ideal', 'failure', 'history'];

        this.systemPrefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;

    }

    config: AppConfig;

    onlyCritical = true;
    // data: any;
    expandedAttribute: SmartAttributeModel | null;

    metadata: { [p: string]: AttributeMetadataModel } | { [p: number]: AttributeMetadataModel };
    device: DeviceModel;
    // tslint:disable-next-line:variable-name
    smart_results: SmartModel[];

    commonSparklineOptions: Partial<ApexOptions>;
    smartAttributeDataSource: MatTableDataSource<SmartAttributeModel>;
    smartAttributeTableColumns: string[];

    @ViewChild('smartAttributeTable', {read: MatSort})
    smartAttributeTableMatSort: MatSort;

    // Private
    private _unsubscribeAll: Subject<void>;
    private systemPrefersDark: boolean;

    readonly humanizeDuration = humanizeDuration;

    deviceStatusForModelWithThreshold = DeviceStatusPipe.deviceStatusForModelWithThreshold
    // -----------------------------------------------------------------------------------------------------
    // @ Lifecycle hooks
    // -----------------------------------------------------------------------------------------------------

    /**
     * On init
     */
    ngOnInit(): void {
        // Subscribe to config changes
        this._configService.config$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((config: AppConfig) => {

                this.config = config;
            });

        // Get the data
        this._detailService.data$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((respWrapper) => {

                // Store the data
                // this.data = data;
                this.device = respWrapper.data.device;
                this.smart_results = respWrapper.data.smart_results
                this.metadata = respWrapper.metadata;


                // Store the table data
                this.smartAttributeDataSource.data = this._generateSmartAttributeTableDataSource(this.smart_results);

                // Prepare the chart data
                this._prepareChartData();
            });
    }

    /**
     * After view init
     */
    ngAfterViewInit(): void {
        // Make the data source sortable
        this.smartAttributeDataSource.sort = this.smartAttributeTableMatSort;
    }

    /**
     * On destroy
     */
    ngOnDestroy(): void {
        // Unsubscribe from all subscriptions
        this._unsubscribeAll.next();
        this._unsubscribeAll.complete();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------

    getAttributeStatusName(attributeStatus: number): string {
        // tslint:disable:no-bitwise

        if (attributeStatus === AttributeStatusPassed) {
            return 'passed'

        } else if ((attributeStatus & AttributeStatusFailedScrutiny) !== 0 || (attributeStatus & AttributeStatusFailedSmart) !== 0) {
            return 'failed'
        } else if ((attributeStatus & AttributeStatusWarningScrutiny) !== 0) {
            return 'warn'
        }
        return ''
        // tslint:enable:no-bitwise
    }

    getAttributeScrutinyStatusName(attributeStatus: number): string {
        // tslint:disable:no-bitwise
        if ((attributeStatus & AttributeStatusFailedScrutiny) !== 0) {
            return 'failed'
        } else if ((attributeStatus & AttributeStatusWarningScrutiny) !== 0) {
            return 'warn'
        } else {
            return 'passed'
        }
        // tslint:enable:no-bitwise
    }

    getAttributeSmartStatusName(attributeStatus: number): string {
        // tslint:disable:no-bitwise
        if ((attributeStatus & AttributeStatusFailedSmart) !== 0) {
            return 'failed'
        } else {
            return 'passed'
        }
        // tslint:enable:no-bitwise
    }


    getAttributeName(attributeData: SmartAttributeModel): string {
        const attributeMetadata = this.metadata[attributeData.attribute_id]
        if (!attributeMetadata) {
            return 'Unknown Attribute Name'
        } else {
            return attributeMetadata.display_name
        }
    }

    getAttributeDescription(attributeData: SmartAttributeModel): string {
        const attributeMetadata = this.metadata[attributeData.attribute_id]
        if (!attributeMetadata) {
            return 'Unknown'
        } else {
            return attributeMetadata.description
        }
    }

    getAttributeValue(attributeData: SmartAttributeModel): number {
        if (this.isAta()) {
            const attributeMetadata = this.metadata[attributeData.attribute_id]
            if (!attributeMetadata) {
                return attributeData.value
            } else if (attributeMetadata.display_type === 'raw') {
                // Device statistics (devstat_*) don't have raw_value, use value instead
                return attributeData.raw_value !== undefined ? attributeData.raw_value : attributeData.value
            } else if (attributeMetadata.display_type === 'transformed' && attributeData.transformed_value) {
                return attributeData.transformed_value
            } else {
                return attributeData.value
            }
        } else {
            return attributeData.value
        }
    }

    getAttributeValueType(attributeData: SmartAttributeModel): string {
        if (this.isAta()) {
            const attributeMetadata = this.metadata[attributeData.attribute_id]
            if (!attributeMetadata) {
                return ''
            } else {
                return attributeMetadata.display_type
            }
        } else {
            return ''
        }
    }

    getAttributeIdeal(attributeData: SmartAttributeModel): string {
        if (this.isAta()) {
            return this.metadata[attributeData.attribute_id]?.display_type === 'raw' ? this.metadata[attributeData.attribute_id]?.ideal : ''
        } else {
            return this.metadata[attributeData.attribute_id]?.ideal
        }
    }

    getAttributeWorst(attributeData: SmartAttributeModel): number | string {
        const attributeMetadata = this.metadata[attributeData.attribute_id]
        if (!attributeMetadata) {
            return attributeData.worst
        } else {
            return attributeMetadata?.display_type === 'normalized' ? attributeData.worst : ''
        }
    }

    getAttributeThreshold(attributeData: SmartAttributeModel): number | string {
        if (this.isAta()) {
            const attributeMetadata = this.metadata[attributeData.attribute_id]
            if (!attributeMetadata || attributeMetadata.display_type === 'normalized') {
                return attributeData.thresh
            } else {
                // if(this.data.metadata[attribute_data.attribute_id].observed_thresholds){
                //
                // } else {
                // }
                // return ''
                return attributeData.thresh
            }
        } else {
            return (attributeData.thresh === -1 ? '' : attributeData.thresh)
        }
    }

    getAttributeCritical(attributeData: SmartAttributeModel): boolean {
        return this.metadata[attributeData.attribute_id]?.critical
    }

    getHiddenAttributes(): number {
        if (!this.smart_results || this.smart_results.length === 0) {
            return 0
        }

        let attributesLength = 0
        const attributes = this.smart_results[0]?.attrs
        if (attributes) {
            attributesLength = Object.keys(attributes).length
        }

        return attributesLength - this.smartAttributeDataSource.data.length
    }

    isAta(): boolean {
        return this.device.device_protocol === 'ATA'
    }

    isScsi(): boolean {
        return this.device.device_protocol === 'SCSI'
    }

    isNvme(): boolean {
        return this.device.device_protocol === 'NVMe'
    }

    private _generateSmartAttributeTableDataSource(smartResults: SmartModel[]): SmartAttributeModel[] {
        const smartAttributeDataSource: SmartAttributeModel[] = [];

        if (smartResults.length === 0) {
            return smartAttributeDataSource
        }
        const latestSmartResult = smartResults[0];
        let attributes: { [p: string]: SmartAttributeModel } = {}
        if (this.isScsi()) {
            this.smartAttributeTableColumns = ['status', 'name', 'value', 'thresh', 'history'];
            attributes = latestSmartResult.attrs
        } else if (this.isNvme()) {
            this.smartAttributeTableColumns = ['status', 'name', 'value', 'thresh', 'ideal', 'history'];
            attributes = latestSmartResult.attrs
        } else {
            // ATA
            attributes = latestSmartResult.attrs
            this.smartAttributeTableColumns = ['status', 'id', 'name', 'value', 'thresh', 'ideal', 'failure', 'history'];
        }

        for (const attrId in attributes) {
            const attr = attributes[attrId]

            // chart history data
            if (!attr.chartData) {


                const attrHistory = []
                for (const smartResult of smartResults) {
                    // attrHistory.push(this.getAttributeValue(smart_result.attrs[attrId]))

                    const chartDatapoint = {
                        x: formatDate(smartResult.date, 'MMMM dd, yyyy - HH:mm', this.locale),
                        y: this.getAttributeValue(smartResult.attrs[attrId])
                    }
                    const attributeStatusName = this.getAttributeStatusName(smartResult.attrs[attrId].status)
                    if (attributeStatusName === 'failed') {
                        chartDatapoint['strokeColor'] = '#F05252'
                        chartDatapoint['fillColor'] = '#F05252'
                    } else if (attributeStatusName === 'warn') {
                        chartDatapoint['strokeColor'] = '#C27803'
                        chartDatapoint['fillColor'] = '#C27803'
                    }
                    attrHistory.push(chartDatapoint)
                }

                // var rawHistory = (attr.history || []).map(hist_attr => this.getAttributeValue(hist_attr)).reverse()
                // rawHistory.push(this.getAttributeValue(attr))

                attributes[attrId].chartData = [
                    {
                        name: 'chart-line-sparkline',
                        // attrHistory needs to be reversed, so the newest data is on the right
                        // fixes #339
                        data: attrHistory.reverse()
                    }
                ]
            }
            // determine when to include the attributes in table.

            if (!this.onlyCritical || this.onlyCritical && this.metadata[attr.attribute_id]?.critical || attr.value < attr.thresh) {
                smartAttributeDataSource.push(attr)
            }
        }
        return smartAttributeDataSource
    }

    /**
     * Prepare the chart data from the data
     *
     * @private
     */
    private _prepareChartData(): void {

        // Account balance
        this.commonSparklineOptions = {
            chart: {
                type: 'bar',
                width: 100,
                height: 25,
                sparkline: {
                    enabled: true
                },
                animations: {
                    enabled: false
                },
                events: {
                    mouseMove: () => {
                        const wrapper = this.elementRef.nativeElement.querySelector('.smart-table-wrapper');
                        wrapper?.classList.add('sparkline-hover');
                    },
                    mouseLeave: () => {
                        const wrapper = this.elementRef.nativeElement.querySelector('.smart-table-wrapper');
                        wrapper?.classList.remove('sparkline-hover');
                    }
                }
            },
            // theme:{
            //     // @ts-ignore
            //     // mode:
            //     mode: 'dark',
            // },
            tooltip: {
                fixed: {
                    enabled: false
                },
                x: {
                    show: true
                },
                y: {
                    title: {
                        formatter: (seriesName) => {
                            return '';
                        }
                    }
                },
                marker: {
                    show: false
                },
                theme: this.determineTheme(this.config)

            },
            stroke: {
                width: 2,
                colors: ['#667EEA']
            }
        };
    }

    private determineTheme(config: AppConfig): string {
        if (config.theme === 'system') {
            return this.systemPrefersDark ? 'dark' : 'light'
        } else {
            return config.theme
        }
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    toHex(decimalNumb: number | string): string {
        // Device statistics use string-based IDs like "devstat_7_8"
        // Only convert numeric values to hex
        const num = Number(decimalNumb);
        if (isNaN(num)) {
            return '';
        }
        return '0x' + num.toString(16).padStart(2, '0').toUpperCase()
    }

    formatAttributeId(attributeId: number | string): string {
        // For string-based IDs (device statistics), just return the ID
        if (typeof attributeId === 'string' && isNaN(Number(attributeId))) {
            return attributeId;
        }
        // For numeric IDs, show both decimal and hex
        const hex = this.toHex(attributeId);
        return hex ? `${attributeId} (${hex})` : `${attributeId}`;
    }

    toggleOnlyCritical(): void {
        this.onlyCritical = !this.onlyCritical
        this.smartAttributeDataSource.data = this._generateSmartAttributeTableDataSource(this.smart_results);
    }

    openSettingsDialog(): void {
        const dialogRef = this.dialog.open(DetailSettingsComponent, {
            width: '600px',
            data: {
                curMuted: this.device.muted,
                curLabel: this.device.label
            },
        });

        dialogRef.afterClosed().subscribe((result: undefined | null | { muted: boolean, label: string }) => {
            if (!result) return;

            const promises: Promise<any>[] = [];

            if (result.muted !== this.device.muted) {
                promises.push(this._detailService.setMuted(this.device.wwn, result.muted).toPromise());
            }

            if (result.label !== this.device.label) {
                promises.push(this._detailService.setLabel(this.device.wwn, result.label).toPromise());
            }

            if (promises.length > 0) {
                Promise.all(promises).then(() => {
                    return this._detailService.getData(this.device.wwn).toPromise();
                });
            }
        });
    }

    /**
     * Track by function for ngFor loops
     *
     * @param index
     * @param item
     */
    trackByFn(index: number, item: any): any {
        return index;
        // return item.id || index;
    }

}
