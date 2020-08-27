import {AfterViewInit, Component, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {ApexOptions} from "ng-apexcharts";
import {MatTableDataSource} from "@angular/material/table";
import {MatSort} from "@angular/material/sort";
import {Subject} from "rxjs";
import {DetailService} from "../detail/detail.service";
import {takeUntil} from "rxjs/operators";
import {fadeOut} from "../../../../@treo/animations/fade";

@Component({
  selector: 'detail',
  templateUrl: './detail.component.html',
  styleUrls: ['./detail.component.scss']
})

export class DetailComponent implements OnInit, AfterViewInit, OnDestroy {

    onlyCritical: boolean = true;
    data: any;
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
        private _detailService: DetailService
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
                this.data = data;

                // Store the table data
                this.smartAttributeDataSource.data = this._generateSmartAttributeTableDataSource(data.data.smart_results);

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
        return this.data.lookup[attribute_data.attribute_id]?.description
    }

    getAttributeValue(attribute_data){
        let attribute_metadata = this.data.lookup[attribute_data.attribute_id]
        if(attribute_metadata.display_type == "raw"){
            return  attribute_data.raw_value
        }
        else if(attribute_metadata.display_type == "transformed" && attribute_data.transformed_value ){
            return attribute_data.transformed_value
        } else {
            return attribute_data.value
        }
    }

    getAttributeValueType(attribute_data){
        let attribute_metadata = this.data.lookup[attribute_data.attribute_id]
        return this.data.lookup[attribute_data.attribute_id].display_type
    }

    getAttributeIdeal(attribute_data){
        return this.data.lookup[attribute_data.attribute_id]?.display_type == "raw" ? this.data.lookup[attribute_data.attribute_id]?.ideal : ''
    }

    getAttributeWorst(attribute_data){
        return attribute_data.worst
    }
    getAttributeThreshold(attribute_data){
        return attribute_data.thresh
    }

    private _generateSmartAttributeTableDataSource(smart_results){
        var smartAttributeDataSource = [];

        if(smart_results.length == 0){
            return smartAttributeDataSource
        }

        var latest_smart_result = smart_results[0];
        for(let attr of latest_smart_result.ata_attributes){
            //chart history data
            if (!attr.chartData) {
                var rawHistory = (attr.history || []).map(hist_attr => this.getAttributeValue(hist_attr)).reverse()
                rawHistory.push(this.getAttributeValue(attr))
                attr.chartData = [
                    {
                        name: "chart-line-sparkline",
                        // data: Array.from({length: 40}, () => Math.floor(Math.random() * 40))
                        data: rawHistory
                    }
                ]
            }
            //determine when to include the attributes in table.
            if(!this.onlyCritical || this.onlyCritical && this.data.lookup[attr.attribute_id]?.critical || attr.value <= attr.thresh){
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
        this.smartAttributeDataSource.data = this._generateSmartAttributeTableDataSource(this.data.data.smart_results);

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

    humanizeHours(hours: number): string {
        if(!hours){
            return '--'
        }

        var value: number
        let unit = ""
        if(hours > (24*365)){ //more than a year
            value = Math.round((hours/(24*365)) * 10)/10
            unit = "years"
        } else if (hours > 24){
            value = Math.round((hours/24) *10 )/10
            unit = "days"
        } else{
            value = hours
            unit = "hours"
        }
        return `${value} ${unit}`
    }
}
