# ioFog Go SDK documentation

This directory is the production documentation for module `github.com/eclipse-iofog/iofog-go-sdk/v3`. It describes the public API in the current source tree (Controller REST **v3.9** client, EdgeletAPI **v1**, YAML `apiVersion` default **`iofog.org/v3`**). Release notes live in [`CHANGELOG.md`](../CHANGELOG.md).

The SDK is a library. It does not ship a CLI. `iofogctl` and `iofog-operator` call these packages.

## Start here

| If you want to | Read |
| --- | --- |
| Install the module and call the Controller once | [Getting started](getting-started.md) |
| Understand which package to import | [Packages](explanation/packages.md) |
| Authenticate, TLS, and token handling | [Security model](explanation/security.md) |
| See why routes were removed | [Messaging](explanation/messaging.md) |

## How-to guides

| Task | Guide |
| --- | --- |
| Log in and call Controller REST | [Use the Controller client](howto/controller-client.md) |
| Deploy Application or Microservice YAML | [Deploy YAML](howto/deploy-yaml.md) |
| Read config inside an edge container | [Write an edge microservice](howto/edge-microservice.md) |
| Exec into a container or tail logs | [Exec and logs](howto/exec-and-logs.md) |
| Manage models, Knowledge, and RuntimeClasses | [Fleet resources](howto/fleet-resources.md) |

## Reference

| Package | Import path | Reference |
| --- | --- | --- |
| Controller REST client | `.../pkg/client` | [Client](reference/client.md) |
| YAML deploy executors | `.../pkg/apps` | [Apps](reference/apps.md) |
| EdgeletAPI v1 | `.../pkg/microservices` | [Microservices](reference/microservices.md) |
| NATS control plane | methods on `client.Client` | [NATS](reference/nats.md) |
| Kubernetes helpers | `.../pkg/k8s`, `.../pkg/k8s/operator` | [Kubernetes](reference/kubernetes.md) |
| Architecture codes and registry ids | `.../pkg/arch`, `client` maps | [Codes](reference/codes.md) |
| Errors, retries, HTTP status | `client`, `apps`, `microservices` | [Errors](reference/errors.md) |

## Operations

Build, test, and release gates: [Operations](operations.md).
