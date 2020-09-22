# Contributing

There are multiple ways to develop on the scrutiny codebase locally. The two most popular are:
- Docker Development Container - only requires docker
- Run Components Locally - requires smartmontools, golang & nodejs installed locally

## Docker Development
```
docker build -f docker/Dockerfile . -t analogj/scrutiny
docker run -it --rm -p 9090:8080 -v /run:/run  -v /dev/disk:/dev/disk --privileged analogj/scrutiny
/scrutiny/bin/scrutiny-collector-metrics run
```


## Local Development

### Frontend
The frontend is written in Angular.
If you're working on the frontend and can use mocked data rather than a real backend, you can use
```
cd webapp/frontend && ng serve
```

However, if you need to also run the backend, and use real data, you'll need to run the following command:
```
cd webapp/frontend && ng build --watch --output-path=../../dist --deploy-url="/web/" --base-href="/web/" --prod
```

> Note: if you do not add `--prod` flag, app will display mocked data for api calls.

### Backend
```
go run webapp/backend/cmd/scrutiny/scrutiny.go start --config ./example.scrutiny.yaml
```
Now visit http://localhost:8080


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
