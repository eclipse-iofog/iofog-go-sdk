# ioFog Go SDK

[![CI](https://github.com/eclipse-iofog/iofog-go-sdk/actions/workflows/ci.yml/badge.svg)](https://github.com/eclipse-iofog/iofog-go-sdk/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.26.4-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-EPL--2.0-blue.svg)](LICENSE)

**Upstream:** [eclipse-iofog/iofog-go-sdk](https://github.com/eclipse-iofog/iofog-go-sdk) · **Datasance distribution:** [Datasance/iofog-go-sdk](https://github.com/Datasance/iofog-go-sdk)

This SDK contains a set of Golang packages developers can use for the purposes of:

* Developing your own edge microservices that run on ioFog Edge Compute Networks
* Interacting with ioFog Controller REST API

## Go

The SDK requires **Go 1.26.4+** (see `go.mod`).

```bash
go get github.com/eclipse-iofog/iofog-go-sdk/v3
```

## Packages

The following is a high-level overview of the functionality provided by each package.

Each package contains its own README.md so please refer to those for further details.

#### Microservices

The `microservices` package (`github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/microservices`) contains functionality required to implement edge microservices that run on ioFog Edge Compute Networks. It is aligned with **EdgeletAPI v1** and includes microservice configuration retrieval and control websocket signal handling over HTTPS/WSS using mounted service-account token and CA material.

#### Client

The `client` package (`github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client`) contains an HTTP client to use with ioFog Controller's REST API v3. You can view the full REST API specification at [iofog.org](https://iofog.org/docs/1.3.0/controllers/rest-api.html).

#### Deploy applications

The `apps` package (`github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/apps`) contains executors to deploy ioFog applications and microservices using the `client` package. This package is used by `iofogctl` and `iofog-operator` to deploy applications and microservices based on YAML configuration files with default `apiVersion: iofog.org/v3`.

## Quality gates

Quality gates before contributing:

```bash
make lint
make security-code
make vulncheck
```
