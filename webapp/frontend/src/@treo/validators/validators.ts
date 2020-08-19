import { FormGroup, ValidatorFn } from '@angular/forms';

export class TreoValidators
{
    /**
     * Check for empty (optional fields) values
     *
     * @param value
     */
    static isEmptyInputValue(value: any): boolean
    {
        return value == null || value.length === 0;
    }

    /**
     * Must match validator
     *
     * @param controlPath A dot-delimited string values that define the path to the control.
     * @param matchingControlPath A dot-delimited string values that define the path to the matching control.
     */
    static mustMatch(controlPath: string, matchingControlPath: string): ValidatorFn
    {
        return (formGroup: FormGroup): null => {

            // Get the control and matching control
            const control = formGroup.get(controlPath);
            const matchingControl = formGroup.get(matchingControlPath);

            // Return if control or matching control doesn't exist
            if ( !control || !matchingControl )
            {
                return;
            }

            // Delete the mustMatch error to reset the error on the matching control
            if ( matchingControl.hasError('mustMatch') )
            {
                delete matchingControl.errors.mustMatch;
                matchingControl.updateValueAndValidity();
            }

            // Don't validate empty values on the matching control
            // Don't validate if values are matching
            if ( this.isEmptyInputValue(matchingControl.value) || control.value === matchingControl.value )
            {
                return;
            }

            // Set the validation error on the matching control
            matchingControl.setErrors({mustMatch: true});
        };
    }
}
