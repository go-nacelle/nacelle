<p align="center">
    <a href="https://godoc.org/github.com/efritz/nacelle"><img src="https://godoc.org/github.com/efritz/nacelle?status.svg" alt="GoDoc"></a>
    <a href="http://travis-ci.org/efritz/nacelle"><img src="https://secure.travis-ci.org/efritz/nacelle.png" alt="Build Status"></a>
    <a href="https://codeclimate.com/github/efritz/nacelle/maintainability"><img src="https://api.codeclimate.com/v1/badges/8118b324f3b7ac9b442a/maintainability" alt="Maintainability"></a>
    <a href="https://codeclimate.com/github/efritz/nacelle/test_coverage"><img src="https://api.codeclimate.com/v1/badges/8118b324f3b7ac9b442a/test_coverage" alt="Test Coverage"></a>
</p>

<p align="center">
    <img width="200" src="https://github.com/efritz/nacelle/blob/master/images/nacelle.png" alt="Nacelle logo">
</p>

<h2 align="center">Nacelle</h2>

Nacelle is a configuration and dependency injection framework for services written
in Go. For example usage, see the examples directory.

---

## Concepts

The following sections outline general concepts required for using Nacelle.

### Process

A **process** is an interface that has an `Init`, `Start`, and `Stop` method.
The `Init` method is called by Nacelle after the configuration has been loaded
from the environment and passes the hydrated config as a parameter. The `Start`
method should begin the process - some examples include serving content over some
port or reading work from a distributed queue. This method should block until the
process has ended either gracefully or due to a fatal error. The `Stop` method
should inform the process that it should begin shutting down.

The `process` package provides abstract processes for an HTTP server, a gRPC
server, and a worker which does some kind of work on a timed interval.

An **initializer** is similar to a process, but only has an `Init` method. These
initializers generally prep some data in a package or make a connection to an
external service required by a process before they are started.

A program can consist of a number of processes all executing concurrently. They
are registered and supervised through a *ProcessRunner* instance. The order that
initializers are registered will be the same order in which they are executed. A
process may be registered with a priority such that the `Init` and `Start` methods
of priority *n* are executed before looking at processes with priority *n+1*.

### Config

Nacelle provides a **Config** object - the default implementation of which reads its
values from the environment. This config object consists of a series of zero-value
initialized structs whose fields are tagged with environment value names. Once the
config is loaded from an external source, the *hydrated* structs registered to the
config object can be re-retrieved via a registration token.

We use the the following configuration object as an example.

```go
type (
    Config struct {
        A int    `env:"A"`
        B string `env:"B"`
        C bool   `env:"C"`
    }

    configToken string
)

var ConfigToken = configToken("my-config-name")
```

An zero-value instance of this config can be registered with the singleton token on
startup (see `main.go` in the example repo).

```go
var configs = map[interface{}]interface{}{
    ConfigToken: &Config{},
    // ...
}
```

Then, a process that requires these config values can retrieve them in its `Init`
method.

```go
func (p *Process) Init(config nacelle.Config) error {
    c := &Config{}
    if err := config.Fetch(ConfigToken, c); err != nil {
        // ...
    }

    c.A // use
    c.B // these
    c.C // values

    // ...
}
```

Config value fields can also be tagged with `required:"true"` if the value
must be supplied and `default:"val"` if a default value should be used when
the associated environment value is not set.

### Services

A **service** is a dependency for an initializer or a process. This can be
something like a database handle, a shared cache object, or a logger tagged
with specific attributes. Services are generally attached to a container by
an initializer and later injected into a process instance.

For example, we can create an initializer as follows to create a shared cache
object. Initializers can access the service container object created by the
Nacelle boot process either by injection (tagging a field with `service:""`)
or by wrapping an initializer function with the `WrapServiceInitializerFunc`
method. This example uses the former.

```go
type CacheInitializer struct{
    Container *nacelle.ServiceContainer `service:"container"`
}

func (i *CacheInitializer) Init(config nacelle.Config) error {
    return i.Container.Set("cache", NewCache())
}
```

Later, once this process begins its initialization phase, the fields tagged
with service will be set with the matching values in the container.

```go
type Process struct {
    Cache SharedCache `service:"cache"`
}
```

Two or more processes can share this value so that the same values are cached
across services.

## License

Copyright (c) 2017 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
