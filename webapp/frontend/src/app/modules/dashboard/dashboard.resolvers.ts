import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Observable} from 'rxjs';
import {DashboardService} from 'app/modules/dashboard/dashboard.service';
import {DeviceSummaryModel} from 'app/core/models/device-summary-model';

@Injectable({
    providedIn: 'root'
})
export class DashboardResolver implements Resolve<any> {
    /**
     * Constructor
     *
     * @param {FinanceService} _dashboardService
     */
    constructor(
        private _dashboardService: DashboardService
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
    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<{ [p: string]: DeviceSummaryModel }> {
        return this._dashboardService.getSummaryData();
    }
}
