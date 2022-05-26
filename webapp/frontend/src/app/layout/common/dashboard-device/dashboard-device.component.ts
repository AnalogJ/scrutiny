import {ChangeDetectorRef, Component, Input, OnInit} from '@angular/core';
import * as moment from "moment";
import {takeUntil} from "rxjs/operators";
import {AppConfig} from "app/core/config/app.config";
import {TreoConfigService} from "@treo/services/config";
import {Subject} from "rxjs";
import  humanizeDuration from 'humanize-duration'
import {MatDialog} from '@angular/material/dialog';
import {DashboardDeviceDeleteDialogComponent} from "app/layout/common/dashboard-device-delete-dialog/dashboard-device-delete-dialog.component";

@Component({
  selector: 'app-dashboard-device',
  templateUrl: './dashboard-device.component.html',
  styleUrls: ['./dashboard-device.component.scss']
})
export class DashboardDeviceComponent implements OnInit {
    @Input() deviceSummary: any;
    @Input() deviceWWN: string;
    deleted = false;

    config: AppConfig;

    private _unsubscribeAll: Subject<any>;

    constructor(
        private _configService: TreoConfigService,
        public dialog: MatDialog,
        private cdRef: ChangeDetectorRef,
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

    deviceTitle(disk){

        console.log(`Displaying Device ${disk.wwn} with: ${this.config.dashboardDisplay}`)
        let titleParts = []
        if (disk.host_id) titleParts.push(disk.host_id)

        //add device identifier (fallback to generated device name)
        titleParts.push(deviceDisplayTitle(disk, this.config.dashboardDisplay) || deviceDisplayTitle(disk, 'name'))

        return titleParts.join(' - ')
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
            data: {wwn: this.deviceWWN, title: this.deviceTitle(this.deviceSummary.device)}
        });

        dialogRef.afterClosed().subscribe(result => {
            console.log('The dialog was closed', result);
            this.deleted = result.success
            this.cdRef.detectChanges()
        });
    }
}

export function deviceDisplayTitle(disk, titleType: string){
    let titleParts = []
    switch(titleType){
        case 'name':
            titleParts.push(`/dev/${disk.device_name}`)
            if (disk.device_type && disk.device_type != 'scsi' && disk.device_type != 'ata'){
                titleParts.push(disk.device_type)
            }
            titleParts.push(disk.model_name)

            break;
        case 'serial_id':
            if(!disk.device_serial_id) return ''
            titleParts.push(`/by-id/${disk.device_serial_id}`)
            break;
        case 'uuid':
            if(!disk.device_uuid) return ''
            titleParts.push(`/by-uuid/${disk.device_uuid}`)
            break;
        case 'label':
            if(disk.label){
                titleParts.push(disk.label)
            } else if(disk.device_label){
                titleParts.push(`/by-label/${disk.device_label}`)
            }
            break;
    }
    return titleParts.join(' - ')
}
