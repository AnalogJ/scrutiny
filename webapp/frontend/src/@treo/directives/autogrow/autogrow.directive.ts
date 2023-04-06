import { Directive, ElementRef, HostBinding, HostListener, Input, OnDestroy, OnInit, Renderer2 } from '@angular/core';
import { Subject } from 'rxjs';

@Directive({
    selector: 'textarea[treoAutogrow]',
    exportAs: 'treoAutogrow'
})
export class TreoAutogrowDirective implements OnInit, OnDestroy
{
    @HostBinding('rows')
    rows: number;

    // Private
    private _padding: number;
    private _unsubscribeAll: Subject<void>;

    /**
     * Constructor
     *
     * @param {ElementRef} _elementRef
     * @param {Renderer2} _renderer2
     */
    constructor(
        private _elementRef: ElementRef,
        private _renderer2: Renderer2
    )
    {
        // Set the private defaults
        this._unsubscribeAll = new Subject();

        // Set the defaults
        this.padding = 8;
        this.rows = 1;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Setter and getter for padding
     *
     * @param value
     */
    @Input('treoAutogrowVerticalPadding')
    set padding(value)
    {
        // Store the value
        this._padding = value;
    }

    get padding(): number
    {
        return this._padding;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Lifecycle hooks
    // -----------------------------------------------------------------------------------------------------

    /**
     * On init
     */
    ngOnInit(): void
    {
        // Set base styles
        this._renderer2.setStyle(this._elementRef.nativeElement, 'resize', 'none');
        this._renderer2.setStyle(this._elementRef.nativeElement, 'overflow', 'hidden');

        // Set the height for the first time
        setTimeout(() => {
            this._resize();
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
     * Resize on 'input' and 'ngModelChange' events
     *
     * @private
     */
    @HostListener('input')
    @HostListener('ngModelChange')
    private _resize(): void
    {
        // Set the height to 'auto' so we can correctly read the scrollHeight
        this._renderer2.setStyle(this._elementRef.nativeElement, 'height', 'auto');

        // Get the scrollHeight and subtract the vertical padding
        const height = this._elementRef.nativeElement.scrollHeight - this.padding + 'px';
        this._renderer2.setStyle(this._elementRef.nativeElement, 'height', height);
    }
}
