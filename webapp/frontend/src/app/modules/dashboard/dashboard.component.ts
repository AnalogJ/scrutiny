import { AfterViewInit, ChangeDetectionStrategy, Component, OnDestroy, OnInit, ViewChild, ViewEncapsulation } from '@angular/core';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ApexOptions } from 'ng-apexcharts';
import { DashboardService } from 'app/modules/dashboard/dashboard.service';
import * as moment from "moment";
import {MatDialog} from '@angular/material/dialog';
import { DashboardSettingsComponent } from 'app/layout/common/dashboard-settings/dashboard-settings.component';
import  humanizeDuration from 'humanize-duration'

@Component({
    selector       : 'example',
    templateUrl    : './dashboard.component.html',
    styleUrls      : ['./dashboard.component.scss'],
    encapsulation  : ViewEncapsulation.None,
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class DashboardComponent implements OnInit, AfterViewInit, OnDestroy
{
    data: any;
    temperatureOptions: ApexOptions;

    // Private
    private _unsubscribeAll: Subject<any>;

    /**
     * Constructor
     *
     * @param {SmartService} _smartService
     */
    constructor(
        private _smartService: DashboardService,
        public dialog: MatDialog
    )
    {
        // Set the private defaults
        this._unsubscribeAll = new Subject();

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
        this._smartService.data$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((data) => {

                // Store the data
                this.data = data;

                // Prepare the chart data
                this._prepareChartData();
            });
    }

    /**
     * After view init
     */
    ngAfterViewInit(): void
    {}

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
    private _deviceDataTemperatureSeries() {
        var deviceTemperatureSeries = []

        console.log("DEVICE DATA SUMMARY", this.data)

        for(const wwn in this.data.data.summary){
            var deviceSummary = this.data.data.summary[wwn]
            if (!deviceSummary.temp_history){
                continue
            }
            var deviceSeriesMetadata = {
                name: `/dev/${deviceSummary.device.device_name}`,
                data: []
            }

            for(let tempHistory of deviceSummary.temp_history){
                let newDate = new Date(tempHistory.date);
                deviceSeriesMetadata.data.push({
                    x: newDate,
                    y: tempHistory.temp
                })
            }
            deviceTemperatureSeries.push(deviceSeriesMetadata)
        }
        return deviceTemperatureSeries
    }
    /**
     * Prepare the chart data from the data
     *
     * @private
     */
    private _prepareChartData(): void
    {
        // Account balance
        this.temperatureOptions = {
            chart  : {
                animations: {
                    speed           : 400,
                    animateGradually: {
                        enabled: false
                    }
                },
                fontFamily: 'inherit',
                foreColor : 'inherit',
                width     : '100%',
                height    : '100%',
                type      : 'area',
                sparkline : {
                    enabled: true
                }
            },
            colors : ['#A3BFFA', '#667EEA'],
            fill   : {
                colors : ['#CED9FB', '#AECDFD'],
                opacity: 0.5,
                type   : 'solid'
            },
            series : this._deviceDataTemperatureSeries(),
            stroke : {
                curve: 'straight',
                width: 2
            },
            tooltip: {
                theme: 'dark',
                x    : {
                    format: 'MMM dd, yyyy HH:mm:ss'
                },
                y    : {
                    formatter: (value) => {
                        return value + '°C';
                    }
                }
            },
            xaxis  : {
                type: 'datetime'
            }
        };
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    openDialog() {
        const dialogRef = this.dialog.open(DashboardSettingsComponent);

        dialogRef.afterClosed().subscribe(result => {
            console.log(`Dialog result: ${result}`);
        });
    }

    deviceTitle(disk){
        let title = []

        if (disk.host_id) title.push(disk.host_id)

        title.push(`/dev/${disk.device_name}`)

        if (disk.device_type && disk.device_type != 'scsi' && disk.device_type != 'ata'){
            title.push(disk.device_type)
        }

        title.push(disk.model_name)

        return title.join(' - ')
    }

    deviceStatusString(deviceStatus){
        if(deviceStatus == 0){
            return "passed"
        } else {
            return "failed"
        }
    }

    /**
     * Track by function for ngFor loops
     *
     * @param index
     * @param item
     */
    trackByFn(index: number, item: any): any
    {
        return item.id || index;
    }

    readonly humanizeDuration = humanizeDuration;

}
