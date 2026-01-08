import {Inject, Injectable} from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {TREO_APP_CONFIG} from '@treo/services/config/config.constants';
import {BehaviorSubject, Observable} from 'rxjs';
import {getBasePath} from '../../app.routing';
import {map, tap} from 'rxjs/operators';
import {AppConfig} from './app.config';
import {merge} from 'lodash';

@Injectable({
    providedIn: 'root'
})
export class ScrutinyConfigService {
    // Private
    private _config: BehaviorSubject<AppConfig>;
    private _defaultConfig: AppConfig;

    constructor(
        private _httpClient: HttpClient,
        @Inject(TREO_APP_CONFIG) defaultConfig: AppConfig
    ) {
        // Set the private defaults
        this._defaultConfig = defaultConfig
        this._config = new BehaviorSubject(null);
    }


    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Setter & getter for config
     */
    set config(value: AppConfig) {
        // get the current config, merge the new values, and then submit. (setTheme only sets a single key, not the whole obj)
        const mergedSettings = merge({}, this._config.getValue(), value);

        this._httpClient.post(getBasePath() + '/api/settings', mergedSettings).pipe(
            map((response: any) => {
                return response.settings
            }),
            tap((settings: AppConfig) => {
                this._config.next(settings);
                return settings
            })
        ).subscribe()
    }

    get config$(): Observable<AppConfig> {
        if (this._config.getValue()) {
            return this._config.asObservable()
        } else {
            return this._httpClient.get(getBasePath() + '/api/settings').pipe(
                map((response: any) => {
                    return response.settings
                }),
                tap((settings: AppConfig) => {
                    this._config.next(settings);
                    return this._config.asObservable()
                })
            );
        }

    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Resets the config to the default
     */
    reset(): void {
        // Set the config
        this.config = this._defaultConfig
    }
}
