import { Component, ElementRef, EventEmitter, HostBinding, HostListener, Input, OnDestroy, OnInit, Output, Renderer2, ViewEncapsulation } from '@angular/core';
import { animate, AnimationBuilder, AnimationPlayer, style } from '@angular/animations';
import { TreoDrawerMode, TreoDrawerPosition } from '@treo/components/drawer/drawer.types';
import { TreoDrawerService } from '@treo/components/drawer/drawer.service';

@Component({
    selector: 'treo-drawer',
    templateUrl: './drawer.component.html',
    styleUrls: ['./drawer.component.scss'],
    encapsulation: ViewEncapsulation.None,
    exportAs: 'treoDrawer',
    standalone: false
})
export class TreoDrawerComponent implements OnInit, OnDestroy
{
    // Name
    @Input()
    name: string;

    // Private
    private _fixed: boolean;
    private _mode: TreoDrawerMode;
    private _opened: boolean | '';
    private _overlay: HTMLElement | null;
    private _player: AnimationPlayer;
    private _position: TreoDrawerPosition;
    private _transparentOverlay: boolean | '';

    // On fixed changed
    @Output()
    readonly fixedChanged: EventEmitter<boolean>;

    // On mode changed
    @Output()
    readonly modeChanged: EventEmitter<TreoDrawerMode>;

    // On opened changed
    @Output()
    readonly openedChanged: EventEmitter<boolean | ''>;

    // On position changed
    @Output()
    readonly positionChanged: EventEmitter<TreoDrawerPosition>;

    @HostBinding('class.treo-drawer-animations-enabled')
    protected _animationsEnabled: boolean;

    /**
     * Constructor
     *
     * @param {AnimationBuilder} _animationBuilder
     * @param {TreoDrawerService} _treoDrawerService
     * @param {ElementRef} _elementRef
     * @param {Renderer2} _renderer2
     */
    constructor(
        private _animationBuilder: AnimationBuilder,
        private _treoDrawerService: TreoDrawerService,
        private _elementRef: ElementRef,
        private _renderer2: Renderer2
    )
    {
        // Set the private defaults
        this._animationsEnabled = false;
        this._overlay = null;

        // Set the defaults
        this.fixedChanged = new EventEmitter<boolean>();
        this.modeChanged = new EventEmitter<TreoDrawerMode>();
        this.openedChanged = new EventEmitter<boolean | ''>();
        this.positionChanged = new EventEmitter<TreoDrawerPosition>();

        this.fixed = false;
        this.mode = 'side';
        this.opened = false;
        this.position = 'left';
        this.transparentOverlay = false;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Setter & getter for fixed
     *
     * @param value
     */
    @Input()
    set fixed(value: boolean)
    {
        // If the value is the same, return...
        if ( this._fixed === value )
        {
            return;
        }

        // Store the fixed value
        this._fixed = value;

        // Update the class
        if ( this.fixed )
        {
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-drawer-fixed');
        }
        else
        {
            this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-drawer-fixed');
        }

        // Execute the observable
        this.fixedChanged.next(this.fixed);
    }

    get fixed(): boolean
    {
        return this._fixed;
    }

    /**
     * Setter & getter for mode
     *
     * @param value
     */
    @Input()
    set mode(value: TreoDrawerMode)
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
            // If the drawer is opened
            if ( this.opened )
            {
                // Show the overlay
                this._showOverlay();
            }
        }

        let modeClassName;

        // Remove the previous mode class
        modeClassName = 'treo-drawer-mode-' + this.mode;
        this._renderer2.removeClass(this._elementRef.nativeElement, modeClassName);

        // Store the mode
        this._mode = value;

        // Add the new mode class
        modeClassName = 'treo-drawer-mode-' + this.mode;
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

    get mode(): TreoDrawerMode
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

        // If the drawer opened, and the mode
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

        // Update opened classes
        if ( this.opened )
        {
            this._renderer2.setStyle(this._elementRef.nativeElement, 'visibility', 'visible');
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-drawer-opened');
        }
        else
        {
            this._renderer2.setStyle(this._elementRef.nativeElement, 'visibility', 'hidden');
            this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-drawer-opened');
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
    set position(value: TreoDrawerPosition)
    {
        // If the value is the same, return...
        if ( this._position === value )
        {
            return;
        }

        let positionClassName;

        // Remove the previous position class
        positionClassName = 'treo-drawer-position-' + this.position;
        this._renderer2.removeClass(this._elementRef.nativeElement, positionClassName);

        // Store the position
        this._position = value;

        // Add the new position class
        positionClassName = 'treo-drawer-position-' + this.position;
        this._renderer2.addClass(this._elementRef.nativeElement, positionClassName);

        // Execute the observable
        this.positionChanged.next(this.position);
    }

    get position(): TreoDrawerPosition
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
        // Register the drawer
        this._treoDrawerService.registerComponent(this.name, this);
    }

    /**
     * On destroy
     */
    ngOnDestroy(): void
    {
        // Deregister the drawer from the registry
        this._treoDrawerService.deregisterComponent(this.name);
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
     * Show the backdrop
     *
     * @private
     */
    private _showOverlay(): void
    {
        // Create the backdrop element
        this._overlay = this._renderer2.createElement('div');

        // Add a class to the backdrop element
        this._overlay.classList.add('treo-drawer-overlay');

        // Add a class depending on the fixed option
        if ( this.fixed )
        {
            this._overlay.classList.add('treo-drawer-overlay-fixed');
        }

        // Add a class depending on the transparentOverlay option
        if ( this.transparentOverlay )
        {
            this._overlay.classList.add('treo-drawer-overlay-transparent');
        }

        // Append the backdrop to the parent of the drawer
        this._renderer2.appendChild(this._elementRef.nativeElement.parentElement, this._overlay);

        // Create the enter animation and attach it to the player
        this._player =
            this._animationBuilder
                .build([
                    animate('300ms cubic-bezier(0.25, 0.8, 0.25, 1)', style({opacity: 1}))
                ]).create(this._overlay);

        // Play the animation
        this._player.play();

        // Add an event listener to the overlay
        this._overlay.addEventListener('click', () => {
            this.close();
        });
    }

    /**
     * Hide the backdrop
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

            // If the backdrop still exists...
            if ( this._overlay )
            {
                // Remove the backdrop
                this._overlay.parentNode.removeChild(this._overlay);
                this._overlay = null;
            }
        });
    }

    /**
     * On mouseenter
     *
     * @protected
     */
    @HostListener('mouseenter')
    protected _onMouseenter(): void
    {
        // Enable the animations
        this._enableAnimations();

        // Add a class
        this._renderer2.addClass(this._elementRef.nativeElement, 'treo-drawer-hover');
    }

    /**
     * On mouseleave
     *
     * @protected
     */
    @HostListener('mouseleave')
    protected _onMouseleave(): void
    {
        // Enable the animations
        this._enableAnimations();

        // Remove the class
        this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-drawer-hover');
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Open the drawer
     */
    open(): void
    {
        // Enable the animations
        this._enableAnimations();

        // Open
        this.opened = true;
    }

    /**
     * Close the drawer
     */
    close(): void
    {
        // Enable the animations
        this._enableAnimations();

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
}
