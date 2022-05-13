# Scrutiny <-> SmartMonTools 

Scrutiny uses `smartctl --scan` to detect devices/drives. If your devices are not being detected by Scrutiny, or some 
data is missing, this is probably due to a `smartctl` issue. 
The following page will document commonly asked questions and troubleshooting steps for the Scrutiny S.M.A.R.T. data collector.

## WWN vs Device name
As discussed in [`#117`](https://github.com/AnalogJ/scrutiny/issues/117), `/dev/sd*` device paths are ephemeral. 

> Device paths in Linux aren't guaranteed to be consistent across restarts. Device names consist of major numbers (letters) and minor numbers. When the Linux storage device driver detects a new device, the driver assigns major and minor numbers from the available range to the device. When a device is removed, the device numbers are freed for reuse.
>
> The problem occurs because device scanning in Linux is scheduled by the SCSI subsystem to happen asynchronously. As a result, a device path name can vary across restarts.
>
> https://docs.microsoft.com/en-us/troubleshoot/azure/virtual-machines/troubleshoot-device-names-problems

While the Docker Scrutiny collector does require devices to attached to the docker container by device name (using `--device=/dev/sd..`), internally 
Scrutiny stores and references the devices by their `WWN` which is globally unique, and never changes. 

As such, passing devices to the Scrutiny collector container using `/dev/disk/by-id/`, `/dev/disk/by-label/`, `/dev/disk/by-path/` and `/dev/disk/by-uuid/`
paths are unnecessary, unless you'd like to ensure the docker run command never needs to change.


## Device Detection By Smartctl

The first thing you'll want to do is run `smartctl` locally (not in Docker) and make sure the output shows all your drives as expected.
See the `Drive Types` section below for what this output should look like for `NVMe`/`ATA`/`RAID` drives.

```bash
smartctl --scan

/dev/sda -d scsi # /dev/sda, SCSI device
/dev/sdb -d scsi # /dev/sdb, SCSI device
/dev/sdc -d scsi # /dev/sdc, SCSI device
/dev/sdd -d scsi # /dev/sdd, SCSI device
```

Once you've verified that `smartctl` correctly detects your drives, make sure scrutiny is correctly detecting them as well.
> NOTE: make sure you specify all the devices you'd like scrutiny to process using `--device=` flags.

```bash
docker run -it --rm \
  -v /run/udev:/run/udev:ro \
  --cap-add SYS_RAWIO \
  --device=/dev/sda \
  --device=/dev/sdb \
  ghcr.io/analogj/scrutiny:master-collector smartctl --scan
```

If the output is the same, your devices will be processed by Scrutiny.

### Collector Config File
In some cases `--scan` does not correctly detect the device type, returning [incomplete SMART data](https://github.com/AnalogJ/scrutiny/issues/45).
Scrutiny will supports overriding the detected device type via the config file.

### RAID Controllers (Megaraid/3ware/HBA/Adaptec/HPE/etc)
Smartctl has support for a large number of [RAID controllers](https://www.smartmontools.org/wiki/Supported_RAID-Controllers), however this 
support is not automatic, and may require some additional device type hinting. You can provide this information to the Scrutiny collector
using a collector config file. See [example.collector.yaml](/example.collector.yaml)

> NOTE: If you use docker, you **must** pass though the RAID virtual disk to the container using `--device` (see below)
>
> This device may be in `/dev/*` or `/dev/bus/*`. 
>
> If you're unsure, run `smartctl --scan` on your host, and pass all listed devices to the container.

```yaml
# /opt/scrutiny/config/collector.yaml
devices:
  # Dell PERC/Broadcom Megaraid example: https://github.com/AnalogJ/scrutiny/issues/30
  - device: /dev/bus/0
    type:
      - megaraid,14
      - megaraid,15
      - megaraid,18
      - megaraid,19
      - megaraid,20
      - megaraid,21

  - device: /dev/twa0
    type:
      - 3ware,0
      - 3ware,1
      - 3ware,2
      - 3ware,3
      - 3ware,4
      - 3ware,5
  
  # Adapec RAID: https://github.com/AnalogJ/scrutiny/issues/189
  - device: /dev/sdb
    type:
      - aacraid,0,0,0
      - aacraid,0,0,1
  
  # HPE Smart Array example:  https://github.com/AnalogJ/scrutiny/issues/213
  - device: /dev/sda
    type:
      - 'cciss,0'
      - 'cciss,1'
```

### NVMe Drives
As mentioned in the [README.md](/README.md), NVMe devices require both `--cap-add SYS_RAWIO` and `--cap-add SYS_ADMIN` 
to allow smartctl permission to query your NVMe device SMART data [#26](https://github.com/AnalogJ/scrutiny/issues/26)

When attaching NVMe devices using `--device=/dev/nvme..`, make sure to provide the device controller (`/dev/nvme0`) 
instead of the block device (`/dev/nvme0n1`). See [#209](https://github.com/AnalogJ/scrutiny/issues/209).

> The character device /dev/nvme0 is the NVME device controller, and block devices like /dev/nvme0n1 are the NVME storage namespaces: the devices you use for actual storage, which will behave essentially as disks.
>
> In enterprise-grade hardware, there might be support for several namespaces, thin provisioning within namespaces and other features. For now, you could think namespaces as sort of meta-partitions with extra features for enterprise use.

### ATA

### Standby/Sleeping Disks
- https://github.com/AnalogJ/scrutiny/issues/221
- https://github.com/AnalogJ/scrutiny/issues/157

### Volume Mount All Devices (`/dev`) - Privileged


## Scrutiny detects Failure but SMART Passed?

There's 2 different mechanisms that Scrutiny uses to detect failures.

The first is simple SMART failures. If SMART thinks an attribute is in a failed state, Scrutiny will display it as failed as well.

The second is using BackBlaze failure data: [https://backblaze.com/blog-smart-stats-2014-8.html](https://backblaze.com/blog-smart-stats-2014-8.html) 
If Scrutiny detects that an attribute corresponds with a high rate of failure using BackBlaze's data, it will also mark that attribute (and disk) as failed (even though SMART may think the device is still healthy).

This can cause some confusion when comparing Scrutiny's dashboard against other SMART analysis tools. 
If you hover over the "failed" label beside an attribute, Scrutiny will tell you if the failure was due to SMART or Scrutiny/BackBlaze data. 


## Hub & Spoke model, with multiple Hosts.

When deploying Scrutiny in a hub & spoke model, it can be difficult to determine exactly which node a set of devices are associated with.
Thankfully the collector has a special `--host-id` flag (or `COLLECTOR_HOST_ID` env variable) that can be used to associate devices with a friendly host name.

See the [docs/INSTALL_HUB_SPOKE.md](/docs/INSTALL_HUB_SPOKE.md) guide for more information. 

