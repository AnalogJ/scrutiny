# Docker Images `master-omnibus` vs `latest`

> TL;DR; The `master-omnibus` and `latest` tags are almost semantically identical, as I follow a `golden master` 
development process. However if you want to ensure you're only using the latest release, you can change to `latest`

The CI script used to orchestrate the docker image builds can be found here: https://github.com/AnalogJ/scrutiny/blob/master/.github/workflows/docker-build.yaml#L166-L184

In general Scrutiny follows a `golden master` development process, which means that the `master` branch is not directly updated (unless its for documentation changes), 
instead development is done in a feature branch, or committed to the `beta` branch. 

As development progresses, and we're satisfied that a feature is complete, and the quality is acceptable, 
I merge the changes to `master` and trigger the creation of a new release -- ie, when master is updated, a new release
is almost immediately created (and tagged with `latest`)

So changing from `master-omnibus -> latest` will be the same thing for all intents and purposes. 

Having said that -- the one key difference is the `automated cron builds` that run on the `master` and `beta` branches. 
They trigger a `nightly` build, even if nothing has changed on the branch. This has a couple of benefits, but one is to 
ensure that there's no broken external dependencies in our (unchanged) code. 

However, as everyone unfortunately found out recently, I had an error in my CI script, which caused failures to be 
ignored -- https://github.com/AnalogJ/scrutiny/issues/287. That has since been fixed.

Hope that gives you an understanding for how everything is wired up. 

