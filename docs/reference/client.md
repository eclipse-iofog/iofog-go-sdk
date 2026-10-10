# Controller client

Package `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client` is the HTTP and WebSocket client for ioFog Controller REST **v3**. Paths below are relative to `Options.BaseURL`, which must already end at `/api/v3`.

Wire structs and JSON tags live in `pkg/client/types.go`. This page lists every exported method and the fields whose units or defaults are easy to misread. Authentication and TLS are in [Security model](../explanation/security.md).

## Construction

```go
type Options struct {
	BaseURL   *url.URL
	Retries   *Retries
	Timeout   int          // seconds; 0 becomes 10
	TLSConfig *tls.Config  // nil skips certificate verification
}
```

| Function | Behavior |
| --- | --- |
| `New(opt) *Client` | Builds the client and calls `GET /status`. Failure leaves the cached version empty. `BaseURL` must be non-nil. |
| `NewAndLogin(opt, email, password) (*Client, error)` | `New` then `Login`. Prompts for OTP when `Totp` is empty. |
| `SessionLogin(opt, token, email, password)` | `POST /user/refresh`. On failure, `NewAndLogin`. |
| `NewWithToken(opt, token)` | Stores `token` as the access token. |
| `NewWithRefreshToken(opt, refreshToken)` | Stores the argument via `SetAccessToken`. Does not call `/user/refresh`. |

| Method | Behavior |
| --- | --- |
| `GetBaseURL() string` | `baseURL.String()`. |
| `GetRetries` / `SetRetries` | Per-client retry policy. |
| `GetAccessToken` / `SetAccessToken` | Bearer token. Non-empty means `isLoggedIn`. |
| `GetRefreshToken` / `SetRefreshToken` | Stored refresh token. |
| `GetVersion() string` | Cached Controller version from `New`. |
| `GetVersionNumbers() (major, minor, patch int, error)` | Splits `x.y.z` after stripping a `-suffix`. |

`SetVerbosity(true)` and `IsVerbose` print HTTP bodies. `SetGlobalRetries` sets the default copied by `New` when `Options.Retries` is nil.

```go
type Retries struct {
	Timeout       int            // extra attempts after HTTP 408
	CustomMessage map[string]int // substring of error text -> max attempts
}
```

Sleep between attempts is `attempt` seconds. HTTP 404 becomes `*NotFoundError`. Every other non-2xx becomes `*HTTPError` with `Code` set. See [Errors](errors.md).

## Session

| Method | HTTP |
| --- | --- |
| `GetStatus() (ControllerStatus, error)` | `GET /status` |
| `Login(LoginRequest) error` | `POST /user/login` |
| `Refresh(RefreshTokenRequest) error` | `POST /user/refresh` |
| `Profile(WithTokenRequest) (UserProfile, error)` | `GET /user/profile` |
| `ChangePassword(ChangePasswordRequest) error` | `POST /user/change-password` |
| `Logout() error` | `POST /user/logout` |

`LoginRequest` fields: `Email`, `Password`, `Totp`. `ChangePasswordRequest`: `CurrentPassword`, `NewPassword`, `ResetToken`.

`ControllerStatus`: `Status`, `UptimeSeconds` (`uptimeSec`), `Versions.Controller`, `Versions.EcnViewer`.

## Users and groups

Embedded auth users (`/users`) are separate from `Login`.

| Method | HTTP |
| --- | --- |
| `ListAuthUsers() ([]AuthUserResponse, error)` | `GET /users` |
| `CreateAuthUser(AuthUserCreateRequest)` | `POST /users` |
| `GetAuthUser(id)` | `GET /users/{id}` |
| `UpdateAuthUser(id, AuthUserUpdateRequest)` | `PATCH /users/{id}` |
| `DeleteAuthUser(id)` | `DELETE /users/{id}` (soft delete) |
| `ResetAuthUserPassword(id) (AuthUserResetPasswordResponse, error)` | `POST /users/{id}/reset-password` |
| `ResetAuthUserToken(id) (AuthUserResetTokenResponse, error)` | `POST /users/{id}/reset-token` |

Create body: `email`, `password`, optional `groups`. Response includes `id`, `email`, `groups`, `mustChangePassword`, `mfaEnabled`, `isBootstrap`.

