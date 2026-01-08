import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from 'app/shared/shared.module';
import { ZFSPoolsComponent } from 'app/modules/zfs-pools/zfs-pools.component';
import { zfsPoolsRoutes } from 'app/modules/zfs-pools/zfs-pools.routing';
import { MatButtonModule as MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule as MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule as MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule as MatTableModule } from '@angular/material/table';
import { MatTooltipModule as MatTooltipModule } from '@angular/material/tooltip';
import { ZFSPoolCardModule } from 'app/layout/common/zfs-pool-card/zfs-pool-card.module';

@NgModule({
    declarations: [
        ZFSPoolsComponent
    ],
    imports: [
        RouterModule.forChild(zfsPoolsRoutes),
        MatButtonModule,
        MatDividerModule,
        MatTooltipModule,
        MatIconModule,
        MatMenuModule,
        MatProgressBarModule,
        MatSortModule,
        MatTableModule,
        SharedModule,
        ZFSPoolCardModule
    ]
})
export class ZFSPoolsModule {
}
