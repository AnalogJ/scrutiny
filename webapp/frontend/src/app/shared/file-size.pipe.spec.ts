import { FileSizePipe } from './file-size.pipe';

describe('FileSizePipe', () => {
    it('create an instance', () => {
        const pipe = new FileSizePipe();
        expect(pipe).toBeTruthy();
    });

    describe('#transform',() => {
        const testCases = [
            {
                'bytes': 1500,
                'precision': undefined,
                'result': '1 KB'
            },{
                'bytes': 2_100_000_000,
                'precision': undefined,
                'result': '2.0 GB',
            },{
                'bytes': 1500,
                'precision': 2,
                'result': '1.46 KB',
            }
        ]
        testCases.forEach((test, index) => {
            it(`should correctly format bytes ${test.bytes}. (testcase: ${index + 1})`, () => {
                // test
                const pipe = new FileSizePipe();
                const formatted = pipe.transform(test.bytes, test.precision)
                expect(formatted).toEqual(test.result);
            });
        })
    })

});
