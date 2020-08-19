import { NgModule } from '@angular/core';
import { TreoFindByKeyPipe } from '@treo/pipes/find-by-key/find-by-key.pipe';

@NgModule({
    declarations: [
        TreoFindByKeyPipe
    ],
    exports     : [
        TreoFindByKeyPipe
    ]
})
export class TreoFindByKeyPipeModule
{
}
