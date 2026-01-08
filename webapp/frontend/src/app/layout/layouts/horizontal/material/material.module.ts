import { NgModule } from '@angular/core';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { RouterModule } from '@angular/router';
import { MatButtonModule as MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule as MatMenuModule } from '@angular/material/menu';
import { TreoNavigationModule } from '@treo/components/navigation';
import { SearchModule } from 'app/layout/common/search/search.module';
import { SharedModule } from 'app/shared/shared.module';
import { MaterialLayoutComponent } from 'app/layout/layouts/horizontal/material/material.component';

@NgModule({ declarations: [
        MaterialLayoutComponent
    ],
    exports: [
        MaterialLayoutComponent
    ], imports: [RouterModule,
        MatButtonModule,
        MatDividerModule,
        MatIconModule,
        MatMenuModule,
        TreoNavigationModule,
        SearchModule,
        SharedModule], providers: [provideHttpClient(withInterceptorsFromDi())] })
export class MaterialLayoutModule
{
}
