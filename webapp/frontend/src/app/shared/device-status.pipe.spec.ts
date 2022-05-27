import { DeviceStatusPipe } from './device-status.pipe';

describe('DeviceStatusPipe', () => {
  it('create an instance', () => {
    const pipe = new DeviceStatusPipe();
    expect(pipe).toBeTruthy();
  });
});
