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

You have two options for running the collector rootless. An easy setup as well as a more advanced one that thoroughly locks everything down. In both approaches you will need to create a group and udev rule:


#### Setting Up Permissions

This is the group that will be used when running `scrutiny-collector-metrics`. Note this isn't strictly needed for all setups, but is useful from a logging/auditing perspective. Also if you are running your web and influxdb instances via rootless podman, you can skip this as the group was already created when with the podman user.

- Debian-based distros:
    - `sudo addgroup --system scrutiny-svc`
- RHEL-based distros:
    - `sudo groupadd --system scrutiny-svc`

Next, for nvme drives you may need to create a udev rule on many systems, as /dev/nvme* are often owned only by root:

##### add udev rule `/etc/udev/rules.d/99-nvme.rules` with contents:

```
KERNEL=="nvme[0-9]*", GROUP="disk", MODE="0640"
```

then run the following commands to load the udev rule:

```sh
sudo udevadm control --reload-rules
sudo udevadm trigger --subsystem-match=nvme --action=add
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
DynamicUser=yes
SupplementaryGroups=scrutiny-svc disk
ExecStart=/opt/scrutiny/bin/scrutiny-collector-metrics run --api-endpoint "http://localhost:8080"

# --- PRIVILEGE LOCKDOWN ---
## CAP_SYS_RAWIO is needed for SATA drives
AmbientCapabilities=CAP_SYS_RAWIO
CapabilityBoundingSet=CAP_SYS_RAWIO
## unfortunately nvme drives require CAP_SYS_ADMIN
## if you want nvme drives you must do the following:
#AmbientCapabilities=CAP_SYS_RAWIO CAP_SYS_ADMIN
#CapabilityBoundingSet=CAP_SYS_RAWIO CAP_SYS_ADMIN

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


#### Creating a Restricted Systemd Service using smartctl Proxy Via Unix Socket (advanced)

If granting scrutiny `CAP_SYS_RAWIO` and/or `CAP_SYS_ADMIN` exceeds your risk appetite, you have another option, though one more complicated. The payoff, though, is a hardened setup where scrutiny's collector is completely unprivileged.

1. install `socat` via your package manager. This is needed by the smartctl-shim
    * Debian-based distros: `sudo apt install socat`
    * RHEL-based distros: `sudo dnf install socat`
2. run `sudo mkdir -p /opt/smartctl-shim/bin`
3. edit `/opt/smartctl-shim/bin/smartctl` with the following content:

```bash
#!/usr/bin/env bash

# Shim for accounts to use smartctl without being root via unix socket

# don't send json, so it won't be processed by scrutiny
# but make the error informative for logs
send_smartctl_shim_error() {
    local MSG="$1" ARGS="$2" EXIT_CODE="${3:-3}"

    cat <<EOF
smartctl-shim 0.3
=======> ERROR <=======
Args: "${ARGS}"
Message: ${MSG}
Please check journald for proxy errors
and review the docs: https://github.com/AnalogJ/scrutiny/blob/master/docs/INSTALL_MANUAL.md
EOF
    exit $EXIT_CODE
}


# request the smartctl data
RESPONSE=$(echo "$*" | socat -t 10 - UNIX-CONNECT:/run/scrutiny/smartctl.sock 2>/dev/null)

if [[ -z "$RESPONSE" ]]; then
    send_smartctl_shim_error "Shim: Failed to connect to proxy socket" "$*" "3"
fi

# if no exit_status, assume bad.
# If you run into issues with grepping, switch to jq
#EXIT_STATUS=$(echo "$RESPONSE" | jq ".smartctl.exit_status // 2" 2>/dev/null)
EXIT_STATUS=$(echo "$RESPONSE" | grep -m 1 -oP $'[\'"]exit_status[\'"]:\s*\K\d+')
PARSE_STATUS=$?

