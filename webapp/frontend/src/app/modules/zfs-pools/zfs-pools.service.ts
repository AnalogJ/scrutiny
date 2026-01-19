import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';
import { getBasePath } from 'app/app.routing';
import { ZFSPoolSummaryResponseWrapper } from 'app/core/models/zfs-pool-summary-model';
import { ZFSPoolModel } from 'app/core/models/zfs-pool-model';

@Injectable({
    providedIn: 'root',
})
export class ZFSPoolsService {
    // Observables
    private _data: BehaviorSubject<Record<string, ZFSPoolModel>>;

    constructor(private _httpClient: HttpClient) {
        this._data = new BehaviorSubject(null);
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    get data$(): Observable<Record<string, ZFSPoolModel>> {
        return this._data.asObservable();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    getSummaryData(): Observable<Record<string, ZFSPoolModel>> {
        return this._httpClient.get(getBasePath() + '/api/zfs/summary').pipe(
            map((response: ZFSPoolSummaryResponseWrapper) => {
                return response.data.pools;
            }),
            tap((response: Record<string, ZFSPoolModel>) => {
                this._data.next(response);
            })
        );
    }

    archivePool(guid: string): Observable<any> {
        return this._httpClient.post(getBasePath() + `/api/zfs/pool/${guid}/archive`, {});
    }

    unarchivePool(guid: string): Observable<any> {
        return this._httpClient.post(getBasePath() + `/api/zfs/pool/${guid}/unarchive`, {});
    }

    mutePool(guid: string): Observable<any> {
        return this._httpClient.post(getBasePath() + `/api/zfs/pool/${guid}/mute`, {});
    }

    unmutePool(guid: string): Observable<any> {
        return this._httpClient.post(getBasePath() + `/api/zfs/pool/${guid}/unmute`, {});
    }

    deletePool(guid: string): Observable<any> {
        return this._httpClient.delete(getBasePath() + `/api/zfs/pool/${guid}`);
    }

    setLabel(guid: string, label: string): Observable<any> {
        return this._httpClient.post(getBasePath() + `/api/zfs/pool/${guid}/label`, { label });
    }
}
