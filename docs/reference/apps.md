# Apps (YAML deploy)

Package `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/apps` turns Go values into Controller YAML and creates or updates the resource. It depends on `pkg/client`. Deploy behavior is in [Deploy YAML](../howto/deploy-yaml.md).

## Entry points

```go
func DeployApplication(controller IofogController, application any, name string, opts ...DeployOption) error
func DeployMicroservice(controller IofogController, microservice any, appName, name string, opts ...DeployOption) error
func DeployApplicationTemplate(controller IofogController, controllerBaseURL *url.URL, template any, name string, opts ...DeployOption) error
func DeployMicroserviceTemplate(controller IofogController, controllerBaseURL *url.URL, template any, name string, opts ...DeployOption) error
```

```go
const DefaultAPIVersion = "iofog.org/v3"

func WithAPIVersion(v string) DeployOption
```

`IofogController`:

| Field | Use |
| --- | --- |
| `Endpoint` | Absolute base URL for application and microservice deploy, including `/api/v3`. |
| `Token` | When non-empty, `client.NewWithToken`. |
| `RefreshToken`, `Email`, `Password` | Used by `client.SessionLogin` when `Token` is empty. |

Template deploy ignores `Endpoint` and uses `controllerBaseURL`.

`ParseFQMsvcName(fqName)` splits `app/name` into two strings, or a single segment into `("", name)`. Any other shape returns `*InputError`.

## Kinds

```go
ApplicationKind          Kind = "Application"
ApplicationTemplateKind  Kind = "ApplicationTemplate"
MicroserviceKind         Kind = "Microservice"
MicroserviceTemplateKind Kind = "MicroserviceTemplate"
ModelKind                Kind = "Model"
KnowledgeKind            Kind = "Knowledge"
RuntimeClassKind         Kind = "RuntimeClass"
RouteKind                Kind = "Route" // deprecated; unmarshals only
```

```go
type Header struct {
	APIVersion string
	Kind       Kind
	Metadata   HeaderMetadata // Name, Namespace
	Spec       any
}
```

`IofogHeader` is an alias of `Header`. Executors build this header and `yaml.Marshal` it. `RuntimeClass` in this package is not wrapped in `spec`: `handler` sits next to `metadata`. Upload that document with the client, not with a deploy function.

## Application spec

`Application`: `name`, `microservices`, optional `natsConfig`, `id`, optional `template`.

`ApplicationTemplate`: `name`, `description`, `variables`, `application` (`ApplicationTemplateInfo` with `microservices` and `natsConfig`).

`TemplateVariable`: `key`, `description`, `defaultValue`, `value`. The JSON fields are `*apiextensions.JSON`. Deep copy is shallow.

## Microservice spec

`Microservice` field order matches describe/redeploy YAML: `uuid`, `application`, `name`, `agent`, `images`, `natsConfig`, `models`, `knowledge`, `container`, `schedule`, `config`, `serviceAccount`, `template`, `created`, `rebuild`, `status`, `execStatus`.

`agent.name` is the fog name (`MicroserviceAgent`).

`images` (`MicroserviceImages`): `registry` (`RegistryRef`), `arm64`, `amd64`, optional `riscv64`, `arm`, optional `catalogId`.

`config` and `container.annotations` use `ArbitraryJSON`. YAML maps and lists are normalized before JSON encoding so describe/redeploy round-trips. Use `ToMap` / `FromMap` when you want `map[string]any`.

`template` (`MicroserviceTemplateRef`): `name` plus `variables`. Variables unmarshal from a YAML mapping or a list of `{key, value}` pairs and marshal as a mapping.

`serviceAccount.roleRef`: `kind`, `name`, optional `apiGroup`.

`MicroserviceTemplate` spec: `name`, `description`, `variables`, `microservice`.

### Container

`MicroserviceContainer` includes network and privilege flags, `runAsUser`, `runAsGroup`, `readOnlyRootFilesystem`, `ipcMode`, `pidMode`, `platform`, `runtime`, `capAdd`, `capDrop`, `annotations`, `sysctls`, `ulimits` (soft/hard, `-1` means unlimited), `cpuSetCpus`, `cpus`, memory fields, `shmSize`, `cdiDevices`, `devices` (`hostPath`, `containerPath`, `permissions`), `volumes`, `tmpfs` (`size` is MiB, `mode` is an optional octal string), `extraHosts`, `env`, `ports`, `workingDir`, `entrypoint`, `commands`, `healthCheck`.

`MicroservicePortMapping` has no `yaml` tags. `yaml.v2` uses the field name (`Internal`, `External`, `Protocol`) unless you populate the struct from JSON. Prefer a YAML struct you control, or set the fields in Go and let the executor marshal them, if you need lowercase keys. The JSON tags are `internal`, `external`, `protocol`.

`MicroserviceVolumeMapping`: `hostDestination`, `containerDestination`, `accessMode`, optional `type` (`bind`, `serviceAccount`, `volume`, and other Controller kinds), optional `scope` (`private` or `shared`, only when `type` is `volume`).

`MicroserviceEnvironment`: `key`, and one of `value`, `valueFromSecret`, `valueFromConfigMap`.

### Status (describe output)

`MicroserviceStatusInfo` matches the client status fields, including `lastError`, `lastErrorAt`, `restartCount`, and `podId`. `MicroserviceExecStatusInfo`: `status`, `execSessionId`.

## Fleet YAML specs

These types are the `spec` (or, for RuntimeClass, the whole document). Deploy functions do not upload them.

| Type | Fields |
| --- | --- |
| `Model` | `repo`, `revision`, `registryId`, `files`, `format` |
| `Knowledge` | same |
| `RuntimeClass` | `apiVersion`, `kind`, `metadata`, `handler` |

`CatalogItem`: `id`, `registry` (`RegistryRef`), `arm64`, `amd64`, optional `riscv64` and `arm`, `name`, `description`, `configExample`.

## RegistryRef

`RegistryRef` is an integer registry id.

| Input | Result |
| --- | --- |
| empty, null, `0` | `1` (remote) |
| `remote` | `1` |
| `local` | `2` |
| numeric string or number | that id |

`Int()` returns `1` when the value is `0`. `MarshalYAML` and `MarshalJSON` emit the integer from `Int()`, unquoted in YAML.

## Agent configuration in YAML

`apps.AgentConfiguration` is the agent block used by higher-level YAML (container engine, disk and memory limits, GPS, watchdog, router, NATS, `archId`). It does not include `deviceScanFrequency`, `bluetoothEnabled`, or `abstractedHardwareEnabled`.

Router comments in the struct: `routerMode` is `edge`, `interior`, or `none` (default `edge`). `routerPort` default `5671`. `upstreamRouters` is ignored when mode is `none`. `networkRouter` is required when mode is `none`.

NATS: `natsMode` is `none`, `leaf`, or `server`, plus port and JetStream size strings.

## Deep copy

`pkg/apps/deepcopy_generated.go` is generated. Regenerate with `make gen` (Darwin) or `make gen-linux`. `make gen-check` fails on drift. `TemplateVariable`, `ArbitraryJSON`, `MicroserviceTemplateVariables`, and `MicroserviceTemplateRef` are hand-written because deepcopy-gen cannot emit them.

## Errors

`NewError`, `NewNotFoundError`, `NewConflictError`, `NewInputError`, `NewInternalError`, `NewHTTPError` mirror the client constructors. Deploy lookup treats `*client.NotFoundError` as "create". Other client errors propagate.

## Related

- [Deploy YAML](../howto/deploy-yaml.md)
- [Client](client.md)
- [Codes](codes.md)
