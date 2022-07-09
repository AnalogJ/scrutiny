import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BehaviorSubject, Observable} from 'rxjs';
import {tap} from 'rxjs/operators';
import {getBasePath} from 'app/app.routing';
import {DeviceDetailsResponseWrapper} from 'app/core/models/device-details-response-wrapper';

@Injectable({
    providedIn: 'root'
})
export class DetailService {
    // Observables
    private _data: BehaviorSubject<DeviceDetailsResponseWrapper>;

    /**
     * Constructor
     *
     * @param {HttpClient} _httpClient
     */
    constructor(
        private _httpClient: HttpClient
    ) {
        // Set the private defaults
        this._data = new BehaviorSubject(null);
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Getter for data
     */
    get data$(): Observable<DeviceDetailsResponseWrapper> {
        return this._data.asObservable();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Get data
     */
    getData(wwn): Observable<DeviceDetailsResponseWrapper> {
        return this._httpClient.get(getBasePath() + `/api/device/${wwn}/details`).pipe(
            tap((response: DeviceDetailsResponseWrapper) => {
                this._data.next(response);
            })
        );
    }
}
