import {Component, OnInit} from '@angular/core';
import {
    AppConfig,
    DashboardDisplay,
    DashboardSort,
    MetricsStatusFilterAttributes,
    MetricsStatusThreshold,
    TemperatureUnit,
    Theme
} from 'app/core/config/app.config';
import {ScrutinyConfigService} from 'app/core/config/scrutiny-config.service';
import {Subject} from 'rxjs';
import {takeUntil} from 'rxjs/operators';

@Component({
    selector: 'app-dashboard-settings',
    templateUrl: './dashboard-settings.component.html',
    styleUrls: ['./dashboard-settings.component.scss']
})
export class DashboardSettingsComponent implements OnInit {

    dashboardDisplay: string;
    dashboardSort: string;
    temperatureUnit: string;
    theme: string;
    statusThreshold: number;
    statusFilterAttributes: number;

    // Private
    private _unsubscribeAll: Subject<any>;

    constructor(
        private _configService: ScrutinyConfigService,
    ) {
        // Set the private defaults
        this._unsubscribeAll = new Subject();
    }

    ngOnInit(): void {
        // Subscribe to config changes
        this._configService.config$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((config: AppConfig) => {

                // Store the config
                this.dashboardDisplay = config.dashboard_display;
                this.dashboardSort = config.dashboard_sort;
                this.temperatureUnit = config.temperature_unit;
                this.theme = config.theme;

                this.statusFilterAttributes = config.metrics.status_filter_attributes;
                this.statusThreshold = config.metrics.status_threshold;

            });

    }

    saveSettings(): void {
        const newSettings: AppConfig = {
            dashboard_display: this.dashboardDisplay as DashboardDisplay,
            dashboard_sort: this.dashboardSort as DashboardSort,
            temperature_unit: this.temperatureUnit as TemperatureUnit,
            theme: this.theme as Theme,
            metrics: {
                status_filter_attributes: this.statusFilterAttributes as MetricsStatusFilterAttributes,
                status_threshold: this.statusThreshold as MetricsStatusThreshold
            }
        }
        this._configService.config = newSettings
        console.log(`Saved Settings: ${JSON.stringify(newSettings)}`)
    }

    formatLabel(value: number): number {
        return value;
    }
}
