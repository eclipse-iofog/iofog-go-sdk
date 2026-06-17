# Changelog

## [v3.8.0] - 17 June 2026

### Breaking changes

- **Module path:** `github.com/eclipse-iofog/iofog-go-sdk/v3` (neutral upstream). Update imports from `github.com/datasance/iofog-go-sdk/v3`.
- **Go 1.26.4** minimum (see `go.mod`).
- **`pkg/microservices` — EdgeletAPI v1:** greenfield renames with no deprecated aliases:
  - `IoFogClient` → `EdgeletAPIClient`
  - `NewIoFogClient` / `NewIoFogClientV3` → `NewEdgeletAPIClient(id string, opts ...ClientOption)`
  - `NewDefaultIoFogClient` / `NewDefaultIoFogClientV3` → `NewDefaultEdgeletAPIClient`
  - `V3APIError` → `EdgeletAPIError`
  - `PortIoFog` → `PortEdgeletAPI`
  - Defaults and routes target **EdgeletAPI v1** (`/v1/microservices/*` on port 54321).
- **`pkg/apps`:** default deploy YAML `apiVersion` is `iofog.org/v3`. Datasance-flavored callers use `apps.WithAPIVersion("datasance.com/v3")`.
- **`pkg/apps` deploy types:** canonical image keys `amd64`, `arm64`, `riscv64`, `arm` (removed `x86`/`arm`); `flow` → `application`; `dockerUrl` → `containerEngineUrl`.
- **`pkg/client` — Controller REST v3.8:** `fogTypeId`/`FogType`/`agentType` → `archId`/`ArchID`/`arch`; flow APIs removed in favor of application name; auth endpoint updates (`POST /users`, `POST /user/change-password`); `RefreshUserSubscriptionKey` removed.

### Added

- `pkg/apps.DefaultAPIVersion` and `WithAPIVersion` deploy option.
- `pkg/arch` shared architecture codes for Controller v3.8.
- GitHub Actions CI (`.github/workflows/ci.yml`): `make lint`, `make security-code`, `make vulncheck`.
- `SECURITY.md` maintainer security gates and documented gosec exceptions.
- `NOTICE` Eclipse ioFog attribution; per-file Datasance copyright headers removed from `pkg/**`.

### Changed

- Copyright attribution consolidated in `NOTICE`.
- golangci-lint v2 config; gosec runs via `make security-code` (not inside golangci-lint).

### Removed

- Azure Pipelines workflow (`azure-pipelines.yml`).
- Deprecated `IoFogClient` / `V3APIError` aliases and legacy four-argument `NewIoFogClient` constructor.

### Migration

Replace Datasance module imports with the neutral upstream path:

```go
// before
import "github.com/datasance/iofog-go-sdk/v3/pkg/microservices"

// after
import "github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/microservices"
```

Update microservice client construction:

```go
client, err := msvcs.NewDefaultEdgeletAPIClient()
```

For Datasance-flavored deploy YAML, pass the apiVersion override at the caller:

```go
apps.DeployApplication(ctrl, app, name, apps.WithAPIVersion("datasance.com/v3"))
```

The Datasance git mirror (`github.com/Datasance/iofog-go-sdk`) ships the same commit SHA as `eclipse-iofog/iofog-go-sdk`; only the **module import path** changes.

## [v3.0.0-beta1] - 13 Auguest 2021

* No changes since alpha2

## [v3.0.0-alpha2] - 28 July 2021

* REST client initialized with Controller base URL
* Go 1.16

## [v3.0.0-alpha1] - 11 March 2021

* Add k8s/operator pkg for operator SDK helpers
* Decrease default timeout for Controller REST client (120s->4s)
* Allow users to specify timeout seconds for Controller REST client
* Add support for EdgeResources
* Add support for Application Templates
* Add supports for Agent upgrade and rollback in Controller REST client
* Update Controller REST client to be aware of backend version
* Add support for new Application routes in Controller REST client

## [v2.0.0]

* Added type to volume mapping
* Update WaitForLoadBalancer to get hostname too
* Fix error reporting in WaitForLoadBalancer func
* Update GetAgentByName to take system flag
* Update ListAgents and allow users to specify filters and system flag
* Remove iofogctl resource kinds
* Update AgentConfiguration for configuring docker frequency
### Bugs

* Stop passing ListAgentsRequest into body of request

## [v2.0.0-alpha] - 2020-03-10

### Features

* Add PutDefaultProxy function to client pkg
* Set retries as optional on new client
* Add retry policy to client
* Add omitempty to optional PATCH msvc args
* Add GetAllMicroservicePublicPorts function
* Move port validation to be run in application deployment too
* Add PublicLink to msvcPortMapping and update microservice update to detect public port mapping changes
* Add DefaultRouterName constant
* Update publicPort json key
* Add router fields to AgentConfiguration struct
* Update agent info to contain router info 
* Update routerConfig in agent config yaml
* Add networkRouter to AgentConfiguration
* Add isSystem in agent and applications
* Update AgentCreateRequest to allow for configuration
* Add Agent Prune API call
* Add PORT to apps.Microservice.Container
* Make CMD optional on microservice update
* Allows CMD in microservice creation and update

### Bugs

* Fix make gen to update file in $PWD/pkg/apps

## [v1.3.0]

* Add client package to the repo
* Re-organize the repo to maintain multiple packages
  
[Unreleased]: https://github.com/eclipse-iofog/iofog-go-sdk/compare/v3.8.0..HEAD
[v3.8.0]: https://github.com/eclipse-iofog/iofog-go-sdk/compare/v3.8.0-beta.2..v3.8.0
[v2.0.0-beta3]: https://github.com/eclipse-iofog/iofog-go-sdk/compare/v2.0.0-beta2..v2.0.0-beta3
[v2.0.0-beta]: https://github.com/eclipse-iofog/iofog-go-sdk/compare/v2.0.0-alpha..v2.0.0-beta2
[v2.0.0-alpha]: https://github.com/eclipse-iofog/iofog-go-sdk/compare/v1.3.0..v2.0.0-beta
[v1.3.0]: https://github.com/eclipse-iofog/iofog-go-sdk/tree/v1.3.0
