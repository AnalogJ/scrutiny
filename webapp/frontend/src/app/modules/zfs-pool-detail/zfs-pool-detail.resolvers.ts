import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, Resolve, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';
import { ZfsPoolDetailResponseWrapper } from '../../core/models/zfs-pool-model';
import { ZfsService } from '../dashboard/zfs.service';

@Injectable({
    providedIn: 'root'
})
export class ZfsPoolDetailResolver implements Resolve<any> {
    /**
     * Constructor
     */
    constructor(private _zfsService: ZfsService) {
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Resolver
     *
     * @param route
     * @param state
     */
    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<ZfsPoolDetailResponseWrapper> {
        const poolGuid = route.paramMap.get('poolGuid');
        return this._zfsService.getZfsPoolDetails(poolGuid);
    }
}