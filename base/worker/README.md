# Nacelle Base Worker Process

This package contains a base process implementation for a generic worker which performs some
work prompted by an external event (e.g. a remote work queue on a schedule).

## Usage

To use a worker, initailize a worker by passing a Worker Spec to the `NewWorker` constructor.
A worker spec is an object with an `Init` method that takes a ncelle config object (as all
process initializer methods do) as well as a reference to the worker, and a `Tick` method. On
a timer, the worker will prompt the worker to perform a single unit of work.

The worker spec will have services injected and will receive the nacelle config object on
initialization as if it were a process.

To get a better understanding of the full usage, see the
[example](https://github.com/efritz/nacelle/tree/master/examples/worker).

## Configuration

The default process behavior can be configured by the following environment variables.

| Environment Variable | Default | Description |
| -------------------- | ------- | ----------- |
| WORKER_TICK_INTERVAL | 0       | The duration (in seconds) between calls to the spec's tick method. |
