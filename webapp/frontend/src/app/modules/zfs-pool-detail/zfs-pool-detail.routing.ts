import { Route } from '@angular/router';
import { ZFSPoolDetailComponent } from 'app/modules/zfs-pool-detail/zfs-pool-detail.component';
import { ZFSPoolDetailResolver } from './zfs-pool-detail.resolvers';

export const zfsPoolDetailRoutes: Route[] = [
    {
        path: '',
        component: ZFSPoolDetailComponent,
        resolve: {
            pool: ZFSPoolDetailResolver
        }
    }
];
