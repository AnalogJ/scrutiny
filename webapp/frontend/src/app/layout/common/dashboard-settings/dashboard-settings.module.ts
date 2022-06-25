import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { Overlay } from '@angular/cdk/overlay';
import { MAT_AUTOCOMPLETE_SCROLL_STRATEGY, MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { SharedModule } from 'app/shared/shared.module';
import {DashboardSettingsComponent} from 'app/layout/common/dashboard-settings/dashboard-settings.component'
import { MatDialogModule } from '@angular/material/dialog';
import { MatButtonToggleModule} from '@angular/material/button-toggle';
import {MatTabsModule} from '@angular/material/tabs';
import {MatSliderModule} from '@angular/material/slider';
import {MatSlideToggleModule} from '@angular/material/slide-toggle';
import {MatTooltipModule} from '@angular/material/tooltip';

@NgModule({
    declarations: [
        DashboardSettingsComponent
    ],
    imports     : [
        RouterModule.forChild([]),
        MatAutocompleteModule,
        MatDialogModule,
        MatButtonModule,
        MatSelectModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatButtonToggleModule,
        MatTabsModule,
        MatTooltipModule,
        MatSliderModule,
        MatSlideToggleModule,
        SharedModule
    ],
    exports     : [
        DashboardSettingsComponent
    ],
    providers   : []
})
export class DashboardSettingsModule
{
}
