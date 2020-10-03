<div align="center"><img width="160" src="https://raw.githubusercontent.com/go-nacelle/nacelle/master/images/nacelle.png" alt="Nacelle logo"></div>

# Nacelle [![GoDoc](https://godoc.org/github.com/go-nacelle/nacelle?status.svg)](https://godoc.org/github.com/go-nacelle/nacelle) [![CircleCI](https://circleci.com/gh/go-nacelle/nacelle.svg?style=svg)](https://circleci.com/gh/go-nacelle/nacelle) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/nacelle/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/nacelle?branch=master)

Service framework written in Go.

---

For more details, see [the website](https://nacelle.dev) and the [getting started gide](https://nacelle.dev/getting-started).

## Goals

Core goals:

- Provide a common convention for application organization so that developers can quickly dive into the meaningful logic of an application.
- Support a common convention for declaring, reading, and validating configuration values from the runtime environment.
- Support a common convention for registering, declaring, and injecting struct and interface dependencies.
- Support a common convention for structured logging.

Additional goals:

- Provide additional non-core functionality via separate opt-in libraries. Keep the dependencies for the core-functionality minimal.
- Operate within existing infrastructures and do not require tools or technologies outside of what this project provides.

## Non-goals

- Impose opinions on service discovery.
- Impose opinions on inter-process or inter-service communication.
- Impose opinions on runtime environment, deployment, orchestration.services.
