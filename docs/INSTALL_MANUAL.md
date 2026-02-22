# Manual Install

While the easiest way to get started with [Scrutiny is using Docker](https://github.com/AnalogJ/scrutiny#docker),
it is possible to run it manually without much work. You can even mix and match, using Docker for one component and
a manual installation for the other. There's also [an installer](INSTALL_ANSIBLE.md) which automates this manual installation procedure.

Scrutiny is made up of three components: an influxdb Database, a collector and a webapp/api. Here's how each component can be deployed manually.

> Note: the `/opt/scrutiny` directory is not hardcoded, you can use any directory name/path.

## InfluxDB

Please follow the official InfluxDB installation guide. Note, you'll need to install v2.8.0+. 

https://docs.influxdata.com/influxdb/v2/install/

## Webapp/API

### Dependencies

Since the webapp is packaged as a stand alone binary, there isn't really any software you need to install other than `glibc`
which is included by most linux OS's already.


### Directory Structure

Now let's create a directory structure to contain the Scrutiny files & binary.

```
mkdir -p /opt/scrutiny/config
mkdir -p /opt/scrutiny/web
mkdir -p /opt/scrutiny/bin
```

### Config file

While it is possible to run the webapp/api without a config file, the defaults are designed for use in a container environment,
and so will need to be overridden. So the first thing you'll need to do is create a config file that looks like the following:

```
# stored in /opt/scrutiny/config/scrutiny.yaml

version: 1

web:
  database:
    # The Scrutiny webapp will create a database for you, however the parent directory must exist.
    location: /opt/scrutiny/config/scrutiny.db
  src:
    frontend:
      # The path to the Scrutiny frontend files (js, css, images) must be specified.
      # We'll populate it with files in the next section
      path: /opt/scrutiny/web
  
  # if you're runnning influxdb on a different host (or using a cloud-provider) you'll need to update the host & port below. 
  # token, org, bucket are unnecessary for a new InfluxDB installation, as Scrutiny will automatically run the InfluxDB setup, 
  # and store the information in the config file. If you 're re-using an existing influxdb installation, you'll need to provide
  # the `token`
  influxdb:
    host: localhost
    port: 8086
#    token: 'my-token'
#    org: 'my-org'
#    bucket: 'bucket'
```

> Note: for a full list of available configuration options, please check the [example.scrutiny.yaml](https://github.com/AnalogJ/scrutiny/blob/master/example.scrutiny.yaml) file.

### Download Files

Next, we'll download the Scrutiny API binary and frontend files from the [latest Github release](https://github.com/analogj/scrutiny/releases).
The files you need to download are named:

- **scrutiny-web-linux-amd64** - save this file to `/opt/scrutiny/bin`
- **scrutiny-web-frontend.tar.gz** - save this file to `/opt/scrutiny/web`

### Prepare Scrutiny

Now that we have downloaded the required files, let's prepare the filesystem.

```
# Let's make sure the Scrutiny webapp is executable.
chmod +x /opt/scrutiny/bin/scrutiny-web-linux-amd64

# Next, lets extract the frontend files.
# NOTE: after extraction, there **should not** be a `dist` subdirectory in `/opt/scrutiny/web` directory.
cd /opt/scrutiny/web
tar xvzf scrutiny-web-frontend.tar.gz --strip-components 1 -C .


# Cleanup
rm -rf scrutiny-web-frontend.tar.gz
```

### Start Scrutiny Webapp

Finally, we start the Scrutiny webapp:

```
/opt/scrutiny/bin/scrutiny-web-linux-amd64 start --config /opt/scrutiny/config/scrutiny.yaml
```

The webapp listens for traffic on `http://0.0.0.0:8080` by default.


## Collector

### Dependencies

Unlike the webapp, the collector does have some dependencies:

- `smartctl`, v7+
- `cron` (or an alternative process scheduler)

Unfortunately the version of `smartmontools` (which contains `smartctl`) available in some of the base OS repositories is ancient.
So you'll need to install the v7+ version using one of the following commands:

- **Ubuntu (22.04/Jammy/LTS):** `apt-get install -y smartmontools`
- **Ubuntu (18.04/Bionic):** `apt-get install -y smartmontools=7.0-0ubuntu1~ubuntu18.04.1`
- **Centos8:**
    - `dnf install https://extras.getpagespeed.com/release-el8-latest.rpm`
    - `dnf install smartmontools`
- **FreeBSD:** `pkg install smartmontools`

The following additional dependencies are needed if you want to run the collector as an unprivileged user:

- systemd version > 235
- a restricted user account

### Directory Structure

Now let's create a directory structure to contain the Scrutiny collector binary.

```
mkdir -p /opt/scrutiny/bin
```


### Download Files

Next, we'll download the Scrutiny collector binary from the [latest Github release](https://github.com/analogj/scrutiny/releases). You are looking for the one titled **scrutiny-collector-metrics-linux-amd64** unless you know you are on arm.

```sh
wget -O /tmp/scrutiny-collector-metrics https://github.com/AnalogJ/scrutiny/releases/latest/download/scrutiny-collector-metrics-linux-amd64
```

Optional, but recommended: Before continuing it's recommended you compare the sha from the release page with the downloaded file to ensure it's the same file and not corrupted/tampered with. The command to do this is:

`echo "SHA_GOES_HERE /tmp/scrutiny-collector-metrics" | sha256sum -c`

example for the v0.8.6 release:

`echo "4c163645ce24e5487f4684a25ec73485d77a82a57f084808ff5aad0c11499ad2 /tmp/scrutiny-collector-metrics" | sha256sum -c`

followed by:

`sudo mv /tmp/scrutiny-collector-metrics /opt/scrutiny/bin/`

to move the binary to its final resting place


### Prepare Scrutiny

Now that we have downloaded the required files, let's prepare the filesystem.

```sh
# Let's make sure the Scrutiny collector is executable.
chmod +x /opt/scrutiny/bin/scrutiny-collector-metrics
```

if you are using SELinux, you may need to also do the following:

```sh
# tell SELinux to allow these binaries
sudo semanage fcontext -a -t bin_t "/opt/scrutiny/bin(/.*)?"
# update labels
sudo restorecon -Rv /opt/scrutiny/bin
```


### Start Scrutiny Collector, Populate Webapp

Next, we will manually trigger the collector, to populate the Scrutiny dashboard:

> NOTE: if you need to pass a config file to the scrutiny collector, you can provide it using the `--config` flag.

```sh
/opt/scrutiny/bin/scrutiny-collector-metrics run --api-endpoint "http://localhost:8080"
```

### Schedule Collector with (root) Cron

Finally you need to schedule the collector to run periodically.
This may be different depending on your OS/environment, but it may look something like this:

```sh
# open crontab
sudo crontab -e

# add a line for Scrutiny
*/15 * * * * . /etc/profile; /opt/scrutiny/bin/scrutiny-collector-metrics run --api-endpoint "http://localhost:8080"
```

### Schedule Collector with Systemd (rootless)

Alternatively you can run `scrutiny-collector-metrics` as non-root so long as the relevant capabilities and permissions are granted.


#### Creating a Restricted Service Account

This is the account that will run `scrutiny-collector-metrics`. Note this isn't strictly needed for all setups, but is useful from a logging/auditing perspective.

- Debian-based distros:
    - `sudo adduser --system scrutiny-svc --group --home /opt/scrutiny-svc`
- RHEL-based distros:
    - `sudo useradd --system --home-dir /opt/scrutiny-svc --shell /sbin/nologin scrutiny-svc`

Next, add the user to the `disk` group:

```sh
sudo usermod -aG disk scrutiny-svc
```


#### Creating a Restricted Systemd Service using AmbientCapabilities (easier)

This is the simpler setup, which allows you to run scrutiny rootless, but depending on what you want, may require granting more permissions to scrutiny than you would like to.

1. go to `/etc/systemd/system`
2. create scrutiny-collector.service with the following contents:


```ini
[Unit]
Description=Daily Restricted Scrutiny Collector
After=network.target

[Service]
[Unit]
Description=Daily Restricted Scrutiny Collector
After=network.target

[Service]
Type=oneshot
User=scrutiny-svc
Group=disk
ExecStart=/opt/scrutiny/bin/scrutiny-collector-metrics run --api-endpoint "http://localhost:8080"

# --- PRIVILEGE LOCKDOWN ---
## CAP_SYS_RAWIO is needed for SATA drives
AmbientCapabilities=CAP_SYS_RAWIO
CapabilityBoundingSet=CAP_SYS_RAWIO
## unfortunately nvme drives require CAP_SYS_ADMIN
## if you want nvme drives you must do the following:
#AmbientCapabilities=CAP_SYS_RAWIO CAP_SYS_ADMIN
#CapabilityBoundingSet=

NoNewPrivileges=yes

# Security/sandboxing settings
KeyringMode=private
LockPersonality=yes
MemoryDenyWriteExecute=yes
ProtectSystem=strict
ProtectHome=yes
PrivateDevices=no
## you can restrict devices using:
#DevicePolicy=closed
#DeviceAllow=/dev/sda r
#DeviceAllow=/dev/nvme0 r
ProtectKernelModules=yes
ProtectKernelTunables=yes
ProtectControlGroups=yes
ProtectClock=yes
ProtectHostname=yes
ProtectKernelLogs=yes
RemoveIPC=yes
RestrictSUIDSGID=true


# --- NETWORK LOCKDOWN
## use these to restrict what scrutiny can talk to over the network
## if using a hub on a different host you will need to change the values accordingly
RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX
IPAddressDeny=any
IPAddressAllow=localhost

[Install]
WantedBy=multi-user.target

```

Additionally, for nvme drives you may need to create a udev rule on many systems, as /dev/nvme* is often owned only by root:

##### add udev rule `/etc/udev/rules.d/99-nvme.rules` with contents:

```
KERNEL=="nvme[0-9]*", GROUP="disk", MODE="0640"
```

then run the following commands to load the udev rule:

```sh
sudo udevadm control --reload-rules
sudo udevadm trigger --subsystem-match=nvme --action=add
```


##### Pros:

- easy to maintain
- much better than running as root (especially if you don't need nvme drives)
- there are no privilege escalations needed


##### Cons:

NOTE: These cons basically only apply if a major supply-chain attack happens against scrutiny, and reflect a worst-case scenario that is unlikely to ever occur:

- CAP_SYS_RAWIO allows for data exfiltration/modification from SATA drives (ssh keys, /etc/shadow, etc)
- CAP_SYS_ADMIN would theoretically allow for significant system compromise
- nvme drives requires a udev rule for reliable access


If you are happy with that, you can jump to [Create a Systemd Timer to run scrutiny-collector.service](#create-a-systemd-timer-to-run-scrutiny-collectorservice)


#### Creating a Restricted Systemd Service using sudo and Shim Script

If granting scrutiny `CAP_SYS_RAWIO` and/or `CAP_SYS_ADMIN` exceeds your risk appetite, you have another option, though one more complicated and with its own set of pros/cons

1. run `sudo mkdir -p /opt/smartctl-shim/bin`
2. edit `/opt/smartctl-shim/bin/smartctl` with the following content:

```sh
#!/bin/bash
# Shim for accounts to use smartctl without being root
# for automation requires the account be in sudoers
exec /usr/bin/sudo /usr/sbin/smartctl "$@"
```

3. create a new `scrutiny-collector` file in `/etc/sudoers.d/`
4. inside `/etc/sudoers.d/scrutiny-collector` add the following:

```sh
scrutiny-svc ALL=(root) NOPASSWD: /usr/sbin/smartctl *
```

5. go to `/etc/systemd/system`
6. create scrutiny-collector.service with the following contents:


```ini
[Unit]
Description=Daily Restricted Scrutiny Collector
After=network.target

[Service]
Type=oneshot
User=scrutiny-svc
Environment="PATH=/opt/smartctl-shim/bin:/usr/bin:/bin"
ExecStart=/opt/scrutiny/bin/scrutiny-collector-metrics run --api-endpoint "http://localhost:8080"

# --- PRIVILEGE LOCKDOWN ---
## we use sudo to elevate privileges for smartctl only, so no Ambient Capabilities are needed
AmbientCapabilities=
## CAP_SYS_RAWIO is needed for SATA drives
CapabilityBoundingSet=CAP_SETUID CAP_SETGID CAP_AUDIT_WRITE CAP_SYS_RAWIO CAP_SYS_RESOURCE
## unfortunately nvme drives require CAP_SYS_ADMIN
## if you want nvme drives you must do the following:
# CapabilityBoundingSet=CAP_SETUID CAP_SETGID CAP_AUDIT_WRITE CAP_SYS_RAWIO CAP_SYS_ADMIN CAP_SYS_RESOURCE

## since sudo needs to be used to elevate permissions in this setup, we need to allow new privileges
NoNewPrivileges=no

# Security/sandboxing settings
KeyringMode=private
LockPersonality=yes
MemoryDenyWriteExecute=yes
ProtectSystem=strict
ProtectHome=yes
PrivateDevices=no
ProtectKernelModules=yes
ProtectKernelTunables=yes
ProtectControlGroups=yes
ProtectClock=yes
ProtectHostname=yes
ProtectKernelLogs=yes
RemoveIPC=yes
RestrictSUIDSGID=true


# --- NETWORK LOCKDOWN
## use these to restrict what scrutiny can talk to over the network
## if using a hub on a different host you will need to change the values accordingly
RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX
IPAddressDeny=any
IPAddressAllow=localhost

[Install]
WantedBy=multi-user.target
```


##### Pros:

- the scrutiny binary itself will not have permissions like CAP_SYS_ADMIN
- much better than running as root (especially if you don't need nvme drives)
- `sudo` restricts privilege escalation to just `smartctl`
- no udev rule needed


##### Cons:

NOTE: These cons basically only apply if a major supply-chain attack happens against scrutiny, and reflect a worst-case scenario that is unlikely to ever occur:

- Any sort of privilege escalation attack in sudo could theoretically allow a compromised scrutiny to gain additional privileges, since the process has permission to escelate privileges in general
- Even though sudo only allows `smartctl`, it still has `CAP_SYS_RAWIO` and `CAP_SYS_ADMIN` so in theory the same attacks from the first method are possible, though now only with an exploit using smartctl instead of scrutiny directly
- even though you don't need a udev rule, this adds a lot of additional administrative overhead
- while the scrutiny binary itself isn't elevated, it has a sub-process that is running as root (systemctl)

#### Create a Systemd Timer to run scrutiny-collector.service

First, lets test our service. It doesn't matter which method you used above, as either way you need to load and run it.

```sh
# reload changes for systemd services
sudo systemctl daemon-reload

# enable the service
sudo systemctl enable scrutiny-collector.service

# now run the service
sudo systemctl start scrutiny-collector.service
```

You should see the data in your hub instance of scrutiny now. If your run into issues I recommend turning on debug logging for scrutiny and checking your system logs using journalctl. It may be a permission is missing or wrong.

Now that things have been validated, lets create the systemd timer to run the service for us on a schedule:

1. if you are not still there, go to `/etc/systemd/system`
2. create scrutiny-collector.timer with the following contents:

```ini
[Unit]
Description=Run Scruitiny Collector daily at 2am

[Timer]
# Standard calendar trigger
OnCalendar=*-*-* 02:00:00
# Ensures the job runs if the computer was off at 2am
Persistent=true
# Minimizes I/O spikes by staggering start time
RandomizedDelaySec=30

[Install]
WantedBy=timers.target

```

Update the schedule as you see fit for your needs

Once you are satisfied with our timer, you'll need to load and enable it:

```sh
# reload changes for systemd services
sudo systemctl daemon-reload

# now enable the timer
sudo systemctl enable --now scrutiny-collector.timer
```

That's it! you're done. You can check the status of the timer using `sudo systemctl status scrutiny-collector.timer
`