import {
    AfterViewInit,
    ChangeDetectionStrategy,
    Component,
    OnDestroy,
    OnInit,
    ViewChild,
    ViewEncapsulation
} from '@angular/core';
import {Subject} from 'rxjs';
import {takeUntil} from 'rxjs/operators';
import {ApexOptions, ChartComponent} from 'ng-apexcharts';
import {DashboardService} from 'app/modules/dashboard/dashboard.service';
import {MatDialog} from '@angular/material/dialog';
import {DashboardSettingsComponent} from 'app/layout/common/dashboard-settings/dashboard-settings.component';
import {AppConfig} from 'app/core/config/app.config';
import {ScrutinyConfigService} from 'app/core/config/scrutiny-config.service';
import {Router} from '@angular/router';
import {TemperaturePipe} from 'app/shared/temperature.pipe';
import {DeviceTitlePipe} from 'app/shared/device-title.pipe';
import {DeviceSummaryModel} from 'app/core/models/device-summary-model';

@Component({
    selector       : 'example',
    templateUrl    : './dashboard.component.html',
    styleUrls      : ['./dashboard.component.scss'],
    encapsulation  : ViewEncapsulation.None,
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class DashboardComponent implements OnInit, AfterViewInit, OnDestroy
{
    summaryData: { [key: string]: DeviceSummaryModel };
    hostGroups: { [hostId: string]: string[] } = {}
    temperatureOptions: ApexOptions;
    tempDurationKey = 'forever'
    config: AppConfig;

    // Private
    private _unsubscribeAll: Subject<any>;
    @ViewChild('tempChart', { static: false }) tempChart: ChartComponent;

    /**
     * Constructor
     *
     * @param {DashboardService} _dashboardService
     * @param {ScrutinyConfigService} _configService
     * @param {MatDialog} dialog
     * @param {Router} router
     */
    constructor(
        private _dashboardService: DashboardService,
        private _configService: ScrutinyConfigService,
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

                // check if the old config and the new config do not match.
                const oldConfig = JSON.stringify(this.config)
                const newConfig = JSON.stringify(config)

                if(oldConfig !== newConfig){
                    console.log(`Configuration updated: ${newConfig} vs ${oldConfig}`)
                    // Store the config
                    this.config = config;

                    if(oldConfig){
                        console.log('reloading component...')
                        this.refreshComponent()
                    }
                }
            });

        // Get the data
        this._dashboardService.data$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((data) => {

                // Store the data
                this.summaryData = data;

                // generate group data.
                for (const wwn in this.summaryData) {
                    const hostid = this.summaryData[wwn].device.host_id
                    const hostDeviceList = this.hostGroups[hostid] || []
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
    private refreshComponent(): void {

        const currentUrl = this.router.url;
        this.router.routeReuseStrategy.shouldReuseRoute = () => false;
        this.router.onSameUrlNavigation = 'reload';
        this.router.navigate([currentUrl]);
    }

    private _deviceDataTemperatureSeries(): any[] {
        const deviceTemperatureSeries = []

        console.log('DEVICE DATA SUMMARY', this.summaryData)

        for (const wwn in this.summaryData) {
            const deviceSummary = this.summaryData[wwn]
            if (!deviceSummary.temp_history) {
                continue
            }

            const deviceName = DeviceTitlePipe.deviceTitleWithFallback(deviceSummary.device, this.config.dashboard_display)

            const deviceSeriesMetadata = {
                name: deviceName,
                data: []
            }

            for(const tempHistory of deviceSummary.temp_history){
                const newDate = new Date(tempHistory.date);
                deviceSeriesMetadata.data.push({
                    x: newDate,
                    y: TemperaturePipe.formatTemperature(tempHistory.temp, this.config.temperature_unit, false)
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
                        return TemperaturePipe.formatTemperature(value, this.config.temperature_unit, true) as string;
                    }
                }
            },
            xaxis: {
                type: 'datetime'
            }
        };
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    deviceSummariesForHostGroup(hostGroupWWNs: string[]): DeviceSummaryModel[] {
        const deviceSummaries: DeviceSummaryModel[] = []
        for (const wwn of hostGroupWWNs) {
            if (this.summaryData[wwn]) {
                deviceSummaries.push(this.summaryData[wwn])
            }
        }
        return deviceSummaries
    }

    openDialog(): void {
        const dialogRef = this.dialog.open(DashboardSettingsComponent, {width: '600px',});

        dialogRef.afterClosed().subscribe(result => {
            console.log(`Dialog result: ${result}`);
        });
    }

    onDeviceDeleted(wwn: string): void {
        delete this.summaryData[wwn] // remove the device from the summary list.
    }

    /*

    DURATION_KEY_WEEK    = "week"
	DURATION_KEY_MONTH   = "month"
	DURATION_KEY_YEAR    = "year"
	DURATION_KEY_FOREVER = "forever"
     */

    changeSummaryTempDuration(durationKey: string): void {
        this.tempDurationKey = durationKey

        this._dashboardService.getSummaryTempData(durationKey)
            .subscribe((tempHistoryData) => {

                // given a list of device temp history, override the data in the "summary" object.
                for (const wwn in this.summaryData) {
                    // console.log(`Updating ${wwn}, length: ${this.data.data.summary[wwn].temp_history.length}`)
                    this.summaryData[wwn].temp_history = tempHistoryData[wwn] || []
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
