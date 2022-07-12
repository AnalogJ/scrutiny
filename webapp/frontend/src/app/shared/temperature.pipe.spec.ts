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
                'f': -460,
            },{
                'c': -34.44,
                'f': -30,
            },{
                'c': -23.33,
                'f': -10,
            },{
                'c': -17.78,
                'f': -0,
            },{
                'c': 0,
                'f': 32,
            },{
                'c': 10,
                'f': 50,
            },{
                'c': 26.67,
                'f': 80,
            },{
                'c': 37,
                'f': 99,
            },{
                'c': 60,
                'f': 140,
            }
        ]
        testCases.forEach((test, index) => {
            it(`should correctly convert ${test.c}, Celsius to Fahrenheit (testcase: ${index + 1})`, () => {
                // test
                const numb = TemperaturePipe.celsiusToFahrenheit(test.c)
                const roundNumb = Math.round(numb);
                expect(roundNumb).toEqual(test.f);
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
                'c': 26.67,
                'unit': 'celsius',
                'includeUnits': false,
                'result': '26.67',
            },{
                'c': 26.67,
                'unit': 'fahrenheit',
                'includeUnits': true,
                'result': '80.006°F',
            },{
                'c': 26.67,
                'unit': 'fahrenheit',
                'includeUnits': false,
                'result': '80.006',
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
