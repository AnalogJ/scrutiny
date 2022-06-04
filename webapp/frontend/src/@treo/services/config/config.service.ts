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
    private systemPrefersDark: boolean;

    /**
     * Constructor
     */
    constructor(@Inject(TREO_APP_CONFIG) defaultConfig: any)
    {
        this.systemPrefersDark = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;

        let currentScrutinyConfig = defaultConfig

        let localConfigStr = localStorage.getItem(SCRUTINY_CONFIG_LOCAL_STORAGE_KEY)
        if (localConfigStr){
            //check localstorage for a value
            let localConfig = JSON.parse(localConfigStr)
            currentScrutinyConfig = Object.assign({}, localConfig, currentScrutinyConfig) // make sure defaults are available if missing from localStorage.
        }

        currentScrutinyConfig.theme = this.determineTheme(currentScrutinyConfig);

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
        let config = _.merge({}, this._config.getValue(), value);

        //Store the config in localstorage
        localStorage.setItem(SCRUTINY_CONFIG_LOCAL_STORAGE_KEY, JSON.stringify(config));
        
        config.theme = this.determineTheme(config);

        // Execute the observable
        this._config.next(config);
    }

    //Getter
    get config$(): Observable<any>
    {
        return this._config.asObservable();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Checks if theme should be set to dark based on config & system settings
     */
    private determineTheme(config:AppConfig): string {
        return (config.darkModeUseSystem && this.systemPrefersDark) ? "dark" :  config.theme;
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
