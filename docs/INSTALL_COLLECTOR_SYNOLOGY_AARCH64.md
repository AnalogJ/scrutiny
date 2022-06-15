# Install collector on Synology

### Install Entware

This will allow you to install a newer version of smartmontools on your Synology. Follow the instructions here (This is tested on DSM7) - https://github.com/Entware/Entware/wiki/Install-on-Synology-NAS

**PLEASE NOTE THAT IF YOU UPDATE DSM FIRMWARE YOU MAY BORK THE EXISTING ENTWARE INSTALLATION, FOR ANYTHING THAT MAY RELATE TO ENTWARE PLEASE VISIT THEIR REPO**

## Collector Setup

**1. Run an update**

`sudo opkg update`

**2. Run an upgrade**

`sudo opkg upgrade`

**3. Install smartmontools**

`sudo opkg install smartmontools`

*It should install v7.2-2*

`Installing smartmontools (7.2-2) to root...`

**4. We will now create the directories.**

```
mkdir -p /volume1/\@Entware/scrutiny/bin
mkdir -p /volume1/\@Entware/scrutiny/conf
```

**5. change into the bin directory**

`cd /volume1/\@Entware/scrutiny/bin`

**6. Download the collector binary for your architecture and make it executable**

`wget https://github.com/AnalogJ/scrutiny/releases/download/v0.4.12/scrutiny-collector-metrics-linux-arm64`

`chmod +x /volume1/\@Entware/scrutiny/bin/scrutiny-collector-metrics-linux-arm64`

**7. Create a config file for the collector**

```
cd /volume1/\@Entware/scrutiny/conf
wget https://raw.githubusercontent.com/AnalogJ/scrutiny/master/example.collector.yaml
mv example.collector.yaml collector.yaml
```

**8. Lets make some changes in the config file, these are what i uncommented/added, please tweak the device paths to your needs**

```
host:
  id: 'Server_Name'


devices:
#  # example for forcing device type detection for a single disk
 - device: /dev/sda
   type: 'sat'
 - device: /dev/sdb
   type: 'sat'
 - device: /dev/sdc
   type: 'sat'
 - device: /dev/sdd
   type: 'sat'
    
api:
 endpoint: 'http://<url>:8080'
```

**9. Let's update the smartd db**

```
cd /volume1/\@Entware/scrutiny/bin/
wget https://raw.githubusercontent.com/smartmontools/smartmontools/master/smartmontools/drivedb.h
```

**10. I ran it like this but you can tweak to your liking, the most important part is the --drivedb, as this loads it into the aplication for future use**

`smartctl -d sat --all /dev/sda  --drivedb=/volume1/\@Entware/scrutiny/bin/drivedb.h`

**11. Now lets create a small bash script, this will be used for the scheduled task inside Synology**

`vim /volume1/\@Entware/scrutiny/bin/run_collect.sh`

**The contents are below, copy and paste them in**

```
#!/bin/bash

/volume1/\@Entware/scrutiny/bin/scrutiny-collector-metrics-linux-arm64 run --config /volume1/\@Entware/scrutiny/config/collector.yaml
```

## Set up Synology to run a scheduled task. 

Log in to DSM and do the following:

Goto: DSM > Control Panel > Task Scheduler

Create > Scheduled Task > User Defined Script

###### General

```
Task: Scrutiny_Collector
User: root
Enabled: yes
```

###### Schedule
```
Run on the following days: Daily
```
###### Time:

```
Frequency: <Your desired frequency>
```

###### Task Settings

**Run Command**

```
. /opt/etc/profile; /volume1/\@Entware/scrutiny/bin/run_collect.sh
```
