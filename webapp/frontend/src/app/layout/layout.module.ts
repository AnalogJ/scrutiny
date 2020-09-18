import { NgModule } from '@angular/core';
import { TreoDrawerModule } from '@treo/components/drawer';
import { LayoutComponent } from 'app/layout/layout.component';
import { EmptyLayoutModule } from 'app/layout/layouts/empty/empty.module';
import { MaterialLayoutModule } from 'app/layout/layouts/horizontal/material/material.module';

import { SharedModule } from 'app/shared/shared.module';

const modules = [
    // Empty
    EmptyLayoutModule,

    // Horizontal navigation
    MaterialLayoutModule,
];

@NgModule({
    declarations: [
        LayoutComponent,
    ],
    imports     : [
        TreoDrawerModule,
        SharedModule,
        ...modules
    ],
    exports     : [
        ...modules
    ]
})
export class LayoutModule
{
}
