# Manual Install

While the easiest way to get started with [Scrutiny is using Docker](https://github.com/AnalogJ/scrutiny#docker),
it is possible to run it manually without much work. You can even mix and match, using Docker for one component and
a manual installation for the other. There's also [an installer](INSTALL_ANSIBLE.md) which automates this manual installation procedure.

Scrutiny is made up of three components: an influxdb Database, a collector and a webapp/api. Here's how each component can be deployed manually.

> Note: the `/opt/scrutiny` directory is not hardcoded, you can use any directory name/path.

## InfluxDB

Please follow the official InfluxDB installation guide. Note, you'll need to install v2.2.0+. 

https://docs.influxdata.com/influxdb/v2.2/install/

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
    host: 0.0.0.0
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

- **Ubuntu:** `apt-get install -y smartmontools=7.0-0ubuntu1~ubuntu18.04.1`
- **Centos8:**
    - `dnf install https://extras.getpagespeed.com/release-el8-latest.rpm`
    - `dnf install smartmontools`
- **FreeBSD:** `pkg install smartmontools`

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

### Start Scrutiny Collector, Populate Webapp

Next, we will manually trigger the collector, to populate the Scrutiny dashboard:

> NOTE: if you need to pass a config file to the scrutiny collector, you can provide it using the `--config` flag.

```
/opt/scrutiny/bin/scrutiny-collector-metrics-linux-amd64 run --api-endpoint "http://localhost:8080"
```

### Schedule Collector with Cron

Finally you need to schedule the collector to run periodically.
This may be different depending on your OS/environment, but it may look something like this:

```
# open crontab
crontab -e

# add a line for Scrutiny
*/15 * * * * . /etc/profile; /opt/scrutiny/bin/scrutiny-collector-metrics-linux-amd64 run --api-endpoint "http://localhost:8080"
```
