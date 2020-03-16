# Getting Started

---

### What is Nacelle?

Nacelle is a Golang framework for setting up and monitoring the behavioral components of your application. It is intended to provide a common and convention-guided way to bootstrap an application or set of microsevices. Broadly, boostrapping your application with Nacelle gives you:

- a convention for code organization via processes and initializers
- a convention for declaring, reading, and validating configuration values from the runtime environment
- a convention for registering, declaring, and injecting structure interface dependencies
- a convention for structured logging

Concretely, bootstrapping your application with Nacelle gives you:

- an expanding set of libraries that use the conventions outlined above
- a process initializer that reads and validates declared component configuration
- a process initializer that injects declared component dependencies
- a process runner that invokes each process in a dedicated goroutine
- a process monitor that watches for error and cleanly shutdown your application

### Installation

```bash
go get -u github.com/go-nacelle/nacelle
```

### Fifteen-Minute Walkthrough

The core ideas of Nacelle revolve around *processes* and *initializers*. A process is a behavioral component of your application which does some work over the process lifetime. An initializer is a component of your application which does some work at application startup.

An application can be composed of one or more processes, which are commonly long-running such as a server or a worker that accepts work from an external system. An application may also have zero or more initializers, which generally create or initialize a resource (such as a database connection) used by an application process.

#### Setup

Applications using Nacelle to bootstrap will have the following minimal `main` function. This hands control off to the bootrapper, which will invoke the registered `setup` function in order to populate the process and service containers. The bootrapper will then initialize each process and monitor it for the lifetime of the application.

```go
package main

import "github.com/go-nacelle/nacelle"

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
    // Register processes and initializers here
    return nil
}

func main() {
    nacelle.NewBootstrapper("app-name", setup).BootAndExit()
}
```


If you were to run this application, you would see Nacelle trying to initialize each registered initializer (of which there are none), and initialize and start each registered process (of which there are none).

#### Registering a Process

Let's create simple HTTP server process that responds with a 200 OK/Hello World response for each request.

Each process is initialized by calling its `Init` function with a configuration container (more on this in the next section). On initialization success, the `Start` method is invoked in a dedicated go-routine. This method is expected to block for long-running processes such as servers. On application shutdown, the `Stop` method is invoked which should unblock any active work being done in the `Start` method.

Here, the `Init` method creates a server and configures its handler. The `Start` creates a TCP listener and starts serving HTTP traffic. The `Stop` method signals for the server to stop accepting new connections and shutdown.

```go
import (
    "context"
    "fmt"
    "net"
    "net/http"
)

type server struct {
    server *http.Server
    port   int
}

func (s *server) Init(config nacelle.Config) error {
    s.server = &http.Server{}
    s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hello, World!\n"))
    })
    return nil
}

func (s *server) Start() error {
    addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", s.port))
    if err != nil {
        return err
    }

    listener, err := net.ListenTCP("tcp", addr)
    if err != nil {
        return err
    }

    defer listener.Close()
    defer s.server.Close()

    // Run server, block until shutdown (do not return ErrServerClosed)
    if err := s.server.Serve(listener); err != http.ErrServerClosed {
        return err
    }

    return nil
}

func (s *server) Stop() error {
    s.server.Shutdown(context.Background())
    return nil
}
```

Now, we modify the setup function to register this server process.

```go
func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
    processes.RegisterProcess(&server{port: 5000}, nacelle.WithProcessName("hw-server"))
    return nil
}
```

If you were to run this application, you would see Nacelle initialize and start the *hw-server* process. Curl-ing any path at `http://localhost:5000` will return the same 200-level response.

#### Process Configuration

The application above creates a process with hard-coded port of 5000. This is problematic in the case you need to change the port when running on a different environment, or run two servers on the same host.

We can instead accept this value from the environment (environment variable, file, configmap, etc) at runtime so that no code change is required to configure this value.

