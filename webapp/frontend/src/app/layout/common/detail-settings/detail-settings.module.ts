import {NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import {MatAutocompleteModule as MatAutocompleteModule} from '@angular/material/autocomplete';
import {MatButtonModule as MatButtonModule} from '@angular/material/button';
import {MatSelectModule as MatSelectModule} from '@angular/material/select';
import {MatFormFieldModule as MatFormFieldModule} from '@angular/material/form-field';
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule as MatInputModule} from '@angular/material/input';
import {SharedModule} from 'app/shared/shared.module';
import {DetailSettingsComponent} from 'app/layout/common/detail-settings/detail-settings.component'
import {MatDialogModule as MatDialogModule} from '@angular/material/dialog';
import {MatButtonToggleModule} from '@angular/material/button-toggle';
import {MatTabsModule as MatTabsModule} from '@angular/material/tabs';
import {MatSliderModule as MatSliderModule} from '@angular/material/slider';
import {MatSlideToggleModule as MatSlideToggleModule} from '@angular/material/slide-toggle';
import {MatTooltipModule as MatTooltipModule} from '@angular/material/tooltip';

@NgModule({
    declarations: [
        DetailSettingsComponent
    ],
    imports: [
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
    exports: [
        DetailSettingsComponent
    ],
    providers: []
})
export class DetailSettingsModule {
}