# Invalid JSON (smartctl-proxy error) or smartctl error that shouldn't be reported up as json
if [[ ($PARSE_STATUS -gt 0) || ($EXIT_STATUS -gt 0 && $EXIT_STATUS -lt 4) ]]; then
    send_smartctl_shim_error "$RESPONSE" "$*" "${EXIT_STATUS:-3}"
else
    echo "$RESPONSE"
    exit $EXIT_STATUS
fi

```

4. make the shim executable `sudo chmod +x /opt/smartctl-shim/bin/smartctl`
5. create a new directory `sudo mkdir -p /opt/smartctl-proxy/bin`
6. create inside there a proxy script. See below examples using a bash script and python script


<details><summary>Click to expand: smartctl-proxy.sh</summary>

```bash
#!/usr/bin/env bash
SMARTCTL="/usr/sbin/smartctl"

send_smartctl_proxy_error() {
    local EXIT_CODE="${3:-3}"
    printf "%s\n" "=======> PROXYERR <=======" "'exit_status': $EXIT_CODE" "MSG: $1" "ARGS: $2"
    # for logging
    printf "%s\n" "=======> PROXYERR <=======" "'exit_status': $EXIT_CODE" "MSG: $1" "ARGS: $2" >&2
    # still should exit with status code zero so systemd doesn't consider the service failed
    exit 0
}

# Get Args
if ! read -t 2 -r -a ARGS_ARRAY; then
    send_smartctl_proxy_error "No Args... Did you send a newline?"
fi
FULL_LINE="${ARGS_ARRAY[*]}"

# Prevent malicoius forking/chaining to prevent attacks to the privileged socket
# only allow alphanumeric, spaces, dots, underscores, and dashes
if [[ "$FULL_LINE" =~ [^a-zA-Z0-9\ ./\_-] ]]; then
    send_smartctl_proxy_error "Security violation: Illegal characters detected" "$FULL_LINE"
fi

# Must be a scan, info, or xall smartctl action
ACTION="${ARGS_ARRAY[0]}"
if [[ "$ACTION" == "--scan" || "$ACTION" == "--info" || "$ACTION" == "--xall" ]]; then
    # We good
    : 
else
    send_smartctl_proxy_error "Forbidden command action: $ACTION" "$FULL_LINE"
fi

# and confirm json
if [[ "$FULL_LINE" != *"--json"* ]]; then
    send_smartctl_proxy_error "Protocol error: --json flag is mandatory" "$FULL_LINE"
fi

COMMAND_OUTPUT=$($SMARTCTL "${ARGS_ARRAY[@]}")
SMARTCTL_EXIT=$?

# keep json for exit codes for SMART errors, but 
if [[ $SMARTCTL_EXIT -eq 0 || $SMARTCTL_EXIT -gt 3 ]]; then
    echo "$COMMAND_OUTPUT"
    exit 0
else
    send_smartctl_proxy_error "$COMMAND_OUTPUT" "$FULL_LINE" "$SMARTCTL_EXIT"
fi

```

7. make it executable `sudo chmod +x /opt/smartctl-proxy/bin/smartctl-proxy.sh`

</details>

<details><summary>Click to expand: smartctl-proxy.py</summary>

```py
#!/usr/bin/env python3
import sys, subprocess, re, select, shlex

smartctl_bin = '/usr/sbin/smartctl'

smartctl_proxy_err_template = """
=======> PROXYERR <=======
"exit_status": {}
MSG: {}
ARGS: {}
"""

def smartctl_proxy_err(msg, cmd_args, exit_code = 3):
    err = smartctl_proxy_err_template.format(exit_code, msg, cmd_args)
    sys.stderr.write(err)
    sys.stderr.flush()
    print(err, flush=True)
    sys.exit()

