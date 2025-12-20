import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'app-detail-settings',
  templateUrl: './detail-settings.component.html',
  styleUrls: ['./detail-settings.component.scss']
})
export class DetailSettingsComponent implements OnInit {

  muted: boolean;
  label: string;

  constructor(
      @Inject(MAT_DIALOG_DATA) public data: { curMuted: boolean, curLabel: string }
  ) {
      this.muted = data.curMuted;
      this.label = data.curLabel || '';
  }

  ngOnInit(): void {
  }
}
