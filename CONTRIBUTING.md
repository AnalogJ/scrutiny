# Contributing

There are multiple ways to develop on the scrutiny codebase locally. The two most popular are:
- Docker Development Container - only requires docker
- Run Components Locally - requires smartmontools, golang & nodejs installed locally

## Docker Development
```
docker build -f docker/Dockerfile . -t analogj/scrutiny
docker run -it --rm -p 8080:8080 \
-v /run/udev:/run/udev:ro \
--cap-add SYS_RAWIO \
--device=/dev/sda \
--device=/dev/sdb \
analogj/scrutiny
/scrutiny/bin/scrutiny-collector-metrics run
```


## Local Development

### Frontend
The frontend is written in Angular.
If you're working on the frontend and can use mocked data rather than a real backend, you can use
```
cd webapp/frontend
npm install
ng serve
```

However, if you need to also run the backend, and use real data, you'll need to run the following command:
```
cd webapp/frontend && ng build --watch --output-path=../../dist --deploy-url="/web/" --base-href="/web/" --prod
```

> Note: if you do not add `--prod` flag, app will display mocked data for api calls.

### Backend

If you're using the `ng build` command above to generate your frontend, you'll need to create a custom config file and
override the `web.src.frontend.path` value.

```
# config file for local development. store as scrutiny.yaml
version: 1

web:
  listen:
    port: 8080
    host: 0.0.0.0
  database:
    # can also set absolute path here
    location: ./scrutiny.db
  src:
    frontend:
      path: ./dist
  influxdb:
    retention_policy: false

log:
  file: 'web.log' #absolute or relative paths allowed, eg. web.log
  level: DEBUG

```

Once you've created a config file, you can pass it to the scrutiny binary during startup.

```
go run webapp/backend/cmd/scrutiny/scrutiny.go start --config ./scrutiny.yaml
```

Now visit http://localhost:8080


If you'd like to populate the database with some test data,  you can run the following commands:

> NOTE: you may need to update the `local_time` key within the JSON file, any timestamps older than ~3 weeks will be automatically ignored
> (since the downsampling & retention policy takes effect at 2 weeks)
> This is done automatically by the `webapp/backend/pkg/models/testdata/helper.go` script

```
docker run -p 8086:8086 --rm influxdb:2.2


docker run --rm -p 8086:8086 \
      -e DOCKER_INFLUXDB_INIT_MODE=setup \
      -e DOCKER_INFLUXDB_INIT_USERNAME=admin \
      -e DOCKER_INFLUXDB_INIT_PASSWORD=password12345 \
      -e DOCKER_INFLUXDB_INIT_ORG=scrutiny \
      -e DOCKER_INFLUXDB_INIT_BUCKET=metrics \
      influxdb:2.2


# curl -X POST -H "Content-Type: application/json" -d @webapp/backend/pkg/web/testdata/register-devices-req.json localhost:8080/api/devices/register
# curl -X POST -H "Content-Type: application/json" -d @webapp/backend/pkg/models/testdata/smart-ata.json localhost:8080/api/device/0x5000cca264eb01d7/smart
# curl -X POST -H "Content-Type: application/json" -d @webapp/backend/pkg/models/testdata/smart-ata-date.json localhost:8080/api/device/0x5000cca264eb01d7/smart
# curl -X POST -H "Content-Type: application/json" -d @webapp/backend/pkg/models/testdata/smart-ata-date2.json localhost:8080/api/device/0x5000cca264eb01d7/smart
# curl -X POST -H "Content-Type: application/json" -d @webapp/backend/pkg/models/testdata/smart-fail2.json localhost:8080/api/device/0x5000cca264ec3183/smart
# curl -X POST -H "Content-Type: application/json" -d @webapp/backend/pkg/models/testdata/smart-nvme.json localhost:8080/api/device/0x5002538e40a22954/smart
# curl -X POST -H "Content-Type: application/json" -d @webapp/backend/pkg/models/testdata/smart-scsi.json localhost:8080/api/device/0x5000cca252c859cc/smart
# curl -X POST -H "Content-Type: application/json" -d @webapp/backend/pkg/models/testdata/smart-scsi2.json localhost:8080/api/device/0x5000cca264ebc248/smart
go run webapp/backend/pkg/models/testdata/helper.go

curl localhost:8080/api/summary

```

### Collector
```
brew install smartmontools
go run collector/cmd/collector-metrics/collector-metrics.go run --debug
```


## Debugging

If you need more verbose logs for debugging, you can use the following environmental variables:

- `DEBUG=true` - enables debug level logging on both the `collector` and `webapp`
- `COLLECTOR_DEBUG=true` - enables debug level logging on the `collector`
- `SCRUTINY_DEBUG=true` - enables debug level logging on the `webapp`

In addition, you can instruct scrutiny to write its logs to a file using the following environmental variables:

- `COLLECTOR_LOG_FILE=/tmp/collector.log` - write the `collector` logs to a file
- `SCRUTINY_LOG_FILE=/tmp/web.log` - write the `webapp` logs to a file

Finally, you can copy the files from the scrutiny container to your host using the following command(s)

```
docker cp scrutiny:/tmp/collector.log collector.log
docker cp scrutiny:/tmp/web.log web.log
```
