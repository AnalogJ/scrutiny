import {NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import {MatButtonModule as MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {SharedModule} from 'app/shared/shared.module';
import {DashboardDeviceComponent} from 'app/layout/common/dashboard-device/dashboard-device.component'
import {dashboardRoutes} from '../../../modules/dashboard/dashboard.routing';
import {MatMenuModule as MatMenuModule} from '@angular/material/menu';
import {DashboardDeviceDeleteDialogModule} from 'app/layout/common/dashboard-device-delete-dialog/dashboard-device-delete-dialog.module';
import {DashboardDeviceArchiveDialogModule} from '../dashboard-device-archive-dialog/dashboard-device-archive-dialog.module';

@NgModule({
    declarations: [
        DashboardDeviceComponent
    ],
    imports: [
        RouterModule.forChild([]),
        RouterModule.forChild(dashboardRoutes),
        MatButtonModule,
        MatIconModule,
        MatMenuModule,
        SharedModule,
        DashboardDeviceDeleteDialogModule,
        DashboardDeviceArchiveDialogModule
    ],
    exports: [
        DashboardDeviceComponent,
    ],
    providers: []
})
export class DashboardDeviceModule {
}
