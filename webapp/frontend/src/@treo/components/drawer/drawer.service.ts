import { Injectable } from '@angular/core';
import { TreoDrawerComponent } from '@treo/components/drawer/drawer.component';

@Injectable({
    providedIn: 'root'
})
export class TreoDrawerService
{
    // Private
    private _componentRegistry: Map<string, TreoDrawerComponent>;

    /**
     * Constructor
     */
    constructor()
    {
        // Set the defaults
        this._componentRegistry = new Map<string, TreoDrawerComponent>();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Register drawer component
     *
     * @param name
     * @param component
     */
    registerComponent(name: string, component: TreoDrawerComponent): void
    {
        this._componentRegistry.set(name, component);
    }

    /**
     * Deregister drawer component
     *
     * @param name
     */
    deregisterComponent(name: string): void
    {
        this._componentRegistry.delete(name);
    }

    /**
     * Get drawer component from the registry
     *
     * @param name
     */
    getComponent(name: string): TreoDrawerComponent
    {
        return this._componentRegistry.get(name);
    }
}
