# Exec and logs

You will open a WebSocket exec session and a log tail. Both dial the Controller directly. There is no REST attach call for microservice exec.

Import `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client`.

## Prerequisites

- A logged-in client (`accessToken` non-empty). Base URL scheme must be `http` or `https`. The dialer rewrites those to `ws` or `wss`.
- The microservice UUID, or the fog UUID for node logs.
- For exec, the microservice must be running. The session waits until stderr contains `Agent connected`, unless you disable that wait.

Frames are MessagePack.

## Steps

### 1. Exec into a microservice

```go
session, err := clt.DialMicroserviceExecWithOptions(microserviceUUID, &client.DialExecOptions{
	OnStatusLine: func(line string) {
		fmt.Fprintln(os.Stderr, line)
	},
})
if err != nil {
	return err
}
defer session.Close()

go func() {
	_, _ = session.WriteStdin([]byte("ls\n"))
}()

for {
	frame, err := session.Read()
	if err != nil {
		return err
	}
	switch frame.Type {
	case client.ExecMessageStdout:
		os.Stdout.Write(frame.Data)
	case client.ExecMessageStderr:
		os.Stderr.Write(frame.Data)
	case client.ExecMessageClose:
		return nil
	}
}
```

| Method | Path under `/api/v3` |
| --- | --- |
| `DialMicroserviceExec` / `DialMicroserviceExecWithOptions` | `/microservices/exec/{uuid}` |
| `DialSystemMicroserviceExec` / `DialSystemMicroserviceExecWithOptions` | `/microservices/system/exec/{uuid}` |

`WaitForAgentReady` defaults to true. Set the pointer to `false` to return as soon as the activation frame arrives.

`WriteStdin` sends type `0`. `WriteControl` sends type `3` (keepalive). `Close` sends type `4`, then closes the socket. `Close` is idempotent.

Exec types: `ExecMessageStdin` 0, `Stdout` 1, `Stderr` 2, `Control` 3, `Close` 4, `Activation` 5. `Read` does not return the activation frame; dial consumes it and fills `session.SessionID`.

### 2. Tail logs

```go
logs, err := clt.DialMicroserviceLogs(microserviceUUID, &client.LogTailOptions{
	Tail:   100,
	Follow: true,
})
if err != nil {
	return err
}
defer logs.Close()

for {
	frame, err := logs.Read()
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

| Method | Path |
| --- | --- |
| `DialMicroserviceLogs` | `/microservices/{uuid}/logs` |
| `DialSystemMicroserviceLogs` | `/microservices/system/{uuid}/logs` |
| `DialFogLogs` | `/iofog/{uuid}/logs` |

Nil options send `tail=100` and `follow=true`. A positive `Tail` replaces 100. `Follow` is sent as given when options are non-nil, so the zero value `false` turns follow off. `Since` and `Until` are optional query strings.

Log types: `LogMessageLine` 6, `Start` 7, `Stop` 8, `Error` 9. `LogSession.Close` is idempotent and does not send an application close frame.

### 3. Fog debug exec

Node debug still provisions a debug microservice, then dials that microservice:

```go
err := clt.AttachExecToAgent(&client.AttachExecToAgentRequest{UUID: fogUUID})
// DialSystemMicroserviceExec on the debug microservice UUID
err = clt.DetachExecFromAgent(&client.DetachExecFromAgentRequest{UUID: fogUUID})
```

`AttachExecMicroservice` and `DetachExecMicroservice` were removed. Do not call a REST enable step before `DialMicroserviceExec`.

## Verification

Exec: stderr shows `Agent connected`, then stdout frames carry command output. Logs: a `LogMessageStart` or `LogMessageLine` frame arrives, and `LogMessageStop` ends a non-follow stream.

## Troubleshooting

Close code 1008:

| Sentinel | When |
| --- | --- |
| `ErrExecSessionQuotaExceeded` | Reason contains `Maximum of 3 concurrent exec sessions`. |
| `ErrWsAgentTimeout` | Reason contains `Timeout waiting for agent connection`. |
| `ErrMicroserviceNotRunning` | Reason contains `not running`. |

Shared by exec and logs:

| Sentinel | Close code |
| --- | --- |
| `ErrWsRelayUnavailable` | 1013, cross-replica relay |
| `ErrWsServerDraining` | 1001 |
| `ErrWsAgentTimeout` | 1008 with the timeout reason |

Log-only sentinels include `ErrLogSessionUnavailable`, `ErrLogAuthenticationFailed`, `ErrAgentNotRunning`, `ErrLogInsufficientPermissions`, `ErrLogPolicyViolation`, `ErrLogConnectionLost`, `ErrLogMessageTooLarge`, and `ErrLogServerError`. Use `errors.Is`.

## Related

- [Client reference](../reference/client.md)
- [Errors](../reference/errors.md)
