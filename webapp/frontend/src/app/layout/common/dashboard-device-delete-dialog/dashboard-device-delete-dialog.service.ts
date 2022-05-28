import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { getBasePath } from 'app/app.routing';

@Injectable({
    providedIn: 'root'
})
export class DashboardDeviceDeleteDialogService
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


    deleteDevice(wwn: string): Observable<any>
    {
        return this._httpClient.delete( `${getBasePath()}/api/device/${wwn}`, {});
    }
}
