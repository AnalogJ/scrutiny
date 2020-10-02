# Manual Install

While the easiest way to get started with [Scrutiny is using Docker](https://github.com/AnalogJ/scrutiny#docker),
it is possible to run it manually without much work. You can even mix and match, using Docker for one component and
a manual installation for the other.

Scrutiny is made up of two components: a collector and a webapp/api. Here's how each component can be deployed manually.

> Note: the `/opt/scrutiny` directory is not hardcoded, you can use any directory name/path.

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
cd /opt/scrutiny/web
tar xvzf scrutiny-web-frontend.tar.gz --strip-components 1 -C .

# Cleanup
rm -rf scrutiny-web-frontend.tar.gz
```

### Install Scrutiny systemd service

You may install Scrutiny as a systemd service:

Create user and group for the service:

```
groupadd -r scrutiny
useradd -m -d /opt/scrutiny -s /sbin/nologin -r -g scrutiny scrutiny
```

Change file permissions:

```
chown -R scrutiny\:scrutiny /opt/scrutiny
```

Create the service unit file:

```
cat > /etc/systemd/system/scrutiny.service <<EOF
[Unit]
Description=Scrutiny disk health monitor
After=network.target

[Service]
Type=idle
User=scrutiny
Group=scrutiny
ExecStart=/opt/scrutiny/bin/scrutiny-web-linux-amd64 start --config /opt/scrutiny/config/scrutiny.yaml
TimeoutStartSec=600
TimeoutStopSec=600

[Install]
WantedBy=multi-user.target
EOF
```

Enable and start the service:

```
systemctl daemon-reload
systemctl enable scrutiny
systemctl start scrutiny
```

Check log for any issues:

```
journalctl -u scrutiny
```

### Start Scrutiny Webapp (non-systemd)

If you are not using the systemd service, you can start the Scrutiny webapp in the
foreground:

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

- **Ubuntu:** `apt-get install -y smartmontools=7.0-0ubuntu1~ubuntu18.04.1`
- **Centos8:**
    - `dnf install https://extras.getpagespeed.com/release-el8-latest.rpm`
    - `dnf install smartmontools`


### Directory Structure

Now let's create a directory structure to contain the Scrutiny collector binary.

```
mkdir -p /opt/scrutiny/bin
```


### Download Files

Next, we'll download the Scrutiny collector binary from the [latest Github release](https://github.com/analogj/scrutiny/releases).
The file you need to download is named:

- **scrutiny-collector-metrics-linux-amd64** - save this file to `/opt/scrutiny/bin`


### Prepare Scrutiny

Now that we have downloaded the required files, let's prepare the filesystem.

```
# Let's make sure the Scrutiny collector is executable.
chmod +x /opt/scrutiny/bin/scrutiny-collector-metrics-linux-amd64
```

### Install Scrutiny Collector systemd service

If you are using systemd, you may configure collector as a systemd service and
create a timer to invoke it periodically:

```
sudo cat > /etc/systemd/system/scrutiny-collector.service <<EOF
[Unit]
Description=Scrutiny disk health data collector

[Service]
Type=idle
ExecStart=/opt/scrutiny/bin/scrutiny-collector-metrics-linux-amd64 run --api-endpoint "http://localhost:8080"
TimeoutStartSec=600
TimeoutStopSec=600
EOF
```

Configure collector systemd timer:

```
sudo cat > /etc/systemd/system/scrutiny-collector.timer <<EOF
[Unit]
Description=Scrutiny disk health data collector timer

[Timer]
OnCalendar=*:0/15
Unit=scrutiny-collector.service
Persistent=true

[Install]
WantedBy=timers.target
EOF
```

Enable timer:

```
systemctl daemon-reload
systemctl enable scrutiny-collector.timer
systemctl start scrutiny-collector.timer
```

Check timer:

```
systemctl is-active scrutiny-collector.timer
```

Should say `active`

Check the next scheduled run:

```
systemctl list-timers scrutiny*
```

Manually run the first data collection:

```
systemctl start scrutiny-collector.service
```

Check log for any issues:

```
journalctl -u scrutiny-collector.service
```

### Non-systemd configuration

#### Start Scrutiny Collector, Populate Webapp

Next, we will manually trigger the collector, to populate the Scrutiny dashboard:

```
/opt/scrutiny/bin/scrutiny-collector-metrics-linux-amd64 run --api-endpoint "http://localhost:8080"
```

#### Schedule Collector with Cron

Finally you need to schedule the collector to run periodically.
This may be different depending on your OS/environment, but it may look something like this:

```
# open crontab
crontab -e

# add a line for Scrutiny
*/15 * * * * /opt/scrutiny/bin/scrutiny-collector-metrics-linux-amd64 run --api-endpoint "http://localhost:8080"
```
