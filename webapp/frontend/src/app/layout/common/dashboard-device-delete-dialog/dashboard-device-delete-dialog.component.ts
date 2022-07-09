import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {DashboardDeviceDeleteDialogService} from 'app/layout/common/dashboard-device-delete-dialog/dashboard-device-delete-dialog.service';

@Component({
  selector: 'app-dashboard-device-delete-dialog',
  templateUrl: './dashboard-device-delete-dialog.component.html',
  styleUrls: ['./dashboard-device-delete-dialog.component.scss']
})
export class DashboardDeviceDeleteDialogComponent implements OnInit {

    constructor(
        public dialogRef: MatDialogRef<DashboardDeviceDeleteDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: {wwn: string, title: string},
        private _deleteService: DashboardDeviceDeleteDialogService,
    ) {
    }

  ngOnInit(): void {
  }

  onDeleteClick(): void {
      this._deleteService.deleteDevice(this.data.wwn)
          .subscribe((data) => {
              this.dialogRef.close(data);
          });

  }
}
