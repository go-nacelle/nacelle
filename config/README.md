# Nacelle Config

Out of the box, nacelle provides loading application configuration from a
configurable source and injecting the values into an (initially zero-valued)
struct requested by the application.

We use the following configuration struct as an example.

```go
type Config struct {
    A string        `env:"X"`
    B bool          `env:"Y" default:"true"`
    C string        `env:"Z" required:"true"`
    D []string      `env:"W" default:"[\"foo\", \"bar\", \"baz\"]"`
}
```

When pulling values from a variable source, a missing value (or empty string)
will use the default value, if provided. If no value is set for a required
configuration value, a fatal error will occur. String values will retrieve
the variable value unaltered. All other field types will attempt to unmarshal
the variable value as JSON.

Then, an initializer or a process that requires these config values can
retrieve them in its `Init` method.

```go
func (p *Process) Init(config nacelle.Config) error {
    c := &Config{}
    if err := config.Load(c); err != nil {
        // ...
    }

    // c is hydrated
    // ...
}
```

## PostLoading Configuration Structs

After hydration, the `PostLoad` method will be invoked on all registered
configuration structs (where such a method exists). This allows additional
non-type validation to occur, and to create any types which are not
directly/easily encodable as JSON.

```go
func (c *Config) PostLoad() error {
    if c.Field != "foo" && c.Field != "bar" {
        return fmt.Errorf("field value must be foo or bar")
    }

    return nil
}
```

## Embedded Configs

It is possible to embed anonymous configuration structs in order to get
configuration reusability. Embedded config structs have the same set of
struct tags.

```go
type (
    BaseConfig struct {
        X string `env:"X"`
        Y string `env:"Y"`
        Z string `env:"Z"`
    }

    ProducerConfig struct {
        BaseConfig
        W string `env:"W"`
    }

    ConsumerConfig struct {
        BaseConfig
        Q string `env:"Q"`

    }
)
```

## Sources

By default, the nacelle bootstrapper uses the environment as the source for
configuration data. You can alternatively supply a Sourcer which pulls from
values from a different source (file, network, etc).

The following struct loads a variable `X` from the environment or loads the
path `a.b.c` from a configuration file (this assumes the configuration file
contains a nested dictionary structure with the path `a.b.c`).

```go
type Config struct {
    X string `env:"x" file:"a.b.c"`
}
```

This scenario also assumes the following option is given to the bootstrapper.

```go
WithConfigSourcer(NewMultiSourcer(
    NewYAMLFileSourcer("config.yaml"),
    NewEnvSourcer("APP"),
))
```

## Config Tags

In some circumstances, it may be necessary to dynamically alter the tags
on a configuration struct. This has become an issue in two circumstances
so far. First, two instances of the same configuration struct may need to
be registered but must be configured separately (i.e. they need to look at
distinct environment variables). This is a particular problem when running
two HTTP servers with the same base, for example. Second, the default value
of a field may need to be altered. This issue can also arise when two
instances of the same configuration struct are registered but shouldn't get
clashing defaults (e.g. a default listening port).

Nacelle provides two tag modifiers which can be applied at configuration
registration time. In the following, the configuration struct is loaded
such that the environment variables used to hydrate the object are `Q_X`,
`Q_Y`, `Q_Z`, `Q_W`, instead of `X`, `Y`, `Z`, and `W` the default value
of the struct field `B` (loaded through the environment variable `Q_Y`) is
false instead of true.

```go
func (p *Process) Init(config nacelle.Config) error {
    c := &Config{}

    if err := config.Load(
        c,
        nacelle.NewEnvTagPrefixer("Q")
        nacelle.NewDefaultTagSetter("B", "false"),
    ); err != nil {
        // ...
    }

    // c is hydrated
    // ...
}
```

If other dynamic modifications of a configuration struct is necessary,
simply implement the `TagModifier` interface and use it in the same way.
