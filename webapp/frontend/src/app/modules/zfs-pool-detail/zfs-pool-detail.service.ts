import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { getBasePath } from 'app/app.routing';
import { ZFSPoolDetailsResponseWrapper } from 'app/core/models/zfs-pool-summary-model';

@Injectable({
    providedIn: 'root'
})
export class ZFSPoolDetailService {
    // Observables
    private _data: BehaviorSubject<ZFSPoolDetailsResponseWrapper>;

    constructor(
        private _httpClient: HttpClient
    ) {
        this._data = new BehaviorSubject(null);
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    get data$(): Observable<ZFSPoolDetailsResponseWrapper> {
        return this._data.asObservable();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    getData(guid: string): Observable<ZFSPoolDetailsResponseWrapper> {
        return this._httpClient.get(getBasePath() + `/api/zfs/pool/${guid}/details`).pipe(
            tap((response: ZFSPoolDetailsResponseWrapper) => {
                this._data.next(response);
            })
        );
    }

    setMuted(guid: string, muted: boolean): Observable<any> {
        const action = muted ? 'mute' : 'unmute';
        return this._httpClient.post(getBasePath() + `/api/zfs/pool/${guid}/${action}`, {});
    }

    setLabel(guid: string, label: string): Observable<any> {
        return this._httpClient.post(getBasePath() + `/api/zfs/pool/${guid}/label`, { label });
    }
}
