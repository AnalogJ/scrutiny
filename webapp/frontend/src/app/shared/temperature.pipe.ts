import { Pipe, PipeTransform } from '@angular/core';
import {formatNumber} from '@angular/common';

@Pipe({
  name: 'temperature'
})
export class TemperaturePipe implements PipeTransform {
    static celsiusToFahrenheit(celsiusTemp: number): number {
        return celsiusTemp * 9/5 + 32;
    }
    static formatTemperature(temp: number, unit: string, includeUnits: boolean): number|string {
        let unitSuffix
        switch (unit) {
            case 'celsius':
                unitSuffix = '°C'
                break
            case 'fahrenheit':
                unitSuffix = '°F'
                break
        }
        if(includeUnits){
            return formatNumber(temp, 'en-US') + unitSuffix
        } else {
            return formatNumber(temp, 'en-US',)
        }
    }

  transform(celsiusTemp: number, unit = 'celsius', includeUnits = false): number|string {
        let temperature;
        switch (unit) {
            case 'celsius':
                temperature = celsiusTemp;
                break
            case 'fahrenheit':
                temperature = TemperaturePipe.celsiusToFahrenheit(celsiusTemp)
                break
        }
        return TemperaturePipe.formatTemperature(celsiusTemp, unit, includeUnits)
  }

}
