# Nacelle Base HTTP Process

This package contains a base process implementation for an HTTP server. For a more
full-featured HTTP framework, see [chevron](https://github.com/efritz/chevron).

## Usage

To use the server, initialize a process by passing a Server Initializer to the `NewServer`
constructor. A server initializer is an object with an `Init` method that takes a nacelle
config object (as all process initializer methods do) as well as an `*http.Server`. This
*hook* is provided so the the application can independently configure the server's HTTP
handler.

The server initializer will have services injected and will receive the nacelle config
object on initialization as if it were a process.

To get a better understanding of the full usage, see the
[example](https://github.com/efritz/nacelle/tree/master/examples/http).

## Configuration

The default process behavior can be configured by the following environment variables.

| Environment Variable  | Default | Description |
| --------------------- | ------- | ----------- |
| HTTP_HOST             | 0.0.0.0 | The host on which the server accepts clients. |
| HTTP_PORT             | 5000    | The port on which the server accepts clients. |
| HTTP_CERT_FILE        |         | An absolute path to a cert file. |
| HTTP_KEY_FILE         |         | An absolute path to a key file. |
| HTTP_SHUTDOWN_TIMEOUT | 5       | The duration (in seconds) to allow for a graceful shutdown. |

If both a cert file and key file paths are provided, the server will serve TLS. It is an
error to provide a path to a cert file or a path to a key file but not both.

## Using Multiple Servers

In order to run multiple HTTP servers, tag modifiers can be applied during config
registration. For more details on how to do this, see the
[example](https://github.com/efritz/nacelle/tree/master/examples/multi-http).
