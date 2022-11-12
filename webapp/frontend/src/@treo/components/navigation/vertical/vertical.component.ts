import { AfterViewInit, ChangeDetectionStrategy, ChangeDetectorRef, Component, ElementRef, EventEmitter, HostBinding, HostListener, Input, OnDestroy, OnInit, Output, QueryList, Renderer2, ViewChild, ViewChildren, ViewEncapsulation } from '@angular/core';
import { animate, AnimationBuilder, AnimationPlayer, style } from '@angular/animations';
import { NavigationEnd, Router } from '@angular/router';
import { ScrollStrategy, ScrollStrategyOptions } from '@angular/cdk/overlay';
import { BehaviorSubject, merge, Subject, Subscription } from 'rxjs';
import { delay, filter, takeUntil } from 'rxjs/operators';
import { TreoAnimations } from '@treo/animations';
import { TreoVerticalNavigationAppearance, TreoNavigationItem, TreoVerticalNavigationMode, TreoVerticalNavigationPosition } from '@treo/components/navigation/navigation.types';
import { TreoNavigationService } from '@treo/components/navigation/navigation.service';
import { TreoScrollbarDirective } from '@treo/directives/scrollbar/scrollbar.directive';

@Component({
    selector       : 'treo-vertical-navigation',
    templateUrl    : './vertical.component.html',
    styleUrls      : ['./vertical.component.scss'],
    animations     : TreoAnimations,
    encapsulation  : ViewEncapsulation.None,
    changeDetection: ChangeDetectionStrategy.OnPush,
    exportAs       : 'treoVerticalNavigation'
})
export class TreoVerticalNavigationComponent implements OnInit, AfterViewInit, OnDestroy
{
    activeAsideItemId: null | string;
    onCollapsableItemCollapsed: BehaviorSubject<TreoNavigationItem | null>;
    onCollapsableItemExpanded: BehaviorSubject<TreoNavigationItem | null>;
    onRefreshed: BehaviorSubject<boolean | null>;

    // Auto collapse
    @Input()
    autoCollapse: boolean;

    // Name
    @Input()
    name: string;

    // On appearance changed
    @Output()
    readonly appearanceChanged: EventEmitter<TreoVerticalNavigationAppearance>;

    // On mode changed
    @Output()
    readonly modeChanged: EventEmitter<TreoVerticalNavigationMode>;

    // On opened changed
    @Output()
    readonly openedChanged: EventEmitter<boolean | ''>;

    // On position changed
    @Output()
    readonly positionChanged: EventEmitter<TreoVerticalNavigationPosition>;

    // Private
    private _appearance: TreoVerticalNavigationAppearance;
    private _asideOverlay: HTMLElement | null;
    private _treoScrollbarDirectives: QueryList<TreoScrollbarDirective>;
    private _treoScrollbarDirectivesSubscription: Subscription;
    private _handleAsideOverlayClick: any;
    private _handleOverlayClick: any;
    private _inner: boolean;
    private _mode: TreoVerticalNavigationMode;
    private _navigation: TreoNavigationItem[];
    private _opened: boolean | '';
    private _overlay: HTMLElement | null;
    private _player: AnimationPlayer;
    private _position: TreoVerticalNavigationPosition;
    private _scrollStrategy: ScrollStrategy;
    private _transparentOverlay: boolean | '';
    private _unsubscribeAll: Subject<void>;

    @HostBinding('class.treo-vertical-navigation-animations-enabled')
    private _animationsEnabled: boolean;

    @ViewChild('navigationContent')
    private _navigationContentEl: ElementRef;