| Method | HTTP |
| --- | --- |
| `ListAuthGroups() ([]AuthGroupResponse, error)` | `GET /groups` |
| `CreateAuthGroup(AuthGroupCreateRequest)` | `POST /groups` |
| `GetAuthGroup(name)` | `GET /groups/{name}` (path-escaped) |
| `UpdateAuthGroup(name, AuthGroupUpdateRequest)` | `PATCH /groups/{name}` |
| `DeleteAuthGroup(name)` | `DELETE /groups/{name}` |

After a rename, use `response.Name` as the canonical name.

## Agents (fog nodes)

Several agent methods return `*Error` (`Controller client must be logged into...`) when the access token is empty: `CreateAgent`, `GetAgentProvisionKey`, `ListAgents`, `GetAgentByID`, `DeleteAgent`, `ReconcileAgent`.

| Method | HTTP | Notes |
| --- | --- | --- |
| `CreateAgent(*CreateAgentRequest) (CreateAgentResponse, error)` | `POST /iofog` | Returns `UUID`. Platform reconcile is async. |
| `GetAgentProvisionKey(uuid) (GetAgentProvisionKeyResponse, error)` | `GET /iofog/{uuid}/provisioning-key` | `Key`, `CaCert`, `ExpireTimeMsUTC`. |
| `ListAgents(ListAgentsRequest) (ListAgentsResponse, error)` | `GET /iofog-list` | Body is `fogs`. Filters are appended as `&filters[n][key\|value\|condition]`. |
| `GetAgentByID(uuid) (*AgentInfo, error)` | `GET /iofog/{uuid}` | Includes `platformStatus`. List responses omit it. |
| `GetAgentByName(name) (*AgentInfo, error)` | list + scan | `*NotFoundError` when missing. |
| `UpdateAgent(*AgentUpdateRequest) (*AgentInfo, error)` | `PATCH /iofog/{uuid}` | Then GET. Platform changes are async. |
| `DeleteAgent(uuid) error` | `DELETE /iofog/{uuid}` | Teardown continues in phase `Deleting`. |
| `RebootAgent(uuid) error` | `POST /iofog/{uuid}/reboot` | |
| `PruneAgent(uuid) error` | `POST /iofog/{uuid}/prune` | |
| `ReconcileAgent(uuid) (CreateAgentResponse, error)` | `POST /iofog/{uuid}/reconcile` | Resets `Failed`, clears backoff. |
| `WaitForAgentPlatformReady(uuid, timeout) error` | polls GET every 2s | Ready returns nil. Failed, Deleting, or timeout return `*Error`. |
| `SetNodeVersionCommand(uuid, versionCommand, *SetNodeVersionCommandRequest) error` | `POST /iofog/{uuid}/version/{upgrade\|rollback}` | |
| `UpgradeNode` / `RollbackNode` | same | `semver == nil` omits the body. |
| `UpgradeAgent` / `RollbackAgent` | lookup by name, then the node call | Second argument is `*string`. `nil` keeps the Controller default. |

`CreateAgentRequest` embeds `AgentUpdateRequest` (name, location, coordinates, description, tags, and `AgentConfiguration`).

### Platform phase

`PlatformPending`, `PlatformProgressing`, `PlatformReady`, `PlatformFailed`, `PlatformDeleting`.

`PlatformStatus` carries `generation`, `observedGeneration`, `lastError`, `lastTransitionAt`, and `conditions`.

### Host metrics on `AgentInfo`

Present on GET and list when the agent reports them. Older Controllers omit keys and the fields stay zero.

| Field | JSON | Unit |
| --- | --- | --- |
| `SystemCpus` | `systemCpus` | Logical CPU count. |
| `SystemTotalCPU` | `systemTotalCpu` | Busy percent, 0–100. Not the core count. |
| `SystemTotalMemory`, `SystemAvailableMemory` | `systemTotalMemory`, `systemAvailableMemory` | Bytes. |
| `SystemTotalDisk`, `SystemAvailableDisk` | `systemTotalDisk`, `systemAvailableDisk` | Bytes on the `diskDirectory` filesystem. |
| `DiskUsage` | `diskUsage` | Edgelet data-directory usage in **GiB**. Not host disk totals. |
| `SystemOs`, `SystemOsVersion`, `SystemKernelVersion` | matching camel case | Kernel release is empty on non-Linux. |

v3.9 status strings (parse them yourself): `runtimeClasses`, `availableCdiDevices`, `modelStatus`, `knowledgeStatus`. Counts: `activeModels`, `activeKnowledge`. Timestamps `modelLastUpdate` and `knowledgeLastUpdate` are Unix milliseconds, or `0` when the managed list is empty.

