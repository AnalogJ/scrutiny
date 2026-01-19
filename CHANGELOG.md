# Changelog

All notable changes to Scrutiny will be documented in this file.

## [1.8.0](https://github.com/Starosdev/scrutiny/compare/v1.7.2...v1.8.0) (2026-01-17)

### Features

* **thresholds:** add metadata for all ATA Device Statistics ([6602bf8](https://github.com/Starosdev/scrutiny/commit/6602bf83aea4fbc625a71e1232b732078b511694))
* **thresholds:** add metadata for all remaining unknown attributes ([163284c](https://github.com/Starosdev/scrutiny/commit/163284c196eca32a4ed25beaefb9c6bc536f91be))
* **thresholds:** add metadata for Page 3 and Page 5 device statistics ([72a7ea8](https://github.com/Starosdev/scrutiny/commit/72a7ea84ad8d60a1c66d7d9adf42c254bc4dd863))
* **thresholds:** add metadata for remaining unknown attributes ([7b0acd0](https://github.com/Starosdev/scrutiny/commit/7b0acd0b572a392c9c74c5f61f2860e671b44b3e))

### Bug Fixes

* **frontend:** handle device statistics display in detail view ([1f81689](https://github.com/Starosdev/scrutiny/commit/1f816898da5c628bdd29ef7825e1e38c5e3c9ecc))
* **smart:** add support for ATA Device Statistics (enterprise SSD metrics) ([79d7841](https://github.com/Starosdev/scrutiny/commit/79d784140d7faee7c979047843d4825316bf3603)), closes [#7](https://github.com/Starosdev/scrutiny/issues/7)

## [1.7.2](https://github.com/Starosdev/scrutiny/compare/v1.7.1...v1.7.2) (2026-01-17)

### Bug Fixes

* **mock:** add ZFS pool management methods to MockDeviceRepo ([af2d4bd](https://github.com/Starosdev/scrutiny/commit/af2d4bdb78c9b932a5ba63b73e274212b4386d8e))

## [1.7.1](https://github.com/Starosdev/scrutiny/compare/v1.7.0...v1.7.1) (2026-01-09)

### Bug Fixes

* **deps:** security audit and dependency inventory ([b42e940](https://github.com/Starosdev/scrutiny/commit/b42e94059ac19db10794517ec3bef027558d03e8)), closes [#69](https://github.com/Starosdev/scrutiny/issues/69) [#70](https://github.com/Starosdev/scrutiny/issues/70) [#36](https://github.com/Starosdev/scrutiny/issues/36)

## [1.7.0](https://github.com/Starosdev/scrutiny/compare/v1.6.2...v1.7.0) (2026-01-08)

### Features

* **zfs:** add ZFS pool monitoring support ([6df294a](https://github.com/Starosdev/scrutiny/commit/6df294a8c208f2c2db3a8fecb49d764d47704bbf)), closes [#66](https://github.com/Starosdev/scrutiny/issues/66)

## [1.6.2](https://github.com/Starosdev/scrutiny/compare/v1.6.1...v1.6.2) (2026-01-08)

### Bug Fixes

* **docker:** correct Angular 21 frontend build paths ([18d464b](https://github.com/Starosdev/scrutiny/commit/18d464bad430cb7e3d36f97cd55a7829a79b040f)), closes [#59](https://github.com/Starosdev/scrutiny/issues/59)

## [1.6.1](https://github.com/Starosdev/scrutiny/compare/v1.6.0...v1.6.1) (2026-01-08)

### Bug Fixes

* **ci:** correct frontend tarball path in release workflow ([d46b8f0](https://github.com/Starosdev/scrutiny/commit/d46b8f0696c71ac64af9fd8c3fb7890e44e9db3d)), closes [#59](https://github.com/Starosdev/scrutiny/issues/59)
* **ci:** make frontend coverage upload optional ([5cc5ed1](https://github.com/Starosdev/scrutiny/commit/5cc5ed134044a6dc6d2d2df521e662083bb0b696))

## [1.6.0](https://github.com/Starosdev/scrutiny/compare/v1.5.0...v1.6.0) (2026-01-08)

### Features

* **frontend:** Upgrade Angular 13 to Angular 21 ([d9e4b6a](https://github.com/Starosdev/scrutiny/commit/d9e4b6ad5753a9c1d343a7a44b1bd145dafb92ba)), closes [#9](https://github.com/Starosdev/scrutiny/issues/9)

## [1.5.0](https://github.com/Starosdev/scrutiny/compare/v1.4.1...v1.5.0) (2026-01-08)

### Features

* **notify:** add device label to notification payload ([#48](https://github.com/Starosdev/scrutiny/issues/48)) ([231cc4c](https://github.com/Starosdev/scrutiny/commit/231cc4c2e11d3d04df4aa1a076f9fe839ca5bc56))

### Bug Fixes

* batch of quick wins from GitHub issues ([5eef50e](https://github.com/Starosdev/scrutiny/commit/5eef50e13ba71f400a9846004fedff05e431afed)), closes [#47](https://github.com/Starosdev/scrutiny/issues/47) [#50](https://github.com/Starosdev/scrutiny/issues/50) [#47](https://github.com/Starosdev/scrutiny/issues/47) [#50](https://github.com/Starosdev/scrutiny/issues/50) [#8](https://github.com/Starosdev/scrutiny/issues/8) [#56](https://github.com/Starosdev/scrutiny/issues/56) [#59](https://github.com/Starosdev/scrutiny/issues/59) [#26](https://github.com/Starosdev/scrutiny/issues/26)
* **collector:** populate DeviceType from smartctl info when not set ([6704245](https://github.com/Starosdev/scrutiny/commit/670424567ab2f9d262fc1f3b04fbbf081fc0267e))
* **tests:** add GetString notify.urls mock for notify.Send() ([a2a9f71](https://github.com/Starosdev/scrutiny/commit/a2a9f7109cfe17dce91b036b5b5ad906ea477a55))
* **tests:** add index.html to all web tests for health check ([e1ee2c3](https://github.com/Starosdev/scrutiny/commit/e1ee2c31c5fc341eb23225dde9733a0602c43ffe))
* **tests:** add missing config mock expectations for GORM logging ([5a74c7c](https://github.com/Starosdev/scrutiny/commit/5a74c7c2b050661a45b5aa80779eea36ff318ea2))
* **tests:** add missing config mocks for GORM logging ([2f312cf](https://github.com/Starosdev/scrutiny/commit/2f312cf52ab0f52b5c4f49ad52fd45eefdf53fa8))
* **tests:** add web.metrics.enabled mock ([6418581](https://github.com/Starosdev/scrutiny/commit/6418581a96b26291788beb938f504e1026b93995))
* **tests:** add web.metrics.enabled mock to all test blocks ([3b297c2](https://github.com/Starosdev/scrutiny/commit/3b297c27eff799aae02e8b86d90f9aa418eba8a8))

### Build

* disable VCS stamping in binary builds ([010d287](https://github.com/Starosdev/scrutiny/commit/010d287a108ec4a9069f5d683b33a05f91ec9e81))

## [1.4.3](https://github.com/Starosdev/scrutiny/compare/v1.4.2...v1.4.3) (2026-01-08)

### Build

* disable VCS stamping in binary builds ([010d287](https://github.com/Starosdev/scrutiny/commit/010d287a108ec4a9069f5d683b33a05f91ec9e81))

## [1.4.2](https://github.com/Starosdev/scrutiny/compare/v1.4.1...v1.4.2) (2026-01-08)

### Bug Fixes

* batch of quick wins from GitHub issues ([5eef50e](https://github.com/Starosdev/scrutiny/commit/5eef50e13ba71f400a9846004fedff05e431afed)), closes [#47](https://github.com/Starosdev/scrutiny/issues/47) [#50](https://github.com/Starosdev/scrutiny/issues/50) [#47](https://github.com/Starosdev/scrutiny/issues/47) [#50](https://github.com/Starosdev/scrutiny/issues/50) [#8](https://github.com/Starosdev/scrutiny/issues/8) [#56](https://github.com/Starosdev/scrutiny/issues/56) [#59](https://github.com/Starosdev/scrutiny/issues/59) [#26](https://github.com/Starosdev/scrutiny/issues/26)
* **collector:** populate DeviceType from smartctl info when not set ([6704245](https://github.com/Starosdev/scrutiny/commit/670424567ab2f9d262fc1f3b04fbbf081fc0267e))
* **tests:** add GetString notify.urls mock for notify.Send() ([a2a9f71](https://github.com/Starosdev/scrutiny/commit/a2a9f7109cfe17dce91b036b5b5ad906ea477a55))
* **tests:** add index.html to all web tests for health check ([e1ee2c3](https://github.com/Starosdev/scrutiny/commit/e1ee2c31c5fc341eb23225dde9733a0602c43ffe))
* **tests:** add missing config mock expectations for GORM logging ([5a74c7c](https://github.com/Starosdev/scrutiny/commit/5a74c7c2b050661a45b5aa80779eea36ff318ea2))
* **tests:** add missing config mocks for GORM logging ([2f312cf](https://github.com/Starosdev/scrutiny/commit/2f312cf52ab0f52b5c4f49ad52fd45eefdf53fa8))
* **tests:** add web.metrics.enabled mock ([6418581](https://github.com/Starosdev/scrutiny/commit/6418581a96b26291788beb938f504e1026b93995))
* **tests:** add web.metrics.enabled mock to all test blocks ([3b297c2](https://github.com/Starosdev/scrutiny/commit/3b297c27eff799aae02e8b86d90f9aa418eba8a8))

## [1.4.1](https://github.com/Starosdev/scrutiny/compare/v1.4.0...v1.4.1) (2026-01-08)

### Bug Fixes

* batch of quick wins from GitHub issues ([#60](https://github.com/Starosdev/scrutiny/issues/60)) ([a11d619](https://github.com/Starosdev/scrutiny/commit/a11d619a893458949e67560ff96ee6881dcf13b5))

## [1.3.0](https://github.com/Starosdev/scrutiny/compare/v1.2.0...v1.3.0) (2025-12-20)

### Features

* add device label editing and API timeout configuration ([75050d5](https://github.com/Starosdev/scrutiny/commit/75050d57fa28fe59e833c417671667f43effc472))

## [1.2.0](https://github.com/Starosdev/scrutiny/compare/v1.1.2...v1.2.0) (2025-12-19)

### Features

* **ci:** add SHA256 checksums to GitHub releases ([367a2dc](https://github.com/Starosdev/scrutiny/commit/367a2dc27e95cf17b95f4ea672154c0f8d871cbf)), closes [#28](https://github.com/Starosdev/scrutiny/issues/28)

### Bug Fixes

* Frontend Demo Mode now loads ([#57](https://github.com/Starosdev/scrutiny/issues/57)) ([462a0c3](https://github.com/Starosdev/scrutiny/commit/462a0c362ce5a7b8f5f04a81fe3076fbce4073a8))

## [1.1.2](https://github.com/Starosdev/scrutiny/compare/v1.1.1...v1.1.2) (2025-12-18)

### Refactoring

* **database:** extract hardcoded time ranges to constants ([deb2df0](https://github.com/Starosdev/scrutiny/commit/deb2df0bc718461c5a9826d6b6c1c1307b7122e8)), closes [#49](https://github.com/Starosdev/scrutiny/issues/49)

## [1.1.1](https://github.com/Starosdev/scrutiny/compare/v1.1.0...v1.1.1) (2025-12-09)

### Bug Fixes

* **collector:** handle large LBA values in SMART data parsing ([7f4bceb](https://github.com/Starosdev/scrutiny/commit/7f4bceb85506606d6318253fd406da4b55921383)), closes [#24](https://github.com/Starosdev/scrutiny/issues/24) [AnalogJ/scrutiny#800](https://github.com/AnalogJ/scrutiny/issues/800)
* **collector:** ignore bit 6 in smartctl exit-code during detect ([735fe2e](https://github.com/Starosdev/scrutiny/commit/735fe2e57d9afc9d32832619d6c3c758ec91eb11))
* **collector:** keep existing device type ([b5bb1a2](https://github.com/Starosdev/scrutiny/commit/b5bb1a232a2e38e6bbffb041ffa397b54999fc02))
* **config:** use structured logging for config file messages ([03513b7](https://github.com/Starosdev/scrutiny/commit/03513b742622b77d27cd08b941147eadf35bec91)), closes [#22](https://github.com/Starosdev/scrutiny/issues/22) [AnalogJ/scrutiny#814](https://github.com/AnalogJ/scrutiny/issues/814)
* **database:** use WAL mode to prevent readonly errors in restricted Docker ([1db337d](https://github.com/Starosdev/scrutiny/commit/1db337d872b655e0c68a4a506f9706f0cb7d4a79)), closes [#25](https://github.com/Starosdev/scrutiny/issues/25) [AnalogJ/scrutiny#772](https://github.com/AnalogJ/scrutiny/issues/772)
* **notify:** try to unmarshal notify.urls as JSON array ([9109fb5](https://github.com/Starosdev/scrutiny/commit/9109fb5447080b5faab3377721b830f1e0266500))
* **thresholds:** add observed threshold for attribute 188 with value 0 ([c86ee89](https://github.com/Starosdev/scrutiny/commit/c86ee894468068830fa9e8cf93cde3ef6df1f5d0))
* **thresholds:** mark wear leveling count (attr 177) as critical ([c072119](https://github.com/Starosdev/scrutiny/commit/c0721199b86b02ae398afcc439f4162a760f1d5e)), closes [#21](https://github.com/Starosdev/scrutiny/issues/21) [AnalogJ/scrutiny#818](https://github.com/AnalogJ/scrutiny/issues/818)
* **ui:** display temperature graph times in local timezone ([6123347](https://github.com/Starosdev/scrutiny/commit/6123347165794a5de177248802229c9ea0ea4a9f)), closes [#30](https://github.com/Starosdev/scrutiny/issues/30)

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
