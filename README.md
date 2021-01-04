<div align="center"><img width="160" src="https://raw.githubusercontent.com/go-nacelle/nacelle/master/images/nacelle.png" alt="Nacelle logo"></div>

# Nacelle service framework

[![PkgGoDev](https://pkg.go.dev/badge/badge/github.com/go-nacelle/nacelle.svg)](https://pkg.go.dev/github.com/go-nacelle/nacelle) [![CircleCI status](https://circleci.com/gh/go-nacelle/nacelle.svg?style=svg)](https://circleci.com/gh/go-nacelle/nacelle) [![Coverage status](https://coveralls.io/repos/github/go-nacelle/nacelle/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/nacelle?branch=master) ![Sonarcloud bugs count](https://sonarcloud.io/api/project_badges/measure?project=go-nacelle_nacelle&metric=bugs) ![Sonarcloud vulnerabilities count](https://sonarcloud.io/api/project_badges/measure?project=go-nacelle_nacelle&metric=vulnerabilities) ![Sonarcloud maintainability rating](https://sonarcloud.io/api/project_badges/measure?project=go-nacelle_nacelle&metric=sqale_rating) ![Sonarcloud code smells count](https://sonarcloud.io/api/project_badges/measure?project=go-nacelle_nacelle&metric=code_smells)

---

See the package documentation on [nacelle.dev](https://nacelle.dev).

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
- Impose opinions on runtime environment, deployment, or orchestration of services.
