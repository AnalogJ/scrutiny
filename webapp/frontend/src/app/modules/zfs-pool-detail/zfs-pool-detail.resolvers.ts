import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';
import { ZFSPoolDetailService } from 'app/modules/zfs-pool-detail/zfs-pool-detail.service';
import { ZFSPoolDetailsResponseWrapper } from 'app/core/models/zfs-pool-summary-model';

@Injectable({
    providedIn: 'root'
})
export class ZFSPoolDetailResolver {
    constructor(
        private _zfsPoolDetailService: ZFSPoolDetailService
    ) {
    }

    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<ZFSPoolDetailsResponseWrapper> {
        return this._zfsPoolDetailService.getData(route.params.guid);
    }
}
