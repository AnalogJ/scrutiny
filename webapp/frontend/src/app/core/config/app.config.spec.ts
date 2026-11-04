import {normalizeTemperatureDuration} from './app.config';

describe('normalizeTemperatureDuration', () => {
    it('keeps known durations', () => {
        expect(normalizeTemperatureDuration('forever')).toBe('forever');
    });

    it('uses week for unknown or missing durations', () => {
        expect(normalizeTemperatureDuration('bogus')).toBe('week');
        expect(normalizeTemperatureDuration('')).toBe('week');
        expect(normalizeTemperatureDuration(undefined)).toBe('week');
    });
});
