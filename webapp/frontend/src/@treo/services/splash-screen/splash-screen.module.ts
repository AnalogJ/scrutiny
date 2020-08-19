import { NgModule } from '@angular/core';
import { TreoSplashScreenService } from '@treo/services/splash-screen/splash-screen.service';

@NgModule({
    providers: [
        TreoSplashScreenService
    ]
})
export class TreoSplashScreenModule
{
    /**
     * Constructor
     *
     * @param {TreoSplashScreenService} _treoSplashScreenService
     */
    constructor(
        private _treoSplashScreenService: TreoSplashScreenService
    )
    {
    }
}
