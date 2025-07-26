import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { RouterModule } from '@angular/router';

import { DashboardZfsPoolComponent } from './dashboard-zfs-pool.component';

@NgModule({
    declarations: [
        DashboardZfsPoolComponent
    ],
    imports: [
        CommonModule,
        MatButtonModule,
        MatIconModule,
        MatMenuModule,
        RouterModule
    ],
    exports: [
        DashboardZfsPoolComponent
    ]
})
export class DashboardZfsPoolModule {}