import { AfterViewInit, ChangeDetectionStrategy, ChangeDetectorRef, Component, ElementRef, EmbeddedViewRef, Input, Renderer2, SecurityContext, TemplateRef, ViewChild, ViewContainerRef, ViewEncapsulation } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { TreoHighlightService } from '@treo/components/highlight/highlight.service';

@Component({
    selector       : 'textarea[treo-highlight]',
    templateUrl    : './highlight.component.html',
    styleUrls      : ['./highlight.component.scss'],
    encapsulation  : ViewEncapsulation.None,
    changeDetection: ChangeDetectionStrategy.OnPush,
    exportAs       : 'treoHighlight'
})
export class TreoHighlightComponent implements AfterViewInit
{
    highlightedCode: string;
    viewRef: EmbeddedViewRef<any>;

    @ViewChild(TemplateRef)
    templateRef: TemplateRef<any>;

    // Private
    private _code: string;
    private _lang: string;

    /**
     * Constructor
     *
     * @param {TreoHighlightService} _treoHighlightService
     * @param {DomSanitizer} _domSanitizer
     * @param {ChangeDetectorRef} _changeDetectorRef
     * @param {ElementRef} _elementRef
     * @param {Renderer2} _renderer2
     * @param {ViewContainerRef} _viewContainerRef
     */
    constructor(
        private _treoHighlightService: TreoHighlightService,
        private _domSanitizer: DomSanitizer,
        private _changeDetectorRef: ChangeDetectorRef,
        private _elementRef: ElementRef,
        private _renderer2: Renderer2,
        private _viewContainerRef: ViewContainerRef
    )
    {
        // Set the private defaults
        this._code = '';
        this._lang = '';
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Accessors
    // -----------------------------------------------------------------------------------------------------

    /**
     * Setter and getter for the code
     */
    @Input()
    set code(value: string)
    {
        // If the value is the same, return...
        if ( this._code === value )
        {
            return;
        }

        // Set the code
        this._code = value;

        // Highlight and insert the code if the
        // viewContainerRef is available. This will
        // ensure the highlightAndInsert method
        // won't run before the AfterContentInit hook.
        if ( this._viewContainerRef.length )
        {
            this._highlightAndInsert();
        }
    }

    get code(): string
    {
        return this._code;
    }

    /**
     * Setter and getter for the language
     */
    @Input()
    set lang(value: string)
    {
        // If the value is the same, return...
        if ( this._lang === value )
        {
            return;
        }

        // Set the language
        this._lang = value;

        // Highlight and insert the code if the
        // viewContainerRef is available. This will
        // ensure the highlightAndInsert method
        // won't run before the AfterContentInit hook.
        if ( this._viewContainerRef.length )
        {
            this._highlightAndInsert();
        }
    }

    get lang(): string
    {
        return this._lang;
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Lifecycle hooks
    // -----------------------------------------------------------------------------------------------------

    /**
     * After view init
     */
    ngAfterViewInit(): void
    {
        // Return, if there is no language set
        if ( !this.lang )
        {
            return;
        }

        // If there is no code input, get the code from
        // the textarea
        if ( !this.code )
        {
            // Get the code
            this.code = this._elementRef.nativeElement.value;
        }

        // Highlight and insert
        this._highlightAndInsert();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Highlight and insert the highlighted code
     *
     * @private
     */
    private _highlightAndInsert(): void
    {
        // Return, if the code or language is not defined
        if ( !this.code || !this.lang )
        {
            return;
        }

        // Destroy the component if there is already one
        if ( this.viewRef )
        {
            this.viewRef.destroy();
        }

        // Highlight and sanitize the code just in case
        this.highlightedCode = this._domSanitizer.sanitize(SecurityContext.HTML, this._treoHighlightService.highlight(this.code, this.lang));

        // Render and insert the template
        this.viewRef = this._viewContainerRef.createEmbeddedView(this.templateRef, {
            highlightedCode: this.highlightedCode,
            lang           : this.lang
        });

        // Detect the changes
        this.viewRef.detectChanges();
    }
}
