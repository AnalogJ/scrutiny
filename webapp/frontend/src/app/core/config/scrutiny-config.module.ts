import {ModuleWithProviders, NgModule} from '@angular/core';
import {ScrutinyConfigService} from 'app/core/config/scrutiny-config.service';
import {TREO_APP_CONFIG} from '@treo/services/config/config.constants';

@NgModule()
export class ScrutinyConfigModule {
    /**
     * Constructor
     *
     * @param {ScrutinyConfigService} _scrutinyConfigService
     */
    constructor(
        private _scrutinyConfigService: ScrutinyConfigService
    ) {
    }

    /**
     * forRoot method for setting user configuration
     *
     * @param config
     */
    static forRoot(config: any): ModuleWithProviders<ScrutinyConfigModule> {
        return {
            ngModule: ScrutinyConfigModule,
            providers: [
                {
                    provide: TREO_APP_CONFIG,
                    useValue: config
                }
            ]
        };
    }
}
