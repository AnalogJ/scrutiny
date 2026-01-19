import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';
import { ZFSPoolsService } from 'app/modules/zfs-pools/zfs-pools.service';
import { ZFSPoolModel } from 'app/core/models/zfs-pool-model';

@Injectable({
    providedIn: 'root',
})
export class ZFSPoolsResolver {
    constructor(private _zfsPoolsService: ZFSPoolsService) {}

    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Record<string, ZFSPoolModel>> {
        return this._zfsPoolsService.getSummaryData();
    }
}
