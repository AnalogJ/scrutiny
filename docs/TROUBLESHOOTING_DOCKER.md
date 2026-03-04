# Docker Images `latest` vs `nightly`

> TL;DR; The `latest-omnibus`, `latest-collector`, and `latest-web` tags point to the most recent release. (`latest` points to `latest-omnibus`)
> The `nightly-omnibus`, `nightly-collector`, and `nightly-web` tags point to builds that are generated every night from the latest commit on the `master` branch.

The CD scripts used to orchestrate the docker image builds can be found here:
* https://github.com/AnalogJ/scrutiny/blob/master/.github/workflows/docker-build.yaml
* https://github.com/AnalogJ/scrutiny/blob/master/.github/workflows/docker-nightly.yaml

In general scrutiny follows a feature branch development process, which means that the `master` branch should ideally always be free of bugs 
This is driven by the requirement that every PR be reviewed and pass all tests.  Unfortunately, bugs do make it through, especially because of the
enormous number of hard drives that scrutiny must support.. 

This means that while the nightly builds should have the latest features and bug fixes, there may be things that sneak through. Unless you need a particular
feature or bug fix, we recommend sticking to releases. Also note that using `latest` tags is generally considered a bad practice; pin a specific version instead.

# Running Docker `rootless`

To avoid that the container(s) restart when you installed Docker as `rootless` you need to isssue the following commands to allow the session to stay alive even after you close your (SSH) sesssion:

`sudo loginctl enable-linger $(whoami)`

`systemctl --user enable docker`
