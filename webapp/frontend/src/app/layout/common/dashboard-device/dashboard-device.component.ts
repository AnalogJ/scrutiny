import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import * as moment from 'moment';
import {takeUntil} from 'rxjs/operators';
import {AppConfig} from 'app/core/config/app.config';
import {ScrutinyConfigService} from 'app/core/config/scrutiny-config.service';
import {Subject} from 'rxjs';
import humanizeDuration from 'humanize-duration'
import {MatDialog} from '@angular/material/dialog';
import {DashboardDeviceDeleteDialogComponent} from 'app/layout/common/dashboard-device-delete-dialog/dashboard-device-delete-dialog.component';
import {DeviceTitlePipe} from 'app/shared/device-title.pipe';
import {DeviceSummaryModel} from 'app/core/models/device-summary-model';

export type deviceStatusName = 'unknown' | 'passed' | 'failed'

@Component({
    selector: 'app-dashboard-device',
    templateUrl: './dashboard-device.component.html',
    styleUrls: ['./dashboard-device.component.scss']
})
export class DashboardDeviceComponent implements OnInit {

    constructor(
        private _configService: ScrutinyConfigService,
        public dialog: MatDialog,
    ) {
        // Set the private defaults
        this._unsubscribeAll = new Subject();
    }

    @Input() deviceSummary: DeviceSummaryModel;
    @Input() deviceWWN: string;
    @Output() deviceDeleted = new EventEmitter<string>();

    config: AppConfig;

    private _unsubscribeAll: Subject<any>;

    readonly humanizeDuration = humanizeDuration;

    ngOnInit(): void {
        // Subscribe to config changes
        this._configService.config$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((config: AppConfig) => {
                this.config = config;
            });
    }


    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    classDeviceLastUpdatedOn(deviceSummary: DeviceSummaryModel): string {
        const deviceStatus = this.deviceStatusString(deviceSummary)
        if (deviceStatus === 'failed') {
            return 'text-red' // if the device has failed, always highlight in red
        } else if (deviceStatus === 'passed') {
            if (moment().subtract(14, 'days').isBefore(deviceSummary.smart.collector_date)) {
                // this device was updated in the last 2 weeks.
                return 'text-green'
            } else if (moment().subtract(1, 'months').isBefore(deviceSummary.smart.collector_date)) {
                // this device was updated in the last month
                return 'text-yellow'
            } else {
                // last updated more than a month ago.
                return 'text-red'
            }
        } else {
            return ''
        }
    }


    deviceStatusString(deviceSummary: DeviceSummaryModel): deviceStatusName {
        // no smart data, so treat the device status as unknown
        if (!deviceSummary.smart) {
            return 'unknown'
        }

        // determine the device status, by comparing it against the allowed threshold
        // tslint:disable-next-line:no-bitwise
        const deviceStatus = deviceSummary.device.device_status & this.config.metrics.status_threshold
        if (deviceStatus === 0) {
            return 'passed'
        } else {
            return 'failed'
        }
    }


    openDeleteDialog(): void {
        const dialogRef = this.dialog.open(DashboardDeviceDeleteDialogComponent, {
            // width: '250px',
            data: {
                wwn: this.deviceWWN,
                title: DeviceTitlePipe.deviceTitleWithFallback(this.deviceSummary.device, this.config.dashboard_display)
            }
        });

        dialogRef.afterClosed().subscribe(result => {
            console.log('The dialog was closed', result);
            if (result.success) {
                this.deviceDeleted.emit(this.deviceWWN)
            }
        });
    }
}
