import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class TreoMessageService
{
    // Private
    private _onDismiss: BehaviorSubject<any>;
    private _onShow: BehaviorSubject<any>;

    /**
     * Constructor
     */
    constructor()
    {
        // Set the private defaults
        this._onDismiss = new BehaviorSubject(null);
        this._onShow = new BehaviorSubject(null);
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Getter for onDismiss
     */
    get onDismiss(): Observable<any>
    {
        return this._onDismiss.asObservable();
    }

    /**
     * Getter for onShow
     */
    get onShow(): Observable<any>
    {
        return this._onShow.asObservable();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Dismiss the message box
     *
     * @param name
     */
    dismiss(name: string): void
    {
        // Return, if the name is not provided
        if ( !name )
        {
            return;
        }

        // Execute the observable
        this._onDismiss.next(name);
    }

    /**
     * Show the dismissed message box
     *
     * @param name
     */
    show(name: string): void
    {
        // Return, if the name is not provided
        if ( !name )
        {
            return;
        }

        // Execute the observable
        this._onShow.next(name);
    }

}
