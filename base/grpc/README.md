# Nacelle Base gRPC Process

This package contains a base process implementation for a gRPC server. For a more
full-featured gRPC framework, see [scarf](https://github.com/efritz/scarf).

## Usage

To use the server, initialize a process by passing a Server Initializer to the `NewServer`
constructor. A server initializer is an object with an `Init` method that takes a nacelle
config object (as all process initializer methods do) as well as a `*grpc.Server`. This
*hook* is provided so that services can be registered to the gRPC server before it begins
accepting clients.

The server initializer will have services injected and will receive the nacelle config
object on initialization as if it were a process.

To get a better understanding of the full usage, see the
[example](https://github.com/efritz/nacelle/tree/master/examples/grpc).

## Configuration

The default process behavior can be configured by the following environment variables.

| Environment Variable | Default | Description |
| -------------------- | ------- | ----------- |
| GRPC_HOST            | 0.0.0.0 | The host on which the server accepts clients. |
| GRPC_PORT            | 6000    | The port on which the server accepts clients. |

## Using Multiple Servers

In order to run multiple gRPC servers, tag modifiers can be applied during config
registration. For more details on how to do this, see the
[example](https://github.com/efritz/nacelle/tree/master/examples/multi-grpc).

Remember that multiple services can be registered to the same grpc.Server instance, so
multiple processes may not even be necessary depending on your use case.
