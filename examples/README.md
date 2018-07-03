# Nacelle Application Examples

This directory serves as repository of usage examples as well as a sanity test to ensure
that app startup is not broken when nacelle is updated.

- [http](https://github.com/efritz/nacelle/tree/master/examples/http) is a basic HTTP server
- [grpc](https://github.com/efritz/nacelle/tree/master/examples/grpc) is a bacic gRPC server
- [worker](https://github.com/efritz/nacelle/tree/master/examples/worker) is a basic worker that simply counts on a timer
- [multi-http](https://github.com/efritz/nacelle/tree/master/examples/multi-http) is an app with two HTTP servers
- [multi-grpc](https://github.com/efritz/nacelle/tree/master/examples/multi-grpc) is an app with two gRPC servers
- [multi-process](https://github.com/efritz/nacelle/tree/master/examples/multi-process) is an *exploding secret* API with an HTTP and gRPC interface backed by Redis
