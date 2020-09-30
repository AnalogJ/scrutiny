import { Injectable } from '@angular/core';
import * as _ from 'lodash';
import { TreoMockApi } from '@treo/lib/mock-api/mock-api.interfaces';
import { TreoMockApiService } from '@treo/lib/mock-api/mock-api.service';
import { sda } from 'app/data/mock/device/details/sda';
import { sdb } from 'app/data/mock/device/details/sdb';
import { sdc } from 'app/data/mock/device/details/sdc';
import { sdd } from 'app/data/mock/device/details/sdd';
import { sde } from 'app/data/mock/device/details/sde';
import { sdf } from 'app/data/mock/device/details/sdf';

@Injectable({
    providedIn: 'root'
})
export class DetailsMockApi implements TreoMockApi
{
    // Private
    private _details: any;

    /**
     * Constructor
     *
     * @param _treoMockApiService
     */
    constructor(
        private _treoMockApiService: TreoMockApiService
    )
    {
        // Register the API endpoints
        this.register();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Register
     */
    register(): void
    {
        this._treoMockApiService
            .onGet('/api/device/0x5002538e40a22954/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sda)
                ];
            });

        this._treoMockApiService
            .onGet('/api/device/0x5000cca264eb01d7/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sdb)
                ];
            });

        this._treoMockApiService
            .onGet('/api/device/0x5000cca264ec3183/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sdc)
                ];
            });

        this._treoMockApiService
            .onGet('/api/device/0x5000cca252c859cc/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sdd)
                ];
            });

        this._treoMockApiService
            .onGet('/api/device/0x5000cca264ebc248/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sdf)
                ];
            });
    }
}