def main():
    command_args = [smartctl_bin]
    has_json = False
    allowed_args = ['--scan', '--info', '--xall', '--device', '--json']
    # only allow alphanumeric, underscore, and slashes, if it's not an allowed arg
    allowed_arg_rex = re.compile(r'^[\w/]+$')

    def is_valid(arg: str):
        nonlocal has_json
        if arg == "--json":
            has_json = True
        return arg in allowed_args or allowed_arg_rex.match(arg)
    
    # wait a maximum of 2 seconds for input
    readable, _, _ = select.select([sys.stdin], [], [], 2)
    if (not readable):
        smartctl_proxy_err("No Args... Did you send a newline?", [])

    try:
        raw_input = sys.stdin.readline()
        args = shlex.split(raw_input)
        command_args.extend(args)
        if (len(args) == 0):
            smartctl_proxy_err("No Args... Did you send a newline?", command_args)
        elif all(is_valid(arg) for arg in args):
            if has_json:
                # shim has socat a timeout of 10s, change accordingly, just some clean up stuff
                smartctl_result = subprocess.run(command_args, capture_output=True, check=False, text=True, timeout=10)
                if (smartctl_result.returncode == 0 or smartctl_result.returncode > 3):
                    print(smartctl_result.stdout, flush=True)
                else:
                    smartctl_proxy_err(smartctl_result.stdout, command_args, smartctl_result.returncode)
            else:
                smartctl_proxy_err("Protocol error: --json flag is mandatory", command_args)
        else:
            smartctl_proxy_err("Security violation: Illegal arguments and/or characters detected", command_args)
    except Exception as e:
        smartctl_proxy_err(f"smartctl Proxy Internal Error: {str(e)}", command_args)


if __name__ == "__main__":
    main()


```

7. make it executable `sudo chmod +x /opt/smartctl-proxy/bin/smartctl-proxy.py`

</details>

8. go to `/etc/systemd/system`
9. create `scrutiny-collector.service` with the following contents:


```ini
[Unit]
Description=Scrutiny Collector Service
After=network.target smartctl-proxy.socket

[Service]
Environment="PATH=/opt/smartctl-shim/bin:/usr/bin"
ExecStart=/opt/scrutiny/bin/scrutiny-collector-metrics run --api-endpoint "http://localhost:8080"

DynamicUser=yes
SupplementaryGroups=scrutiny-svc
CapabilityBoundingSet=

RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX
IPAddressDeny=any
IPAddressAllow=localhost

ProtectProc=invisible
ProcSubset=pid
RestrictNamespaces=yes
LockPersonality=yes
MemoryDenyWriteExecute=yes
SystemCallArchitectures=native
SystemCallFilter=@system-service
RestrictRealtime=yes

InaccessiblePaths=/root /boot /home /etc/shadow /etc/ssh /etc/sudoers /etc/sudoers.d
## NOTE: SELinux users should use the below instead since SELinux will already protect /etc/shadow
# InaccessiblePaths=/root /boot /home /etc/ssh /etc/sudoers /etc/sudoers.d
ReadOnlyPaths=/usr /bin

ProtectSystem=strict
ProtectHome=yes
PrivateTmp=yes
PrivateDevices=yes
DevicePolicy=closed
NoNewPrivileges=yes

ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectControlGroups=yes
ProtectClock=yes
ProtectHostname=yes
ProtectKernelLogs=yes
RestrictSUIDSGID=true
RemoveIPC=yes

```

10. create `smartctl-proxy.socket` with the following contents:


```ini
[Unit]
Description=smartctl proxy socket

[Socket]
ListenStream=/run/scrutiny/smartctl.sock
SocketGroup=scrutiny-svc
SocketMode=0660
Accept=yes

MaxConnections=20
Backlog=10

RemoveOnStop=yes
FlushPending=yes
TriggerLimitIntervalSec=10s
TriggerLimitBurst=30

RuntimeDirectory=smartctl
RuntimeDirectoryMode=0755

[Install]
WantedBy=sockets.target

```

11. create `smartctl-proxy@.service` with the following contents:

```ini
[Unit]
Description=smartctl Proxy Service
Requires=smartctl-proxy.socket
After=smartctl-proxy.socket
CollectMode=inactive-or-failed

