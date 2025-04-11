import { TemperaturePipe } from './temperature.pipe';

describe('TemperaturePipe', () => {
    it('create an instance', () => {
        const pipe = new TemperaturePipe();
        expect(pipe).toBeTruthy();
    });


    describe('#celsiusToFahrenheit', () => {
        const testCases = [
            {
                'c': -273.15,
                'f': -459.66999999999996,
            },{
                'c': -34.44,
                'f': -29.991999999999997,
            },{
                'c': -23.33,
                'f': -9.993999999999993,
            },{
                'c': -17.78,
                'f': -0.0040000000000048885,
            },{
                'c': 0,
                'f': 32,
            },{
                'c': 10,
                'f': 50,
            },{
                'c': 26.67,
                'f': 80.006,
            },{
                'c': 37,
                'f': 98.6,
            },{
                'c': 60,
                'f': 140,
            }
        ]
        testCases.forEach((test, index) => {
            it(`should correctly convert ${test.c}, Celsius to Fahrenheit (testcase: ${index + 1})`, () => {
                // test
                const numb = TemperaturePipe.celsiusToFahrenheit(test.c)
                expect(numb).toEqual(test.f);
            });
        })
    });

    describe('#formatTemperature',() => {
        const testCases = [
            {
                'c': 26.67,
                'unit': 'celsius',
                'includeUnits': true,
                'result': '26.67°C'
            },{
                'c': 26.6767,
                'unit': 'celsius',
                'includeUnits': true,
                'result': '26.677°C'
            },{
                'c': 26.67,
                'unit': 'celsius',
                'includeUnits': false,
                'result': '26.67',
            },{
                'c': 26.67,
                'unit': 'fahrenheit',
                'includeUnits': true,
                'result': '26.67°F',
            },{
                'c': 26.6767,
                'unit': 'fahrenheit',
                'includeUnits': true,
                'result': '26.677°F',
            },{
                'c': 26.67,
                'unit': 'fahrenheit',
                'includeUnits': false,
                'result': '26.67',
            }
        ]
        testCases.forEach((test, index) => {
            it(`should correctly format temperature ${test.c} to ${test.unit} ${test.includeUnits ? 'with' : 'without'} unit. (testcase: ${index + 1})`, () => {
                // test
                const formatted = TemperaturePipe.formatTemperature(test.c, test.unit, test.includeUnits)
                expect(formatted).toEqual(test.result);
            });
        })
    })
});
