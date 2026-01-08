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
        // @ Summary - GET
        // -----------------------------------------------------------------------------------------------------
        this._treoMockApiService
            .onGet('/api/summary')
            .reply(() => {

                return [
                    200,
                    _.cloneDeep(this._summary)
                ];
            });

        // -----------------------------------------------------------------------------------------------------
        // @ Summary Temp History - GET
        // -----------------------------------------------------------------------------------------------------
        this._treoMockApiService
            .onGet('/api/summary/temp')
            .reply(() => {

                // Extract temp_history from summary data for each device
                const tempHistory: { [key: string]: any[] } = {};

                if (this._summary.data && this._summary.data.summary) {
                    for (const wwn in this._summary.data.summary) {
                        const deviceData = this._summary.data.summary[wwn];
                        if (deviceData.temp_history) {
                            tempHistory[wwn] = deviceData.temp_history;
                        } else {
                            tempHistory[wwn] = [];
                        }
                    }
                }

                return [
                    200,
                    {
                        success: true,
                        data: {
                            temp_history: tempHistory
                        }
                    }
                ];
            });
    }
}
