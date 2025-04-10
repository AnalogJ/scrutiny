import {NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {SharedModule} from 'app/shared/shared.module';
import {dashboardRoutes} from 'app/modules/dashboard/dashboard.routing';
import {MatDialogModule} from '@angular/material/dialog';
import {DashboardDeviceArchiveDialogComponent} from './dashboard-device-archive-dialog.component';

@NgModule({
    declarations: [
        DashboardDeviceArchiveDialogComponent
    ],
    imports: [
        RouterModule.forChild([]),
        RouterModule.forChild(dashboardRoutes),
        MatButtonModule,
        MatIconModule,
        SharedModule,
        MatDialogModule
    ],
    exports     : [
        DashboardDeviceArchiveDialogComponent,
    ],
    providers   : []
})
export class DashboardDeviceArchiveDialogModule
{
}
