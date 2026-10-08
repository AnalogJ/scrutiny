import { Route } from '@angular/router';
import { DashboardComponent } from 'app/modules/dashboard/dashboard.component';
import {DashboardResolver} from 'app/modules/dashboard/dashboard.resolvers';

export const dashboardRoutes: Route[] = [
    {
        path     : '',
        component: DashboardComponent,
        // re-fetch the summary when a settings change reloads the dashboard, since
        // its temp_history follows the configured chart duration
        runGuardsAndResolvers: 'always',
        resolve  : {
            sales: DashboardResolver
        }
    }
];
