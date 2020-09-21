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

```
docker exec scrutiny scrutiny-collector-metrics run --debug --log-file /tmp/test.log
# then use docker cp to copy the log file out of the container. 
docker cp scrutiny:/tmp/test.log test.log
```
