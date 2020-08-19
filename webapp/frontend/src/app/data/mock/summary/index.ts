import { Injectable } from '@angular/core';
import * as _ from 'lodash';
import { TreoMockApi } from '@treo/lib/mock-api/mock-api.interfaces';
import { TreoMockApiService } from '@treo/lib/mock-api/mock-api.service';
import { summary as summaryData } from 'app/data/mock/summary/data';

@Injectable({
    providedIn: 'root'
})
export class SummaryMockApi implements TreoMockApi
{
    // Private
    private _summary: any;

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
        this._summary = summaryData;

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
            .onGet('/api/summary')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(this._summary)
                ];
            });
    }
}