Removed and not decoded: `deviceScanFrequency`, `bluetoothEnabled`, `abstractedHardwareEnabled`. This client does not call HAL hardware or USB inventory endpoints. `EdgeGuardFrequency` remains.

`ArchID` is `archId`. `Arch` is the architecture object when the Controller sends it. NATS mode fields are `natsMode` (`none`, `leaf`, `server`) and the `nats*Port` pointers.

## Applications

| Method | HTTP |
| --- | --- |
| `GetAllApplications() (*ApplicationListResponse, error)` | `GET /application` |
| `GetApplicationByName(name) (*ApplicationInfo, error)` | `GET /application/{name}` |
| `CreateApplicationFromYAML(io.Reader) (*ApplicationInfo, error)` | `POST /application/yaml` field `application` |
| `UpdateApplicationFromYAML(name, io.Reader)` | `PUT /application/yaml/{name}` field `application` |
| `PatchApplication(name, *ApplicationPatchRequest)` | `PATCH /application/{name}` |
| `StartApplication(name)` | PATCH `isActivated: true` |
| `StopApplication(name)` | PATCH `isActivated: false` |
| `DeleteApplication(name) error` | `DELETE /application/{name}` |
| `GetAllSystemApplications` | `GET /application/system` |
| `GetSystemApplicationByName(name)` | `GET /application/system/{name}` |
| `DeleteSystemApplication(name)` | `DELETE /application/system/{name}` |

`ApplicationInfo`: `name`, `description`, `isActivated`, `isSystem`, `userId`, `id`, `microservices`, optional `natsConfig` (`natsAccess`, `natsRule`).

YAML create reads `name` from the create response, then GET.

## Application templates

Calls run `HEAD /capabilities/applicationTemplates` first. HTTP 404 becomes `*NotSupportedError` (`Controller API does not support Application Templates`). The client must be logged in.

| Method | HTTP |
| --- | --- |
| `IsApplicationTemplateCapable() error` | `HEAD /capabilities/applicationTemplates` |
| `ListApplicationTemplates` | `GET /applicationTemplates` |
| `GetApplicationTemplate(name)` | `GET /applicationTemplate/{name}` |
| `CreateApplicationTemplateFromYAML(io.Reader)` | `POST /applicationTemplate/yaml` field `template` |
| `UpdateApplicationTemplateFromYAML(name, io.Reader)` | `PUT /applicationTemplate/yaml/{name}` field `template` |
| `UpdateApplicationTemplateMetadata(name, *ApplicationTemplateMetadataUpdateRequest)` | `PATCH /applicationTemplate/{name}` |
| `DeleteApplicationTemplate(name)` | `DELETE /applicationTemplate/{name}` |

## Microservices

| Method | HTTP | Notes |
| --- | --- | --- |
| `GetAllMicroservices` | `GET /microservices` | |
| `GetMicroservicesByApplication(name)` | `GET /microservices?application={name}` | |
| `GetMicroserviceByName(app, name)` | list + scan | `*NotFoundError` if missing. |
| `GetMicroserviceByID(uuid)` | `GET /microservices/{uuid}` | |
| `CreateMicroserviceFromYAML(io.Reader)` | `POST /microservices/yaml` field `microservice` | Then GET by returned UUID. |
| `UpdateMicroserviceFromYAML(uuid, io.Reader)` | `PATCH /microservices/yaml/{uuid}` field `microservice` | |
| `DeleteMicroservice(uuid)` | `DELETE /microservices/{uuid}` | |
| `StartMicroservice(uuid)` | `PATCH /microservices/{uuid}/start` | |
| `StopMicroservice(uuid)` | `PATCH /microservices/{uuid}/stop` | |
| `RebuildsMicroservice(uuid)` | `PATCH /microservices/{uuid}/rebuild` | |
| `GetMicroservicePortMapping(uuid)` | `GET /microservices/{uuid}/port-mapping` | |
| `CreateMicroservicePortMapping(uuid, *MicroservicePortMappingInfo)` | `POST /microservices/{uuid}/port-mapping` | |
| `DeleteMicroservicePortMapping(uuid, *MicroservicePortMappingInfo)` | `DELETE /microservices/{uuid}/port-mapping/{internal}` | Uses `Internal`. |
| `PatchMicroserviceModels(uuid, MicroserviceCatalog)` | `PATCH /microservices/{uuid}/models` | Catalog-only. See [Fleet resources](../howto/fleet-resources.md). |
| `PatchMicroserviceKnowledge(uuid, KnowledgeCatalog)` | `PATCH /microservices/{uuid}/knowledge` | User microservices. Success HTTP 204. |
| `GetAllSystemMicroservices` | `GET /microservices/system` | |
| `GetSystemMicroservicesByApplication` | `GET /microservices/system?application=` | |
| `GetSystemMicroserviceByName` / `GetSystemMicroserviceByID` | system list or `GET /microservices/system/{uuid}` | |
| `UpdateSystemMicroserviceFromYAML` | `PATCH /microservices/system/yaml/{uuid}` | |
| `RebuildsSystemMicroservice` | `PATCH /microservices/system/{uuid}/rebuild` | |

