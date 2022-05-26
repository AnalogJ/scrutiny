import { AfterViewInit, ChangeDetectionStrategy, Component, OnDestroy, OnInit, ViewChild, ViewEncapsulation } from '@angular/core';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import {ApexOptions, ChartComponent} from 'ng-apexcharts';
import { DashboardService } from 'app/modules/dashboard/dashboard.service';
import {MatDialog} from '@angular/material/dialog';
import { DashboardSettingsComponent } from 'app/layout/common/dashboard-settings/dashboard-settings.component';
import {deviceDisplayTitle} from "app/layout/common/dashboard-device/dashboard-device.component";
import {AppConfig} from "app/core/config/app.config";
import {TreoConfigService} from "@treo/services/config";
import {Router} from "@angular/router";

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
    hostGroups: { [hostId: string]: string[] } = {}
    temperatureOptions: ApexOptions;
    tempDurationKey: string = "forever"
    config: AppConfig;

    // Private
    private _unsubscribeAll: Subject<any>;
    @ViewChild("tempChart", { static: false }) tempChart: ChartComponent;

    /**
     * Constructor
     *
     * @param {SmartService} _smartService
     */
    constructor(
        private _smartService: DashboardService,
        private _configService: TreoConfigService,
        public dialog: MatDialog,
        private router: Router,
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

        // Subscribe to config changes
        this._configService.config$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((config: AppConfig) => {

                //check if the old config and the new config do not match.
                let oldConfig = JSON.stringify(this.config)
                let newConfig = JSON.stringify(config)

                if(oldConfig != newConfig){
                    console.log(`Configuration updated: ${newConfig} vs ${oldConfig}`)
                    // Store the config
                    this.config = config;

                    if(oldConfig){
                        console.log("reloading component...")
                        this.refreshComponent()
                    }
                }
            });

        // Get the data
        this._smartService.data$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((data) => {

                // Store the data
                this.data = data;

                //generate group data.
                for(let wwn in this.data.data.summary){
                    let hostid = this.data.data.summary[wwn].device.host_id
                    let hostDeviceList = this.hostGroups[hostid] || []
                    hostDeviceList.push(wwn)
                    this.hostGroups[hostid] = hostDeviceList
                }
                console.log(this.hostGroups)

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
    private refreshComponent(){

        let currentUrl = this.router.url;
        this.router.routeReuseStrategy.shouldReuseRoute = () => false;
        this.router.onSameUrlNavigation = 'reload';
        this.router.navigate([currentUrl]);
    }

    private _deviceDataTemperatureSeries() {
        var deviceTemperatureSeries = []

        console.log("DEVICE DATA SUMMARY", this.data)

        for(const wwn in this.data.data.summary){
            var deviceSummary = this.data.data.summary[wwn]
            if (!deviceSummary.temp_history){
                continue
            }

            let deviceName = this.deviceTitle(deviceSummary.device)

            var deviceSeriesMetadata = {
                name: deviceName,
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

    deviceTitle(disk){

        console.log(`Displaying Device ${disk.wwn} with: ${this.config.dashboardDisplay}`)
        let titleParts = []
        if (disk.host_id) titleParts.push(disk.host_id)

        //add device identifier (fallback to generated device name)
        titleParts.push(deviceDisplayTitle(disk, this.config.dashboardDisplay) || deviceDisplayTitle(disk, 'name'))

        return titleParts.join(' - ')
    }

    deviceSummariesForHostGroup(hostGroupWWNs: string[]) {
        let deviceSummaries = []
        for(let wwn of hostGroupWWNs){
            if(this.data.data.summary[wwn]){
                deviceSummaries.push(this.data.data.summary[wwn])
            }
        }
        return deviceSummaries
    }

    openDialog() {
        const dialogRef = this.dialog.open(DashboardSettingsComponent);

        dialogRef.afterClosed().subscribe(result => {
            console.log(`Dialog result: ${result}`);
        });
    }

    onDeviceDeleted(wwn: string) {
        delete this.data.data.summary[wwn] // remove the device from the summary list.
    }

    /*

    DURATION_KEY_WEEK    = "week"
	DURATION_KEY_MONTH   = "month"
	DURATION_KEY_YEAR    = "year"
	DURATION_KEY_FOREVER = "forever"
     */

    changeSummaryTempDuration(durationKey: string){
        this.tempDurationKey = durationKey

        this._smartService.getSummaryTempData(durationKey)
            .subscribe((data) => {

                // given a list of device temp history, override the data in the "summary" object.
                for(const wwn in this.data.data.summary) {
                    // console.log(`Updating ${wwn}, length: ${this.data.data.summary[wwn].temp_history.length}`)
                    this.data.data.summary[wwn].temp_history = data.data.temp_history[wwn] || []
                }

                // Prepare the chart series data
                this.tempChart.updateSeries(this._deviceDataTemperatureSeries())
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
        return item.id || index;
    }

}
