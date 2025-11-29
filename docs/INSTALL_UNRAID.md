# UnRAID Install

Installation of Scrutiny in UnRAID follows the same process as installing any other docker container, utilizing the Community Applications plugin

## Install the 'Community Applications' Plugin

All docker containers in UnRAID are typically installed utilizing the Community Applications plugin. To get started:
- Navigate to the plugins tab ( <UnRaid_IP_Address>/Plugins )
- Select the 'Install Plugin' tab, and enter the following address into the input field
```
https://raw.githubusercontent.com/Squidly271/community.applications/master/plugins/community.applications.plg
```

You're all set with the pre-requisites!

## Installing the Scrutiny docker image

To install, simply click 'Install'; the configuration parameters should not need modification as the template within CA already defines the necessary parameters.

As a docker image can be created using various OS bases, the image choice is entirely the users choice. Recommendations of a specific image from a specific maintainer is beyond the scope of this guide. However, to provide some context given the number of questions posed regarding the various versions available:

- **ghcr.io/Starosdev/scrutiny:master-omnibus**
    - `Image maintained directly by the application author`
    - `Debian based docker image`
- **linuxserver/scrutiny**
    - `Image maintained by the LinuxServer.io group`
    - `Alpine based docker image`
- **hotio/scrutiny**
    - `Image maintained by hotio`
    - `DETAILS TBD`

The support for a given image is provided by that images maintainers, while support for the application itself remains with the developer - i.e. LinuxServer.io supports the docker image of Scrutiny which they create, to the extent an issue is specific to that image. If an issue/enhancement pertains directly to the source code, support would still come directly from this repository's contributors. 
