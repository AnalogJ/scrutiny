import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from 'app/shared/shared.module';
import { DashboardComponent } from 'app/modules/dashboard/dashboard.component';
import { dashboardRoutes } from 'app/modules/dashboard/dashboard.routing';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { NgApexchartsModule } from 'ng-apexcharts';
import { MatTooltipModule } from '@angular/material/tooltip'
import { DashboardSettingsModule } from 'app/layout/common/dashboard-settings/dashboard-settings.module';
import { DashboardDeviceModule } from 'app/layout/common/dashboard-device/dashboard-device.module';

@NgModule({
    declarations: [
        DashboardComponent
    ],
    imports: [
        RouterModule.forChild(dashboardRoutes),
        MatButtonModule,
        MatDividerModule,
        MatTooltipModule,
        MatIconModule,
        MatMenuModule,
        MatProgressBarModule,
        MatSortModule,
        MatTableModule,
        NgApexchartsModule,
        SharedModule,
        DashboardSettingsModule,
        DashboardDeviceModule
    ]
})
export class DashboardModule
{
}
