# Changelog

## [0.7.1](https://github.com/monkescience/phasor/compare/0.7.0...0.7.1) (2026-01-10)


### Bug Fixes

* use FQDN in analysis template health check URLs ([398ca36](https://github.com/monkescience/phasor/commit/398ca36e378f5a9990c648da3792454e5399e707))

## [0.7.0](https://github.com/monkescience/phasor/compare/0.6.0...0.7.0) (2026-01-10)


### Features

* add HTTPRoute template for Gateway API traffic routing in backend rollout ([15e6723](https://github.com/monkescience/phasor/commit/15e67234130325a37b6c2b86ad5daa579dcc2200))


### Bug Fixes

* **deps:** update github.com/monkescience/vital digest to 8503480 ([#15](https://github.com/monkescience/phasor/issues/15)) ([65a86fd](https://github.com/monkescience/phasor/commit/65a86fda7a84bac48abaa63aa3b5f5be54084db6))

## [0.6.0](https://github.com/monkescience/phasor/compare/0.5.0...0.6.0) (2026-01-05)


### Features

* add Gateway API traffic routing to backend rollout and configure analysis templates with health check count ([d49721d](https://github.com/monkescience/phasor/commit/d49721dcd425a7471f68f8697b60cd55c81924a1))

## [0.5.0](https://github.com/monkescience/phasor/compare/0.4.0...0.5.0) (2026-01-05)


### Features

* add canary and blue-green rollout support for backend and frontend services, including health checks and associated services ([4bc3d9c](https://github.com/monkescience/phasor/commit/4bc3d9c72c013abaaeb733da560c7c7e876601f6))

## [0.4.0](https://github.com/monkescience/phasor/compare/0.3.0...0.4.0) (2025-12-23)


### Features

* add config checksum annotation to frontend and backend deployments ([75e47c4](https://github.com/monkescience/phasor/commit/75e47c416d0c863b795d687980f7d4ee1f117740))
* add Renovate bot with workflow and schema integration ([69f2b86](https://github.com/monkescience/phasor/commit/69f2b86768924224a7c5638248030c83df9b8d69))

## [0.3.0](https://github.com/monkescience/phasor/compare/0.2.0...0.3.0) (2025-12-18)


### Features

* add backend health checker and integrate it into health endpoint ([452c1bd](https://github.com/monkescience/phasor/commit/452c1bd32a652ca04bf7f1cfa1b7f9a3bbc401f1))
* consolidate configuration handling by replacing inline YAML with `toYaml` and commenting out unused config options ([d9e32f6](https://github.com/monkescience/phasor/commit/d9e32f6093c15bb5a84ad080be8909c782205d3e))

## [0.2.0](https://github.com/monkescience/phasor/compare/0.1.0...0.2.0) (2025-12-17)


### Features

* add initial project structure ([16b94ae](https://github.com/monkescience/phasor/commit/16b94ae216080f98970358d12e184917b2a84623))
