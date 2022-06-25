import { Component, OnInit } from '@angular/core';
import {AppConfig} from 'app/core/config/app.config';
import { TreoConfigService } from '@treo/services/config';
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

    // Private
    private _unsubscribeAll: Subject<any>;

    constructor(
        private _configService: TreoConfigService,
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
              this.dashboardDisplay = config.dashboardDisplay;
              this.dashboardSort = config.dashboardSort;
              this.temperatureUnit = config.temperatureUnit;
              this.theme = config.theme;

          });

  }

  saveSettings(): void {


        const newSettings = {
            dashboardDisplay: this.dashboardDisplay,
            dashboardSort: this.dashboardSort,
            temperatureUnit: this.temperatureUnit,
            theme: this.theme
        }
        this._configService.config = newSettings
        console.log(`Saved Settings: ${JSON.stringify(newSettings)}`)
  }

    formatLabel(value: number): number {
        return value;
    }
}
