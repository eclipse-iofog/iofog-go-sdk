# NATS control plane

These methods live on `client.Client` under `/api/v3/nats`. They manage accounts, users, and rules. They do not publish or subscribe. Application data uses a NATS client with the creds this API returns.

`GetNatsBootstrap` returns HTTP 403 unless the Controller is running on the Kubernetes control plane. The response includes operator JWT, public key, and seed, plus system account JWT, public key, and system user creds (`sysUserCredsBase64`). Treat the seed and creds as secrets.

## Operator and hub

| Method | HTTP |
| --- | --- |
| `GetNatsOperator() (NatsOperatorResponse, error)` | `GET /nats/operator` |
| `RotateNatsOperator()` | `POST /nats/operator/rotate` |
| `GetNatsBootstrap() (NatsBootstrapResponse, error)` | `GET /nats/bootstrap` |
| `GetNatsHub() (NatsHubResponse, error)` | `GET /nats/hub` |
| `UpsertNatsHub(*NatsHubRequest)` | `PUT /nats/hub` |

`NatsOperatorResponse`: `name`, `publicKey`, `jwt`.

`NatsHubResponse`: `host`, `serverPort`, `clusterPort`, `leafPort`, `mqttPort`, `httpPort`. The request uses pointers so omitted fields stay out of the JSON.

## Accounts and users

Accounts are keyed by application name.

| Method | HTTP |
| --- | --- |
| `ListNatsAccounts` | `GET /nats/accounts` |
| `GetNatsAccount(appName)` | `GET /nats/accounts/{appName}` |
| `EnsureNatsAccount(appName, *NatsEnsureAccountRequest)` | `POST /nats/accounts/{appName}` |
| `ListNatsUsers` | `GET /nats/users` |
| `ListNatsAccountUsers(appName)` | `GET /nats/accounts/{appName}/users` |
| `CreateNatsUser(appName, *NatsCreateUserRequest)` | `POST /nats/accounts/{appName}/users` |
| `GetNatsUserCreds(appName, userName)` | `GET /nats/accounts/{appName}/users/{userName}/creds` |
| `DeleteNatsUser(appName, userName)` | `DELETE /nats/accounts/{appName}/users/{userName}` |
| `CreateNatsMqttBearer(appName, *NatsCreateMqttBearerRequest)` | `POST /nats/accounts/{appName}/mqtt-bearer` |
| `DeleteNatsMqttBearer(appName, userName)` | `DELETE /nats/accounts/{appName}/mqtt-bearer/{userName}` |

`NatsEnsureAccountRequest` and the user create bodies accept optional `natsRule`. User create also accepts `name` and optional `expiresIn`.

`NatsUserCredsResponse.CredsBase64` is the creds file, base64-encoded.

`NatsAccountInfo`: `id`, `name`, `publicKey`, `jwt`, `isSystem`, optional `isLeafSystem`, optional `applicationId`.

`NatsUserInfo`: `id`, `name`, `publicKey`, `jwt`, `isBearer`, and optional account, application, and `microserviceUuid`.

Application and microservice specs carry the same switch in YAML and JSON:

```yaml
natsConfig:
  natsAccess: true
  natsRule: default-rule
```

Types: `client.ApplicationNatsConfig`, `client.MicroserviceNatsConfig`, and the matching `apps` structs.

Fog-side mode is `AgentInfo.NatsMode` (`none`, `leaf`, `server`) plus port and storage fields. See [Client](client.md).

## Rules

Account rules and user rules share `NatsRuleInfo` on list and write responses. Create and update JSON bodies are `NatsAccountRulePayload` and `NatsUserRulePayload` (`map[string]any`). The Controller validates the schema. Numeric fields on `NatsRuleInfo` use `FlexInt` and `FlexInt64`, which accept a JSON number or a numeric string. `FlexStringSlice` accepts a JSON array or a string that itself contains a JSON array.

| Method | HTTP | YAML field |
| --- | --- | --- |
| `ListNatsAccountRules` | `GET /nats/account-rules` | |
| `CreateNatsAccountRule(payload)` | `POST /nats/account-rules` | |
| `CreateNatsAccountRuleFromYAML(io.Reader)` | `POST /nats/account-rules/yaml` | `natsAccountRule` |
| `UpdateNatsAccountRule(name, payload)` | `PATCH /nats/account-rules/{name}` | |
| `UpdateNatsAccountRuleFromYAML(name, file)` | `PATCH /nats/account-rules/yaml/{name}` | `natsAccountRule` |
| `DeleteNatsAccountRule(name)` | `DELETE /nats/account-rules/{name}` | |
| `ListNatsUserRules` | `GET /nats/user-rules` | |
| `CreateNatsUserRule` | `POST /nats/user-rules` | |
| `CreateNatsUserRuleFromYAML` | `POST /nats/user-rules/yaml` | `natsUserRule` |
| `UpdateNatsUserRule` | `PATCH /nats/user-rules/{name}` | |
| `UpdateNatsUserRuleFromYAML` | `PATCH /nats/user-rules/yaml/{name}` | `natsUserRule` |
| `DeleteNatsUserRule` | `DELETE /nats/user-rules/{name}` | |

`NatsRuleInfo` covers connection caps, import and export documents, JetStream storage, response permissions, pub/sub allow and deny, bearer and proxy flags, and time windows. Fields the Controller omits stay nil. Pass through unknown nested objects with the map payload instead of the struct when you need a field the struct types as `any` and you must resubmit it unchanged.

## Related

- [Messaging](../explanation/messaging.md)
- [Client](client.md)
