package client

import (
	"fmt"

	ws "github.com/gorilla/websocket"
)

// sharedWSCloseErr maps WebSocket close codes shared by exec and log sessions.
// Prefer close code over reason text; reason is included for operator diagnostics.
func sharedWSCloseErr(code int, reason string) error {
	switch code {
	case 1013:
		return fmt.Errorf("%w: %s", ErrWsRelayUnavailable, reason)
	case ws.CloseGoingAway: // 1001 — server draining during rollout
		return fmt.Errorf("%w: %s", ErrWsServerDraining, reason)
	default:
		return nil
	}
}
