import {Component, Inject, OnDestroy, OnInit, ViewEncapsulation} from '@angular/core';
import {DOCUMENT} from '@angular/common';
import {ActivatedRoute, NavigationEnd, Router} from '@angular/router';
import {MatSlideToggleChange} from '@angular/material/slide-toggle';
import {Subject} from 'rxjs';
import {filter, takeUntil} from 'rxjs/operators';
import {ScrutinyConfigService} from 'app/core/config/scrutiny-config.service';
import {TreoDrawerService} from '@treo/components/drawer';
import {Layout} from 'app/layout/layout.types';
import {AppConfig, Theme} from 'app/core/config/app.config';

@Component({
    selector: 'layout',
    templateUrl: './layout.component.html',
    styleUrls: ['./layout.component.scss'],
    encapsulation: ViewEncapsulation.None
})
export class LayoutComponent implements OnInit, OnDestroy {
    config: AppConfig;
    layout: Layout;
    theme: Theme;

    // Private
    private _unsubscribeAll: Subject<void>;
    private systemPrefersDark: boolean;

    /**
     * Constructor
     *
     * @param {ActivatedRoute} _activatedRoute
     * @param {ScrutinyConfigService} _scrutinyConfigService
     * @param {TreoDrawerService} _treoDrawerService
     * @param {DOCUMENT} _document
     * @param {Router} _router
     */
    constructor(
        private _activatedRoute: ActivatedRoute,
        private _scrutinyConfigService: ScrutinyConfigService,
        private _treoDrawerService: TreoDrawerService,
        @Inject(DOCUMENT) private _document: any,
        private _router: Router
    )
    {
        // Set the private defaults
        this._unsubscribeAll = new Subject();

        this.systemPrefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;

    }

    // -----------------------------------------------------------------------------------------------------
    // @ Lifecycle hooks
    // -----------------------------------------------------------------------------------------------------

    /**
     * On init
     */
    ngOnInit(): void
    {
        // Subscribe to config changes
        this._scrutinyConfigService.config$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((config: AppConfig) => {

                // Store the config
                this.config = config;

                // Store the theme
                this.theme = config.theme;

                // Update the selected theme class name on body
                const themeName = 'treo-theme-' + this.determineTheme(config);
                this._document.body.classList.forEach((className) => {
                    if ( className.startsWith('treo-theme-') && className !== themeName )
                    {
                        this._document.body.classList.remove(className);
                        this._document.body.classList.add(themeName);
                        return;
                    }
                });

                // Update the layout
                this._updateLayout();
            });

        // Subscribe to NavigationEnd event
        this._router.events.pipe(
            filter(event => event instanceof NavigationEnd),
            takeUntil(this._unsubscribeAll)
        ).subscribe(() => {

            // Update the layout
            this._updateLayout();
        });
    }

    /**
     * On destroy
     */
    ngOnDestroy(): void
    {
        // Unsubscribe from all subscriptions
        this._unsubscribeAll.next();
        this._unsubscribeAll.complete();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Checks if theme should be set to dark based on config & system settings
     */
    private determineTheme(config:AppConfig): string {
        if (config.theme === 'system') {
            return this.systemPrefersDark ? 'dark' : 'light'
        } else {
            return config.theme
        }
    }

    /**
     * Update the selected layout
     */
    private _updateLayout(): void
    {
        // Get the current activated route
        let route = this._activatedRoute;
        while ( route.firstChild )
        {
            route = route.firstChild;
        }

        // 1. Set the layout from the config
        this.layout = this.config.layout;

        // 2. Get the query parameter from the current route and
        // set the layout and save the layout to the config
        const layoutFromQueryParam = (route.snapshot.queryParamMap.get('layout') as Layout);
        if ( layoutFromQueryParam )
        {
            this.config.layout = this.layout = layoutFromQueryParam;
        }

        // 3. Iterate through the paths and change the layout as we find
        // a config for it.
        //
        // The reason we do this is that there might be empty grouping
        // paths or componentless routes along the path. Because of that,
        // we cannot just assume that the layout configuration will be
        // in the last path's config or in the first path's config.
        //
        // So, we get all the paths that matched starting from root all
        // the way to the current activated route, walk through them one
        // by one and change the layout as we find the layout config. This
        // way, layout configuration can live anywhere within the path and
        // we won't miss it.
        //
        // Also, this will allow overriding the layout in any time so we
        // can have different layouts for different routes.
        const paths = route.pathFromRoot;
        paths.forEach((path) => {

            // Check if there is a 'layout' data
            if ( path.routeConfig && path.routeConfig.data && path.routeConfig.data.layout )
            {
                // Set the layout
                this.layout = path.routeConfig.data.layout;
            }
        });
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Set the layout on the config
     *
     * @param layout
     */
    setLayout(layout: Layout): void {
        // Clear the 'layout' query param to allow layout changes
        this._router.navigate([], {
            queryParams: {
                layout: null
            },
            queryParamsHandling: 'merge'
        }).then(() => {

            // Set the config
            this._scrutinyConfigService.config = {layout};
        });
    }

    /**
     * Set the theme on the config
     *
     * @param change
     */
    setTheme(change: MatSlideToggleChange): void
    {
        this._scrutinyConfigService.config = {theme: change.checked ? 'dark' : 'light'};
    }
}
