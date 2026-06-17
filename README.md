# ioFog Go SDK

This SDK contains a set of Golang packages developers can use for the purposes of:

* Developing your own edge microservices that run on ioFog Edge Compute Networks
* Interacting with ioFog Controller REST API

#### Golang v1.18

The SDK is tested against Golang v1.18 environments.

## Packages

The following is a high-level overview of the functionality provided by each package.

Each package contains its own README.md so please refer to those for further details.

#### Microservices

The `microservices` package contains functionality required to implement edge microservices that run on ioFog Edge
Compute Networks. It is aligned with EdgeletAPI v1 and includes microservice configuration retrieval and control websocket
signal handling over HTTPS/WSS using mounted service-account token and CA material.

#### Client

The `client` package contains an HTTP client to use with ioFog Controller's REST API. You can view see the full REST API
specification at [iofog.org](https://iofog.org/docs/1.3.0/controllers/rest-api.html).

#### Deploy applications

The `deployapps` package contains executors to deploy iofog applications and microservices using the `client` package.
This package is used by `iofogctl` and `iofog-operator` to deploy applications and microservices based on yaml
configuration files.