# Errors and retries

## Controller HTTP mapping

`pkg/client` treats status codes below 200 or at least 300 as errors.

| Status | Type | `Error()` prefix |
| --- | --- | --- |
| 404 | `*NotFoundError` | `Unknown resource error` |
| anything else | `*HTTPError` (`Code` is the status) | `Unexpected HTTP response` |

`errors.As` works for both. There is no special type for 409 or 401. Conflict responses are `*HTTPError` with `Code == 409`. Login treats 400 and 401 as a bad credential and retries the prompt.

| Constructor | Type | When the SDK uses it |
| --- | --- | --- |
| `NewError` | `*Error` | Logged-out guard on some agent and platform calls. Platform wait failures. |
| `NewInputError` | `*InputError` | Exec or log dial without a token or UUID. |
| `NewInternalError` | `*InternalError` | Unexpected decode or a closed session. |
| `NewNotSupportedError` | `*NotSupportedError` | Capability `HEAD` returned 404. |
| `NewConflictError` | `*ConflictError` | Available to callers. The HTTP layer does not return it. |
| `NewNotFoundError` | `*NotFoundError` | Name scans (`GetAgentByName`, catalog, microservice by name) and HTTP 404. |

`ErrRoutesNotSupported` is a plain `errors.New` value. Route methods return it directly.

### WebSocket sentinels

Use `errors.Is`. The formatted error wraps the sentinel and the close reason.

Shared:

| Variable | Close |
| --- | --- |
| `ErrWsRelayUnavailable` | 1013 |
| `ErrWsServerDraining` | 1001 |
| `ErrWsAgentTimeout` | 1008 and reason `Timeout waiting for agent connection` |

Exec (1008):

| Variable | Reason contains |
| --- | --- |
| `ErrExecSessionQuotaExceeded` | `Maximum of 3 concurrent exec sessions` |
| `ErrMicroserviceNotRunning` | `not running` |

Log: `ErrLogSessionUnavailable`, `ErrLogAuthenticationFailed`, `ErrAgentNotRunning`, `ErrLogInsufficientPermissions`, `ErrLogPolicyViolation`, `ErrLogConnectionLost`, `ErrLogMessageTooLarge`, `ErrLogServerError`.

`ExecSession.Close` and `LogSession.Close` ignore `use of closed network connection`, `connection reset by peer`, and `broken pipe`.

## Retries

`client.Retries.Timeout` counts extra attempts after HTTP 408 only. `CustomMessage` counts extra attempts when `err.Error()` contains the map key. The delay is `time.Duration(attempt) * time.Second` using the attempt number after the increment.

The zero `GlobalRetriesPolicy` does not retry. Set it with `SetGlobalRetries` before `New`, or pass `Options.Retries`.

Timeouts that are not HTTP 408 (client deadline, DNS) are not retried unless their error text hits `CustomMessage`.

## Apps

`pkg/apps` defines the same `Error`, `NotFoundError`, `ConflictError`, `InputError`, `InternalError`, and `HTTPError` types. Deploy does not translate Controller errors into the apps types. A missing application is detected with `errors.As` against `*client.NotFoundError`.

`apps.NewInputError` is returned for a bad `ParseFQMsvcName` and for `DeployMicroservice` with an empty application name.

## EdgeletAPI

| Type | Meaning |
| --- | --- |
| `*microservices.AuthMaterialError` | Token or CA file could not be read. `Unwrap` returns the file error. |
| `*microservices.EdgeletAPIError` | Envelope `success` is false, or the HTTP status is outside 2xx. |

Host fallback retries the next host only for retriable dial errors. An EdgeletAPI error body is returned immediately.

## `pkg/util`

`util.AssertInt` and `util.AssertBool` panic when the dynamic value has the wrong type. `AssertInt` accepts `int` and `float64` (truncated). Do not use them on Controller responses. Decode into structs instead.

## Related

- [Security model](../explanation/security.md)
- [Exec and logs](../howto/exec-and-logs.md)
