import {enableProdMode, NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {ExtraOptions, PreloadAllModules, RouterModule} from '@angular/router';
import {APP_BASE_HREF} from '@angular/common';
import {MarkdownModule} from 'ngx-markdown';
import {TreoModule} from '@treo';
import {ScrutinyConfigModule} from 'app/core/config/scrutiny-config.module';
import {TreoMockApiModule} from '@treo/lib/mock-api';
import {CoreModule} from 'app/core/core.module';
import {appConfig} from 'app/core/config/app.config';
import {mockDataServices} from 'app/data/mock';
import {LayoutModule} from 'app/layout/layout.module';
import {AppComponent} from 'app/app.component';
import {appRoutes, getAppBaseHref} from 'app/app.routing';

const routerConfig: ExtraOptions = {
    scrollPositionRestoration: 'enabled',
    preloadingStrategy: PreloadAllModules
};

let dev = [
    TreoMockApiModule.forRoot(mockDataServices),
];

// if production clear dev imports and set to prod mode
// @ts-ignore
if (process.env.NODE_ENV === 'production') {
    dev = [];
    enableProdMode();
}

@NgModule({
    declarations: [
        AppComponent,
    ],
    imports     : [
        BrowserModule,
        BrowserAnimationsModule,
        RouterModule.forRoot(appRoutes, routerConfig),

        // Treo & Treo Mock API
        TreoModule,
        ScrutinyConfigModule.forRoot(appConfig),
        ...dev,

        // Core
        CoreModule,

        // Layout
        LayoutModule,

        // 3rd party modules
        MarkdownModule.forRoot({})
    ],
    bootstrap   : [
        AppComponent
    ],
    providers: [
        {
            provide: APP_BASE_HREF,
            useValue: getAppBaseHref()
        }
    ],
})
export class AppModule
{
}
