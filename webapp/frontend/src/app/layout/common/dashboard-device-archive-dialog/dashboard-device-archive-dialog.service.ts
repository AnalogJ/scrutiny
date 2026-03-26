import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {getBasePath} from 'app/app.routing';

@Injectable({
    providedIn: 'root'
})
export class DashboardDeviceArchiveDialogService
{


    /**
     * Constructor
     *
     * @param {HttpClient} _httpClient
     */
    constructor(
        private _httpClient: HttpClient
    )
    {
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------


    archiveDevice(scrutiny_uid: string): Observable<any>
    {
        return this._httpClient.post( `${getBasePath()}/api/device/${scrutiny_uid}/archive`, {});
    }

    unarchiveDevice(scrutiny_uid: string): Observable<any>
    {
        return this._httpClient.post( `${getBasePath()}/api/device/${scrutiny_uid}/unarchive`, {});
    }
}
