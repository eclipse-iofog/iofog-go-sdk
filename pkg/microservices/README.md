# Microservices Package (EdgeletAPI v1)

This package is the Edgelet microservice SDK for EdgeletAPI v1.

It supports:
- reading microservice config over `GET /v1/microservices/config`
- receiving control signals over `GET /v1/microservices/control` WebSocket

It does not support legacy LocalAPI v2 messagebus APIs.
For data-plane messaging, use NATS.

## Runtime Assumptions

The running microservice has service-account material mounted by Edgelet:

- token: `/var/run/secrets/edgelet.iofog.org/serviceaccount/token`
- CA: `/var/run/secrets/edgelet.iofog.org/serviceaccount/ca.crt`

Default client behavior is HTTPS/WSS to:
- host: `edgelet.default.svc.bridge.local`
- port: `54321`

## Basic Usage

```go
import (
	msvcs "github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/microservices"
)

func run() error {
	client, err := msvcs.NewDefaultEdgeletAPIClient()
	if err != nil {
		return err
	}

	cfg, err := client.GetConfig()
	if err != nil {
		return err
	}
	_ = cfg

	controlCh := client.EstablishControlWsConnection(0)
	for range controlCh {
		updated, err := client.GetConfig()
		if err != nil {
			return err
		}
		_ = updated
	}
	return nil
}
```

## Advanced Configuration

```go
client, err := msvcs.NewEdgeletAPIClient(
	"microservice-id",
	msvcs.WithHost("edgelet.default.svc.bridge.local"),
	msvcs.WithPort(msvcs.PortEdgeletAPI),
	msvcs.WithTLS(true),
	msvcs.WithTokenPath("/var/run/secrets/edgelet.iofog.org/serviceaccount/token"),
	msvcs.WithCAPath("/var/run/secrets/edgelet.iofog.org/serviceaccount/ca.crt"),
)
```

## Migration Notes (iofog-agent v2 -> edgelet v1)

- Removed v2 endpoints (`/v2/...`) and messagebus methods.
- Removed `IoMessage`/`IoMessageReadable` types from the SDK surface.
- `GetConfig` now uses `GET /v1/microservices/config` and EdgeletAPI response envelope parsing.
- `EstablishControlWsConnection` now uses `/v1/microservices/control` with Bearer JWT auth.
- Token and CA are loaded from mounted service-account files.
- Use NATS for message exchange.
