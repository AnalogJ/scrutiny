import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { Overlay } from '@angular/cdk/overlay';
import { MAT_AUTOCOMPLETE_SCROLL_STRATEGY, MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { SharedModule } from 'app/shared/shared.module';
import {DashboardDeviceComponent} from 'app/layout/common/dashboard-device/dashboard-device.component'
import { MatDialogModule } from "@angular/material/dialog";
import { MatButtonToggleModule} from "@angular/material/button-toggle";
import {MatTabsModule} from "@angular/material/tabs";
import {MatSliderModule} from "@angular/material/slider";
import {MatSlideToggleModule} from "@angular/material/slide-toggle";
import {MatTooltipModule} from "@angular/material/tooltip";
import {dashboardRoutes} from "../../../modules/dashboard/dashboard.routing";
import {MatDividerModule} from "@angular/material/divider";
import {MatMenuModule} from "@angular/material/menu";
import {MatProgressBarModule} from "@angular/material/progress-bar";
import {MatSortModule} from "@angular/material/sort";
import {MatTableModule} from "@angular/material/table";
import {NgApexchartsModule} from "ng-apexcharts";
import {DashboardDeviceDeleteDialogModule} from "../dashboard-device-delete-dialog/dashboard-device-delete-dialog.module";

@NgModule({
    declarations: [
        DashboardDeviceComponent
    ],
    imports     : [
        RouterModule.forChild([]),
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
        DashboardDeviceDeleteDialogModule
    ],
    exports     : [
        DashboardDeviceComponent,
    ],
    providers   : []
})
export class DashboardDeviceModule
{
}
