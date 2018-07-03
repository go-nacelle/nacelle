# Nacelle Logging

Nacelle provides rather opinionated structured logging following a few short principles.

(1) Logs must be structured. It is absolutely essential to be able to correlate and
    aggregate log messages to form a view of a running system.
(2) Applications should always log to a standard error. In order to redirect logs to a
    secondary target (such as an ELK stack), the application's output should simply be
    redirected. This keeps the application simple and allows redirection of logs to
    **any** source without requiring an application update. For an example of redirection
    when run in a Docker container, see nacelle's
    [fluentd wrapper](https://github.com/efritz/nacelle-fluentd).

The interfaces provided here are backed by [gomol](https://github.com/aphistic/gomol).

## Interface

The logging interface is simple. Each log method (each one named after the level at which
it will log) provides a printf-like interface. Logging at the fatal level will abort the
application after the log has been flushed.

```go
logger.Debug("A debug message (%#v)", args...)
logger.Info("A info message (%#v)", args...)
logger.Warning("A warning message (%#v)", args...)
logger.Error("A error message (%#v)", args...)
logger.Fatal("A fatal message (%#v)", args...)
```

Each log method also has a WithFields variant, which takes a `Fields` object as its first
parameter. Fields are a map from strings to interfaces. Each field provided to the logger
will be output with (but separately from) the formatted message.

```go
fields:=nacelle.Fields{
    "foo": "bar",
    "baz": 12345,
}

logger.DebugWithFields(fields, "A debug message (%#v)", args...)
logger.InfoWithFields(fields, "A info message (%#v)", args...)
logger.WarningWithFields(fields, "A warning message (%#v)", args...)
logger.ErrorWithFields(fields, "A error message (%#v)", args...)
logger.FatalWithFields(fields, "A fatal message (%#v)", args...)
```

A logger can also be decorated with fields and used later so that multiple log messages
share the same set of fields. This is useful for request correlation in servers where a
logger instance can be given a unique request identifier.

```go
requestLogger := logger.WithFields(fields)

// Same as above
requestLogger.Info("A debug message (%#v)", args...)
```

A logger should **not** be constructed, but should be injected via a service container.
See the [service package documentation](https://github.com/efritz/nacelle/tree/master/service)
for additional formation.

## Configuration

The default logging behavior can be configured by the following environment variables.

| Environment Variable         | Default | Description |
| ---------------------------- | ------- | ----------- |
| LOG_LEVEL                    | info    | The highest level that will be emitted. |
| LOG_ENCODING                 | console | `console` for human-readable output and `json` for JSON-formatted output. |
| LOG_FIELDS                   |         | A JSON-encoded map of fields to include in every log. |
| LOG_FIELD_BLACKLIST          |         | A JSON-encoded list of fields to omit from logs. Works with `console` encoding only. |
| LOG_COLORIZE                 | true    | Colorize log messages by level when true. Works with `console` encoding only. |
| LOG_SHORT_TIME               | false   | Omit date from timestamp when true. Works with `console` encoding only. |
| LOG_DISPLAY_FIELDS           | true    | Omit log fields from output when false. Works with `console` encoding only. |
| LOG_DISPLAY_MULTILINE_FIELDS | true    | Print fields on one line when true, one field per line when false. Works with `console` encoding only. |

## Adapters

Nacelle ships with a handful of logging adapters - objects which wrap a logger
instance and decorates it with some additional behavior or data. A custom adapter
can be created for behavior which is not provided here.

### Replay

The *ReplayAdapter* supports replaying a sequence of log messages but at a higher
log level.

*Example use case:* Each request in an HTTP server has a unique log adapter which
traces the request. This adapter generally logs at the DEBUG level. When a request
encounters an error or is being served slowly, the entire trace can be replayed at
a higher level so the entire context is available for analysis.

```go
adapter := NewReplayAdapter(
    logger,          // base logger
    log.LevelDebug,  // track debug messages for replay
    log.LevelInfo,   // also track info messages
)

// ...

if requestIsTakingLong() {
    // Re-log journaled messages at warning level
    adapter.Replay(log.LevelWarning)
}
```

Messages which are replayed at a higher level will keep the original message timestamp
(if supplied), or use the time the `Log` message was invoked (if not supplied). Each
message will also be sent with an additional field called `replayed-from-level` with a
value equal to the original level of the message.

### Rollup

The *RollupAdapter* supports collapsing similar log messages into a multiplicity. This
is intended to be used with a chatty subsystem that only logs a handful of messages for
which a higher frequency does not provide a benefit (for example, failure to connect to
a Redis cache during a network partition).

```go
adapter := NewRollupAdapter(
    logger,       // base logger
    time.Second,  // rollup window
)

for i:=0; i < 10000; i++ {
    adapter.Debug("Some problem here!")
}
```

A **rollup** begins once two messages with the same format string are seen within the
rollup window period. During a rollup, all log messages (except for the first in the
window) are discarded but counted, and the **first** log message in that window will
be sent at the end of the window period with an additional field called `rollup-multiplicity`
with a value equal to the number of logs in that window.
