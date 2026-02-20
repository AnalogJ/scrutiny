import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {DashboardDeviceArchiveDialogService} from 'app/layout/common/dashboard-device-archive-dialog/dashboard-device-archive-dialog.service';

@Component({
  selector: 'app-dashboard-device-archive-dialog',
  templateUrl: './dashboard-device-archive-dialog.component.html',
  styleUrls: ['./dashboard-device-archive-dialog.component.scss'],
})
export class DashboardDeviceArchiveDialogComponent implements OnInit {

    constructor(
        public dialogRef: MatDialogRef<DashboardDeviceArchiveDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: {scrutiny_uuid: string, title: string},
        private _archiveService: DashboardDeviceArchiveDialogService,
    ) {
    }

  ngOnInit(): void {
  }

  onArchiveClick(): void {
      this._archiveService.archiveDevice(this.data.scrutiny_uuid)
          .subscribe((data) => {
              this.dialogRef.close(data);
          });

  }
}
