# InfluxDB Troubleshooting

## Installation 
InfluxDB is a required dependency for Scrutiny v0.4.0+. 

https://docs.influxdata.com/influxdb/v2.2/install/

## Persistence

To ensure that all data is correctly stored, you must also persist the InfluxDB database directory
    - `/opt/scrutiny/influxdb` (for Docker omnibus image)
    - `/var/lib/influxdb2` (for vanilla Influxdb image `influxdb:2.2`)

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