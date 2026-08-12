import {
    AfterViewInit,
    ChangeDetectionStrategy,
    Component,
    ElementRef,
    NgZone,
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

export interface Box {
    top: number;
    left: number;
    bottom: number;
    right: number;
}

/**
 * How far a tooltip has to move to sit inside the viewport. apexcharts positions the shared tooltip
 * against the plot area and only keeps its top half inside, so once there is a row per drive it is
 * taller than the plot area and hangs off the bottom of the card.
 */
export function tooltipViewportOffset(box: Box, viewport: { width: number; height: number },
                                      margin: number = 8): { dx: number; dy: number } {
    // pull it back from the far edges first, then off the near ones, so a tooltip too big to fit
    // lands against the top left rather than being pushed further off screen
    let dx = Math.min(0, viewport.width - margin - box.right);
    let dy = Math.min(0, viewport.height - margin - box.bottom);
    dx += Math.max(0, margin - (box.left + dx));
    dy += Math.max(0, margin - (box.top + dy));

    return {dx, dy};
}

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
    tempDurationKey = 'week'
    config: AppConfig;
    showArchived: boolean;

    // Private
    private _unsubscribeAll: Subject<void>;
    private _tooltipObserver: MutationObserver;
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
        private _elementRef: ElementRef,
        private _ngZone: NgZone,
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
                for (const scrutiny_uuid in this.summaryData) {
                    const hostid = this.summaryData[scrutiny_uuid].device.host_id
                    const hostDeviceList = this.hostGroups[hostid] || []
                    hostDeviceList.push(scrutiny_uuid)
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
    {
        this._ngZone.runOutsideAngular(() => {
            this._tooltipObserver = new MutationObserver((mutations) => {
                for (const mutation of mutations) {
                    const target = mutation.target as HTMLElement;
                    if (target.classList.contains('apexcharts-tooltip')) {
                        this.nudgeTooltipIntoView(target);
                    }
                }
            });
            this._tooltipObserver.observe(this._elementRef.nativeElement, {
                attributes: true,
                attributeFilter: ['style'],
                subtree: true
            });
        });
    }

    /**
     * On destroy
     */
    ngOnDestroy(): void
    {
        this._tooltipObserver?.disconnect();

        // Unsubscribe from all subscriptions
        this._unsubscribeAll.next();
        this._unsubscribeAll.complete();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------
    private nudgeTooltipIntoView(tooltip: HTMLElement): void {
        const origin = (tooltip.offsetParent as HTMLElement)?.getBoundingClientRect();
        if (!origin || !tooltip.style.left || !tooltip.style.top) {
            return;
        }

        // .apexcharts-tooltip animates every move (`transition: .15s ease all`), so reading it back
        // with getBoundingClientRect() gives a position it is still travelling away from and the
        // correction below never settles. Measure where apexcharts asked for it to go instead.
        const left = parseFloat(tooltip.style.left);
        const top = parseFloat(tooltip.style.top);

        const {dx, dy} = tooltipViewportOffset(
            {
                left: origin.left + left,
                top: origin.top + top,
                right: origin.left + left + tooltip.offsetWidth,
                bottom: origin.top + top + tooltip.offsetHeight
            },
            {width: document.documentElement.clientWidth, height: document.documentElement.clientHeight}
        );

        if (dx === 0 && dy === 0) {
            return;
        }

        tooltip.style.left = `${left + dx}px`;
        tooltip.style.top = `${top + dy}px`;
    }

    private refreshComponent(): void {

        const currentUrl = this.router.url;
        this.router.routeReuseStrategy.shouldReuseRoute = () => false;
        this.router.onSameUrlNavigation = 'reload';
        this.router.navigate([currentUrl]);
    }

    private _deviceDataTemperatureSeries(): any[] {
        const deviceTemperatureSeries = []

        console.log('DEVICE DATA SUMMARY', this.summaryData)

        for (const scrutiny_uuid in this.summaryData) {
            const deviceSummary = this.summaryData[scrutiny_uuid]
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
                let temperature;
                switch (this.config.temperature_unit) {
                    case 'celsius':
                        temperature = tempHistory.temp;
                        break
                    case 'fahrenheit':
                        temperature = TemperaturePipe.celsiusToFahrenheit(tempHistory.temp)
                        break
                }
                deviceSeriesMetadata.data.push({
                    x: newDate,
                    y: temperature
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
            colors : ['#667eea', '#9066ea', '#66c0ea', '#66ead2', '#d266ea', '#66ea90'],
            fill   : {
                colors : ['#b2bef4', '#c7b2f4', '#b2dff4', '#b2f4e8', '#e8b2f4', '#b2f4c7'],
                opacity: 0.5,
                type   : 'gradient'
            },
            series : this._deviceDataTemperatureSeries(),
            stroke : {
                curve: this.config.line_stroke,
                width: 2
            },
            tooltip: {
                theme: 'dark',
                shared: true,
                intersect: false,
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
                type: 'datetime',
                labels: {
                    datetimeUTC: false
                }
            }
        };
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    deviceSummariesForHostGroup(hostGroupScrutinyUUIDs: string[]): DeviceSummaryModel[] {
        const deviceSummaries: DeviceSummaryModel[] = []
        for (const scrutiny_uuid of hostGroupScrutinyUUIDs) {
            if (this.summaryData[scrutiny_uuid]) {
                deviceSummaries.push(this.summaryData[scrutiny_uuid])
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

    onDeviceDeleted(scrutiny_uuid: string): void {
        delete this.summaryData[scrutiny_uuid] // remove the device from the summary list.
    }

    onDeviceArchived(scrutiny_uuid: string): void {
        this.summaryData[scrutiny_uuid].device.archived = true;
    }

    onDeviceUnarchived(scrutiny_uuid: string): void {
        this.summaryData[scrutiny_uuid].device.archived = false;
    }

    /*
    DURATION_KEY_DAY    = "day"
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
                for (const scrutiny_uuid in this.summaryData) {
                    // console.log(`Updating ${scrutiny_uuid}, length: ${this.data.data.summary[scrutiny_uuid].temp_history.length}`)
                    this.summaryData[scrutiny_uuid].temp_history = tempHistoryData[scrutiny_uuid] || []
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
