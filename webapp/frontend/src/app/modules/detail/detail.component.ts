import {AfterViewInit, Component, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {ApexOptions} from "ng-apexcharts";
import {MatTableDataSource} from "@angular/material/table";
import {MatSort} from "@angular/material/sort";
import {Subject} from "rxjs";
import {DetailService} from "./detail.service";
import {takeUntil} from "rxjs/operators";
import {fadeOut} from "../../../@treo/animations/fade";
import {DetailSettingsComponent} from "app/layout/common/detail-settings/detail-settings.component";
import {MatDialog} from "@angular/material/dialog";
import humanizeDuration from 'humanize-duration';

@Component({
  selector: 'detail',
  templateUrl: './detail.component.html',
  styleUrls: ['./detail.component.scss']
})

export class DetailComponent implements OnInit, AfterViewInit, OnDestroy {

    onlyCritical: boolean = true;
    // data: any;

    metadata: any;
    device: any;
    smart_results: any[];

    commonSparklineOptions: Partial<ApexOptions>;
    smartAttributeDataSource: MatTableDataSource<any>;
    smartAttributeTableColumns: string[];


    @ViewChild('smartAttributeTable', {read: MatSort})
    smartAttributeTableMatSort: MatSort;

    // Private
    private _unsubscribeAll: Subject<any>;

    /**
     * Constructor
     *
     * @param {DetailService} _detailService
     */
    constructor(
        private _detailService: DetailService,
        public dialog: MatDialog

    )
    {
        // Set the private defaults
        this._unsubscribeAll = new Subject();

        // Set the defaults
        this.smartAttributeDataSource = new MatTableDataSource();
        // this.recentTransactionsTableColumns = ['status', 'id', 'name', 'value', 'worst', 'thresh'];
        this.smartAttributeTableColumns = ['status', 'id', 'name', 'value', 'worst', 'thresh','ideal', 'failure', 'history'];
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Lifecycle hooks
    // -----------------------------------------------------------------------------------------------------

    /**
     * On init
     */
    ngOnInit(): void
    {
        // Get the data
        this._detailService.data$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((data) => {

                // Store the data
                // this.data = data;
                this.device = data.data.device;
                this.smart_results = data.data.smart_results
                this.metadata = data.metadata;


                // Store the table data
                this.smartAttributeDataSource.data = this._generateSmartAttributeTableDataSource(this.smart_results);

                // Prepare the chart data
                this._prepareChartData();
            });
    }

    /**
     * After view init
     */
    ngAfterViewInit(): void
    {
        // Make the data source sortable
        this.smartAttributeDataSource.sort = this.smartAttributeTableMatSort;
    }

    /**
     * On destroy
     */
    ngOnDestroy(): void
    {
        // Unsubscribe from all subscriptions
        this._unsubscribeAll.next();
        this._unsubscribeAll.complete();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------
    getAttributeDescription(attribute_data){
        let attribute_metadata = this.metadata[attribute_data.attribute_id]
        if(!attribute_metadata){
            return 'Unknown'
        } else {
            return attribute_metadata.description
        }
        return
    }

    getAttributeValue(attribute_data){
        if(this.isAta()) {
            let attribute_metadata = this.metadata[attribute_data.attribute_id]
            if(!attribute_metadata){
                return attribute_data.value
            } else if (attribute_metadata.display_type == "raw") {
                return attribute_data.raw_value
            } else if (attribute_metadata.display_type == "transformed" && attribute_data.transformed_value) {
                return attribute_data.transformed_value
            } else {
                return attribute_data.value
            }
        }
        else{
            return attribute_data.value
        }
    }

    getAttributeValueType(attribute_data){
        if(this.isAta()) {
            let attribute_metadata = this.metadata[attribute_data.attribute_id]
            if(!attribute_metadata){
                return ''
            } else {
                return attribute_metadata.display_type
            }
        } else {
            return ''
        }
    }

    getAttributeIdeal(attribute_data){
        if(this.isAta()){
            return this.metadata[attribute_data.attribute_id]?.display_type == "raw" ? this.metadata[attribute_data.attribute_id]?.ideal : ''
        } else {
            return this.metadata[attribute_data.attribute_id]?.ideal
        }
    }

    getAttributeWorst(attribute_data){
        let attribute_metadata = this.metadata[attribute_data.attribute_id]
        if(!attribute_metadata){
            return attribute_data.worst
        } else {
            return attribute_metadata?.display_type == "normalized" ? attribute_data.worst : ''
        }
    }

    getAttributeThreshold(attribute_data){
        if(this.isAta()){
            let attribute_metadata = this.metadata[attribute_data.attribute_id]
            if(!attribute_metadata || attribute_metadata.display_type == "normalized"){
                return attribute_data.thresh
            } else {
                // if(this.data.metadata[attribute_data.attribute_id].observed_thresholds){
                //
                // } else {
                // }
                // return ''
                return attribute_data.thresh
            }
        } else {
            return (attribute_data.thresh == -1 ? '' : attribute_data.thresh )
        }
    }

    getAttributeCritical(attribute_data){
        return this.metadata[attribute_data.attribute_id]?.critical
    }
    getHiddenAttributes(){
        if (!this.smart_results || this.smart_results.length == 0) {
            return 0
        }

        let attributes_length = 0
        let attributes = this.smart_results[0]?.attrs
        if (attributes) {
            attributes_length = Object.keys(attributes).length
        }

        return attributes_length - this.smartAttributeDataSource.data.length
    }

    isAta(): boolean {
        return this.device.device_protocol == 'ATA'
    }
    isScsi(): boolean {
        return this.device.device_protocol == 'SCSI'
    }
    isNvme(): boolean {
        return this.device.device_protocol == 'NVMe'
    }

    private _generateSmartAttributeTableDataSource(smart_results){
        var smartAttributeDataSource = [];

        if(smart_results.length == 0){
            return smartAttributeDataSource
        }
        var latest_smart_result = smart_results[0];
        let attributes = {}
        if(this.isScsi()) {
            this.smartAttributeTableColumns = ['status', 'name', 'value', 'thresh', 'history'];
            attributes = latest_smart_result.attrs
        } else if(this.isNvme()){
            this.smartAttributeTableColumns = ['status', 'name', 'value', 'thresh', 'ideal', 'history'];
            attributes = latest_smart_result.attrs
        } else {
            //ATA
            attributes = latest_smart_result.attrs
            this.smartAttributeTableColumns = ['status', 'id', 'name', 'value', 'worst', 'thresh','ideal', 'failure', 'history'];
        }

        for(const attrId in attributes){
            var attr = attributes[attrId]

            //chart history data
            if (!attr.chartData) {
                var rawHistory = (attr.history || []).map(hist_attr => this.getAttributeValue(hist_attr)).reverse()
                rawHistory.push(this.getAttributeValue(attr))
                attr.chartData = [
                    {
                        name: "chart-line-sparkline",
                        data: rawHistory
                    }
                ]
            }
            //determine when to include the attributes in table.

            if(!this.onlyCritical || this.onlyCritical && this.metadata[attr.attribute_id]?.critical || attr.value < attr.thresh){
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
    private _prepareChartData(): void
    {

        // Account balance
        this.commonSparklineOptions = {
            chart: {
                type: "bar",
                width: 100,
                height: 25,
                sparkline: {
                    enabled: true
                },
                animations: {
                    enabled: false
                }
            },
            tooltip: {
                fixed: {
                    enabled: false
                },
                x: {
                    show: false
                },
                y: {
                    title: {
                        formatter: function(seriesName) {
                            return "";
                        }
                    }
                },
                marker: {
                    show: false
                }
            },
            stroke: {
                width: 2,
                colors: ['#667EEA']
            }
        };
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    toHex(decimalNumb){
        return "0x" + Number(decimalNumb).toString(16).padStart(2, '0').toUpperCase()
    }
    toggleOnlyCritical(){
        this.onlyCritical = !this.onlyCritical
        this.smartAttributeDataSource.data = this._generateSmartAttributeTableDataSource(this.smart_results);

    }

    openDialog() {
        const dialogRef = this.dialog.open(DetailSettingsComponent);

        dialogRef.afterClosed().subscribe(result => {
            console.log(`Dialog result: ${result}`);
        });
    }

    /**
     * Track by function for ngFor loops
     *
     * @param index
     * @param item
     */
    trackByFn(index: number, item: any): any
    {
        return index;
        // return item.id || index;
    }

    readonly humanizeDuration = humanizeDuration;

}
