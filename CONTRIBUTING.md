# Contributing

The Scrutiny repository is a [monorepo](https://en.wikipedia.org/wiki/Monorepo) containing source code for:
- Scrutiny Backend Server (API)
- Scrutiny Frontend Angular SPA
- S.M.A.R.T Collector

Depending on the functionality you are adding, you may need to setup a development environment for 1 or more projects. 

# Modifying the Scrutiny Backend Server (API)

1. install the [Go runtime](https://go.dev/doc/install) (v1.17+)
2. download the `scrutiny-web-frontend.tar.gz` for the [latest release](https://github.com/AnalogJ/scrutiny/releases/latest). Extract to a folder named `dist`
3. create a `scrutiny.yaml` config file
    ```yaml
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
4. start a InfluxDB docker container. 
    ```bash
    docker run -p 8086:8086 --rm influxdb:2.2
    ```
5. start the scrutiny web server
    ```bash
    go mod vendor
    go run webapp/backend/cmd/scrutiny/scrutiny.go start --config ./scrutiny.yaml
    ```
6. open your browser to [http://localhost:8080/web](http://localhost:8080/web)

# Modifying the Scrutiny Frontend Angular SPA

The frontend is written in Angular. If you're working on the frontend and can use mocked data rather than a real backend, you can follow the instructions below:

1. install [NodeJS](https://nodejs.org/en/download/)
2. start the Angular Frontend Application
    ```bash
    cd webapp/frontend
    npm install
    ng serve --deploy-url="/web/" --base-href="/web/" --port 4200
    ```
3. open your browser and visit [http://localhost:4200/web](http://localhost:4200/web)

# Modifying both Scrutiny Backend and Frontend Applications
If you're developing a feature that requires changes to the backend and the frontend, or a frontend feature that requires real data, 
you'll need to follow the steps below:

1. install the [Go runtime](https://go.dev/doc/install) (v1.17+)
2. install [NodeJS](https://nodejs.org/en/download/)
3. create a `scrutiny.yaml` config file
    ```yaml
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
4. start a InfluxDB docker container.
    ```bash
    docker run -p 8086:8086 --rm influxdb:2.2
    ```
5. build the Angular Frontend Application
    ```bash
    cd webapp/frontend
    npm install
    ng build --watch --output-path=../../dist --prod
    # Note: if you do not add `--prod` flag, app will display mocked data for api calls.
    ```
6. start the scrutiny web server
    ```bash
    go mod vendor
    go run webapp/backend/cmd/scrutiny/scrutiny.go start --config ./scrutiny.yaml
    ```
7. open your browser to [http://localhost:8080/web](http://localhost:8080/web)


If you'd like to populate the database with some test data,  you can run the following commands:

> NOTE: you may need to update the `local_time` key within the JSON file, any timestamps older than ~3 weeks will be automatically ignored
> (since the downsampling & retention policy takes effect at 2 weeks)
> This is done automatically by the `webapp/backend/pkg/models/testdata/helper.go` script

```
docker run -p 8086:8086 --rm influxdb:2.2


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

# Modifying the Collector
```
brew install smartmontools
go run collector/cmd/collector-metrics/collector-metrics.go run --debug
```


# Debugging

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

# Docker Development

```
docker build -f docker/Dockerfile . -t chcr.io/analogj/scrutiny:master-omnibus
docker run -it --rm -p 8080:8080 \
-v /run/udev:/run/udev:ro \
--cap-add SYS_RAWIO \
--device=/dev/sda \
--device=/dev/sdb \
ghcr.io/analogj/scrutiny:master-omnibus
/opt/scrutiny/bin/scrutiny-collector-metrics run
```
