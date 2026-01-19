import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatButtonModule as MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule as MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule as MatTooltipModule } from '@angular/material/tooltip';
import { ZFSPoolCardComponent } from './zfs-pool-card.component';
import { SharedModule } from 'app/shared/shared.module';
import { MatDialogModule } from '@angular/material/dialog';

@NgModule({
    declarations: [ZFSPoolCardComponent],
    imports: [CommonModule, RouterModule, MatButtonModule, MatIconModule, MatMenuModule, MatTooltipModule, SharedModule, MatDialogModule],
    exports: [ZFSPoolCardComponent],
})
export class ZFSPoolCardModule {}
