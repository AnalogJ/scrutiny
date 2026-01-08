import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from 'app/shared/shared.module';
import { ZFSPoolDetailComponent } from 'app/modules/zfs-pool-detail/zfs-pool-detail.component';
import { zfsPoolDetailRoutes } from 'app/modules/zfs-pool-detail/zfs-pool-detail.routing';
import { MatButtonModule as MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule as MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule as MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule as MatTableModule } from '@angular/material/table';
import { MatTooltipModule as MatTooltipModule } from '@angular/material/tooltip';
import { NgApexchartsModule } from 'ng-apexcharts';
import { TreoCardModule } from '@treo/components/card';

@NgModule({
    declarations: [
        ZFSPoolDetailComponent
    ],
    imports: [
        RouterModule.forChild(zfsPoolDetailRoutes),
        MatButtonModule,
        MatDividerModule,
        MatTooltipModule,
        MatIconModule,
        MatMenuModule,
        MatProgressBarModule,
        MatSortModule,
        MatTableModule,
        NgApexchartsModule,
        TreoCardModule,
        SharedModule,
    ]
})
export class ZFSPoolDetailModule {
}
