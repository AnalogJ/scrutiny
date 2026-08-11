import {Injectable} from '@angular/core';
import {HttpRequest} from '@angular/common/http';
import * as _ from 'lodash';
import {TreoMockApi} from '@treo/lib/mock-api/mock-api.interfaces';
import {TreoMockApiService} from '@treo/lib/mock-api/mock-api.service';
import {AppConfig, appConfig} from 'app/core/config/app.config';

@Injectable({
    providedIn: 'root'
})
export class SettingsMockApi implements TreoMockApi {
    private _settings: AppConfig = _.cloneDeep(appConfig);

    constructor(private _treoMockApiService: TreoMockApiService) {
        this.register();
    }

    register(): void {
        this._treoMockApiService
            .onGet('/api/settings')
            .reply(() => [
                200,
                {success: true, settings: _.cloneDeep(this._settings)}
            ]);

        this._treoMockApiService
            .onPost('/api/settings')
            .reply((request: HttpRequest<AppConfig>) => {
                this._settings = _.cloneDeep(request.body);

                return [
                    200,
                    {success: true, settings: _.cloneDeep(this._settings)}
                ];
            });
    }
}