[Service]
ExecStart=/opt/smartctl-proxy/bin/smartctl-proxy.sh
## If you prefer the python script:
# Environment="PYTHONUNBUFFERED=1"
# ExecStart=/opt/smartctl-proxy/bin/smartctl-proxy.py
StandardInput=socket
StandardOutput=socket
StandardError=journal

DynamicUser=yes
SupplementaryGroups=scrutiny-svc disk

## CAP_SYS_RAWIO is needed for SATA drives
AmbientCapabilities=CAP_SYS_RAWIO
CapabilityBoundingSet=CAP_SYS_RAWIO
## unfortunately nvme drives require CAP_SYS_ADMIN
## if you want nvme drives you must do the following:
#AmbientCapabilities=CAP_SYS_RAWIO CAP_SYS_ADMIN
#CapabilityBoundingSet=CAP_SYS_RAWIO CAP_SYS_ADMIN

# --- LOCKDOWN ---
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
PrivateTmp=yes

PrivateDevices=no
DeviceAllow=block-sd r
DeviceAllow=char-nvme r
InaccessiblePaths=/root /boot /home /etc/shadow /etc/ssh /etc/sudoers /etc/sudoers.d
## NOTE: SELinux users should use the below instead since SELinux will already protect /etc/shadow
# InaccessiblePaths=/root /boot /home /etc/ssh /etc/sudoers /etc/sudoers.d

IPAddressDeny=any
RestrictAddressFamilies=AF_UNIX

ProtectControlGroups=yes
ProtectKernelTunables=yes
ProtectControlGroups=yes
ProtectKernelModules=yes
ProtectKernelLogs=yes
ProtectClock=yes
ProtectHostname=yes
LockPersonality=yes
MemoryDenyWriteExecute=yes
ProtectKernelLogs=yes
RemoveIPC=yes
RestrictSUIDSGID=true

MemoryMax=64M
TasksMax=5
TimeoutStopSec=15s

```

Note the at-sign is on purpose. That instructs systemd that it's a template unit file and allows spawning multiple instances as-needed for each socket activation.

12. enable the proxy socket:

```sh
# reload changes for systemd services
sudo systemctl daemon-reload

# enable the proxy socket
sudo systemctl enable smartctl-proxy.socket

# start the socket
sudo systemctl start smartctl-proxy.socket
sudo systemctl status smartctl-proxy.socket

```

13. Test the proxy:

```sh
# should be successful and output json
/opt/smartctl-shim/bin/smartctl --scan --json

# should fail due to proxy error for illegal arguments
/opt/smartctl-shim/bin/smartctl --~xall --json /dev/sda
```

##### Additiona SELinux Considerations

if you are using SELinux, you may need to also do the following:

```sh
# tell SELinux to allow these binaries
sudo semanage fcontext -a -t bin_t "/opt/smartctl-proxy/bin(/.*)?"
sudo semanage fcontext -a -t bin_t "/opt/smartctl-shim/bin(/.*)?"
# update labels
sudo restorecon -Rv /opt/smartctl-proxy/bin
sudo restorecon -Rv /opt/smartctl-shim/bin
```


##### Pros:

- the scrutiny binary itself is extremely locked down with virtually zero privileges
- almost completely immune to the threat of supply-chain attacks
- only smartctl itself runs with elevated privileges (but still non-root)
- the proxy script allows you to restrict what smartctl capabilities are exposed, further limiting the attack surface


##### Cons:

NOTE: These cons basically only apply if a major supply-chain attack happens against scrutiny, and reflect a worst-case scenario that is unlikely to ever occur:

- Much more involved to set up, which increases the risk of something breaking (socket can fail and need to be manually restarted, for example)
- You are responsible for validating and maintaining the proxy scripts. The above examples are a good starting place but may need to be expanded in the future

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

That's it! you're done. You can check the status of the timer using `sudo systemctl status scrutiny-collector.timer`