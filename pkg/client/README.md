# Client Package

This package provides an HTTP client for ioFog Controller's REST API **v3** (`/api/v3/*`).

Import path: `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client`

You can view the full REST API specification at [iofog.org](https://iofog.org/docs/1.3.0/controllers/rest-api.html).

## Usage

Create a client with the Controller base URL (including `/api/v3`), then log in with your credentials.

```go
import (
	"net/url"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
)

baseURL, err := url.Parse("http://localhost:51121/api/v3")
if err != nil {
	return err
}

clt, err := client.NewAndLogin(client.Options{BaseURL: baseURL}, "user@domain.com", "password")
if err != nil {
	return err
}
```

Call any method on the authenticated client:

```go
resp, err := clt.GetStatus()
if err != nil {
	return err
}
println(resp.Status)
```

## HTTPS / TLS

Optional `Options.TLSConfig` for HTTPS Controller endpoints. If omitted, the SDK skips certificate verification (self-signed friendly). CLI tools should set this from their trust store; use `InsecureSkipVerify: true` only in explicit insecure mode.

## Microservice exec (WebSocket)

Interactive exec is WebSocket-first. Dial directly — no REST attach step:

```go
session, err := clt.DialMicroserviceExecWithOptions(microserviceUUID, &client.DialExecOptions{
	OnStatusLine: func(line string) {
		fmt.Fprintln(os.Stdout, line)
	},
})
if err != nil {
	return err
}
defer session.Close()

for {
	frame, err := session.Read()
	if err != nil {
		return err
	}
	if frame.Type == client.ExecMessageStdout {
		os.Stdout.Write(frame.Data)
	}
}
```

`OnStatusLine` receives STDERR status lines emitted while waiting for the agent (for example `Waiting for agent connection...`). `ExecSession.Close()` is idempotent.

System microservices use `DialSystemMicroserviceExec` / `DialSystemMicroserviceExecWithOptions`.

Fog node debug still provisions a debug microservice via `AttachExecToAgent`, then uses `DialSystemMicroserviceExec` on the debug microservice UUID.

## Log streaming (WebSocket)

Remote log tailing uses WebSocket sessions with optional query parameters:

```go
session, err := clt.DialMicroserviceLogs(microserviceUUID, &client.LogTailOptions{
	Tail:   100,
	Follow: true,
})
if err != nil {
	return err
}
defer session.Close()

for {
	frame, err := session.Read()
	if err != nil {
		return err
	}
	switch frame.Type {
	case client.LogMessageLine:
		os.Stdout.Write(frame.Data)
	case client.LogMessageStop:
		return nil
	case client.LogMessageError:
		return fmt.Errorf("%s", frame.Data)
	}
}
```

Use `DialSystemMicroserviceLogs` for system microservices and `DialFogLogs` for Agent (fog node) logs.
