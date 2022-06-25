import { Inject, Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import * as _ from 'lodash';
import { TREO_APP_CONFIG } from '@treo/services/config/config.constants';
import { AppConfig } from 'app/core/config/app.config';

const SCRUTINY_CONFIG_LOCAL_STORAGE_KEY = 'scrutiny';

@Injectable({
    providedIn: 'root'
})
export class TreoConfigService
{
    // Private
    private _config: BehaviorSubject<any>;

    /**
     * Constructor
     */
    constructor(@Inject(TREO_APP_CONFIG) defaultConfig: any)
    {
        let currentScrutinyConfig = defaultConfig

        const localConfigStr = localStorage.getItem(SCRUTINY_CONFIG_LOCAL_STORAGE_KEY)
        if (localConfigStr){
            // check localstorage for a value
            const localConfig = JSON.parse(localConfigStr)
            currentScrutinyConfig = Object.assign({}, currentScrutinyConfig, localConfig) // make sure defaults are available if missing from localStorage.
        }
        // Set the private defaults
        this._config = new BehaviorSubject(currentScrutinyConfig);
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Setter and getter for config
     */
    // Setter
    set config(value: any)
    {
        // Merge the new config over to the current config
        const config = _.merge({}, this._config.getValue(), value);

        // Store the config in localstorage
        localStorage.setItem(SCRUTINY_CONFIG_LOCAL_STORAGE_KEY, JSON.stringify(config));

        // Execute the observable
        this._config.next(config);
    }

    // Getter
    get config$(): Observable<any>
    {
        return this._config.asObservable();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Resets the config to the default
     */
    reset(): void
    {
        // Set the config
        this._config.next(this.config);
    }
}
