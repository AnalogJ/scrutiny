import { Pipe, PipeTransform } from '@angular/core';

/**
 * Finds an object from given source using the given key - value pairs
 */
@Pipe({
    name: 'treoFindByKey',
    pure: false
})
export class TreoFindByKeyPipe implements PipeTransform
{
    /**
     * Constructor
     */
    constructor()
    {
    }

    /**
     * Transform
     *
     * @param value A string or an array of strings to find from source
     * @param key Key of the object property to look for
     * @param source Array of objects to find from
     */
    transform(value: string | string[], key: string, source: any[]): any
    {
        // If the given value is an array of strings...
        if ( Array.isArray(value) )
        {
            return value.map((item) => {
                return source.find((sourceItem) => sourceItem[key] === item);
            });
        }

        // If the value is a string...
        return source.find(sourceItem => sourceItem[key] === value);
    }
}
