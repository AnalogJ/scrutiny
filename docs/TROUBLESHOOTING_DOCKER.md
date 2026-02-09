# Docker Images `nightly` vs `latest`

> TL;DR; The `latest-omnibus`, `latest-collector`, and `latest-web` tags point to the most recent release. (`latest` points to `latest-omnibus`)
> The `nightly-omnibus`, `nightly-collector`, and `nightly-web` tags point to builds that are generated every night from the latest commit.

The CD scripts used to orchestrate the docker image builds can be found here:
* https://github.com/AnalogJ/scrutiny/blob/master/.github/workflows/docker-build.yaml
* https://github.com/AnalogJ/scrutiny/blob/master/.github/workflows/docker-nightly.yaml

In general Scrutiny follows a feature branch development process, which means that the `master` branch should always be fully functional, 
and bug free. This is softly guaranteed by the fact that every PR, and thus any commit, must pass all testing. 

As development progresses, and we're satisfied that a feature is complete, and the quality is acceptable, 
I merge the changes to `master` and trigger the creation of a new release -- ie, when master is updated, a new release
is almost immediately created (and tagged with `latest`)

So changing from `master-omnibus -> latest` will be the same thing for all intents and purposes. 



# Running Docker `rootless`

To avoid that the container(s) restart when you installed Docker as `rootless` you need to isssue the following commands to allow the session to stay alive even after you close your (SSH) sesssion:

`sudo loginctl enable-linger $(whoami)`

`systemctl --user enable docker`
