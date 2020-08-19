import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TreoDrawerComponent } from '@treo/components/drawer/drawer.component';

@NgModule({
    declarations: [
        TreoDrawerComponent
    ],
    imports     : [
        CommonModule
    ],
    exports     : [
        TreoDrawerComponent
    ]
})
export class TreoDrawerModule
{
}
