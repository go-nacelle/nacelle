# Nacelle Processes

Nacelle applications are organized into three distinct categories:

- A **service** is a shared APIs. Services exist in support of a *process*.
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

## Example

TODO
