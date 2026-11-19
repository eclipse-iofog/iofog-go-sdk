# Write an edge microservice

You will read the microservice config from Edgelet and refresh it when the control socket signals a change. This code runs inside the container Edgelet started. It does not run on your laptop against the Controller.

Import `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/microservices`.

## Prerequisites

Edgelet mounts:

- Token: `/var/run/secrets/edgelet.iofog.org/serviceaccount/token`
- CA: `/var/run/secrets/edgelet.iofog.org/serviceaccount/ca.crt`

Environment:

| Variable | Role | Default |
| --- | --- | --- |
| `EDGELET_MICROSERVICE_UID` | Microservice id. Required by `NewDefaultEdgeletAPIClient`. | none |
| `SSL` | `true` or `false` | `true` when unset or malformed |

Default dial target: `edgelet.default.svc.bridge.local:54321`, then fallback host `127.0.0.1`.

## Steps

### 1. Create the client

```go
client, err := msvcs.NewDefaultEdgeletAPIClient()
if err != nil {
	return err
}
```

Explicit options when the mount paths differ:

```go
client, err := msvcs.NewEdgeletAPIClient(microserviceID,
	msvcs.WithHost("edgelet.default.svc.bridge.local"),
	msvcs.WithFallbackHosts("127.0.0.1"),
	msvcs.WithPort(msvcs.PortEdgeletAPI),
	msvcs.WithTLS(true),
	msvcs.WithTokenPath("/var/run/secrets/edgelet.iofog.org/serviceaccount/token"),
	msvcs.WithCAPath("/var/run/secrets/edgelet.iofog.org/serviceaccount/ca.crt"),
)
```

An empty id returns `cannot create client with empty id`.

### 2. Read config

```go
cfg, err := client.GetConfig()
if err != nil {
	return err
}
```

`GetConfig` calls `GET /v1/microservices/config`, unwraps the EdgeletAPI envelope, and returns the `config` object as `map[string]any`. A string payload is JSON-decoded. Use `GetConfigIntoStruct` to decode into your own struct (JSON field names).

### 3. Refresh on control signals

```go
signals := client.EstablishControlWsConnection(0)
for range signals {
	cfg, err = client.GetConfig()
	if err != nil {
		return err
	}
}
```

`0` selects the default buffer of 5. The socket is `GET /v1/microservices/control` over WebSocket (`wss` when TLS is on) with the same bearer token. The reader loops forever and reconnects with exponential backoff from 1s up to 30s. Each signal is one byte. After you observe a signal, call `GetConfig` again.

This channel is not a message bus. Publish and subscribe with NATS. See [Messaging](../explanation/messaging.md).

## Verification

A successful `GetConfig` returns your microservice `config` object. Connection failures log to stderr (`Control ws connection has been established` after the socket connects). Auth or file problems return `*microservices.AuthMaterialError` (`Kind` is the material kind, `Path` is the file). HTTP error envelopes return `*microservices.EdgeletAPIError` with `StatusCode`, `Code`, `Message`, and `Details`.

## Troubleshooting

| Symptom | Cause |
| --- | --- |
| `EDGELET_MICROSERVICE_UID environment variable is not set` | `NewDefaultEdgeletAPIClient` outside an Edgelet container, or the variable was not injected. |
| `token material error` / `ca material error` | Mount path missing or unreadable. Override with `WithTokenPath` / `WithCAPath`. |
| `edgeletapi ...` | Non-success envelope. Read `Code` and `Message`. |
| `missing config payload in response` | HTTP 2xx body had no `config` field. |

## Related

- [Microservices reference](../reference/microservices.md)
