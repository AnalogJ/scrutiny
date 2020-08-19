import { Injectable } from '@angular/core';
import { TreoMockApiRequestHandler } from '@treo/lib/mock-api/mock-api.request-handler';

@Injectable({
    providedIn: 'root'
})
export class TreoMockApiService
{
    requestHandlers: any;

    /**
     * Constructor
     */
    constructor()
    {
        // Set the defaults
        this.requestHandlers = {
            delete: new Map<string, TreoMockApiRequestHandler>(),
            get   : new Map<string, TreoMockApiRequestHandler>(),
            patch : new Map<string, TreoMockApiRequestHandler>(),
            post  : new Map<string, TreoMockApiRequestHandler>(),
            put   : new Map<string, TreoMockApiRequestHandler>()
        };
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Register 'delete' request handler
     *
     * @param url
     * @param delay
     */
    onDelete(url: string, delay: number = 0): TreoMockApiRequestHandler
    {
        return this._registerRequestHandler('delete', url, delay);
    }

    /**
     * Register 'get' request handler
     *
     * @param url
     * @param delay
     */
    onGet(url: string, delay: number = 0): TreoMockApiRequestHandler
    {
        return this._registerRequestHandler('get', url, delay);
    }

    /**
     * Register 'patch' request handler
     *
     * @param url
     * @param delay
     */
    onPatch(url: string, delay: number = 0): TreoMockApiRequestHandler
    {
        return this._registerRequestHandler('patch', url, delay);
    }

    /**
     * Register 'post' request handler
     *
     * @param url
     * @param delay
     */
    onPost(url: string, delay: number = 0): TreoMockApiRequestHandler
    {
        return this._registerRequestHandler('post', url, delay);
    }

    /**
     * Register 'put' request handler
     *
     * @param url
     * @param delay
     */
    onPut(url: string, delay: number = 0): TreoMockApiRequestHandler
    {
        return this._registerRequestHandler('put', url, delay);
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Register a request handler
     *
     * @param requestType
     * @param url
     * @param delay
     * @private
     */
    private _registerRequestHandler(requestType, url, delay): TreoMockApiRequestHandler
    {
        // Create a new instance of TreoMockApiRequestHandler
        const treoMockHttp = new TreoMockApiRequestHandler();

        // Store the url
        treoMockHttp.url = url;

        // Store the delay
        treoMockHttp.delay = delay;

        // Store the request handler to access them from the interceptor
        this.requestHandlers[requestType].set(url, treoMockHttp);

        // Return the instance
        return treoMockHttp;
    }
}
