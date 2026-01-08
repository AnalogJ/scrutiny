import { Route } from '@angular/router';
import { ZFSPoolsComponent } from 'app/modules/zfs-pools/zfs-pools.component';
import { ZFSPoolsResolver } from 'app/modules/zfs-pools/zfs-pools.resolvers';

export const zfsPoolsRoutes: Route[] = [
    {
        path: '',
        component: ZFSPoolsComponent,
        resolve: {
            pools: ZFSPoolsResolver
        }
    }
];
