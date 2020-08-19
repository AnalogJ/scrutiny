import { Injectable } from '@angular/core';
import { BreakpointObserver, BreakpointState } from '@angular/cdk/layout';
import { BehaviorSubject, Observable } from 'rxjs';
import { treoBreakpoints } from '@treo/tailwind/exported/variables';

@Injectable()
export class TreoMediaWatcherService
{
    private _onMediaChange: BehaviorSubject<{ matchingAliases: string[], matchingRules: any }>;

    /**
     * Constructor
     *
     * @param {BreakpointObserver} _breakpointObserver
     */
    constructor(
        private _breakpointObserver: BreakpointObserver
    )
    {
        // Set the defaults
        this._onMediaChange = new BehaviorSubject(null);

        // Initialize
        this._init();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Getter for _onMediaChange
     */
    get onMediaChange$(): Observable<{ matchingAliases: string[], matchingRules: any }>
    {
        return this._onMediaChange.asObservable();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Initialize
     *
     * @private
     */
    private _init(): void
    {
        // Subscribe to the breakpoint observer
        this._breakpointObserver.observe(Object.values(treoBreakpoints))
            .subscribe((state) => {

                const matchingAliases = [];
                const matchingRules = {};

                // If there are no matching rules, execute the observable and bail
                if ( !state.matches )
                {
                    this._onMediaChange.next({
                        matchingAliases,
                        matchingRules
                    });

                    return;
                }

                // Go through the breakpoints and find the ones that match
                for ( const [query, matches] of Object.entries(state.breakpoints) )
                {
                    if ( !matches )
                    {
                        continue;
                    }

                    // Get the alias of the matching query
                    const alias = Object.keys(treoBreakpoints).find(key => treoBreakpoints[key] === query);

                    // Prepare the observable values
                    matchingAliases.push(alias);
                    matchingRules[alias] = query;
                }

                // Execute the observable
                this._onMediaChange.next({
                    matchingAliases,
                    matchingRules
                });
            });
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * On media query change
     *
     * @param query
     */
    onMediaQueryChange$(query: string): Observable<BreakpointState>
    {
        return this._breakpointObserver.observe(query);
    }
}
