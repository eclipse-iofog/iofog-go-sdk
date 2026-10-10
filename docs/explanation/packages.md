# Packages

The module is a set of libraries that talk to two different control planes, plus helpers used by operators.

```text
your program
    |
    +-- pkg/apps          YAML deploy (create-or-update, then start)
    |       |
    |       +-- pkg/client    Controller REST  /api/v3
    |
    +-- pkg/client        direct REST, exec WebSocket, log WebSocket
    |
    +-- pkg/microservices EdgeletAPI v1 inside a running container
            GET /v1/microservices/config
            WSS /v1/microservices/control

pkg/arch          shared architecture integer codes
pkg/k8s           Kubernetes clientset + wait helpers
pkg/k8s/operator  controller-runtime requeue results
pkg/util          panic-on-bad-type assertions (not a product API)
```

## Controller client (`pkg/client`)

Use this when your program runs outside the edge container and needs the Controller: agents, applications, microservices, catalog, registries, secrets, RBAC, NATS, fleet models, Knowledge, RuntimeClasses.

`Options.BaseURL` must be an absolute URL whose path ends at `/api/v3`. Methods append a relative path such as `/iofog` or `/knowledge`. They send `Authorization: Bearer <access token>` on every request, including `GetStatus`.

## YAML deploy (`pkg/apps`)

Use this when you already have an application or microservice spec and want the same create-or-update behavior as `iofogctl`. The executors log in, marshal a header (`apiVersion`, `kind`, `metadata`, `spec`), and upload multipart YAML through `pkg/client`.

Default `apiVersion` is `iofog.org/v3` (`apps.DefaultAPIVersion`). Pass `apps.WithAPIVersion("datasance.com/v3")` only when the Controller expects the Datasance header. Do not hardcode that string in new Eclipse callers.

There is no `DeployModel`, `DeployKnowledge`, or `DeployRuntimeClass`. Those YAML documents are uploaded with client methods (`CreateModelFromYAML`, `CreateKnowledgeFromYAML`, `CreateRuntimeClassFromYAML`, and the matching upsert methods).

## Edge microservice SDK (`pkg/microservices`)

Use this only inside a container that Edgelet started. The client reads a mounted service-account token and CA, then calls EdgeletAPI v1 on port `54321`. It does not talk to the Controller, and it does not publish or subscribe to messages. Data-plane messaging is NATS.

## What this module does not include

- A CLI.
- Legacy LocalAPI v2 message bus types (`IoMessage` was removed).
- Controller route CRUD. Those methods return `client.ErrRoutesNotSupported`.
- HAL hardware and USB inventory endpoints.
- Fog config fields removed in v3.9: `deviceScanFrequency`, `bluetoothEnabled`, `abstractedHardwareEnabled`.

## Related

- [Security model](security.md)
- [Messaging](messaging.md)
- [Client reference](../reference/client.md)
