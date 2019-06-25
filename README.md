<div align="center"><img width="160" src="https://raw.githubusercontent.com/go-nacelle/nacelle/master/images/nacelle.png" alt="Nacelle logo"></div>

# Nacelle [![GoDoc](https://godoc.org/github.com/go-nacelle/nacelle?status.svg)](https://godoc.org/github.com/go-nacelle/nacelle) [![CircleCI](https://circleci.com/gh/go-nacelle/nacelle.svg?style=svg)](https://circleci.com/gh/go-nacelle/nacelle) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/nacelle/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/nacelle?branch=master)

Microservice framework written in Go.

---

For a full example application, see the [example repository](https://github.com/go-nacelle/example) in this project.

This library is a wrapper and bootstrapping function around the following components. Each exported type, constant, variable, and function are also importable from the main nacelle package (the following repositories should not be imported directly).

- See [config](https://github.com/go-nacelle/config) for documentation on defining configuration structs and loading application configuration from the environment.
- See [log](https://github.com/go-nacelle/log) for documentation on structured loggers.
- See [process](https://github.com/go-nacelle/process) for documentation on defining initializers and processes as well as a description of the application initialization order.
- See [service](https://github.com/go-nacelle/service) for documentation on services and dependency injection. Each initializer and process registered via the bootstrapper will have their dependencies injected prior to their initialization.

This library additional provides a `Bootstrapper`. Applications written in with nacelle will have a common entrypoint, as follows.

```go
func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
    // register initializer and process instances
}

func main() {
    nacelle.NewBootstrapper("app-name", setup).BootAndExit()
}
```

The following options can be supplied to the bootstrapper to tune its behavior.

- **WithConfigSourcer** changes the default source for configuration variables. The default sourcer is the application environment using the name given to the bootstrapper as a prefix.
- **WithConfigMaskedKeys** sets the keys to mask from log messages when loading configuration data. This is used to hide sensitive configuration values.
- **WithLoggingInitFunc** sets the factory used to create the base logger. This can be set to supply a different log backend.
- **WithLoggingFields** adds additional fields to every log message. This can be useful to present build information (time, hash, branch), process name, or operating environment.
- **WithRunnerOptions** sets configuration for the process runner. See the [process](https://github.com/go-nacelle/process) library for additional details.


### Frameworks

The following frameworks are built on top of nacelle to provide rich features
to a single *primary* process.

- [chevron](https://github.com/go-nacelle/chevron) - an HTTP server framework
- [scarf](https://github.com/go-nacelle/scarf) - a gRPC server framework

This project also provides a set of abstract base processes for common process types: an [AWS Lambda event listener](https://github.com/go-nacelle/lambdabase), a [gRPC server](https://github.com/go-nacelle/grpcbase), an [HTTP server](https://github.com/go-nacelle/httpbase), and a [generic worker process](https://github.com/go-nacelle/workerbase).
