# Changelog

All notable changes to Scrutiny will be documented in this file.

## [1.1.0](https://github.com/Starosdev/scrutiny/compare/v1.0.0...v1.1.0) (2025-11-30)

### Features

* Add "day" as resolution for temperature graph ([2670af2](https://github.com/Starosdev/scrutiny/commit/2670af216d491c478b36f8ef20497c5cb6002801))
* add day resolution for temperature graph (upstream PR [#823](https://github.com/Starosdev/scrutiny/issues/823)) ([2d6ffa7](https://github.com/Starosdev/scrutiny/commit/2d6ffa732cda4583c0f867540bed87a331fbb6d4))
* add setting to enable/disable SCT temperature history (upstream PR [#557](https://github.com/Starosdev/scrutiny/issues/557)) ([c3692ac](https://github.com/Starosdev/scrutiny/commit/c3692acd17e310e1c5d1470404566ae13e67d9a5))
* Implement device-wise notification mute/unmute ([925e86d](https://github.com/Starosdev/scrutiny/commit/925e86d461fc2bfe4f318851d790a08d99eb6bde))
* implement device-wise notification mute/unmute (upstream PR [#822](https://github.com/Starosdev/scrutiny/issues/822)) ([ea7102e](https://github.com/Starosdev/scrutiny/commit/ea7102e9297aeb011a808f1133fbf03114176900))
* implement Prometheus metrics support (upstream PR [#830](https://github.com/Starosdev/scrutiny/issues/830)) ([7384f7d](https://github.com/Starosdev/scrutiny/commit/7384f7de6ebf8f6c3936fb52d19ffe3b805bae0c))
* support SAS temperature (upstream PR [#816](https://github.com/Starosdev/scrutiny/issues/816)) ([f954cc8](https://github.com/Starosdev/scrutiny/commit/f954cc815f756bef8842f026a5a0e554bfd5ba80))

### Bug Fixes

* better handling of ata_sct_temperature_history (upstream PR [#825](https://github.com/Starosdev/scrutiny/issues/825)) ([d134ad7](https://github.com/Starosdev/scrutiny/commit/d134ad7160b754ad25d10d600a6fc8e56c0d5914))
* **database:** add missing temperature parameter in SCSI migration ([df7da88](https://github.com/Starosdev/scrutiny/commit/df7da8824c3cd3745f66ae426bcec1db7844e840))
* support transient SMART failures (upstream PR [#375](https://github.com/Starosdev/scrutiny/issues/375)) ([601775e](https://github.com/Starosdev/scrutiny/commit/601775e462f6cd56d442386071c6499dfba3cc39))
* **ui:** fix temperature conversion in temperature.pipe.ts (upstream PR [#815](https://github.com/Starosdev/scrutiny/issues/815)) ([e0f2781](https://github.com/Starosdev/scrutiny/commit/e0f27819facc20c6f04c8903f2ebb85035475b47))

### Refactoring

* use limit() instead of tail() for fetching smart attributes (upstream PR [#829](https://github.com/Starosdev/scrutiny/issues/829)) ([2849531](https://github.com/Starosdev/scrutiny/commit/2849531d3893028861cec68f862d4ed32bedbb0c))

## 1.0.0 (2025-11-29)

### Features

* Ability to override commands args ([604dcf3](https://github.com/Starosdev/scrutiny/commit/604dcf355ce387de5b5030473163838c5855fa31))
* create allow-list for filtering down devices to only a subset ([c9429c6](https://github.com/Starosdev/scrutiny/commit/c9429c61b2aa7dbea9ed412bd9d49326cf408e94))
* dynamic line stroke settings ([536b590](https://github.com/Starosdev/scrutiny/commit/536b590080b589a807765b69612990d41ae97773))
* Update dashboard.component.ts ([bb98b8c](https://github.com/Starosdev/scrutiny/commit/bb98b8c45b13d9b01c3a543022608fb746b207d6))

### Bug Fixes

* **collector:** show correct nvme capacity ([db86bac](https://github.com/Starosdev/scrutiny/commit/db86bac9efb10ca11177a1cf00621a8ea91dc6aa)), closes [#466](https://github.com/Starosdev/scrutiny/issues/466)
* https://github.com/AnalogJ/scrutiny/issues/643 ([50561f3](https://github.com/Starosdev/scrutiny/commit/50561f34ead034c118dd7ea5f1d1f067b0d1d97a))
* igeneric types ([e9cf8a9](https://github.com/Starosdev/scrutiny/commit/e9cf8a9180e5d181f62076bb602888e34596885b))
* increase timeout ([222b810](https://github.com/Starosdev/scrutiny/commit/222b8103d635ddfafd29ac93ea110c3d851a3112))
* prod build command ([50321d8](https://github.com/Starosdev/scrutiny/commit/50321d897a21faa515b142f4b2e285ba16815acd))
* remove fullcalendar ([64ad353](https://github.com/Starosdev/scrutiny/commit/64ad3536284f67cb4652a9e83a02f0024b7dcde9))
* remove outdated option ([5518865](https://github.com/Starosdev/scrutiny/commit/5518865bc69f0a9906977facfa4be8895a7b12d9))

### Refactoring

* update dependencies version ([e18a7e9](https://github.com/Starosdev/scrutiny/commit/e18a7e9ce08e9172853f7bd5f6a6388e278ee4e2))
