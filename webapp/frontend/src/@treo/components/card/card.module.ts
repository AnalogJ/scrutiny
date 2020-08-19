import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TreoCardComponent } from '@treo/components/card/card.component';

@NgModule({
    declarations: [
        TreoCardComponent
    ],
    imports     : [
        CommonModule
    ],
    exports     : [
        TreoCardComponent
    ]
})
export class TreoCardModule
{
}
