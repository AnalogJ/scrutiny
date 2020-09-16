import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-dashboard-settings',
  templateUrl: './dashboard-settings.component.html',
  styleUrls: ['./dashboard-settings.component.scss']
})
export class DashboardSettingsComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
  }
    formatLabel(value: number) {
        return value;
    }
}
