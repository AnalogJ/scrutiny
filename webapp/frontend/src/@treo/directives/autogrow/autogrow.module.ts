import { NgModule } from '@angular/core';
import { TreoAutogrowDirective } from '@treo/directives/autogrow/autogrow.directive';

@NgModule({
    declarations: [
        TreoAutogrowDirective
    ],
    exports     : [
        TreoAutogrowDirective
    ]
})
export class TreoAutogrowModule
{
}
