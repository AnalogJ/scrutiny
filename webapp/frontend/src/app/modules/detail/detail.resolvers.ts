import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, Resolve, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';
import { DetailService } from 'app/modules/detail/detail.service';

@Injectable({
    providedIn: 'root'
})
export class DetailResolver implements Resolve<any>
{
    /**
     * Constructor
     *
     * @param {FinanceService} _detailService
     */
    constructor(
        private _detailService: DetailService
    )
    {
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
    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<any>
    {
        return this._detailService.getData(route.params.wwn);
    }
}