Deprecated, all return `ErrRoutesNotSupported`: `CreateMicroserviceRoute`, `DeleteMicroserviceRoute`, `UpdateMicroserviceRoutes`.

`MicroserviceInfo.Commands` is JSON `cmd`. `commands` wins when both keys are present (`UnmarshalJSON`). Container limits: `cpus` (float), `memoryLimit`, `memoryReservation`, `memorySwap`, `shmSize` (int64), `cpuSetCpus`, `sysctls`, `ulimits`, `devices`, `tmpfs`, `entrypoint`, `workingDir`, `runAsUser`, `runAsGroup`, `readOnlyRootFilesystem`, `cdiDevices`, `runtime`, `pidMode`, `ipcMode`. `serviceAccount.roleRef` binds a Role.

`MicroserviceVolumeMappingInfo.Scope` is `private` or `shared`, and only meaningful when `type` is `volume`. Omit on create and the Controller stores `private`. GET always returns the key.

`MicroserviceStatusInfo` crash extras: `lastError` (survives recovery), `lastErrorAt` (Unix ms, 0 if empty), `restartCount` (crashes since the last operator rebuild; the first crash is 1), `podId` (pause/sandbox id on edgelet). Older Controllers omit them and the values stay zero.

## Microservice templates

No fog link API.

| Method | HTTP | Multipart field |
| --- | --- | --- |
| `ListMicroserviceTemplates` | `GET /microserviceTemplates` | |
| `GetMicroserviceTemplate(name)` | `GET /microserviceTemplates/{name}` | |
| `CreateMicroserviceTemplate(*MicroserviceTemplateCreateRequest)` | `POST /microserviceTemplates` | |
| `UpdateMicroserviceTemplate(name, *MicroserviceTemplateUpdateRequest)` | `PATCH /microserviceTemplates/{name}` | Name is immutable. |
| `DeleteMicroserviceTemplate(name)` | `DELETE /microserviceTemplates/{name}` | |
| `CreateMicroserviceTemplateFromYAML(io.Reader)` | `POST /microserviceTemplates/yaml` | `template` |
| `UpdateMicroserviceTemplateFromYAML(name, io.Reader)` | `PUT /microserviceTemplates/yaml/{name}` | `template` |

`MicroserviceTemplate` holds `name`, `description`, `variables` (`key`, `description`, `defaultValue`), and `microservice` as `map[string]any`.

## Models, Knowledge, RuntimeClasses

Shared link body `FogLinkRequest` / `FogLinkSet`: `fogUuids`. The resource name is only in the path. `Get*` does not return links.

| Action | Model | Knowledge | RuntimeClass |
| --- | --- | --- | --- |
| List | `GET /models` | `GET /knowledge` | `GET /runtimeClasses` |
| Get | `GET /models/{name}` | `GET /knowledge/{name}` | `GET /runtimeClasses/{name}` |
| Create | `POST /models` | `POST /knowledge` | `POST /runtimeClasses` |
| Update | `PATCH /models/{name}` | `PATCH /knowledge/{name}` | `PATCH /runtimeClasses/{name}` |
| Delete | `DELETE /models/{name}` | `DELETE /knowledge/{name}` (HTTP 202) | `DELETE /runtimeClasses/{name}` |
| YAML create | `POST /models/yaml` field `model` | `POST /knowledge/yaml` field `knowledge` | `POST /runtimeClasses/yaml` field `runtimeClass` |
| YAML upsert | `PUT /models/yaml/{name}` | `PUT /knowledge/yaml/{name}` | `PUT /runtimeClasses/yaml/{name}` |
| Get link | `GET .../{name}/link` | same | same |
| Link | `POST .../{name}/link` | same | same (HTTP 400 if the fog engine is not `edgelet`) |
| Unlink | `DELETE .../{name}/link` | same | same |