    /**
     * Constructor
     *
     * @param {AnimationBuilder} _animationBuilder
     * @param {TreoNavigationService} _treoNavigationService
     * @param {ChangeDetectorRef} _changeDetectorRef
     * @param {ElementRef} _elementRef
     * @param {Renderer2} _renderer2
     * @param {Router} _router
     * @param {ScrollStrategyOptions} _scrollStrategyOptions
     */
    constructor(
        private _animationBuilder: AnimationBuilder,
        private _treoNavigationService: TreoNavigationService,
        private _changeDetectorRef: ChangeDetectorRef,
        private _elementRef: ElementRef,
        private _renderer2: Renderer2,
        private _router: Router,
        private _scrollStrategyOptions: ScrollStrategyOptions
    )
    {
        // Set the private defaults
        this._animationsEnabled = false;
        this._asideOverlay = null;
        this._handleAsideOverlayClick = () => {
            this.closeAside();
        };
        this._handleOverlayClick = () => {
            this.close();
        };
        this._overlay = null;
        this._scrollStrategy = this._scrollStrategyOptions.block();
        this._unsubscribeAll = new Subject();

        // Set the defaults
        this.appearanceChanged = new EventEmitter<TreoVerticalNavigationAppearance>();
        this.modeChanged = new EventEmitter<TreoVerticalNavigationMode>();
        this.openedChanged = new EventEmitter<boolean | ''>();
        this.positionChanged = new EventEmitter<TreoVerticalNavigationPosition>();

        this.onCollapsableItemCollapsed = new BehaviorSubject(null);
        this.onCollapsableItemExpanded = new BehaviorSubject(null);
        this.onRefreshed = new BehaviorSubject(null);

        this.activeAsideItemId = null;
        this.appearance = 'classic';
        this.autoCollapse = true;
        this.inner = false;
        this.mode = 'side';
        this.opened = false;
        this.position = 'left';
        this.transparentOverlay = false;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Setter & getter for appearance
     *
     * @param value
     */
    @Input()
    set appearance(value: TreoVerticalNavigationAppearance)
    {
        // If the value is the same, return...
        if ( this._appearance === value )
        {
            return;
        }

        let appearanceClassName;

        // Remove the previous appearance class
        appearanceClassName = 'treo-vertical-navigation-appearance-' + this.appearance;
        this._renderer2.removeClass(this._elementRef.nativeElement, appearanceClassName);

        // Store the appearance
        this._appearance = value;

        // Add the new appearance class
        appearanceClassName = 'treo-vertical-navigation-appearance-' + this.appearance;
        this._renderer2.addClass(this._elementRef.nativeElement, appearanceClassName);

        // Execute the observable
        this.appearanceChanged.next(this.appearance);
    }

    get appearance(): TreoVerticalNavigationAppearance
    {
        return this._appearance;
    }

    /**
     * Setter for treoScrollbarDirectives
     */
    @ViewChildren(TreoScrollbarDirective)
    set treoScrollbarDirectives(treoScrollbarDirectives: QueryList<TreoScrollbarDirective>)
    {
        // Store the directives
        this._treoScrollbarDirectives = treoScrollbarDirectives;

        // Return, if there are no directives
        if ( treoScrollbarDirectives.length === 0 )
        {
            return;
        }

        // Unsubscribe the previous subscriptions
        if ( this._treoScrollbarDirectivesSubscription )
        {
            this._treoScrollbarDirectivesSubscription.unsubscribe();
        }

        // Update the scrollbars on collapsable items' collapse/expand
        this._treoScrollbarDirectivesSubscription =
            merge(
                this.onCollapsableItemCollapsed,
                this.onCollapsableItemExpanded
            )
                .pipe(
                    takeUntil(this._unsubscribeAll),
                    delay(250)
                )
                .subscribe(() => {

                    // Loop through the scrollbars and update them
                    treoScrollbarDirectives.forEach((treoScrollbarDirective) => {
                        treoScrollbarDirective.update();
                    });
                });
    }

    /**
     * Setter & getter for data
     */
    @Input()
    set navigation(value: TreoNavigationItem[])
    {
        // Store the data
        this._navigation = value;

        // Mark for check
        this._changeDetectorRef.markForCheck();
    }

    get navigation(): TreoNavigationItem[]
    {
        return this._navigation;
    }

    /**
     * Setter & getter for inner
     *
     * @param value
     */
    @Input()
    set inner(value: boolean)
    {
        // If the value is the same, return...
        if ( this._inner === value )
        {
            return;
        }

        // Set the naked value
        this._inner = value;

        // Update the class
        if ( this.inner )
        {
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-vertical-navigation-inner');
        }
        else
        {
            this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-vertical-navigation-inner');
        }
    }

    get inner(): boolean
    {
        return this._inner;
    }

    /**
     * Setter & getter for mode
     *
     * @param value
     */
    @Input()
    set mode(value: TreoVerticalNavigationMode)
    {
        // If the value is the same, return...
        if ( this._mode === value )
        {
            return;
        }

        // Disable the animations
        this._disableAnimations();

        // If the mode changes: 'over -> side'
        if ( this.mode === 'over' && value === 'side' )
        {
            // Hide the overlay
            this._hideOverlay();
        }

        // If the mode changes: 'side -> over'
        if ( this.mode === 'side' && value === 'over' )
        {
            // Close the aside
            this.closeAside();

            // If the navigation is opened
            if ( this.opened )
            {
                // Show the overlay
                this._showOverlay();
            }
        }

        let modeClassName;

        // Remove the previous mode class
        modeClassName = 'treo-vertical-navigation-mode-' + this.mode;
        this._renderer2.removeClass(this._elementRef.nativeElement, modeClassName);

        // Store the mode
        this._mode = value;

        // Add the new mode class
        modeClassName = 'treo-vertical-navigation-mode-' + this.mode;
        this._renderer2.addClass(this._elementRef.nativeElement, modeClassName);

        // Execute the observable
        this.modeChanged.next(this.mode);

        // Enable the animations after a delay
        // The delay must be bigger than the current transition-duration
        // to make sure nothing will be animated while the mode changing
        setTimeout(() => {
            this._enableAnimations();
        }, 500);
    }

    get mode(): TreoVerticalNavigationMode
    {
        return this._mode;
    }

    /**
     * Setter & getter for opened
     *
     * @param value
     */
    @Input()
    set opened(value: boolean | '')
    {
        // If the value is the same, return...
        if ( this._opened === value )
        {
            return;
        }

        // If the provided value is an empty string,
        // take that as a 'true'
        if ( value === '' )
        {
            value = true;
        }

        // Set the opened value
        this._opened = value;

        // If the navigation opened, and the mode
        // is 'over', show the overlay
        if ( this.mode === 'over' )
        {
            if ( this._opened )
            {
                this._showOverlay();
            }
            else
            {
                this._hideOverlay();
            }
        }

        if ( this.opened )
        {
            // Update styles and classes
            this._renderer2.setStyle(this._elementRef.nativeElement, 'visibility', 'visible');
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-vertical-navigation-opened');
        }
        else
        {
            // Update styles and classes
            this._renderer2.setStyle(this._elementRef.nativeElement, 'visibility', 'hidden');
            this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-vertical-navigation-opened');
        }

        // Execute the observable
        this.openedChanged.next(this.opened);
    }

    get opened(): boolean | ''
    {
        return this._opened;
    }

    /**
     * Setter & getter for position
     *
     * @param value
     */
    @Input()
    set position(value: TreoVerticalNavigationPosition)
    {
        // If the value is the same, return...
        if ( this._position === value )
        {
            return;
        }

        let positionClassName;

        // Remove the previous position class
        positionClassName = 'treo-vertical-navigation-position-' + this.position;
        this._renderer2.removeClass(this._elementRef.nativeElement, positionClassName);

        // Store the position
        this._position = value;

        // Add the new position class
        positionClassName = 'treo-vertical-navigation-position-' + this.position;
        this._renderer2.addClass(this._elementRef.nativeElement, positionClassName);

        // Execute the observable
        this.positionChanged.next(this.position);
    }

    get position(): TreoVerticalNavigationPosition
    {
        return this._position;
    }

    /**
     * Setter & getter for transparent overlay
     *
     * @param value
     */
    @Input()
    set transparentOverlay(value: boolean | '')
    {
        // If the value is the same, return...
        if ( this._opened === value )
        {
            return;
        }

        // If the provided value is an empty string,
        // take that as a 'true' and set the opened value
        if ( value === '' )
        {
            // Set the opened value
            this._transparentOverlay = true;
        }
        else
        {
            // Set the transparent overlay value
            this._transparentOverlay = value;
        }
    }

    get transparentOverlay(): boolean | ''
    {
        return this._transparentOverlay;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Lifecycle hooks
    // -----------------------------------------------------------------------------------------------------

    /**
     * On init
     */
    ngOnInit(): void
    {
        // Register the navigation component
        this._treoNavigationService.registerComponent(this.name, this);

        // Subscribe to the 'NavigationEnd' event
        this._router.events
            .pipe(
                filter(event => event instanceof NavigationEnd),
                takeUntil(this._unsubscribeAll)
            )
            .subscribe(() => {

                if ( this.mode === 'over' && this.opened )
                {
                    // Close the navigation
                    this.close();
                }
            });
    }

    /**
     * After view init
     */
    ngAfterViewInit(): void
    {
        setTimeout(() => {

            // If 'navigation content' element doesn't have
            // perfect scrollbar activated on it...
            if ( !this._navigationContentEl.nativeElement.classList.contains('ps') )
            {
                // Find the active item
                const activeItem = this._navigationContentEl.nativeElement.querySelector('.treo-vertical-navigation-item-active');

                // If the active item exists, scroll it into view
                if ( activeItem )
                {
                    activeItem.scrollIntoView();
                }
            }
            // Otherwise
            else
            {
                // Go through all the scrollbar directives
                this._treoScrollbarDirectives.forEach((treoScrollbarDirective) => {

                    // Skip if not enabled
                    if ( !treoScrollbarDirective.enabled )
                    {
                        return;
                    }

                    // Scroll to the active element
                    treoScrollbarDirective.scrollToElement('.treo-vertical-navigation-item-active', -120, true);
                });
            }
        });
    }

    /**
     * On destroy
     */
    ngOnDestroy(): void
    {
        // Deregister the navigation component from the registry
        this._treoNavigationService.deregisterComponent(this.name);

        // Unsubscribe from all subscriptions
        this._unsubscribeAll.next();
        this._unsubscribeAll.complete();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Enable the animations
     *
     * @private
     */
    private _enableAnimations(): void
    {
        // If the animations are already enabled, return...
        if ( this._animationsEnabled )
        {
            return;
        }

        // Enable the animations
        this._animationsEnabled = true;
    }

    /**
     * Disable the animations
     *
     * @private
     */
    private _disableAnimations(): void
    {
        // If the animations are already disabled, return...
        if ( !this._animationsEnabled )
        {
            return;
        }

        // Disable the animations
        this._animationsEnabled = false;
    }

    /**
     * Show the overlay
     *
     * @private
     */
    private _showOverlay(): void
    {
        // If there is already an overlay, return...
        if ( this._asideOverlay )
        {
            return;
        }

        // Create the overlay element
        this._overlay = this._renderer2.createElement('div');

        // Add a class to the overlay element
        this._overlay.classList.add('treo-vertical-navigation-overlay');

        // Add a class depending on the transparentOverlay option
        if ( this.transparentOverlay )
        {
            this._overlay.classList.add('treo-vertical-navigation-overlay-transparent');
        }

        // Append the overlay to the parent of the navigation
        this._renderer2.appendChild(this._elementRef.nativeElement.parentElement, this._overlay);

        // Enable block scroll strategy
        this._scrollStrategy.enable();

        // Create the enter animation and attach it to the player
        this._player =
            this._animationBuilder
                .build([
                    animate('300ms cubic-bezier(0.25, 0.8, 0.25, 1)', style({opacity: 1}))
                ]).create(this._overlay);

        // Play the animation
        this._player.play();

        // Add an event listener to the overlay
        this._overlay.addEventListener('click', this._handleOverlayClick);
    }

    /**
     * Hide the overlay
     *
     * @private
     */
    private _hideOverlay(): void
    {
        if ( !this._overlay )
        {
            return;
        }

        // Create the leave animation and attach it to the player
        this._player =
            this._animationBuilder
                .build([
                    animate('300ms cubic-bezier(0.25, 0.8, 0.25, 1)', style({opacity: 0}))
                ]).create(this._overlay);

        // Play the animation
        this._player.play();

        // Once the animation is done...
        this._player.onDone(() => {

            // If the overlay still exists...
            if ( this._overlay )
            {
                // Remove the event listener
                this._overlay.removeEventListener('click', this._handleOverlayClick);

                // Remove the overlay
                this._overlay.parentNode.removeChild(this._overlay);
                this._overlay = null;
            }

            // Disable block scroll strategy
            this._scrollStrategy.disable();
        });
    }

    /**
     * Show the aside overlay
     *
     * @private
     */
    private _showAsideOverlay(): void
    {
        // If there is already an overlay, return...
        if ( this._asideOverlay )
        {
            return;
        }

        // Create the aside overlay element
        this._asideOverlay = this._renderer2.createElement('div');

        // Add a class to the aside overlay element
        this._asideOverlay.classList.add('treo-vertical-navigation-aside-overlay');

        // Append the aside overlay to the parent of the navigation
        this._renderer2.appendChild(this._elementRef.nativeElement.parentElement, this._asideOverlay);

        // Create the enter animation and attach it to the player
        this._player =
            this._animationBuilder
                .build([
                    animate('300ms cubic-bezier(0.25, 0.8, 0.25, 1)', style({opacity: 1}))
                ]).create(this._asideOverlay);

        // Play the animation
        this._player.play();

        // Add an event listener to the aside overlay
        this._asideOverlay.addEventListener('click', this._handleAsideOverlayClick);
    }

    /**
     * Hide the aside overlay
     *
     * @private
     */
    private _hideAsideOverlay(): void
    {
        if ( !this._asideOverlay )
        {
            return;
        }

        // Create the leave animation and attach it to the player
        this._player =
            this._animationBuilder
                .build([
                    animate('300ms cubic-bezier(0.25, 0.8, 0.25, 1)', style({opacity: 0}))
                ]).create(this._asideOverlay);

        // Play the animation
        this._player.play();

        // Once the animation is done...
        this._player.onDone(() => {

            // If the aside overlay still exists...
            if ( this._asideOverlay )
            {
                // Remove the event listener
                this._asideOverlay.removeEventListener('click', this._handleAsideOverlayClick);

                // Remove the aside overlay
                this._asideOverlay.parentNode.removeChild(this._asideOverlay);
                this._asideOverlay = null;
            }
        });
    }

    /**
     * On mouseenter
     *
     * @private
     */
    @HostListener('mouseenter')
    private _onMouseenter(): void
    {
        // Enable the animations
        this._enableAnimations();

        // Add a class
        this._renderer2.addClass(this._elementRef.nativeElement, 'treo-vertical-navigation-hover');
    }

    /**
     * On mouseleave
     *
     * @private
     */
    @HostListener('mouseleave')
    private _onMouseleave(): void
    {
        // Enable the animations
        this._enableAnimations();

        // Remove the class
        this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-vertical-navigation-hover');
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Refresh the component to apply the changes
     */
    refresh(): void
    {
        // Mark for check
        this._changeDetectorRef.markForCheck();

        // Execute the observable
        this.onRefreshed.next(true);
    }

    /**
     * Open the navigation
     */
    open(): void
    {
        // Enable the animations
        this._enableAnimations();

        // Open
        this.opened = true;
    }

    /**
     * Close the navigation
     */
    close(): void
    {
        // Enable the animations
        this._enableAnimations();

        // Close the aside
        this.closeAside();

        // Close
        this.opened = false;
    }

    /**
     * Toggle the opened status
     */
    toggle(): void
    {
        // Toggle
        if ( this.opened )
        {
            this.close();
        }
        else
        {
            this.open();
        }
    }

    /**
     * Open the aside
     *
     * @param item
     */
    openAside(item: TreoNavigationItem): void
    {
        // Return if the item is disabled
        if ( item.disabled )
        {
            return;
        }

        // Open
        this.activeAsideItemId = item.id;

        // Show the aside overlay
        this._showAsideOverlay();

        // Mark for check
        this._changeDetectorRef.markForCheck();
    }

    /**
     * Close the aside
     */
    closeAside(): void
    {
        // Close
        this.activeAsideItemId = null;

        // Hide the aside overlay
        this._hideAsideOverlay();

        // Mark for check
        this._changeDetectorRef.markForCheck();
    }

    /**
     * Toggle the aside
     *
     * @param item
     */
    toggleAside(item: TreoNavigationItem): void
    {
        // Toggle
        if ( this.activeAsideItemId === item.id )
        {
            this.closeAside();
        }
        else
        {
            this.openAside(item);
        }
    }

    /**
     * Track by function for ngFor loops
     *
     * @param index
     * @param item
     */
    trackByFn(index: number, item: any): any
    {
        return item.id || index;
    }
}
