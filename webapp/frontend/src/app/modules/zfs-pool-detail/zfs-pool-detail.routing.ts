import { Route } from '@angular/router';
import { ZfsPoolDetailComponent } from './zfs-pool-detail.component';
import { ZfsPoolDetailResolver } from './zfs-pool-detail.resolvers';

export const zfsPoolDetailRoutes: Route[] = [
    {
        path: '',
        component: ZfsPoolDetailComponent,
        resolve: {
            poolDetail: ZfsPoolDetailResolver
        }
    }
];