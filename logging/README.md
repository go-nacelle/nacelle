# Nacelle Logging

TODO

## Replay

TOOD

The *ReplayAdapter* supports replaying a sequence of log messages but at a higher
log level.

*Intended use case:* Each request in an HTTP server has a unique log adapter which
traces the request. This adapter generally logs at the DEBUG level. When a request
encounters an error or is being served slowly, the entire trace can be replayed at
a higher level so the entire context is available for analysis.

## Example

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

## Rollup

The *RollupAdapter* supports collapsing similar log messages into a multiplicity. This
is intended to be used with a chatty subsystem that only logs a handful of messages for
which a higher frequency does not provide a benefit (for example, failure to connect to
a Redis cache during a network partition).

## Example

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
