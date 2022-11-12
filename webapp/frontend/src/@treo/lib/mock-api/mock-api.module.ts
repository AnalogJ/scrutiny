import { APP_INITIALIZER, ModuleWithProviders, NgModule } from '@angular/core';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { TreoMockApiInterceptor } from '@treo/lib/mock-api/mock-api.interceptor';
import { TreoMockApiService } from '@treo/lib/mock-api/mock-api.service';

@NgModule({
    providers: [
        TreoMockApiService,
        {
            provide : HTTP_INTERCEPTORS,
            useClass: TreoMockApiInterceptor,
            multi   : true
        }
    ]
})
export class TreoMockApiModule
{
    /**
     * forRoot method for setting user configuration
     *
     * @param mockDataServices
     */
    static forRoot(mockDataServices: any[]): ModuleWithProviders<TreoMockApiModule>
    {
        return {
            ngModule : TreoMockApiModule,
            providers: [
                {
                    provide   : APP_INITIALIZER,
                    deps      : mockDataServices,
                    useFactory: () => () => null,
                    multi     : true
                },
            ]
        };
    }
}