We declare the configuration values accepted by the server process with a configuration struct. Here, we tag the port field with `env` which indicates the environment variable that should be read to populate this field.

In the `Init` method of the process, we populate an instance of this struct with values and pull the required values inot the process struct for later use.

```go
type serverConfig struct {
    Port int `env:"port" default:"5000"`
}

func (s *server) Init(config nacelle.Config) error {
    serverConfig := &serverConfig{}
    if err := config.Load(serverConfig); err != nil {
        return err
    }

    s.port = serverConfig.Port
    // ...
}

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
    processes.RegisterProcess(&server{}, nacelle.WithProcessName("hw-server"))
    return nil
}
```

Running the application the same way will show the same behavior. Running the application with `PORT=3000` will cause the application to listen to the non-default port 3000. Running the application with a non-integer port value will cause the application to fail on startup with an error message.

#### Process Dependencies

Let's modify the server to return a distinct response for each request. Instead of a canned message, we will print their request count: *Hello #1* for the first request, *Hello #2!* for the second, and so on. We'll store this data in Redis, and atomically increment a request counter each time the handler is invoked.

This creates a dependency for a Redis client in the server process.

```go
import "github.com/go-redis/redis/v7"

type server struct {
    client *redis.Client
    // ...
}
```

We'll change the HTTP handler implementation as follows.

```go
s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    count, err := s.client.Incr("sample").Result()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hello, #%d!\n", count)))
})
```

Finally, we need to supply a concrete Redis client to the process on startup. We'll do this *for now* by initializing the client in the setup function and passing it to the process on creation.

```go
func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	processes.RegisterProcess(&server{client: client}, nacelle.WithProcessName("hw-server"))
	return nil
}
```

The application should now produce HTTP responses with an increasing count in the body. If Redis is not running or accessible on your host, then the server should respond with an internal server error response.

#### Dependency Injection

The last change above creates a few issues, namely:

1. We do not check if the client can reach a remote server.
2. We do not pull the Redis address from the environment. We hard-code the address, which has the same issues as hard-coding the port above. Additionally, the bootstrapper hasn't created the configuration object yet, so we *couldn't* read from the environment at this point in the application lifecycle anyway.
3. We are supplying dependencies manually. Right now, the server process has the dependency on the client, but in a larger application this may be a dependency-of-a-dependency, which requires threading dependencies transitively through your application graph.

We can handle all of these issues by writing an initializer. An initializer is like a process, but only has an `Init` method, called in the same fashion. The following `Init` method reads the Redis address from the environment, constructs a client, pings the remote server, and adds the client to the service container with an application-distinct name.

```go
type redisInitializer struct {
	Services nacelle.ServiceContainer `service:"services"`
}

type redisConfig struct {
	Addr string `env:"redis_addr" default:"localhost:6379"`
}

func (i *redisInitializer) Init(config nacelle.Config) error {
	redisConfig := &redisConfig{}
	if err := config.Load(redisConfig); err != nil {
		return err
	}

	client := redis.NewClient(&redis.Options{Addr: redisConfig.Addr})
	if _, err := client.Ping().Result(); err != nil {
		return err
	}

	return i.Services.Set("redis", client)
}
```

THe redis initializer has field with a `service` tag. This informs the bootstrapper to set the value of that field with the registered service with the same name. The *services* and *logger* services are available to all applications at startup. Similarly, we change the client field of the server process as follows.

```go
type server struct {
	server *http.Server
	port   int
	Client *redis.Client `service:"redis"`
}
```

Note that each injected field must be exported for the bootstrapper to access it. This canged the casing of the field, and will need to be changed within the HTTP handler as well.

Now, we can replace the ad-hoc client creation with the registration of the initializer that replaces it.

```go
func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	processes.RegisterInitializer(&redisInitializer{}, nacelle.WithInitializerName("redis"))
	processes.RegisterProcess(&server{}, nacelle.WithProcessName("hw-server"))
	return nil
}
```
