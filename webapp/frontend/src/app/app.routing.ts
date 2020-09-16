import { Route } from '@angular/router';
import { LayoutComponent } from 'app/layout/layout.component';
import { EmptyLayoutComponent } from 'app/layout/layouts/empty/empty.component';

// @formatter:off
// tslint:disable:max-line-length
export const appRoutes: Route[] = [

    // Redirect empty path to '/example'
    {path: '', pathMatch : 'full', redirectTo: 'dashboard'},


    // Landing routes
    {
        path: '',
        component: EmptyLayoutComponent,
        children   : [
            {path: 'home', loadChildren: () => import('app/modules/landing/home/home.module').then(m => m.LandingHomeModule)},
        ]
    },

    // Admin routes
    {
        path       : '',
        component  : LayoutComponent,
        children   : [

            // Example
            {path: 'dashboard', loadChildren: () => import('app/modules/dashboard/dashboard.module').then(m => m.DashboardModule)},
            {path: 'device/:wwn', loadChildren: () => import('app/modules/detail/detail.module').then(m => m.DetailModule)}

            // 404 & Catch all
            // {path: '404-not-found', pathMatch: 'full', loadChildren: () => import('app/modules/admin/pages/errors/error-404/error-404.module').then(m => m.Error404Module)},
            // {path: '**', redirectTo: '404-not-found'}
        ]
    }
];
