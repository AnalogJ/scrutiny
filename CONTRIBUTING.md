# Contributing to scrutiny

This document describes the process of contributing to scrutiny. It is intended
for anyone considering opening an **issue**, **discussion** or **pull request**.

> [!NOTE]
>
> The intention of these policies is not to be difficult, and
> contributions are greatly appreciated. The goal is to streamline
> and simplify the efforts of both contributers and maintainers.

## AI Usage

scrutiny has strict rules for AI usage. Please see
the [AI Usage Policy](AI_POLICY.md). **This is very important.**

## Quick Guide

### I'd like to contribute

[All issues are actionable](#issues-are-actionable). Pick one and start
working on it. Thank you. If you need help or guidance, comment on the issue.
Issues that are extra friendly to new contributors are tagged with
["contributor friendly"].

["contributor friendly"]: https://github.com/AnalogJ/scrutiny/issues?q=is%3Aissue%20is%3Aopen%20label%3A%22contributor%20friendly%22

### I have a bug! / Something isn't working

First, search the issue tracker and discussions for similar issues. Tip: also
search for [closed issues] and [discussions] — your issue might have already
been fixed!

> [!NOTE]
>
> If there is an _open_ issue or discussion that matches your problem,
> **please do not comment on it unless you have valuable insight to add**.
>
> GitHub has a very _noisy_ set of default notification settings which
> sends an email to _every participant_ in an issue/discussion every time
> someone adds a comment. Instead, use the handy upvote button for discussions,
> and/or emoji reactions on both discussions and issues, which are a visible
> yet non-disruptive way to show your support.

If your issue hasn't been reported already, open an ["Issue Triage"] discussion
and make sure to fill in the template **completely**. They are vital for
maintainers to figure out important details about your setup.

> [!WARNING]
>
> A _very_ common mistake is to file a bug report either as a Q&A or a Feature
> Request. **Please don't do this.** Otherwise, maintainers would have to ask
> for your system information again manually, and sometimes they will even ask
> you to create a new discussion because of how few detailed information is
> required for other discussion types compared to Issue Triage.
>
> Because of this, please make sure that you _only_ use the "Issue Triage"
> category for reporting bugs — thank you!

[closed issues]: https://github.com/AnalogJ/scrutiny/issues?q=is%3Aissue%20state%3Aclosed
[discussions]: https://github.com/AnalogJ/scrutiny/discussions?discussions_q=is%3Aclosed
["Issue Triage"]: https://github.com/AnalogJ/scrutiny/discussions/new?category=issue-triage

### I have an idea for a feature

Like bug reports, first search through both issues and discussions and try to
find if your feature has already been requested. Otherwise, open a discussion
in the ["Feature Requests, Ideas"] category.

["Feature Requests, Ideas"]: https://github.com/AnalogJ/scrutiny/discussions/new?category=feature-requests-ideas

### I've implemented a feature

1. If there is an issue for the feature, open a pull request straight away.
2. If there is no issue, open a discussion and link to your branch.
3. If you want to live dangerously, open a pull request and
   [hope for the best](#pull-requests-implement-an-issue).

### I have a question which is neither a bug report nor a feature request

Open a [Q&A discussion].

> [!NOTE]
> If your question is about a missing feature, please open a discussion under
> the ["Feature Requests, Ideas"] category. If scrutiny is behaving
> unexpectedly, use the ["Issue Triage"] category.
>
> The "Q&A" category is strictly for other kinds of discussions and do not
> require detailed information unlike the two other categories, meaning that
> maintainers would have to spend the extra effort to ask for basic information
> if you submit a bug report under this category.
>
> Therefore, please **pay attention to the category** before opening
> discussions to save us all some time and energy. Thank you!

[Q&A discussion]: https://github.com/AnalogJ/scrutiny/discussions/new?category=q-a

## General Patterns

### Issues are Actionable

The scrutiny [issue tracker](https://github.com/AnalogJ/scrutiny/issues)
is for _actionable items_.

Unlike some other projects, scrutiny **does not use the issue tracker for
discussion or feature requests**. Instead, we use GitHub
[discussions](https://github.com/AnalogJ/scrutiny/discussions) for that.
Once a discussion reaches a point where a well-understood, actionable
item is identified, it is moved to the issue tracker. **This pattern
makes it easier for maintainers or contributors to find issues to work on
since _every issue_ is ready to be worked on.**

If you are experiencing a bug and have clear steps to reproduce it, please
open an issue. If you are experiencing a bug but you are not sure how to
reproduce it or aren't sure if it's a bug, please open a discussion.
If you have an idea for a feature, please open a discussion.

### Pull Requests Implement an Issue

Pull requests should be associated with a previously accepted issue.
**If you open a pull request for something that wasn't previously discussed,**
it may be closed or remain stale for an indefinite period of time. I'm not
saying it will never be accepted, but the odds are stacked against you.

Issues tagged with "feature" represent accepted, well-scoped feature requests.
If you implement an issue tagged with feature as described in the issue, your
pull request will be accepted with a high degree of certainty.

> [!NOTE]
>
> **Pull requests are NOT a place to discuss feature design.** Please do
> not open a WIP pull request to discuss a feature. Instead, use a discussion
> and link to your branch.

# Developer Guide

> [!NOTE]
>
> **The remainder of this file is dedicated to developers actively
> working on scrutiny.** If you're a user reporting an issue, you can
> ignore the rest of this document.

The Scrutiny repository is a [monorepo](https://en.wikipedia.org/wiki/Monorepo) containing source code for:
- Scrutiny Backend Server (API)
- Scrutiny Frontend Angular SPA
- S.M.A.R.T Collector

Depending on the functionality you are adding, you may need to setup a development environment for 1 or more projects.

# Modifying the Scrutiny Backend Server (API)

1. install the [Go runtime](https://go.dev/doc/install) (v1.25)
2. download the `scrutiny-web-frontend.tar.gz` for
   the [latest release](https://github.com/AnalogJ/scrutiny/releases/latest). Extract to a folder named `dist`
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
    npm run start -- --serve-path="/web/" --port 4200
    ```
3. open your browser and visit [http://localhost:4200/web](http://localhost:4200/web)

# Modifying both Scrutiny Backend and Frontend Applications
If you're developing a feature that requires changes to the backend and the frontend, or a frontend feature that requires real data,
you'll need to follow the steps below:

1. install the [Go runtime](https://go.dev/doc/install) (v1.20+)
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
    npm run build:prod -- --watch --output-path=../../dist
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
docker build -f docker/Dockerfile . -t ghcr.io/analogj/scrutiny:master-omnibus
docker run -it --rm -p 8080:8080 \
-v /run/udev:/run/udev:ro \
--cap-add SYS_RAWIO \
--device=/dev/sda \
--device=/dev/sdb \
ghcr.io/analogj/scrutiny:master-omnibus
/opt/scrutiny/bin/scrutiny-collector-metrics run
```


# Running Tests

```bash
docker run -p 8086:8086 -d --rm \
-e DOCKER_INFLUXDB_INIT_MODE=setup \
-e DOCKER_INFLUXDB_INIT_USERNAME=admin \
-e DOCKER_INFLUXDB_INIT_PASSWORD=password12345 \
-e DOCKER_INFLUXDB_INIT_ORG=scrutiny \
-e DOCKER_INFLUXDB_INIT_BUCKET=metrics \
-e DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=my-super-secret-auth-token \
influxdb:2.2
go test ./...

```
