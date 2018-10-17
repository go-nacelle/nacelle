# Nacelle Processes

Nacelle applications are organized into three distinct categories:

- A **service** is a shared API. Services exist in support of a *process*.
Services do not respond to user requests, are not externally accessible, and
are generally *inactive* - they do not process something continuously and
only have observable behavior when they are

- An **initializer** is something that runs once on application startup.
Initializers generally instantiate a service and insert it into a service
container for use by other services, initializers, and processes.

- A **process** is the meat of the application. All processes do something
*actively* - this may be listening for incoming socket connections or reading
messages from a remote work queue and processing them. A process should generally
do a single thing. Multiple processes can communicate directly or through a shared
service.

Several common low-level processes (an HTTP server, a gRPC server, and a generic
worker process) implementations are available in the
[base package](https://github.com/efritz/nacelle/tree/master/base).

## Setup

A nacelle bootstrapper instance is created with a single function pointer which is called
with a reference to a process container and a service container and should register the
application initializers and processes. The structure of an application will look similar
to the following.

```go
func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
    // ...
    return nil
}

func main() {
    nacelle.NewBootstrapper("app-name", setup).BootAndExit()
}
```

In the following examples, we assume this layout and will simply denote in which functions
additional code must be added.

For usage of the configuration object, see the
[config](https://github.com/efritz/nacelle/tree/master/config) package.

An initializer can be registered to a process container with the `RegisterInitializer` function.
Additional options can be provided along with the initializer instance (e.g. provide a name for
the initializer used in logs, set a maximum execution duration, etc). Initializers are run by
nacelle in the order in which they are registered.

A process can be registered to a process containerw ith the `RegisterProcess` function. Additional
options can be provided here as well (e.g. provide a name, set the process priority, whether or not
it should be allowed to exit, etc). Processes are run by their registered priority, then by their
registration order. The default process priority is zero.

## Example

We give a small inline example here. For additional, fully working examples, see the
[examples](https://github.com/efritz/nacelle/tree/master/examples) directory.

### Service Example

First, we define a service that implements a simple cache backed by Redis. At this point everything
is plain Go - there is no nacelle secret sauce.

```go
import (
	"fmt"
	"io"
	"net/http"

	"github.com/efritz/deepjoy"
	"github.com/efritz/nacelle"
	"github.com/garyburd/redigo/redis"
)

type (
    Cache interface {
        Get(string) (string, error)
        Set(string, string) error
    }

    RedisCache struct {
        client deepjoy.Client
    }
)

func NewRedisCache(addr string) Cache{
    return &RedisCache{
        client: deepjoy.NewClient(addr),
    }
}

func (c *RedisCache) Get(key string) (string, error) {
    return redis.String(c.client.Do("GET", key))
}

func (c *RedisCache) Set(key, value string) error {
    _, err := c.client.Do("SET", key, value)
    return err
}
```

### Initializer Example

Next, we create an initializer that adds an instance of the Redis cache defined into a
shared service container. This will allow multiple processes to use the same cache
instance. We define a configuration object that allows the remote Redis address to be
set via an environment variable.

```go
type (
    CacheInitializer struct {
        Services nacelle.ServiceContainer `service:"container"`
    }

    RedisConfig struct {
        CacheAddr string `env:"cache_addr" required:"true"`
    }
)

func NewCacheInitializer() nacelle.Initializer {
    return &CacheInitializer{}
}

func (m *CacheInitializer) Init(config nacelle.Config) error {
    redisConfig := &RedisConfig{}
    if err := config.Load(redisConfig); err != nil {
        return err
    }

    return m.Services.Set("cache", NewRedisCache(redisConfig.CacheAddr))
}
```

Add the following to the `setup` method to register the initializer.

```go
processes.RegisterInitializer(NewCacheInitializer())
```

### Process Example

Next, we create an HTTP process. This process is implemented in a very simple way to highlight
the process structure itself and isn't the recommended way to lay out a server process. For a
more correct way, see the base HTTP server
[implementation](https://github.com/efritz/nacelle/tree/master/base/http).

This server defines its own config and requests an instance of the registered cache service from
the nacelle on init. By the time the server's `Init` method is called, the config will have been
loaded from the environment services will have been injected into the struct.

```go
type (
    Server struct {
        Logger nacelle.Logger `service:"logger"`
        Cache  Cache          `service:"cache"`
        srv    *http.Server
    }

    ServerConfig struct {
        Port int `env:"port" default:"8080"`
    }
)

func NewServer() nacelle.Process {
    return &Server{}
}

func (s *Server) Init(config nacelle.Config) (err error) {
    serverConfig := &ServerConfig{}
    if err := config.Load(serverConfig); err != nil {
        return err
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        cacheKey := r.RequestURI

        // Try to read payload from cache
        payload, err := s.Cache.Get(cacheKey)
        if err != nil {
            s.Logger.Error("cache read failure (%s)", err.Error())
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        if payload == "" {
            // Not cached, do the work
            payload = generatePayload(r)

            // Write payload back to cache
            if err := s.Cache.Set(cacheKey, payload); err != nil {
                s.Logger.Error("cache write failure (%s)", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
        }

        // Write to client
        io.WriteString(w, payload)
    })

    s.srv = &http.Server{
        Addr: fmt.Sprintf(":%d", serverConfig.Port),
    }

    return nil
}

func (s *Server) Start() error {
    s.Logger.Info("Server is now accepting clients")
    return s.srv.ListenAndServe()
}

func (s *Server) Stop() (err error) {
    return s.srv.Close()
}
```

Add the following to the `setup` method to register the server process itself.

```go
processes.RegisterProcess(NewServer())
```
