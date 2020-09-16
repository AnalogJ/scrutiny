import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from 'app/shared/shared.module';
import { DetailComponent } from 'app/modules/detail/detail.component';
import { detailRoutes } from 'app/modules/detail/detail.routing';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip'
import { NgApexchartsModule } from 'ng-apexcharts';
import { TreoCardModule } from '@treo/components/card';

@NgModule({
    declarations: [
        DetailComponent
    ],
    imports     : [
        RouterModule.forChild(detailRoutes),
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
export class DetailModule
{
}
