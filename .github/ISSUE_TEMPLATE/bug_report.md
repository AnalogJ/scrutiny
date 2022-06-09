---
name: Bug report
about: Create a report to help us improve
title: "[BUG]"
labels: bug
assignees: ''

---

**Describe the bug**
A clear and concise description of what the bug is.

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Log Files**
If related to missing devices or SMART data, please run the `collector` in DEBUG mode, and attach the log file.
See [/docs/TROUBLESHOOTING_DEVICE_COLLECTOR.md](docs/TROUBLESHOOTING_DEVICE_COLLECTOR.md) for other troubleshooting tips. 

```
docker run -it --rm -p 8080:8080 \
-v `pwd`/config:/opt/scrutiny/config \
-v /run/udev:/run/udev:ro \
--cap-add SYS_RAWIO \
--device=/dev/sda \
--device=/dev/sdb \
-e DEBUG=true \
-e COLLECTOR_LOG_FILE=/opt/scrutiny/config/collector.log \
-e SCRUTINY_LOG_FILE=/opt/scrutiny/config/web.log \
--name scrutiny \
ghcr.io/analogj/scrutiny:master-omnibus

# in another terminal trigger the collector
docker exec scrutiny scrutiny-collector-metrics run
```

The log files will be available on your host in the `config` directory. Please attach them to this issue. 

Please also provide the output of `docker info`