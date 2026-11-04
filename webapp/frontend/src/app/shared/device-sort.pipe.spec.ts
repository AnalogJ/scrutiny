import { DeviceSortPipe } from './device-sort.pipe';

describe('DeviceSortPipe', () => {
  it('create an instance', () => {
    const pipe = new DeviceSortPipe();
    expect(pipe).toBeTruthy();
  });

  it('sorts a device passing by override with the passing devices', () => {
    const summary = (device_status: number) => ({smart: {}, device: {device_status}});
    const failedByOverride = summary(6);
    const passed = summary(0);
    const passedByOverride = summary(8);

    const sorted = new DeviceSortPipe().transform([passedByOverride, passed, failedByOverride]);

    expect(sorted).toEqual([failedByOverride, passedByOverride, passed]);
  });
});
