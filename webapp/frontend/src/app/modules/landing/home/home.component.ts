import { Component, ViewEncapsulation } from '@angular/core';

@Component({
    selector: 'landing-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.scss'],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class LandingHomeComponent
{
    /**
     * Constructor
     */
    constructor()
    {
    }
}
