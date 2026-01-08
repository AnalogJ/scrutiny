import * as _ from 'lodash';

import { Injectable } from '@angular/core';
import { TreoMockApi } from '@treo/lib/mock-api/mock-api.interfaces';
import { TreoMockApiService } from '@treo/lib/mock-api/mock-api.service';
import { settings as settingsData } from 'app/data/mock/settings/data';

@Injectable({
    providedIn: 'root'
})
export class SettingsMockApi implements TreoMockApi
{
    private _settings: any;

    constructor(private _treoMockApiService: TreoMockApiService)
    {
        this._settings = settingsData;
        this.register();
    }

    register(): void
    {
        this._treoMockApiService.onGet('/api/settings')
            .reply(() => [200, _.cloneDeep(this._settings)]);

        this._treoMockApiService.onPost('/api/settings')
            .reply(() => [200, { success: true }]);
    }
}
