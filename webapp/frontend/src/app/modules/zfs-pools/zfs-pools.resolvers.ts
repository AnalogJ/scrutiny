import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';
import { ZFSPoolsService } from 'app/modules/zfs-pools/zfs-pools.service';
import { ZFSPoolSummaryModel } from 'app/core/models/zfs-pool-summary-model';

@Injectable({
    providedIn: 'root'
})
export class ZFSPoolsResolver {
    constructor(
        private _zfsPoolsService: ZFSPoolsService
    ) {
    }

    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<{ [guid: string]: ZFSPoolSummaryModel }> {
        return this._zfsPoolsService.getSummaryData();
    }
}
