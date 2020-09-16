import { Route } from '@angular/router';
import { DetailComponent } from 'app/modules/detail/detail.component';
import {DetailResolver} from "./detail.resolvers";

export const detailRoutes: Route[] = [
    {
        path     : '',
        component: DetailComponent,
        resolve  : {
            sales: DetailResolver
        }
    }
];
