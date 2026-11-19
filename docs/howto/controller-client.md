# Use the Controller client

You will authenticate and list fog nodes. The client is `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client`.

## Prerequisites

- Controller base URL with scheme and `/api/v3`, for example `https://controller.example:51121/api/v3`.
- An access token, or email, password, and TOTP.

## Steps

### 1. Build `Options`

```go
base, err := url.Parse("https://controller.example:51121/api/v3")
if err != nil {
	return err
}

opts := client.Options{
	BaseURL: base,
	Timeout: 30, // seconds; 0 means 10
	TLSConfig: &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    pool, // your trust store
	},
}
```

Leave `TLSConfig` nil only when you accept skipped certificate verification. See [Security model](../explanation/security.md).

Optional retries:

```go
client.SetGlobalRetries(client.Retries{
	Timeout: 3, // extra attempts after HTTP 408
	CustomMessage: map[string]int{
		"connection refused": 5,
	},
})
```

`New` copies `GlobalRetriesPolicy` unless `Options.Retries` is set. Retries sleep `attempt` seconds between tries. Only HTTP 408 and errors whose text contains a `CustomMessage` key are retried.

`SetVerbosity(true)` prints request method, URL, JSON body, and response body to stdout. Do not enable it when bodies contain tokens or secrets.

### 2. Authenticate

Token (non-interactive):

```go
clt, err := client.NewWithToken(opts, accessToken)
```

Password login (interactive if `Totp` is empty):

```go
clt, err := client.NewAndLogin(opts, "user@example.com", "password")
```

Refresh token with password fallback:

```go
clt, err := client.SessionLogin(opts, refreshToken, "user@example.com", "password")
```

Direct login when you already have a TOTP and must not touch the terminal:

```go
clt := client.New(opts)
err := clt.Login(client.LoginRequest{
	Email:    "user@example.com",
	Password: "password",
	Totp:     totp,
})
```

### 3. Call an API

```go
agent, err := clt.GetAgentByName("edge-1")
if err != nil {
	return err
}
fmt.Println(agent.UUID, agent.DaemonStatus)
```

`GetAgentByName` lists every fog (`GET /iofog-list`) and returns the first name match, or `*client.NotFoundError`.

## Verification

`GetStatus` returns `ControllerStatus`. `GetVersion` returns the version string cached at `New` time (empty if that first status call failed). `GetVersionNumbers` parses `major.minor.patch` and ignores a `-suffix`.

## Troubleshooting

| Symptom | Cause |
| --- | --- |
| Process waits on `Enter OTP:` | `Login` / `NewAndLogin` saw an empty `Totp`. |
| `User input error` / `client is not logged in` | Exec, logs, or another guarded method ran with an empty access token. |
| `*client.NotFoundError` | HTTP 404. Message is `Unknown resource error` plus the body. |
| `*client.HTTPError` with `Code` | Any other non-2xx. `Error()` text starts with `Unexpected HTTP response`. |
| TLS handshake failure | `TLSConfig` does not trust the Controller certificate. The default config skips verify; a custom config does not. |

## Related

- [Client reference](../reference/client.md)
- [Errors](../reference/errors.md)
