import {FileSizePipe} from './file-size.pipe';

describe('FileSizePipe', () => {
    it('create an instance', () => {
        const pipe = new FileSizePipe();
        expect(pipe).toBeTruthy();
    });

    describe('#transform',() => {
        const testCases = [
            {
                'bytes': 1500,
                'si': false,
                'result': '1.5 KiB'
            },
            {
                'bytes': 1500,
                'si': true,
                'result': '1.5 kB'
            },
            {
                'bytes': 5000,
                'si': false,
                'result': '4.9 KiB',
            },
            {
                'bytes': 5000,
                'si': true,
                'result': '5.0 kB',
            },
            {
                'bytes': 999_949,
                'si': false,
                'result': '976.5 KiB',
            },
            {
                'bytes': 999_949,
                'si': true,
                'result': '999.9 kB',
            },
            {
                'bytes': 999_950,
                'si': true,
                'result': '1.0 MB',
            },
            {
                'bytes': 1_551_859_712,
                'si': false,
                'result': '1.4 GiB',
            },
            {
                'bytes': 2_100_000_000,
                'si': false,
                'result': '2.0 GiB',
            },
            {
                'bytes': 2_100_000_000,
                'si': true,
                'result': '2.1 GB',
            }
        ]

        testCases.forEach((test, index) => {
            it(`should correctly format bytes ${test.bytes}. (testcase: ${index + 1})`, () => {
                // test
                const pipe = new FileSizePipe();
                const formatted = pipe.transform(test.bytes, test.si)
                expect(formatted).toEqual(test.result);
            });
        })
    })

});
