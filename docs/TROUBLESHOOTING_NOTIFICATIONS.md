# Notifications

As documented in [example.scrutiny.yaml](https://github.com/AnalogJ/scrutiny/blob/master/example.scrutiny.yaml#L59-L75)
there are multiple ways to configure notifications for Scrutiny.

Under the hood we use a library called [Shoutrrr](https://github.com/containrrr/shoutrrr) to send our notifications, and you should use their documentation if you run into
any issues: https://containrrr.dev/shoutrrr/services/overview/


# Script Notifications

While the Shoutrrr library supports many popular providers for sending notifications Scrutiny also supports a "script" based
notification system, allowing you to execute a custom script whenever a notification needs to be sent. 
Data is provided to this script using the following environmental variables:

```
SCRUTINY_SUBJECT - 	eg. "Scrutiny SMART error (%s) detected on device: %s"
SCRUTINY_DATE 
SCRUTINY_FAILURE_TYPE - EmailTest, SmartFail, ScrutinyFail
SCRUTINY_DEVICE_NAME - eg. /dev/sda
SCRUTINY_DEVICE_TYPE - ATA/SCSI/NVMe
SCRUTINY_DEVICE_SERIAL - eg. WDDJ324KSO
SCRUTINY_MESSAGE - eg. "Scrutiny SMART error notification for device: %s\nFailure Type: %s\nDevice Name: %s\nDevice Serial: %s\nDevice Type: %s\nDate: %s"
SCRUTINY_HOST_ID - (optional) eg. "my-custom-host-id"
```

# Special Characters

`Shoutrrr` supports special characters in the username and password fields, however you'll need to url-encode the
username and the password separately.

- if your username is: `myname@example.com`
- if your password is `124@34$1`

Then your `shoutrrr` url will look something like:

- `smtp://myname%40example%2Ecom:124%4034%241@ms.my.domain.com:587`

# Testing Notifications

You can test that your notifications are configured correctly by posting an empty payload to the notifications health
check API.

```
curl -X POST http://localhost:8080/api/health/notify
```
