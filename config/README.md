# Nacelle Config

Nacelle provides loading application configuration from the environment out of the
box. Parts of an application, organized into services, initializers, and processes,
can associate an instance of their own configuration structs (initially with zero
valued fields) with a unique token with nacelle on startup. These structs are then
hydrated with values from the environment. The hydrated struct can then be requested
with the unique token used during startup.

We use the following configuration struct as an example, defining a unique token in
the idiomatic way.

```go
type (
    Config struct {
        A string        `env:"X"`
        B bool          `env:"Y" default:"true"`
        C string        `env:"Z" required:"true"`
        D []string      `env:"W" default:"[\"foo\", \"bar\", \"baz\"]"`
    }

    configToken string
)

var ConfigToken = configToken("my-config-name")
```

When loading values from the environment, a missing value (or empty string) will use
the default value, if provided. If no value is set for a required configuration value,
a fatal error will occur. String values will retrieve the environment value unaltered.
All other field types will attempt to deserialize the environment value as JSON.

An zero-value instance of this config can be registered with the singleton token on
startup (see `main.go` in any of the provided sample applications).

```go
func setupConfigs(config nacelle.Config) error {
    config.MustRegister(ConfigToken, &Config{})
    // ...
    return nil
}
```

Then, an initializer or a process that requires these config values can retrieve them in
its `Init` method.

```go
func (p *Process) Init(config nacelle.Config) error {
    c := &Config{}
    if err := config.Fetch(ConfigToken, c); err != nil {
        // ...
    }

    // c is hydrated
    // ...
}
```

## PostLoading Configuration Structs

After hydration, the `PostLoad` method will be invoked on all registered configuration
structs (where such a method exists). This allows additional non-type validation to
occur, and to create any types which are not directly/easily encodable as JSON.

```go
func (c *Config) PostLoad() error {
    if c.Field != "foo" && c.Field != "bar" {
        return fmt.Errorf("field value must be foo or bar")
    }

    return nil
}
```

## Config Tags

In some circumstances, it may be necessary to dynamically alter the tags on a configuration
struct. This has become an issue in two circumstances so far. First, two instances of the
same configuration struct may need to be registered but must be configured separately
(i.e. they need to look at distinct environment variables). This is a particular problem
when running two HTTP servers with the same base, for example. Second, the default value of
a field may need to be altered. This issue can also arise when two instances of the same
configuration struct are registered but shouldn't get clashing defaults (e.g. a default
listening port).

Nacelle provides two tag modifiers which can be applied at configuration registration time.
In the following, the instance registered is modified so that the environment variables used
to hydrate the object are `Q_X`, `Q_Y`, `Q_Z`, `Q_W`, instead of `X`, `Y`, `Z`, and `W` the
default value of the struct field `B` (loaded through the environment variable `Q_Y`) is false
instead of true.

```go
func setupConfigs(config nacelle.Config) error {
    config.MustRegister(
        ConfigToken,
        &Config{},
        nacelle.NewEnvTagPrefixer("Q"),
        nacelle.NewDefaultTagSetter("B", "false"),
    )

    // ...
    return nil
}
```

If other dynamic modifications of a configuration struct is necessary, simply implement the
`TagModifier` interface and use it in the same way.
