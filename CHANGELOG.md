# Changelog

This changelog tracks updates to this repository as well as [go-nacelle/config](https://github.com/go-nacelle/config), [go-nacelle/log](https://github.com/go-nacelle/log), [go-nacelle/process](https://github.com/go-nacelle/process), and [go-nacelle/service](https://github.com/go-nacelle/service).

## v1.1.4

- Add `HasReason` to `Health` struct (from [go-nacelle/process](https://github.com/go-nacelle/process)).

## v1.1.3

- Add `NewFlagSourcer`, `NewFlagTagPrefixer`, and `NewFlagTagSetter` constructors (from [go-nacelle/config](https://github.com/go-nacelle/config)).

## v1.1.2

#### Definition updates

- Add `FileSystem` type (from [go-nacelle/config](https://github.com/go-nacelle/config)).
- Add `With{File,Directory,Glob}SourcerFS` constructors (from [go-nacelle/config](https://github.com/go-nacelle/config)).

## v1.1.1

#### Update [go-nacelle/log](https://github.com/go-nacelle/log) to [v1.1.1](https://github.com/go-nacelle/log/releases/tag/v1.1.1)

- Rewrite base logger to remove dependency on [gomol](https://github.com/aphistic/gomol).
- The config variable `log_field_blacklist` now accepts a JSON-encoded array instead of a comma-separated list.
- Add `WithIndirectCaller` to `Logger` interface (closes [log#1](https://github.com/go-nacelle/log/issues/1)).
- Fix malformed templates in console output.

## v1.1.0

#### Update [go-nacelle/config](https://github.com/go-nacelle/config) to [v1.1.0](https://github.com/go-nacelle/config/releases/tag/v1.1.0)

- Add `Filesystem` interface (closes [config#1](https://github.com/go-nacelle/config/issues/1)). An instances of this interface can now be (optionally) supplied to the following constructors:
  - `NewFileSourcer`
  - `NewOptionalFileSourcer`
  - `NewDirectorySourcer`
  - `NewOptionalDirectorySourcer`
  - `NewGlobSourcer`

## v1.0.2

#### Definition updates

- Add `Finalizer` type (from [go-nacelle/process](https://github.com/go-nacelle/process)).

## v1.0.1

#### Update [go-nacelle/log](https://github.com/go-nacelle/log) to [v1.0.1](https://github.com/go-nacelle/log/releases/tag/v1.0.1)

- Regenerate mocks.

#### Definition updates

- Replace `MockLogger` interface and `NewMockLogger` constructor with imports from [go-nacelle/log](https://github.com/go-nacelle/log) (previously imported from [go-nacelle/config](https://github.com/go-nacelle/config)).
- Add `NewTestEnvSourcer` constructor (from [go-nacelle/config](https://github.com/go-nacelle/config)).

## v1.0.0

Initial stable release.
