import { Directive, ElementRef, Input, OnDestroy, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Platform } from '@angular/cdk/platform';
import { fromEvent, Subject } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';
import PerfectScrollbar from 'perfect-scrollbar';
import * as _ from 'lodash';
import { ScrollbarGeometry, ScrollbarPosition } from '@treo/directives/scrollbar/scrollbar.interfaces';

// -----------------------------------------------------------------------------------------------------
// Wrapper directive for the Perfect Scrollbar: https://github.com/mdbootstrap/perfect-scrollbar
// Based on https://github.com/zefoy/ngx-perfect-scrollbar
// -----------------------------------------------------------------------------------------------------
@Directive({
    selector: '[treoScrollbar]',
    exportAs: 'treoScrollbar',
    standalone: false
})
export class TreoScrollbarDirective implements OnInit, OnDestroy
{
    isMobile: boolean;
    ps: PerfectScrollbar | any;

    // Private
    private _animation: number | null;
    private _enabled: boolean;
    private _options: any;
    private _unsubscribeAll: Subject<void>;

    /**
     * Constructor
     *
     * @param {ElementRef} _elementRef
     * @param {Platform} _platform
     * @param {Router} _router
     */
    constructor(
        private _elementRef: ElementRef,
        private _platform: Platform,
        private _router: Router
    )
    {
        // Set the private defaults
        this._animation = null;
        this._options = {};
        this._unsubscribeAll = new Subject();

        // Set the defaults
        this.enabled = true;
        this.isMobile = false;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Scrollbar options
     *
     * @param value
     */
    @Input()
    set treoScrollbarOptions(value: any)
    {
        // Merge the options
        this._options = _.merge({}, this._options, value);

        // Destroy and re-init the PerfectScrollbar to update its options
        setTimeout(() => {
            this._destroy();
        });

        setTimeout(() => {
            this._init();
        });
    }

    get treoScrollbarOptions(): any
    {
        // Return the options
        return this._options;
    }

    /**
     * Is enabled
     *
     * @param value
     */
    @Input('treoScrollbar')
    set enabled(value: boolean | '')
    {
        // If the value is an empty string, interpret it as 'true'
        if ( value === '' )
        {
            value = true;
        }

        // If the value is the same, return...
        if ( this._enabled === value )
        {
            return;
        }

        // Store the value
        this._enabled = value;

        // If enabled...
        if ( this.enabled )
        {
            // Init the directive
            this._init();
        }
        else
        {
            // Otherwise destroy it
            this._destroy();
        }
    }

    get enabled(): boolean | ''
    {
        // Return the enabled status
        return this._enabled;
    }

    /**
     * Getter for _elementRef
     */
    get elementRef(): ElementRef
    {
        return this._elementRef;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Lifecycle hooks
    // -----------------------------------------------------------------------------------------------------

    /**
     * On init
     */
    ngOnInit(): void
    {
        // Subscribe to window resize event
        fromEvent(window, 'resize')
            .pipe(
                takeUntil(this._unsubscribeAll),
                debounceTime(150)
            )
            .subscribe(() => {

                // Update the PerfectScrollbar
                this.update();
            });
    }

    /**
     * On destroy
     */
    ngOnDestroy(): void
    {
        this._destroy();

        // Unsubscribe from all subscriptions
        this._unsubscribeAll.next();
        this._unsubscribeAll.complete();
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
        // Return, if already initialized
        if ( this.ps )
        {
            return;
        }

        // Check if is mobile
        if ( this._platform.ANDROID || this._platform.IOS )
        {
            this.isMobile = true;
        }

        // Return if it's mobile or the platform is not a browser
        if ( this.isMobile || !this._platform.isBrowser )
        {
            // Silently set the enabled to false
            this._enabled = false;

            return;
        }

        // Initialize the PerfectScrollbar
        this.ps = new PerfectScrollbar(this._elementRef.nativeElement, {...this.treoScrollbarOptions});
    }

    /**
     * Destroy
     *
     * @private
     */
    private _destroy(): void
    {
        if ( !this.ps )
        {
            return;
        }

        // Destroy the PerfectScrollbar
        this.ps.destroy();

        // Clean up
        this.ps = null;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Update the scrollbar
     */
    update(): void
    {
        if ( !this.ps )
        {
            return;
        }

        // Update the PerfectScrollbar
        this.ps.update();
    }

    /**
     * Destroy the scrollbar
     */
    destroy(): void
    {
        this.ngOnDestroy();
    }

    /**
     * Returns the geometry of the scrollable element
     *
     * @param prefix
     */
    geometry(prefix: string = 'scroll'): ScrollbarGeometry
    {
        const scrollbarGeometry = new ScrollbarGeometry(
            this._elementRef.nativeElement[prefix + 'Left'],
            this._elementRef.nativeElement[prefix + 'Top'],
            this._elementRef.nativeElement[prefix + 'Width'],
            this._elementRef.nativeElement[prefix + 'Height']);

        return scrollbarGeometry;
    }

    /**
     * Returns the position of the scrollable element
     *
     * @param absolute
     */
    position(absolute: boolean = false): ScrollbarPosition
    {
        let scrollbarPosition;

        if ( !absolute && this.ps )
        {
            scrollbarPosition = new ScrollbarPosition(
                this.ps.reach.x || 0,
                this.ps.reach.y || 0
            );
        }
        else
        {
            scrollbarPosition = new ScrollbarPosition(
                this._elementRef.nativeElement.scrollLeft,
                this._elementRef.nativeElement.scrollTop
            );
        }

        return scrollbarPosition;
    }

    /**
     * Scroll to
     *
     * @param x
     * @param y
     * @param speed
     */
    scrollTo(x: number, y?: number, speed?: number): void
    {
        if ( y == null && speed == null )
        {
            this.animateScrolling('scrollTop', x, speed);
        }
        else
        {
            if ( x != null )
            {
                this.animateScrolling('scrollLeft', x, speed);
            }

            if ( y != null )
            {
                this.animateScrolling('scrollTop', y, speed);
            }
        }
    }

    /**
     * Scroll to X
     *
     * @param {number} x
     * @param {number} speed
     */
    scrollToX(x: number, speed?: number): void
    {
        this.animateScrolling('scrollLeft', x, speed);
    }

    /**
     * Scroll to Y
     *
     * @param {number} y
     * @param {number} speed
     */
    scrollToY(y: number, speed?: number): void
    {
        this.animateScrolling('scrollTop', y, speed);
    }

    /**
     * Scroll to top
     *
     * @param {number} offset
     * @param {number} speed
     */
    scrollToTop(offset: number = 0, speed?: number): void
    {
        this.animateScrolling('scrollTop', offset, speed);
    }

    /**
     * Scroll to bottom
     *
     * @param {number} offset
     * @param {number} speed
     */
    scrollToBottom(offset: number = 0, speed?: number): void
    {
        const top = this._elementRef.nativeElement.scrollHeight - this._elementRef.nativeElement.clientHeight;
        this.animateScrolling('scrollTop', top - offset, speed);
    }

    /**
     * Scroll to left
     *
     * @param {number} offset
     * @param {number} speed
     */
    scrollToLeft(offset: number = 0, speed?: number): void
    {
        this.animateScrolling('scrollLeft', offset, speed);
    }

    /**
     * Scroll to right
     *
     * @param {number} offset
     * @param {number} speed
     */
    scrollToRight(offset: number = 0, speed?: number): void
    {
        const left = this._elementRef.nativeElement.scrollWidth - this._elementRef.nativeElement.clientWidth;
        this.animateScrolling('scrollLeft', left - offset, speed);
    }

    /**
     * Scroll to element
     *
     * @param {string} qs
     * @param {number} offset
     * @param {boolean} ignoreVisible If true, scrollToElement won't happen if element is already inside the current viewport
     * @param {number} speed
     */
    scrollToElement(qs: string, offset: number = 0, ignoreVisible: boolean = false, speed?: number): void
    {
        const element = this._elementRef.nativeElement.querySelector(qs);

        if ( !element )
        {
            return;
        }

        const elementPos = element.getBoundingClientRect();
        const scrollerPos = this._elementRef.nativeElement.getBoundingClientRect();

        if ( this._elementRef.nativeElement.classList.contains('ps--active-x') )
        {
            if ( ignoreVisible && elementPos.right <= (scrollerPos.right - Math.abs(offset)) )
            {
                return;
            }

            const currentPos = this._elementRef.nativeElement['scrollLeft'];
            const position = elementPos.left - scrollerPos.left + currentPos;

            this.animateScrolling('scrollLeft', position + offset, speed);
        }

        if ( this._elementRef.nativeElement.classList.contains('ps--active-y') )
        {
            if ( ignoreVisible && elementPos.bottom <= (scrollerPos.bottom - Math.abs(offset)) )
            {
                return;
            }

            const currentPos = this._elementRef.nativeElement['scrollTop'];
            const position = elementPos.top - scrollerPos.top + currentPos;

            this.animateScrolling('scrollTop', position + offset, speed);
        }
    }

    /**
     * Animate scrolling
     *
     * @param target
     * @param value
     * @param speed
     */
    animateScrolling(target: string, value: number, speed?: number): void
    {
        if ( this._animation )
        {
            window.cancelAnimationFrame(this._animation);
            this._animation = null;
        }

        if ( !speed || typeof window === 'undefined' )
        {
            this._elementRef.nativeElement[target] = value;
        }
        else if ( value !== this._elementRef.nativeElement[target] )
        {
            let newValue = 0;
            let scrollCount = 0;

            let oldTimestamp = performance.now();
            let oldValue = this._elementRef.nativeElement[target];

            const cosParameter = (oldValue - value) / 2;

            const step = (newTimestamp: number) => {
                scrollCount += Math.PI / (speed / (newTimestamp - oldTimestamp));
                newValue = Math.round(value + cosParameter + cosParameter * Math.cos(scrollCount));

                // Only continue animation if scroll position has not changed
                if ( this._elementRef.nativeElement[target] === oldValue )
                {
                    if ( scrollCount >= Math.PI )
                    {
                        this.animateScrolling(target, value, 0);
                    }
                    else
                    {
                        this._elementRef.nativeElement[target] = newValue;

                        // On a zoomed out page the resulting offset may differ
                        oldValue = this._elementRef.nativeElement[target];
                        oldTimestamp = newTimestamp;

                        this._animation = window.requestAnimationFrame(step);
                    }
                }
            };

            window.requestAnimationFrame(step);
        }
    }
}
