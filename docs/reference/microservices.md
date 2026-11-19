# EdgeletAPI client

Package `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/microservices` talks to EdgeletAPI **v1** from inside a microservice container.

A step-by-step program is in [Write an edge microservice](../howto/edge-microservice.md).

## Types

```go
func NewEdgeletAPIClient(id string, opts ...ClientOption) (*EdgeletAPIClient, error)
func NewDefaultEdgeletAPIClient() (*EdgeletAPIClient, error)

func (c *EdgeletAPIClient) GetConfig() (map[string]any, error)
func (c *EdgeletAPIClient) GetConfigIntoStruct(config any) error
func (c *EdgeletAPIClient) EstablishControlWsConnection(signalBufSize int) <-chan byte
```

`NewEdgeletAPIClient` rejects an empty `id`. `NewDefaultEdgeletAPIClient` reads `EDGELET_MICROSERVICE_UID` and `SSL`.

There is no close method. The control-socket goroutine runs until the process exits. It reconnects forever.

## Options

| Option | Field | Default |
| --- | --- | --- |
| `WithHost` | `Host` | `edgelet.default.svc.bridge.local` (`HostDefault`) |
| `WithFallbackHosts` | `FallbackHosts` | `127.0.0.1` (`FallbackHostLocal`) |
| `WithPort` | `Port` | `54321` (`PortEdgeletAPI`). `<= 0` resets to the default. |
| `WithTLS` | `UseTLS` | `true` (`SSLDefault`) |
| `WithTokenPath` | `TokenPath` | `/var/run/secrets/edgelet.iofog.org/serviceaccount/token` |
| `WithCAPath` | `CAPath` | `/var/run/secrets/edgelet.iofog.org/serviceaccount/ca.crt` |
| `WithRequestTimeout` | `RequestTimeout` | 15s |
| `WithWSHandshakeTimeout` | `WSHandshakeTimeout` | 10s |
| `WithWSReconnectDelays` | base and max | 1s and 30s |

Empty host, token path, or CA path are replaced with the defaults. Fallback hosts are trimmed, de-duplicated, and the primary host is removed from the fallback list. If max delay is less than base delay, max is raised to base.

HTTP tries the primary host, then each fallback, when the error is a retriable network failure. The same order applies to the control WebSocket dial.

TLS uses the mounted CA as `RootCAs` and `MinVersion: tls.VersionTLS12`. Certificate verification is not skipped.

## Routes

| Constant | Value |
| --- | --- |
| `URLGetConfigV1` | `/v1/microservices/config` |
| `URLGetControlWSV1` | `/v1/microservices/control` |

Config responses are an envelope:

```json
{"success": true, "data": {"config": { }}}
```

`GetConfig` returns the inner `config` object. A JSON string is decoded. Non-success or non-2xx becomes `*EdgeletAPIError`. A 2xx body that is not an envelope is accepted as a raw object.

Control frames: the server sends `CodeControlSignal` (`0xC`). The client writes `CodeAck` (`0xB`) and forwards the byte on the channel. The send is blocking, so a full buffer stalls the socket reader until the consumer receives a byte. Buffer size `0` becomes `DefaultSignalBufferSize` (5). Reconnect delay doubles per attempt and caps at `WSReconnectMaxDelay`. Empty payloads and any first byte other than `0xC` are ignored.

## Errors

```go
type AuthMaterialError struct {
	Kind string
	Path string
	Err  error
}

type EdgeletAPIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    map[string]any
}
```

`AuthMaterialError` unwraps `Err`. `EdgeletAPIError.Error()` formats as `edgeletapi <code> (<status>): <message>` when `Code` is set.

## Environment

| Name | Constant |
| --- | --- |
| `EDGELET_MICROSERVICE_UID` | `MicroserviceUID` |
| `SSL` | `SSL` |

`SSL` is parsed with `strconv.ParseBool`. A missing or malformed value logs a warning and uses `true`.

## Related

- [Packages](../explanation/packages.md)
- [Messaging](../explanation/messaging.md)
