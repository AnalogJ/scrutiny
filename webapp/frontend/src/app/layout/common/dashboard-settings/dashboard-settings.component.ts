import {Component, OnInit} from '@angular/core';
import {
    AppConfig,
    DashboardDisplay,
    DashboardSort,
    MetricsStatusFilterAttributes,
    MetricsStatusThreshold,
    TemperatureUnit,
    LineStroke,
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
    fileSizeSIUnits: boolean;
    lineStroke: string;
    theme: string;
    statusThreshold: number;
    statusFilterAttributes: number;
    repeatNotifications: boolean;

    // Private
    private _unsubscribeAll: Subject<void>;

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
                this.fileSizeSIUnits = config.file_size_si_units;
                this.lineStroke = config.line_stroke;
                this.theme = config.theme;

                this.statusFilterAttributes = config.metrics.status_filter_attributes;
                this.statusThreshold = config.metrics.status_threshold;
                this.repeatNotifications = config.metrics.repeat_notifications;

            });

    }

    saveSettings(): void {
        const newSettings: AppConfig = {
            dashboard_display: this.dashboardDisplay as DashboardDisplay,
            dashboard_sort: this.dashboardSort as DashboardSort,
            temperature_unit: this.temperatureUnit as TemperatureUnit,
            file_size_si_units: this.fileSizeSIUnits,
            line_stroke: this.lineStroke as LineStroke,
            theme: this.theme as Theme,
            metrics: {
                status_filter_attributes: this.statusFilterAttributes as MetricsStatusFilterAttributes,
                status_threshold: this.statusThreshold as MetricsStatusThreshold,
                repeat_notifications: this.repeatNotifications
            }
        }
        this._configService.config = newSettings
        console.log(`Saved Settings: ${JSON.stringify(newSettings)}`)
    }

    formatLabel(value: number): number {
        return value;
    }
}
