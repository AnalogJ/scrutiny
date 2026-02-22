# Rootless Podman Quadlet Install

Note: These instructions are written with Podman 4.9 in mind, as that's what's available on Ubuntu 24.04. Podman 5+ can simplify the process using a .pod file to run both the hub and influxdb instance in the same pod, sharing localhost. This is a fairly trivial change should anyone want to add the documentation for it. While this document isn't Ubuntu-specific, this is being purposefully done to allow it to apply to the vast majority of Podman users, regardless of what Linux distro they use.


### Dependencies

- Podman > 4.9
- Systemd > 250 (for quadlet support)
- a restricted service account


### Creating a Service Account

See [Creating a Restricted Service Account](INSTALL_MANUAL.md#creating-a-restricted-service-account) for instructions.

While you do not need to use the same account as the collector, this guide will assume you will be for all its examples.

In addition to those steps, you will need to create sub ids and enable lingering for the user:

```sh
# add sub-uids and sub-gids, you may need to adjust numbers if you have other rootless quadlets running for other users already
# it is not recommended to go below 100000
# we choose to start at 500000 in the event you have some other podman accounts
sudo usermod --add-subuids 500000-565535 scrutiny-svc
sudo usermod --add-subgids 500000-565535 scrutiny-svc

# We want the quadlets to stay running even if the user isn't logged in
sudo loginctl enable-linger scrutiny-svc
```


### Directory Structure

Once the account is created, you will need to grab its id to create a few drectories for the data files and rootless quadlet files:

```sh
# create folders for config and influxdb
sudo mkdir -p /opt/scrutiny-svc/scrutiny/{config,influxdb}

# get the config file for scrutiny hub
sudo wget -O /opt/scrutiny-svc/scrutiny/config/scrutiny.yaml https://raw.githubusercontent.com/AnalogJ/scrutiny/refs/heads/master/example.scrutiny.yaml

# set permissions on everything
sudo chown -R scrutiny-svc:scrutiny-svc /opt/scrutiny-svc

# Get the ID of scrutiny-svc so you know it for your own record-keeping
id -u scrutiny-svc

# create a directory
sudo mkdir -p /etc/containers/systemd/users/$(id -u scrutiny-svc)

## go into the directory you just created for the rest of the guide
cd /etc/containers/systemd/users/$(id -u scrutiny-svc)
```


### Quadlet Files

Now that everything is set up and configured for the account to run quadlets, we just need to create a few quadlet files.

All remaining system actions will take place in `/etc/containers/systemd/users/$(id -u scrutiny-svc)` which is why we had you cd into it.


#### Networking

We need the hub and influxdb instances to be able to talk to each other, and in the case of Podman 4.9, they will run separately not sharing a localhost, and as such we need to configure a network for them to share. The file is pretty simple:


##### scrutiny-net.network

```ini
[Network]
NetworkName=scrutiny-net
```


#### Containers

Now we're ready for creating the containers


##### influxdb.container

```ini
[Unit]
Description=influxdb

[Container]
ContainerName=influxdb
Image=docker.io/library/influxdb:2.8
AutoUpdate=registry
Timezone=local
## not strictly necessary, but keeps file permission sane for influxdb
PodmanArgs=--group-add keep-groups
## versions of podman after 5.1 should do the below instead
#GroupAdd=keep-groups
Volume=/opt/scrutiny-svc/scrutiny/influxdb:/var/lib/influxdb2:Z
Network=scrutiny-net

[Service]
Restart=on-failure

[Install]
# Start by default on boot
WantedBy=default.target
```


##### scrutiny-web.container

```ini
[Unit]
Description=scrutiny-web
After=influxdb.service
Requires=influxdb.service

[Container]
ContainerName=scrutiny-web
Image=ghcr.io/analogj/scrutiny:latest-web
AutoUpdate=registry
Timezone=local
Volume=/opt/scrutiny-svc/scrutiny/config:/opt/scrutiny/config:Z
Network=scrutiny-net
PublishPort=8080:8080/tcp

[Service]
Restart=on-failure

[Install]
# Start by default on boot
WantedBy=default.target
```

#### Update scrutiny config

Since our containers are running separately, we need to update `/opt/scrutiny-svc/scrutiny/config/scrutiny.yaml` to the new influxdb host:

1. edit `/opt/scrutiny-svc/scrutiny/config/scrutiny.yaml`
2. under `influxdb` section, change `host: 0.0.0.0` to `host: influxdb` -- remember that yaml is whitespace-sensitive! so be mindful of the indents

```yaml
  influxdb:
#    scheme: 'http'
    host: influxdb
    port: 8086
```

# Running the hub and doing the 

With that done, we're now ready to start up the services:

```sh
# reload all the systemd user files for scrutiny-svc
sudo systemctl --user -M scrutiny-svc@ daemon-reload

# start the scrutiny-net network:
sudo systemctl --user -M scrutiny-svc@ start scrutiny-net-network.service

# start influxdb first and wait for it to come up
sudo systemctl --user -M scrutiny-svc@ start influxdb.service

# check if it's fully up
sudo systemctl --user -M scrutiny-svc@ status influxdb.service

# now start scrutiny
sudo systemctl --user -M scrutiny-svc@ start scrutiny-web.service
```

You are now ready to run the collector, if you would like to run that rootless as well, see the guide at [Schedule Collector with Systemd (rootless)](INSTALL_MANUAL.md#schedule-collector-with-systemd-rootless)