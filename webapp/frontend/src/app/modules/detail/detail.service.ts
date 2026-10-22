import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BehaviorSubject, Observable} from 'rxjs';
import {map, tap} from 'rxjs/operators';
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
     * Set the device's friendly name, an empty name clears it
     */
    setFriendlyName(scrutiny_uuid: string, friendlyName: string): Observable<string> {
        return this._httpClient.post<{ data: { friendly_name: string } }>(getBasePath() + `/api/device/${scrutiny_uuid}/friendly_name`, {friendly_name: friendlyName}).pipe(
            map(response => response.data.friendly_name)
        );
    }

    /**
     * Get data
     */
    getData(scrutiny_uuid): Observable<DeviceDetailsResponseWrapper> {
        return this._httpClient.get(getBasePath() + `/api/device/${scrutiny_uuid}/details`).pipe(
            tap((response: DeviceDetailsResponseWrapper) => {
                this._data.next(response);
            })
        );
    }
}
