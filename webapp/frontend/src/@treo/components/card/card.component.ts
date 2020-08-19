import { Component, ElementRef, Input, Renderer2, ViewEncapsulation } from '@angular/core';
import { TreoAnimations } from '@treo/animations';

@Component({
    selector     : 'treo-card',
    templateUrl  : './card.component.html',
    styleUrls    : ['./card.component.scss'],
    encapsulation: ViewEncapsulation.None,
    animations   : TreoAnimations,
    exportAs     : 'treoCard'
})
export class TreoCardComponent
{
    expanded: boolean;
    flipped: boolean;

    // Private
    private _flippable: boolean;

    /**
     * Constructor
     *
     * @param {Renderer2} _renderer2
     * @param {ElementRef} _elementRef
     */
    constructor(
        private _renderer2: Renderer2,
        private _elementRef: ElementRef
    )
    {
        // Set the defaults
        this.expanded = false;
        this.flippable = false;
        this.flipped = false;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Setter and getter for flippable
     *
     * @param value
     */
    @Input()
    set flippable(value: boolean)
    {
        // If the value is the same, return...
        if ( this._flippable === value )
        {
            return;
        }

        // Update the class name
        if ( value )
        {
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-card-flippable');
        }
        else
        {
            this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-card-flippable');
        }

        // Store the value
        this._flippable = value;
    }

    get flippable(): boolean
    {
        return this._flippable;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Expand the details
     */
    expand(): void
    {
        this.expanded = true;
    }

    /**
     * Collapse the details
     */
    collapse(): void
    {
        this.expanded = false;
    }

    /**
     * Toggle the expand/collapse status
     */
    toggleExpanded(): void
    {
        this.expanded = !this.expanded;
    }

    /**
     * Flip the card
     */
    flip(): void
    {
        // Return if not flippable
        if ( !this.flippable )
        {
            return;
        }

        this.flipped = !this.flipped;

        // Update the class name
        if ( this.flipped )
        {
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-card-flipped');
        }
        else
        {
            this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-card-flipped');
        }
    }
}
