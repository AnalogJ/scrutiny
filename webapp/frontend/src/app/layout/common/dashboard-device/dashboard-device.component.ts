import { Component, Input, Output, OnInit, EventEmitter} from '@angular/core';
import * as moment from "moment";
import {takeUntil} from "rxjs/operators";
import {AppConfig} from "app/core/config/app.config";
import {TreoConfigService} from "@treo/services/config";
import {Subject} from "rxjs";
import  humanizeDuration from 'humanize-duration'
import {MatDialog} from '@angular/material/dialog';
import {DashboardDeviceDeleteDialogComponent} from "app/layout/common/dashboard-device-delete-dialog/dashboard-device-delete-dialog.component";
import {DeviceTitlePipe} from "app/shared/device-title.pipe";

@Component({
  selector: 'app-dashboard-device',
  templateUrl: './dashboard-device.component.html',
  styleUrls: ['./dashboard-device.component.scss']
})
export class DashboardDeviceComponent implements OnInit {
    @Input() deviceSummary: any;
    @Input() deviceWWN: string;
    @Output() deviceDeleted = new EventEmitter<string>();

    config: AppConfig;

    private _unsubscribeAll: Subject<any>;

    constructor(
        private _configService: TreoConfigService,
        public dialog: MatDialog,
    ) {
        // Set the private defaults
        this._unsubscribeAll = new Subject();
    }

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

    classDeviceLastUpdatedOn(deviceSummary){
        if (deviceSummary.device.device_status !== 0) {
            return 'text-red' // if the device has failed, always highlight in red
        } else if(deviceSummary.device.device_status === 0 && deviceSummary.smart){
            if(moment().subtract(14, 'd').isBefore(deviceSummary.smart.collector_date)){
                // this device was updated in the last 2 weeks.
                return 'text-green'
            } else if(moment().subtract(1, 'm').isBefore(deviceSummary.smart.collector_date)){
                // this device was updated in the last month
                return 'text-yellow'
            } else{
                // last updated more than a month ago.
                return 'text-red'
            }

        } else {
            return ''
        }
    }

    deviceStatusString(deviceStatus){
        if(deviceStatus == 0){
            return "passed"
        } else {
            return "failed"
        }
    }

    readonly humanizeDuration = humanizeDuration;



    openDeleteDialog(): void {
        const dialogRef = this.dialog.open(DashboardDeviceDeleteDialogComponent, {
            // width: '250px',
            data: {wwn: this.deviceWWN, title: DeviceTitlePipe.deviceTitleWithFallback(this.deviceSummary.device, this.config.dashboardDisplay)}
        });

        dialogRef.afterClosed().subscribe(result => {
            console.log('The dialog was closed', result);
            if(result.success){
                this.deviceDeleted.emit(this.deviceWWN)
            }
        });
    }
}
