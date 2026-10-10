# Security model

Two trust domains exist in this SDK. Mixing them produces clients that cannot authenticate.

## Controller REST

The Controller client stores an access token and a refresh token in memory on `Client`. Request helpers set:

```text
Authorization: Bearer <access token>
```

`isLoggedIn` is true when the access token string is non-empty. Several methods (`CreateAgent`, exec dial, log dial, edge resources, application templates, platform reconcile) refuse to run when that string is empty and return `client.Error` or `client.InputError`. Most other methods still send the request and surface the Controller's HTTP error.

### How tokens get onto the client

| Constructor or method | What it does |
| --- | --- |
| `New` | No token. Calls `GET /status` and caches the version. |
| `NewWithToken` | Stores the argument as the access token. |
| `NewAndLogin` | Calls `Login`. See the terminal behavior below. |
| `SessionLogin` | `POST /user/refresh` with the refresh token. On failure, falls back to `NewAndLogin`. |
| `NewWithRefreshToken` | Stores the argument with `SetAccessToken`. It does not call the refresh endpoint. Prefer `SessionLogin` or `Refresh` when you hold a refresh token. |
| `Login` | `POST /user/login`. On success, stores `accessToken` and `refreshToken` from the body. |
| `Refresh` | `POST /user/refresh`. Replaces both tokens from the response. |
| `SetAccessToken` / `SetRefreshToken` | Assign the strings directly. |

`Login` is interactive when fields are empty:

- Empty `Email` prompts `User E-mail:` on stdin.
- Empty, `"null"`, or `"NULL"` password prompts with asterisk masking.
- Empty `Totp` prompts `Enter OTP:`.

A wrong email, password, or TOTP (HTTP 400 or 401, or an error containing `failed to login`) retries up to three times and clears the fields so the prompts run again. Any other error returns immediately.

Non-interactive callers must either use `NewWithToken` or call `Login` with `Email`, `Password`, and `Totp` all set. If your Controller does not use TOTP, `Login` still prompts unless `Totp` is non-empty, and that value is sent in the JSON body.

`Logout` is `POST /user/logout`. It does not clear the in-memory tokens.

`Profile` replaces the client's access token with `WithTokenRequest.AccessToken`, then calls `GET /user/profile`. The shared request helper overwrites `Authorization` to `Bearer` plus the access token, so the profile call uses the same bearer header as every other method.

### TLS

If `Options.TLSConfig` is nil, HTTP and WebSocket dials use `InsecureSkipVerify: true`. That default exists because Controller installs often use a private CA. Production programs should pass a `*tls.Config` built from their trust store. Set `InsecureSkipVerify` only for an explicit insecure mode.

`Options.Timeout` is the HTTP client timeout in seconds. `0` becomes `10`.

WebSocket exec and log dials reuse `Options.TLSConfig` the same way. Handshake timeout is 45 seconds.

## EdgeletAPI

`pkg/microservices` does not use Controller credentials. Edgelet mounts:

| Material | Default path |
| --- | --- |
| Bearer JWT | `/var/run/secrets/edgelet.iofog.org/serviceaccount/token` |
| CA certificate | `/var/run/secrets/edgelet.iofog.org/serviceaccount/ca.crt` |

TLS defaults to on (`SSL` env, default `true`). The HTTP transport trusts only that CA and sets TLS 1.2 as the minimum. It does not skip verification.

The microservice id comes from `EDGELET_MICROSERVICE_UID`. `NewDefaultEdgeletAPIClient` fails if that variable is empty.

## YAML deploy credentials

`apps.IofogController` accepts `Token`, `RefreshToken`, `Email`, and `Password`.

- Non-empty `Token` uses `client.NewWithToken`.
- Otherwise the executor calls `client.SessionLogin` with the refresh token, then email and password.

`Endpoint` must be a full Controller base URL, including scheme and `/api/v3`. `url.Parse` rejects a bare `host:port`.

Application and microservice deploy parse `Endpoint` themselves. Application-template and microservice-template deploy take the parsed `*url.URL` as a separate argument.

## Related

- [Errors and retries](../reference/errors.md)
- [Use the Controller client](../howto/controller-client.md)
