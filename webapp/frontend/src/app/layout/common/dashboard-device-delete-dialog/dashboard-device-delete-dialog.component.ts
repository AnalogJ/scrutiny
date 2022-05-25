import { Component, OnInit, Inject } from '@angular/core';
import {MAT_DIALOG_DATA} from '@angular/material/dialog';

@Component({
  selector: 'app-dashboard-device-delete-dialog',
  templateUrl: './dashboard-device-delete-dialog.component.html',
  styleUrls: ['./dashboard-device-delete-dialog.component.scss']
})
export class DashboardDeviceDeleteDialogComponent implements OnInit {

    constructor(@Inject(MAT_DIALOG_DATA) public data: {wwn: string, title: string}) { }

  ngOnInit(): void {
  }

}
