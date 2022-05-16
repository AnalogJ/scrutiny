import { Inject, Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import * as _ from 'lodash';
import { TREO_APP_CONFIG } from '@treo/services/config/config.constants';

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

        let localConfigStr = localStorage.getItem(SCRUTINY_CONFIG_LOCAL_STORAGE_KEY)
        if(localConfigStr){
            //check localstorage for a value
            let localConfig = JSON.parse(localConfigStr)
            currentScrutinyConfig = localConfig
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
    //Setter
    set config(value: any)
    {
        // Merge the new config over to the current config
        const config = _.merge({}, this._config.getValue(), value);

        //Store the config in localstorage
        localStorage.setItem(SCRUTINY_CONFIG_LOCAL_STORAGE_KEY, JSON.stringify(config));

        // Execute the observable
        this._config.next(config);
    }

    //Getter
    get config$(): Observable<any>
    {
        return this._config.asObservable();
    }

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
