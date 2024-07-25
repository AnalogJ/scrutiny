import { DeviceHoursPipe } from "./device-hours.pipe";

describe("DeviceHoursPipe", () => {
  it("create an instance", () => {
    const pipe = new DeviceHoursPipe();
    expect(pipe).toBeTruthy();
  });

  describe("#transform", () => {
    const testCases = [
      {
        input: 12345,
        configuration: "device_hours",
        result: "12345 hours",
      },
      {
        input: 15273,
        configuration: "humanize",
        result: "1 year, 8 months, 3 weeks, 6 days, 15 hours",
      },
      {
        input: 48,
        configuration: null,
        result: "2 days",
      },
      {
        input: 168,
        configuration: "scrutiny",
        result: "1 week",
      },
      {
        input: null,
        configuration: "device_hours",
        result: "Unknown",
      },
      {
        input: null,
        configuration: "humanize",
        result: "Unknown",
      },
    ];

    testCases.forEach((test, index) => {
      it(`format input '${test.input}' with configuration '${test.configuration}', should be '${test.result}' (testcase: ${index + 1})`, () => {
        // test
        const pipe = new DeviceHoursPipe();
        const formatted = pipe.transform(test.input, test.configuration);
        expect(formatted).toEqual(test.result);
      });
    });
  });
});
