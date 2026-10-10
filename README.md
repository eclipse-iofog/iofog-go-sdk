# ioFog Go SDK

[![CI](https://github.com/eclipse-iofog/iofog-go-sdk/actions/workflows/ci.yml/badge.svg)](https://github.com/eclipse-iofog/iofog-go-sdk/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.26.9-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-EPL--2.0-blue.svg)](LICENSE)
[![govulncheck](https://github.com/eclipse-iofog/iofog-go-sdk/actions/workflows/govulncheck.yml/badge.svg)](https://github.com/eclipse-iofog/iofog-go-sdk/actions/workflows/govulncheck.yml)

**Upstream:** [eclipse-iofog/iofog-go-sdk](https://github.com/eclipse-iofog/iofog-go-sdk) · **Datasance distribution:** [Datasance/iofog-go-sdk](https://github.com/Datasance/iofog-go-sdk)

This SDK contains a set of Golang packages developers can use for the purposes of:

* Developing your own edge microservices that run on ioFog Edge Compute Networks
* Interacting with ioFog Controller REST API

## Go

The SDK requires **Go 1.26.9+** (see `go.mod`).

```bash
go get github.com/eclipse-iofog/iofog-go-sdk/v3
```

## Documentation

Production documentation is in [`docs/`](docs/README.md): install and first call, how-to guides, and a reference for every public package (`client`, `apps`, `microservices`, NATS, Kubernetes helpers).

Package READMEs under `pkg/` are short entry points. The `docs/` tree is the contract for current behavior.

## Packages

The following is a high-level overview of the functionality provided by each package.

#### Microservices

The `microservices` package (`github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/microservices`) contains functionality required to implement edge microservices that run on ioFog Edge Compute Networks. It is aligned with **EdgeletAPI v1** and includes microservice configuration retrieval and control websocket signal handling over HTTPS/WSS using mounted service-account token and CA material.

#### Client

The `client` package (`github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client`) contains an HTTP and WebSocket client for ioFog Controller REST API v3 (`BaseURL` ends at `/api/v3`). Method list, field units, and auth behavior: [`docs/reference/client.md`](docs/reference/client.md).

#### Deploy applications

The `apps` package (`github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/apps`) contains executors to deploy ioFog applications and microservices using the `client` package. This package is used by `iofogctl` and `iofog-operator` to deploy applications and microservices based on YAML configuration files with default `apiVersion: iofog.org/v3`.

## Quality gates

Quality gates before contributing:

```bash
make lint
make security-code
make vulncheck
```
