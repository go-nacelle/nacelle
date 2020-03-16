# github.com/go-nacelle/nacelle

This package provides a common [bootstrapper](https://godoc.org/github.com/go-nacelle/nacelle#Bootstrapper) object that initializes and supervises the core framework behaviors.

---

Applications written with nacelle should have a common entrypoint, as follows. The application-specific functionality is passed to a boostrapper on construction as a reference to a function that populates a [process container](https://nacelle.dev/docs/core/process). The `BootAndExit` function initializes and supervises the application, blocks until the application shut down, then calls `os.Exit` with the appropriate status code. A symmetric function called `Boot` will perform the same behavior, but will return the integer status code instead of calling `os.Exit`.

```go
func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
    // register initializer and process instances
}

func main() {
    nacelle.NewBootstrapper("app-name", setup).BootAndExit()
}
```

You can see additional examples of the bootstrapper in the [example repository](https://github.com/go-nacelle/example). Specifically, the main function of the [HTTP API](https://github.com/go-nacelle/example/blob/843979aaa86786784a1ca3646e8d0d1f69e29c65/cmd/http-api/main.go#L17), the [gRPC API](https://github.com/go-nacelle/example/blob/843979aaa86786784a1ca3646e8d0d1f69e29c65/cmd/grpc-api/main.go#L16), and the [worker](https://github.com/go-nacelle/example/blob/843979aaa86786784a1ca3646e8d0d1f69e29c65/cmd/worker/main.go#L17).

The following options can be supplied to the bootstrapper to tune its behavior.

<dl>
  <dt>WithConfigSourcer</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/nacelle#WithConfigSourcer">WithConfigSourcer</a> changes the default source for configuration variables. The default sourcer is the application environment using the name given to the bootstrapper as a prefix.</dd>

  <dt>WithConfigMaskedKeys</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/nacelle#WithConfigMaskedKeys">WithConfigMaskedKeys</a> sets the keys to mask from log messages when loading configuration data. This is used to hide sensitive configuration values.</dd>

  <dt>WithLoggingInitFunc</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/nacelle#WithLoggingInitFunc">WithLoggingInitFunc</a> sets the factory used to create the base logger. This can be set to supply a different log backend.</dd>

  <dt>WithLoggingFields</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/nacelle#WithLoggingFields">WithLoggingFields</a> adds additional fields to every log message. This can be useful to present build information (time, hash, branch), process name, or operating environment.</dd>

  <dt>WithRunnerOptions</dt>
  <dd>
    <a href="https://godoc.org/github.com/go-nacelle/nacelle#WithRunnerOptions">WithRunnerOptions</a> accepts additional options specific to the process runner. The following options can be supplied to tune its behavior.
    <!---->
    <dl>
      <dt>WithHealthCheckBackoff</dt>
      <dd><a href="https://godoc.org/github.com/go-nacelle/process#WithHealthCheckBackoff">WithHealthCheckBackoff</a> sets the <a href="https://github.com/efritz/backoff">backoff</a> instance used to check the health of processes during startup. </dd>
      <!---->
      <dt>WithShutdownTimeout</dt>
      <dd><a href="https://godoc.org/github.com/go-nacelle/process#WithShutdownTimeout">WithShutdownTimeout</a> sets the maximum time that the application can spend shutting down.</dd>
      <!---->
      <dt>WithStartTimeout</dt>
      <dd><a href="https://godoc.org/github.com/go-nacelle/process#WithStartTimeout">WithStartTimeout</a> sets the maximum time that the application can spend in startup.</dd>
    </dl>
  </dd>
</dl>
