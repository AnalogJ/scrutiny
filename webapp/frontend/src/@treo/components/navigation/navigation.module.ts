import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TreoScrollbarModule } from '@treo/directives/scrollbar/public-api';
import { TreoHorizontalNavigationBasicItemComponent } from '@treo/components/navigation/horizontal/components/basic/basic.component';
import { TreoHorizontalNavigationBranchItemComponent } from '@treo/components/navigation/horizontal/components/branch/branch.component';
import { TreoHorizontalNavigationDividerItemComponent } from '@treo/components/navigation/horizontal/components/divider/divider.component';
import { TreoHorizontalNavigationSpacerItemComponent } from '@treo/components/navigation/horizontal/components/spacer/spacer.component';
import { TreoHorizontalNavigationComponent } from '@treo/components/navigation/horizontal/horizontal.component';
import { TreoVerticalNavigationAsideItemComponent } from '@treo/components/navigation/vertical/components/aside/aside.component';
import { TreoVerticalNavigationBasicItemComponent } from '@treo/components/navigation/vertical/components/basic/basic.component';
import { TreoVerticalNavigationCollapsableItemComponent } from '@treo/components/navigation/vertical/components/collapsable/collapsable.component';
import { TreoVerticalNavigationDividerItemComponent } from '@treo/components/navigation/vertical/components/divider/divider.component';
import { TreoVerticalNavigationGroupItemComponent } from '@treo/components/navigation/vertical/components/group/group.component';
import { TreoVerticalNavigationSpacerItemComponent } from '@treo/components/navigation/vertical/components/spacer/spacer.component';
import { TreoVerticalNavigationComponent } from '@treo/components/navigation/vertical/vertical.component';

@NgModule({
    declarations: [
        TreoHorizontalNavigationBasicItemComponent,
        TreoHorizontalNavigationBranchItemComponent,
        TreoHorizontalNavigationDividerItemComponent,
        TreoHorizontalNavigationSpacerItemComponent,
        TreoHorizontalNavigationComponent,
        TreoVerticalNavigationAsideItemComponent,
        TreoVerticalNavigationBasicItemComponent,
        TreoVerticalNavigationCollapsableItemComponent,
        TreoVerticalNavigationDividerItemComponent,
        TreoVerticalNavigationGroupItemComponent,
        TreoVerticalNavigationSpacerItemComponent,
        TreoVerticalNavigationComponent
    ],
    imports     : [
        CommonModule,
        RouterModule,
        MatButtonModule,
        MatDividerModule,
        MatIconModule,
        MatMenuModule,
        MatTooltipModule,
        TreoScrollbarModule
    ],
    exports     : [
        TreoHorizontalNavigationComponent,
        TreoVerticalNavigationComponent
    ]
})
export class TreoNavigationModule
{
}
