import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from 'app/shared/shared.module';
import { ZfsPoolDetailComponent } from './zfs-pool-detail.component';
import { zfsPoolDetailRoutes } from './zfs-pool-detail.routing';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatCardModule } from '@angular/material/card';
import { MatMenuModule } from '@angular/material/menu';
import { TreoCardModule } from '@treo/components/card';
import { NgApexchartsModule } from 'ng-apexcharts';

@NgModule({
    declarations: [
        ZfsPoolDetailComponent
    ],
    imports: [
        RouterModule.forChild(zfsPoolDetailRoutes),
        MatButtonModule,
        MatCardModule,
        MatDividerModule,
        MatIconModule,
        MatMenuModule,
        MatProgressBarModule,
        MatTableModule,
        MatTabsModule,
        MatTooltipModule,
        NgApexchartsModule,
        TreoCardModule,
        SharedModule
    ]
})
export class ZfsPoolDetailModule {
}