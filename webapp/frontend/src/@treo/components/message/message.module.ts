import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TreoMessageComponent } from '@treo/components/message/message.component';

@NgModule({
    declarations: [
        TreoMessageComponent
    ],
    imports     : [
        CommonModule,
        MatButtonModule,
        MatIconModule
    ],
    exports     : [
        TreoMessageComponent
    ]
})
export class TreoMessageModule
{
}