Delete and unlink return HTTP 409 when a microservice still references the name.

Model and Knowledge fields: `uuid`, `name`, `repo` (upstream path without host), `revision` (empty means latest for OCI, `main` for Hugging Face), `registryId`, `files` (Hugging Face only), `format`. Knowledge `format` is normalized to `unknown` when the Controller does not recognize it. RuntimeClass fields: `name`, `handler`. Update changes `handler` only.

## Catalog and registries

| Method | HTTP |
| --- | --- |
| `GetCatalog` | `GET /catalog/microservices` |
| `GetCatalogItem(id)` | `GET /catalog/microservices/{id}` |
| `GetCatalogItemByName(name)` | list + scan, `*NotFoundError` if missing |
| `CreateCatalogItem(*CatalogItemCreateRequest)` | `POST /catalog/microservices` |
| `UpdateCatalogItem(*CatalogItemUpdateRequest)` | `PATCH /catalog/microservices/{id}` |
| `DeleteCatalogItem(id)` | `DELETE /catalog/microservices/{id}` |

`CreateCatalogItem` sets `RegistryID` to `1` when it is `0`.

| Method | HTTP |
| --- | --- |
| `ListRegistries` | `GET /registries` |
| `GetRegistry(id)` | `GET /registries/{id}` |
| `CreateRegistry(*RegistryCreateRequest) (int, error)` | `POST /registries` returns the id |
| `UpdateRegistry(RegistryUpdateRequest)` | `PATCH /registries/{id}` (`ID` is not in the JSON body) |
| `DeleteRegistry(id)` | `DELETE /registries/{id}` |

`RegistryInfo.Type` is `oci` or `hf`. `CA` is an optional base64 PEM bundle (extra trust). `Insecure` allows HTTP and skips TLS verify. Image-pull `ca` / `insecure` apply on edgelet only. Removed: `isSecure`, `certificate`, `requiresCert`.

`RegistryTypeRegistryTypeIDDict` maps YAML aliases `remote` → 1 and `local` → 2. `RegistryTypeIDRegistryTypeDict` is the inverse. These aliases are not the registry `type` field.

`GetArchitectures() (*ArchitecturesListResponse, error)` is `GET /architectures/`. Integer codes are in [Codes](codes.md).

## Registries of config: secrets, services, config maps, volume mounts, certificates

List methods accept either `{ "<plural>": [ ... ] }` or a bare JSON array.

| Resource | Create | YAML field | Get | Update | YAML update field | Delete |
| --- | --- | --- | --- | --- | --- | --- |
| Secret | `POST /secrets` | `secret` | `GET /secrets/{name}` | `PATCH /secrets/{name}` | `secret` | `DELETE /secrets/{name}` |
| Service | `POST /services` | `service` | `GET /services/{name}` | `PATCH /services/{name}` | `service` | `DELETE /services/{name}` |
| ConfigMap | `POST /configmaps` | `configMap` | `GET /configmaps/{name}` | `PATCH /configmaps/{name}` | `configMap` | `DELETE /configmaps/{name}` |
| Volume mount | `POST /volumeMounts` | `volumeMount` | `GET /volumeMounts/{name}` | `PATCH /volumeMounts/{name}` | `volumeMount` | `DELETE /volumeMounts/{name}` |
| Certificate | `POST /certificates` | `certificate` | `GET /certificates/{name}` | | | `DELETE /certificates/{name}` |
| CA | `POST /certificates/ca` | | `GET /certificates/ca/{name}` | | | `DELETE /certificates/ca/{name}` |

Also: `ListSecrets`, `ListServices`, `ListConfigMaps`, `ListVolumeMounts`, `ListCertificates`, `ListCAs`, `ListExpiringCertificates` (`GET /certificates/expiring`), `RenewCertificate` (`POST /certificates/{name}/renew`).

`CreateService` and `UpdateService` provision the hub asynchronously. Poll `GetService` or `WaitForServiceProvisioningReady` (2s interval). `DeleteService` enqueues teardown. `ReconcileService` is `POST /services/{name}/reconcile`. When status is `failed`, the Controller resets it to `pending` and clears `provisioningError`.

