import {NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {SharedModule} from 'app/shared/shared.module';
import {DashboardDeviceDeleteDialogComponent} from 'app/layout/common/dashboard-device-delete-dialog/dashboard-device-delete-dialog.component'
import {dashboardRoutes} from 'app/modules/dashboard/dashboard.routing';
import {MatDialogModule} from '@angular/material/dialog';

@NgModule({
    declarations: [
        DashboardDeviceDeleteDialogComponent
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
        DashboardDeviceDeleteDialogComponent,
    ],
    providers   : []
})
export class DashboardDeviceDeleteDialogModule
{
}
