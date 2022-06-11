# InfluxDB Troubleshooting

## Installation 
InfluxDB is a required dependency for Scrutiny v0.4.0+. 

https://docs.influxdata.com/influxdb/v2.2/install/

## Persistence

To ensure that all data is correctly stored, you must also persist the InfluxDB database directory

- If you're using the Official Scrutiny Omnibus image (`ghcr.io/analogj/scrutiny:master-omnibus`), the path is `/opt/scrutiny/influxdb`
- If you're deploying in Hub/Spoke mode with the InfluxDB maintained image (`influxdb:2.2`), the path is `/var/lib/influxdb2`

If you attempt to restart Scrutiny but you forgot to persist the InfluxDB directory, you will get an error message like follows:

```
scrutiny    | time="2022-05-12T22:54:12Z" level=info msg="Trying to connect to scrutiny sqlite db: /opt/scrutiny/config/scrutiny.db\n"
scrutiny    | time="2022-05-12T22:54:12Z" level=info msg="Successfully connected to scrutiny sqlite db: /opt/scrutiny/config/scrutiny.db\n"
scrutiny    | ts=2022-05-12T22:54:12.240791Z lvl=info msg=Unauthorized log_id=0aQcVlOW000 error="authorization not found"
scrutiny    | panic: unauthorized: unauthorized access
```

Unfortunately this may mean that your database is lost, and the previous Scrutiny data is unavailable. 
You should fix the docker-compose/docker run command that you're using to ensure that your database folder is persisted correctly, 
then delete the `web.influxdb.token` field in your `scrutiny.yaml` file, and then restart Scrutiny.


## First Start
The web/api service will trigger an InfluxDB onboarding process automatically when it first starts. After that, it will store the newly generated influxdb api token in the Scrutiny config file. 

If this Credential is not correctly stored in the scrutiny config file, Scrutiny will fail to start (with an authentication error)

```
scrutiny    | time="2022-05-12T22:52:55Z" level=info msg="Successfully connected to scrutiny sqlite db: /opt/scrutiny/config/scrutiny.db\n"
scrutiny    | ts=2022-05-12T22:52:55.235753Z lvl=error msg="failed to onboard user admin" log_id=0aQcRnc0000 handler=onboard error="onboarding has already been completed" took=0.038ms
scrutiny    | ts=2022-05-12T22:52:55.235816Z lvl=error msg="api error encountered" log_id=0aQcRnc0000 error="onboarding has already been completed"
scrutiny    | panic: conflict: onboarding has already been completed
```

You can fix this issue by authenticating to the InfluxDB admin portal (the default credentials are username: `admin`, password: `password12345`),
then retrieving the API token, and writing it to your `scrutiny.yaml` config file under the `web.influxdb.token` field:

![influx db admin token](./influxdb-admin-token.png)

## Upgrading from v0.3.x to v0.4.x

When upgrading from v0.3.x to v0.4.x, some users have noticed problems such as:

```
2022/05/13 14:38:05 Loading configuration file: /opt/scrutiny/config/scrutiny.yaml
time="2022-05-13T14:38:05Z" level=info msg="Trying to connect to scrutiny sqlite db:"
time="2022-05-13T14:38:05Z" level=info msg="Successfully connected to scrutiny sqlite db:"
panic: a username and password is required for a setup
```

or 

```
Start the scrutiny server
time="2022-06-11T10:35:04-04:00" level=info msg="Trying to connect to scrutiny sqlite db: \n"
time="2022-06-11T10:35:04-04:00" level=info msg="Successfully connected to scrutiny sqlite db: \n"
panic: failed to check influxdb setup status - parse "://:": missing protocol scheme
```

As discussed in [#248](https://github.com/AnalogJ/scrutiny/issues/248) and [#234](https://github.com/AnalogJ/scrutiny/issues/234),
this usually related to either:

- Upgrading from the LSIO Scrutiny image to the Official Scrutiny image, without removing LSIO specific environmental variables
  - remove the `SCRUTINY_WEB=true` and `SCRUTINY_COLLECTOR=true` environmental variables. They were used by the LSIO image, but are unnecessary and cause issues with the official Scrutiny image.
- Updated versions of the [LSIO Scrutiny images are broken](https://github.com/linuxserver/docker-scrutiny/issues/22), as they have not installed InfluxDB which is a required dependency of Scrutiny v0.4.x
  - You can revert to an earlier version of the LSIO image (`lscr.io/linuxserver/scrutiny:060ac7b8-ls34`), or just change to the official Scrutiny image (`ghcr.io/analogj/scrutiny:master-omnibus`)

Here's a couple of confirmed working docker-compose files that you may want to look at:

- https://github.com/AnalogJ/scrutiny/blob/master/docker/example.hubspoke.docker-compose.yml
- https://github.com/AnalogJ/scrutiny/blob/master/docker/example.omnibus.docker-compose.yml
