import {DeviceStatusPipe, deviceStatusHasOverride} from './device-status.pipe';
import {MetricsStatusThreshold} from '../core/config/app.config';
import {DeviceModel} from '../core/models/device-model';

describe('DeviceStatusPipe', () => {
    it('create an instance', () => {
        const pipe = new DeviceStatusPipe();
        expect(pipe).toBeTruthy();
    });

    describe('#deviceStatusForModelWithThreshold', () => {
        it('if healthy device, should be passing', () => {
            expect(DeviceStatusPipe.deviceStatusForModelWithThreshold(
                {device_status: 0} as DeviceModel,
                true,
                MetricsStatusThreshold.Both
            )).toBe('passed')
        });

        it('if device with no smart data, should be unknown', () => {
            expect(DeviceStatusPipe.deviceStatusForModelWithThreshold(
                {device_status: 0} as DeviceModel,
                false,
                MetricsStatusThreshold.Both
            )).toBe('unknown')
        });

        const testCases = [
            // Override bits (4 = failed-by-override, 8 = passed-by-override) explain *why* a
            // verdict was reached and are not failure bits. They must be masked out before the
            // status lookup, otherwise the value falls outside the table and renders undefined.
            {
                'deviceStatus': 8, // passing, decided by an override
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Both,
                'includeReason': false,
                'result': 'passed'
            },
            {
                'deviceStatus': 8,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Both,
                'includeReason': true,
                'result': 'passed'
            },
            {
                'deviceStatus': 6, // 2 (scrutiny) | 4 (failed by override)
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Both,
                'includeReason': false,
                'result': 'failed'
            },
            {
                'deviceStatus': 6,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Both,
                'includeReason': true,
                'result': 'failed: scrutiny'
            },
            {
                'deviceStatus': 6,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Smart, // scrutiny failure filtered out
                'includeReason': false,
                'result': 'passed'
            },
            {
                'deviceStatus': 5, // 1 (smart) | 4 (failed by override)
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Both,
                'includeReason': true,
                'result': 'failed: smart'
            },
            {
                'deviceStatus': 10000, // invalid status
                'hasSmartResults': false,
                'threshold': MetricsStatusThreshold.Smart,
                'includeReason': false,
                'result': 'unknown'
            },

            {
                'deviceStatus': 1,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Smart,
                'includeReason': false,
                'result': 'failed'
            },
            {
                'deviceStatus': 1,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Scrutiny,
                'includeReason': false,
                'result': 'passed'
            },
            {
                'deviceStatus': 1,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Both,
                'includeReason': false,
                'result': 'failed'
            },

            {
                'deviceStatus': 2,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Smart,
                'includeReason': false,
                'result': 'passed'
            },
            {
                'deviceStatus': 2,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Scrutiny,
                'includeReason': false,
                'result': 'failed'
            },
            {
                'deviceStatus': 2,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Both,
                'includeReason': false,
                'result': 'failed'
            },

            {
                'deviceStatus': 3,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Smart,
                'includeReason': false,
                'result': 'failed'
            },
            {
                'deviceStatus': 3,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Scrutiny,
                'includeReason': false,
                'result': 'failed'
            },
            {
                'deviceStatus': 3,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Both,
                'includeReason': false,
                'result': 'failed'
            },

            {
                'deviceStatus': 3,
                'hasSmartResults': false,
                'threshold': MetricsStatusThreshold.Smart,
                'includeReason': true,
                'result': 'unknown'
            },
            {
                'deviceStatus': 3,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Smart,
                'includeReason': true,
                'result': 'failed: smart'
            },
            {
                'deviceStatus': 3,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Scrutiny,
                'includeReason': true,
                'result': 'failed: scrutiny'
            },
            {
                'deviceStatus': 3,
                'hasSmartResults': true,
                'threshold': MetricsStatusThreshold.Both,
                'includeReason': true,
                'result': 'failed: both'
            }


        ]

        testCases.forEach((test, index) => {
            it(`if device with status (${test.deviceStatus}), hasSmartResults(${test.hasSmartResults}) and threshold (${test.threshold}), should be ${test.result}`, () => {
                expect(DeviceStatusPipe.deviceStatusForModelWithThreshold(
                    {device_status: test.deviceStatus} as DeviceModel,
                    test.hasSmartResults,
                    test.threshold,
                    test.includeReason
                )).toBe(test.result)
            });
        });
    });
});

describe('deviceStatusHasOverride', () => {
    it('is false when no override bits are set', () => {
        expect(deviceStatusHasOverride(0)).toBe(false);
        expect(deviceStatusHasOverride(1)).toBe(false);
        expect(deviceStatusHasOverride(3)).toBe(false);
    });

    it('is true when an override decided the verdict', () => {
        expect(deviceStatusHasOverride(4)).toBe(true);  // failed by override
        expect(deviceStatusHasOverride(8)).toBe(true);  // passed by override
        expect(deviceStatusHasOverride(6)).toBe(true);  // scrutiny failure, via override
    });
});
