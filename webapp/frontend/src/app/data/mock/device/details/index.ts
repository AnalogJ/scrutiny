import { Injectable } from '@angular/core';
import * as _ from 'lodash';
import { TreoMockApi } from '@treo/lib/mock-api/mock-api.interfaces';
import { TreoMockApiService } from '@treo/lib/mock-api/mock-api.service';
import { details as detailsData } from 'app/data/mock/device/details/data';

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
        // Set the data
        this._details = detailsData;

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
        // -----------------------------------------------------------------------------------------------------
        // @ Sales - GET
        // -----------------------------------------------------------------------------------------------------
        this._treoMockApiService
            .onGet('/api/device/:wwn/details')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(this._details)
                ];
            });
    }
}
