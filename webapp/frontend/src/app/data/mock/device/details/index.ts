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
            .onGet('/api/device/ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sda)
                ];
            });

        this._treoMockApiService
            .onGet('/api/device/3ea22b35-682b-49fb-a655-abffed108e48/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sdb)
                ];
            });

        this._treoMockApiService
            .onGet('/api/device/42caca8a-9b95-5c75-b059-305771a2a193/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sdc)
                ];
            });

        this._treoMockApiService
            .onGet('/api/device/d8796fe7-2422-520c-8991-e970993dad3e/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sdd)
                ];
            });

        this._treoMockApiService
            .onGet('/api/device/00328b73-9f8a-53ad-8f20-8d0b1be00f47/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sde)
                ];
            });

        this._treoMockApiService
            .onGet('/api/device/e5ccc378-24fc-5a9d-b1ce-8732096a9ea5/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(sdf)
                ];
            });
    }
}
