import { ChangeDetectionStrategy, ChangeDetectorRef, Component, ElementRef, EventEmitter, Input, OnDestroy, OnInit, Output, Renderer2, ViewEncapsulation } from '@angular/core';
import { Subject } from 'rxjs';
import { filter, takeUntil } from 'rxjs/operators';
import { TreoAnimations } from '@treo/animations';
import { TreoMessageAppearance, TreoMessageType } from '@treo/components/message/message.types';
import { TreoMessageService } from '@treo/components/message/message.service';

@Component({
    selector       : 'treo-message',
    templateUrl    : './message.component.html',
    styleUrls      : ['./message.component.scss'],
    encapsulation  : ViewEncapsulation.None,
    changeDetection: ChangeDetectionStrategy.OnPush,
    animations     : TreoAnimations,
    exportAs       : 'treoMessage'
})
export class TreoMessageComponent implements OnInit, OnDestroy
{
    // Name
    @Input()
    name: string;

    @Output()
    readonly afterDismissed: EventEmitter<boolean>;

    @Output()
    readonly afterShown: EventEmitter<boolean>;

    // Private
    private _appearance: TreoMessageAppearance;
    private _dismissed: null | boolean;
    private _showIcon: boolean;
    private _type: TreoMessageType;
    private _unsubscribeAll: Subject<void>;

    /**
     * Constructor
     *
     * @param {TreoMessageService} _treoMessageService
     * @param {ChangeDetectorRef} _changeDetectorRef
     * @param {ElementRef} _elementRef
     * @param {Renderer2} _renderer2
     */
    constructor(
        private _treoMessageService: TreoMessageService,
        private _changeDetectorRef: ChangeDetectorRef,
        private _elementRef: ElementRef,
        private _renderer2: Renderer2
    )
    {
        // Set the private defaults
        this._unsubscribeAll = new Subject();

        // Set the defaults
        this.afterDismissed = new EventEmitter<boolean>();
        this.afterShown = new EventEmitter<boolean>();
        this.appearance = 'fill';
        this.dismissed = null;
        this.showIcon = true;
        this.type = 'primary';
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Setter and getter for appearance
     *
     * @param value
     */
    @Input()
    set appearance(value: TreoMessageAppearance)
    {
        // If the value is the same, return...
        if ( this._appearance === value )
        {
            return;
        }

        // Update the class name
        this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-message-appearance-' + this.appearance);
        this._renderer2.addClass(this._elementRef.nativeElement, 'treo-message-appearance-' + value);

        // Store the value
        this._appearance = value;
    }

    get appearance(): TreoMessageAppearance
    {
        return this._appearance;
    }

    /**
     * Setter and getter for dismissed
     *
     * @param value
     */
    @Input()
    set dismissed(value: null | boolean)
    {
        // If the value is the same, return...
        if ( this._dismissed === value )
        {
            return;
        }

        // Update the class name
        if ( value === null )
        {
            this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-message-dismissible');
        }
        else if ( value === false )
        {
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-message-dismissible');
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-message-dismissed');
        }
        else
        {
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-message-dismissible');
            this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-message-dismissed');
        }

        // Store the value
        this._dismissed = value;
    }

    get dismissed(): null | boolean
    {
        return this._dismissed;
    }

    /**
     * Setter and getter for show icon
     *
     * @param value
     */
    @Input()
    set showIcon(value: boolean)
    {
        // If the value is the same, return...
        if ( this._showIcon === value )
        {
            return;
        }

        // Update the class name
        if ( value )
        {
            this._renderer2.addClass(this._elementRef.nativeElement, 'treo-message-show-icon');
        }
        else
        {
            this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-message-show-icon');
        }

        // Store the value
        this._showIcon = value;
    }

    get showIcon(): boolean
    {
        return this._showIcon;
    }

    /**
     * Setter and getter for type
     *
     * @param value
     */
    @Input()
    set type(value: TreoMessageType)
    {
        // If the value is the same, return...
        if ( this._type === value )
        {
            return;
        }

        // Update the class name
        this._renderer2.removeClass(this._elementRef.nativeElement, 'treo-message-type-' + this.type);
        this._renderer2.addClass(this._elementRef.nativeElement, 'treo-message-type-' + value);

        // Store the value
        this._type = value;
    }

    get type(): TreoMessageType
    {
        return this._type;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Lifecycle hooks
    // -----------------------------------------------------------------------------------------------------

    /**
     * On init
     */
    ngOnInit(): void
    {
        // Subscribe to the service calls if only
        // a name provided for the message box
        if ( this.name )
        {
            // Subscribe to the dismiss calls
            this._treoMessageService.onDismiss
                .pipe(
                    filter((name) => this.name === name),
                    takeUntil(this._unsubscribeAll)
                )
                .subscribe(() => {

                    // Dismiss the message box
                    this.dismiss();
                });

            // Subscribe to the show calls
            this._treoMessageService.onShow
                .pipe(
                    filter((name) => this.name === name),
                    takeUntil(this._unsubscribeAll)
                )
                .subscribe(() => {

                    // Show the message box
                    this.show();
                });
        }
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
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Dismiss the message box
     */
    dismiss(): void
    {
        // Return, if already dismissed
        if ( this.dismissed )
        {
            return;
        }

        // Dismiss
        this.dismissed = true;

        // Execute the observable
        this.afterDismissed.next(true);

        // Notify the change detector
        this._changeDetectorRef.markForCheck();
    }

    /**
     * Show the dismissed message box
     */
    show(): void
    {
        // Return, if not dismissed
        if ( !this.dismissed )
        {
            return;
        }

        // Show
        this.dismissed = false;

        // Execute the observable
        this.afterShown.next(true);

        // Notify the change detector
        this._changeDetectorRef.markForCheck();
    }
}
