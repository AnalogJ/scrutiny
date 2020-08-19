export class TreoMockApiUtils
{
    /**
     * Constructor
     */
    constructor()
    {

    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Generate a globally unique id
     */
    static guid(): string
    {
        /* tslint:disable */

        let d = new Date().getTime();

        // Use high-precision timer if available
        if ( typeof performance !== 'undefined' && typeof performance.now === 'function' )
        {
            d += performance.now();
        }

        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
            const r = (d + Math.random() * 16) % 16 | 0;
            d = Math.floor(d / 16);
            return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
        });

        /* tslint:enable */
    }
}
