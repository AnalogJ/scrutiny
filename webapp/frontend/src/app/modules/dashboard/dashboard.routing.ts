import { Route } from '@angular/router';
import { DashboardComponent } from 'app/modules/dashboard/dashboard.component';
import {DashboardResolver} from 'app/modules/dashboard/dashboard.resolvers';

export const dashboardRoutes: Route[] = [
    {
        path     : '',
        component: DashboardComponent,
        resolve  : {
            sales: DashboardResolver
        }
    }
];