`ProvisioningPending` (`pending`), `ProvisioningReady` (`ready`), `ProvisioningFailed` (`failed`). `ready` means the hub connector, listener, and Kubernetes Service finished. Edge fog bridge updates still converge through fog platform reconcile.

`ServiceUpdateRequest.DefaultBridge` is JSON `defaultBridgePort`. Create uses `defaultBridge`.

`ConfigMapCreateRequest` and `ConfigMapUpdateRequest` use `*bool` for `immutable` and `useVault` so an explicit `false` is encoded. `ConfigMapInfo` timestamps are `created_at` and `updated_at`. `useVault` is a non-pointer `bool` on the GET struct.

Volume mount link: `LinkVolumeMount` `POST /volumeMounts/{name}/link`, `UnlinkVolumeMount` `DELETE` on that path. Body includes `name` and `fogUuids`.

## RBAC

Roles: `ListRoles`, `GetRole`, `CreateRole`, `CreateRoleFromYaml` (field `role`), `UpdateRole`, `UpdateRoleFromYaml`, `DeleteRole` under `/roles`.

Role bindings: the same set under `/rolebindings` with YAML field `rolebinding`.

Service accounts are scoped by application name and account name.

| Method | HTTP |
| --- | --- |
| `ListServiceAccounts(applicationName)` | `GET /serviceaccounts` or `?applicationName=` when non-empty |
| `GetServiceAccount(app, name)` | `GET /serviceaccounts/{app}/{name}` |
| `CreateServiceAccount` | `POST /serviceaccounts` |
| `CreateServiceAccountFromYaml` | `POST /serviceaccounts/yaml` field `serviceaccount` |
| `UpdateServiceAccount(app, name, *ServiceAccountUpdateRequest)` | `PATCH /serviceaccounts/{app}/{name}` |
| `UpdateServiceAccountFromYaml(app, name, file)` | `PATCH /serviceaccounts/yaml/{app}/{name}` |
| `DeleteServiceAccount(app, name)` | `DELETE /serviceaccounts/{app}/{name}` |

`RBACRule`: `apiGroups`, `resources`, `verbs`, optional `resourceNames`. `RoleRef` and `Subject`: `kind`, `name`, optional `apiGroup`.

## Edge resources

Guarded by `HEAD /capabilities/edgeResources` and a non-empty access token. HTTP 404 becomes `*NotSupportedError` for `Edge Resources`.

| Method | HTTP |
| --- | --- |
| `IsEdgeResourceCapable` | `HEAD /capabilities/edgeResources` |
| `CreateHTTPEdgeResource(*EdgeResourceMetadata)` | `POST /edgeResource` |
| `GetHTTPEdgeResourceByName(name, version)` | `GET /edgeResource/{name}/{version}` |
| `ListEdgeResources` | `GET /edgeResources` |
| `UpdateHTTPEdgeResource(name, *EdgeResourceMetadata)` | `PUT /edgeResource/{name}/{version}` |
| `DeleteEdgeResource(name, version)` | `DELETE /edgeResource/{name}/{version}` |
| `LinkEdgeResource(LinkEdgeResourceRequest)` | `POST /edgeResource/{name}/{version}/link` |
| `UnlinkEdgeResource(LinkEdgeResourceRequest)` | `DELETE` on that path |

`LinkEdgeResourceRequest.AgentUUID` is JSON `uuid`. Name and version are path-only.

The interface payload is HTTP-specific (`HTTPEdgeResource.Endpoints`: `name`, `method`, `url`).

## Router

| Method | HTTP |
| --- | --- |
| `GetDefaultRouter() (Router, error)` | `GET /router` |
| `PutDefaultRouter(Router) error` | `PUT /router` |

`DefaultRouterName` is `default-router`. `DefaultNatsServerName` is `default-nats-hub`.

## Exec, logs, routes

Exec and log sessions: [Exec and logs](../howto/exec-and-logs.md).

Fog debug provisioning: `AttachExecToAgent` `POST /iofog/{uuid}/exec` (`uuid`, optional `image`), `DetachExecFromAgent` `DELETE /iofog/{uuid}/exec`.

Route methods (`ListRoutes`, `GetRoute`, `CreateRoute`, `UpdateRoute`, `PatchRoute`, `DeleteRoute`) return `ErrRoutesNotSupported`.

NATS methods: [NATS](nats.md).

## Related

- [Getting started](../getting-started.md)
- [Deploy YAML](../howto/deploy-yaml.md)
- [Fleet resources](../howto/fleet-resources.md)
