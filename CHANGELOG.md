# Changelog

## [Unreleased]

### Added

- Added config object to service container with the key `config`. This will cause a registration error in applications that previous used this key for a custom service. [#7](https://github.com/go-nacelle/nacelle/pull/7)
- Added `WithContextFilter`. [#10](https://github.com/go-nacelle/nacelle/pull/10)
- Imported `WithInitializerContextFilter`, `WithProcessContextFilter`, `WithInitializerPriority`, and `WithProcessPriority` from [go-nacelle/process](https://github.com/go-nacelle/process). [#11](https://github.com/go-nacelle/nacelle/pull/11)
- Registered configurations are now dumped when application is invoked with `--help` flag. [#12](https://github.com/go-nacelle/nacelle/pull/12)
- Added `Configurable` and `ConfigurationRegistry` interfaces. [#13](https://github.com/go-nacelle/nacelle/pull/13)

### Changed

- Changed signature of `ServiceInitializerFunc`. [#11](https://github.com/go-nacelle/nacelle/pull/11)
- Changed signature of `AppInitFunc`, `ServiceInitializerFunc`, and `WrapServiceInitializerFunc`. [#17](https://github.com/go-nacelle/nacelle/pull/17)

### Removed

- Removed mocks package. [#14](https://github.com/go-nacelle/nacelle/pull/14)
- Removed `Overlay` import. [#17](https://github.com/go-nacelle/nacelle/pull/17)

## [v1.2.0] - 2020-10-04

### Added

- Imported `Overlay` struct from [go-nacelle/service](https://github.com/go-nacelle/service). [abb708b](https://github.com/go-nacelle/nacelle/commit/abb708b780370823c35ce654c8feb79611a7f29e)
- Imported `WithInitializerLogFields` and `WithProcessLogFields` options from [go-nacelle/process](https://github.com/go-nacelle/process). [abb708b](https://github.com/go-nacelle/nacelle/commit/abb708b780370823c35ce654c8feb79611a7f29e)
- [go-nacelle/process@v1.0.1] -> [go-nacelle/process@v1.1.0]
  - Added `WithInitializerLogFields` and `WithProcessLogFields`. [#2](https://github.com/go-nacelle/process/pull/2)
- [go-nacelle/service@v1.0.0] -> [go-nacelle/service@v1.0.2]
  - Added overlay container. [#1](https://github.com/go-nacelle/service/pull/1)

### Changed

- Replaced `WithHealthCheckBackoff` options with `WithHealthCheckInterval`. [c6b9130](https://github.com/go-nacelle/nacelle/commit/c6b91304d1e7c258889109e4ed763dff04764fb6)

### Removed

- Removed dependency on [aphistic/sweet](https://github.com/aphistic/sweet) by rewriting tests to use [testify](https://github.com/stretchr/testify). [#2](https://github.com/go-nacelle/nacelle/pull/2)
- [go-nacelle/config@v1.2.0] -> [go-nacelle/config@v1.2.1]
  - Removed dependency on [aphistic/sweet](https://github.com/aphistic/sweet) by rewriting tests to use [testify](https://github.com/stretchr/testify). [#5](https://github.com/go-nacelle/config/pull/5)
- [go-nacelle/log@v1.1.1] -> [go-nacelle/log@v1.1.2]
  - Removed dependency on [aphistic/sweet](https://github.com/aphistic/sweet) by rewriting tests to use [testify](https://github.com/stretchr/testify). [#3](https://github.com/go-nacelle/log/pull/3)
- [go-nacelle/process@v1.0.1] -> [go-nacelle/process@v1.1.0]
  - Removed dependency on [efritz/backoff](https://github.com/efritz/backoff). [bd4092d](https://github.com/go-nacelle/process/commit/bd4092d39078bba1e9cdce0e3187560fbfb172bc)
  - Removed dependency on [efritz/watchdog](https://github.com/efritz/watchdog). [4121898](https://github.com/go-nacelle/process/commit/41218985f4849dc0e89c26e0fe2b274a31af61fb)
  - Removed dependency on [aphistic/sweet](https://github.com/aphistic/sweet) by rewriting tests to use [testify](https://github.com/stretchr/testify). [#3](https://github.com/go-nacelle/process/pull/3)
- [go-nacelle/service@v1.0.0] -> [go-nacelle/service@v1.0.2]
  - Removed dependency on [aphistic/sweet](https://github.com/aphistic/sweet) by rewriting tests to use [testify](https://github.com/stretchr/testify). [#2](https://github.com/go-nacelle/service/pull/2)

## [v1.1.4] - 2020-06-08

### Added

- [go-nacelle/process@v1.0.0] -> [go-nacelle/process@v1.0.1]
  - Added `HasReason` to `Health`. [#1](https://github.com/go-nacelle/process/pull/1)

## [v1.1.3] - 2020-04-02

### Added

- Imported `NewFlagSourcer`, `NewFlagTagPrefixer`, and `NewFlagTagSetter` from [go-nacelle/config](https://github.com/go-nacelle/config). [bc39689](https://github.com/go-nacelle/nacelle/commit/bc396890f965e3b359cb707c0ff2840d058a2504)
- [go-nacelle/config@v1.1.0] -> [go-nacelle/config@v1.2.0]
  - Added `FlagSourcer` that reads configuration values from the command line. [#3](https://github.com/go-nacelle/config/pull/3)
  - Added `Init` method to `Config` and `Sourcer`. [#4](https://github.com/go-nacelle/config/pull/4)

## [v1.1.2] - 2020-01-02

### Added

- Imported `WithDirectorySourcerFS`, `WithFileSourcerFS`, and `WithGlobSourcerFS` from [go-nacelle/config](https://github.com/go-nacelle/config). [4575828](https://github.com/go-nacelle/nacelle/commit/4575828f9c7dbb2821dc585faf369432dbed4086)

## [v1.1.1] - 2019-11-19

### Fixed

- [go-nacelle/log@v1.0.1] -> [go-nacelle/log@v1.1.1]
  - Fixed bad console output. [db6e246](https://github.com/go-nacelle/log/commit/db6e24657334615a099e39bae0359179778016e4), [45875f1](https://github.com/go-nacelle/log/commit/45875f173a0db48fc3f615d96a4f83e015cdf130)

### Added

- [go-nacelle/log@v1.0.1] -> [go-nacelle/log@v1.1.1]
  - Added `WithIndirectCaller` to control the number of stack frames to omit. [#2](https://github.com/go-nacelle/log/pull/2)

### Removed

- [go-nacelle/log@v1.0.1] -> [go-nacelle/log@v1.1.1]
  - Removed dependency on [aphistic/gomol](https://github.com/aphistic/gomol) by rewriting base logger internally. [4e537aa](https://github.com/go-nacelle/log/commit/4e537aa0e5a08638bfb45f5153e8deccf6e1d00d)

### Changed

- [go-nacelle/log@v1.0.1] -> [go-nacelle/log@v1.1.1]
  - Changed log field blacklist from a comma-separated list to a json-encoded array. [96b9d53](https://github.com/go-nacelle/log/commit/96b9d53baff25f7c0436799f520c3d4a5970941e)

## [v1.1.0] - 2019-09-05

### Added

- [go-nacelle/config@v1.0.0] -> [go-nacelle/config@v1.1.0]
  - Added options to supply a filesystem adapter to glob, file, and directory sourcers. [#2](https://github.com/go-nacelle/config/pull/2)

## [v1.0.2] - 2019-06-24

### Added

- Imported `Finalizer` from [go-nacelle/process](https://github.com/go-nacelle/process). [000fcd8](https://github.com/go-nacelle/nacelle/commit/000fcd833621e1ef0a2bdec44afbe8cd15a3644d)

## [v1.0.1] - 2019-06-20

### Added

- Imported `NewTestEnvSourcer` from [go-nacelle/config](https://github.com/go-nacelle/config). [c577ab0](https://github.com/go-nacelle/nacelle/commit/c577ab075bede49ea8151e2f945472cb6228bfd0)
- [go-nacelle/log@v1.0.0] -> [go-nacelle/log@v1.0.1]
  - Added mocks package. [d24aad2](https://github.com/go-nacelle/log/commit/d24aad20df4c5b24dbdff3860c348af82abed169)

### Changed

- Import logger mocks from [go-nacelle/log](https://github.com/go-nacelle/log). [b3a0df4](https://github.com/go-nacelle/nacelle/commit/b3a0df415b7bb1d2261bed9b57f423cca45ad455)

## [v1.0.0] - 2019-06-17

### Changed

- Migrated from [efritz/nacelle](https://github.com/efritz/nacelle).

[Unreleased]: https://github.com/go-nacelle/nacelle/compare/v1.2.0...HEAD
[go-nacelle/config@v1.0.0]: https://github.com/go-nacelle/config/releases/tag/v1.0.0
[go-nacelle/config@v1.1.0]: https://github.com/go-nacelle/config/releases/tag/v1.1.0
[go-nacelle/config@v1.2.0]: https://github.com/go-nacelle/config/releases/tag/v1.2.0
[go-nacelle/config@v1.2.1]: https://github.com/go-nacelle/config/releases/tag/v1.2.1
[go-nacelle/log@v1.0.0]: https://github.com/go-nacelle/log/releases/tag/v1.0.0
[go-nacelle/log@v1.0.1]: https://github.com/go-nacelle/log/releases/tag/v1.0.1
[go-nacelle/log@v1.1.1]: https://github.com/go-nacelle/log/releases/tag/v1.1.1
[go-nacelle/log@v1.1.2]: https://github.com/go-nacelle/log/releases/tag/v1.1.2
[go-nacelle/process@v1.0.0]: https://github.com/go-nacelle/process/releases/tag/v1.0.0
[go-nacelle/process@v1.0.1]: https://github.com/go-nacelle/process/releases/tag/v1.0.1
[go-nacelle/process@v1.1.0]: https://github.com/go-nacelle/process/releases/tag/v1.1.0
[go-nacelle/service@v1.0.0]: https://github.com/go-nacelle/service/releases/tag/v1.0.0
[go-nacelle/service@v1.0.2]: https://github.com/go-nacelle/service/releases/tag/v1.0.2
[v1.0.0]: https://github.com/go-nacelle/nacelle/releases/tag/v1.0.0
[v1.0.1]: https://github.com/go-nacelle/nacelle/compare/v1.0.0...v1.0.1
[v1.0.2]: https://github.com/go-nacelle/nacelle/compare/v1.0.1...v1.0.2
[v1.1.0]: https://github.com/go-nacelle/nacelle/compare/v1.0.2...v1.1.0
[v1.1.1]: https://github.com/go-nacelle/nacelle/compare/v1.1.0...v1.1.1
[v1.1.2]: https://github.com/go-nacelle/nacelle/compare/v1.1.1...v1.1.2
[v1.1.3]: https://github.com/go-nacelle/nacelle/compare/v1.1.2...v1.1.3
[v1.1.4]: https://github.com/go-nacelle/nacelle/compare/v1.1.3...v1.1.4
[v1.2.0]: https://github.com/go-nacelle/nacelle/compare/v1.1.4...v1.2.0
