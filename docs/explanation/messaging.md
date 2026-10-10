# Messaging

ioFog moved the data plane to NATS. This SDK follows that split: the Controller client manages NATS accounts and rules, and the edge SDK does not carry messages.

## What was removed

Controller no longer exposes `/api/v3/routes` or `/api/v3/microservices/:uuid/routes`. These methods return `client.ErrRoutesNotSupported` and do not open a connection:

- `ListRoutes`, `GetRoute`, `CreateRoute`, `UpdateRoute`, `PatchRoute`, `DeleteRoute`
- `CreateMicroserviceRoute`, `DeleteMicroserviceRoute`, `UpdateMicroserviceRoutes`

`apps.Route` and `apps.RouteKind` remain so old YAML still unmarshals. `validateRoutes` and `createRoutes` are deprecated and unused by deploy. `createRoutes` returns `ErrRoutesNotSupported`.

`pkg/microservices` has no message-bus publish or subscribe API and no `IoMessage` type. A microservice that needs to exchange data opens a NATS connection with credentials the Controller issued.

## What you use instead

| Need | API |
| --- | --- |
| Turn NATS on for an application | `ApplicationNatsConfig` (`natsAccess`, optional `natsRule`) on create/patch, or `EnsureNatsAccount` |
| Issue a user and download creds | `CreateNatsUser`, `GetNatsUserCreds` (`credsBase64`) |
| MQTT bearer user | `CreateNatsMqttBearer`, `DeleteNatsMqttBearer` |
| Account and user limit documents | `CreateNatsAccountRule` / `CreateNatsUserRule` (JSON or YAML) |
| Operator bootstrap on Kubernetes | `GetNatsBootstrap` (HTTP 403 when the Controller is not on the Kubernetes control plane) |

Fog nodes report NATS mode on `AgentInfo` (`natsMode`: `none`, `leaf`, or `server`, plus ports). Agent configuration uses the embedded `NatsConfig` struct.

## Edge control signals are not the data plane

`EstablishControlWsConnection` receives a single byte from Edgelet when the microservice config should be re-read (`CodeControlSignal` is `0xC`). The client answers with `CodeAck` (`0xB`). That channel is not a message bus.

## Related

- [NATS reference](../reference/nats.md)
- [Edge microservice guide](../howto/edge-microservice.md)
