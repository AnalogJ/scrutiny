import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BehaviorSubject, Observable} from 'rxjs';
import {map, tap} from 'rxjs/operators';
import {getBasePath} from 'app/app.routing';
import {DeviceSummaryResponseWrapper} from 'app/core/models/device-summary-response-wrapper';
import {DeviceSummaryModel} from 'app/core/models/device-summary-model';
import {SmartTemperatureModel} from 'app/core/models/measurements/smart-temperature-model';
import {DeviceSummaryTempResponseWrapper} from 'app/core/models/device-summary-temp-response-wrapper';

@Injectable({
    providedIn: 'root'
})
export class DashboardService {
    // Observables
    private _data: BehaviorSubject<{ [p: string]: DeviceSummaryModel }>;

    /**
     * Constructor
     *
     * @param {HttpClient} _httpClient
     */
    constructor(
        private _httpClient: HttpClient
    )
    {
        // Set the private defaults
        this._data = new BehaviorSubject(null);
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Getter for data
     */
    get data$(): Observable<{ [p: string]: DeviceSummaryModel }> {
        return this._data.asObservable();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Get data
     */
    getSummaryData(): Observable<{ [key: string]: DeviceSummaryModel }> {
        return this._httpClient.get(getBasePath() + '/api/summary').pipe(
            map((response: DeviceSummaryResponseWrapper) => {
                // console.log("FILTERING=----", response.data.summary)
                return response.data.summary
            }),
            tap((response: { [key: string]: DeviceSummaryModel }) => {
                this._data.next(response);
            })
        );
    }

    getSummaryTempData(durationKey: string): Observable<{ [key: string]: SmartTemperatureModel[] }> {
        const params = {}
        if (durationKey) {
            params['duration_key'] = durationKey
        }

        return this._httpClient.get(getBasePath() + '/api/summary/temp', {params: params}).pipe(
            map((response: DeviceSummaryTempResponseWrapper) => {
                return response.data.temp_history
            })
        );
    }
}
