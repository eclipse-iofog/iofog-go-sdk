package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
)

const defaultExecHandshakeTimeout = 45 * time.Second

// ExecSessionInfo is parsed from the ACTIVATION frame after a successful exec WebSocket upgrade.
type ExecSessionInfo struct {
	SessionID        string `json:"sessionId"`
	MicroserviceUUID string `json:"microserviceUuid"`
}

// DialExecOptions configures exec WebSocket dial behavior.
type DialExecOptions struct {
	// WaitForAgentReady blocks until STDERR indicates the agent is connected.
	// Defaults to true when nil.
	WaitForAgentReady *bool
	// OnStatusLine is invoked for each STDERR status line emitted while waiting
	// for the agent (for example "Waiting for agent connection...").
	OnStatusLine func(string)
}

// ExecSession is an interactive exec WebSocket session to a microservice container.
type ExecSession struct {
	SessionID        string
	MicroserviceUUID string
	conn             *ws.Conn
	closeOnce        sync.Once
	closeErr         error
}

// DialMicroserviceExec opens a WebSocket exec session to an application microservice.
func (clt *Client) DialMicroserviceExec(microserviceUUID string) (*ExecSession, error) {
	return clt.dialExec(microserviceUUID, fmt.Sprintf("/microservices/exec/%s", microserviceUUID), nil)
}

// DialSystemMicroserviceExec opens a WebSocket exec session to a system microservice.
func (clt *Client) DialSystemMicroserviceExec(microserviceUUID string) (*ExecSession, error) {
	return clt.dialExec(microserviceUUID, fmt.Sprintf("/microservices/system/exec/%s", microserviceUUID), nil)
}

// DialMicroserviceExecWithOptions opens a WebSocket exec session with dial options.
func (clt *Client) DialMicroserviceExecWithOptions(microserviceUUID string, opts *DialExecOptions) (*ExecSession, error) {
	return clt.dialExec(microserviceUUID, fmt.Sprintf("/microservices/exec/%s", microserviceUUID), opts)
}

// DialSystemMicroserviceExecWithOptions opens a system microservice exec session with dial options.
func (clt *Client) DialSystemMicroserviceExecWithOptions(microserviceUUID string, opts *DialExecOptions) (*ExecSession, error) {
	return clt.dialExec(microserviceUUID, fmt.Sprintf("/microservices/system/exec/%s", microserviceUUID), opts)
}

func (clt *Client) dialExec(microserviceUUID, requestPath string, opts *DialExecOptions) (*ExecSession, error) {
	if !clt.isLoggedIn() {
		return nil, NewInputError("client is not logged in")
	}
	if microserviceUUID == "" {
		return nil, NewInputError("microservice UUID is required")
	}

	wsURL, err := clt.execWSURL(requestPath)
	if err != nil {
		return nil, err
	}

	dialer := ws.Dialer{
		HandshakeTimeout: defaultExecHandshakeTimeout,
		TLSClientConfig:  clt.execDialTLSConfig(),
	}

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+clt.accessToken)

	conn, resp, err := dialer.Dial(wsURL, headers)
	if err != nil {
		if resp != nil && resp.StatusCode != http.StatusSwitchingProtocols {
			return nil, NewHTTPError(fmt.Sprintf("exec WebSocket upgrade failed with status %d", resp.StatusCode), resp.StatusCode)
		}
		return nil, fmt.Errorf("exec WebSocket dial: %w", err)
	}

	session := &ExecSession{
		MicroserviceUUID: microserviceUUID,
		conn:             conn,
	}

	info, err := session.readActivation()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	session.SessionID = info.SessionID
	session.MicroserviceUUID = info.MicroserviceUUID

	waitForAgent := true
	if opts != nil && opts.WaitForAgentReady != nil {
		waitForAgent = *opts.WaitForAgentReady
	}
	var onStatusLine func(string)
	if opts != nil {
		onStatusLine = opts.OnStatusLine
	}
	if waitForAgent {
		if err := session.waitForAgentReady(onStatusLine); err != nil {
			_ = session.Close()
			return nil, err
		}
	}

	return session, nil
}

func (clt *Client) execWSURL(requestPath string) (string, error) {
	u, err := url.Parse(clt.baseURL.String())
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	default:
		return "", fmt.Errorf("unsupported base URL scheme %q for WebSocket", u.Scheme)
	}
	u.Path = path.Join(u.Path, strings.TrimPrefix(requestPath, "/"))
	return u.String(), nil
}

func (clt *Client) execDialTLSConfig() *tls.Config {
	if clt.tlsConfig != nil {
		return clt.tlsConfig
	}
	return defaultTLSConfig()
}

// Read reads the next exec frame from the Controller.
func (s *ExecSession) Read() (*ExecFrame, error) {
	if s.conn == nil {
		return nil, NewInternalError("exec session is closed")
	}

	_, data, err := s.conn.ReadMessage()
	if err != nil {
		return nil, mapExecCloseError(err)
	}

	msg, err := decodeExecWSMessage(data)
	if err != nil {
		return nil, err
	}
	return execFrameFromMessage(msg), nil
}

// WriteStdin sends stdin data to the remote container.
func (s *ExecSession) WriteStdin(data []byte) error {
	return s.writeFrame(ExecMessageStdin, data)
}

// WriteControl sends a keepalive control frame.
func (s *ExecSession) WriteControl(data []byte) error {
	return s.writeFrame(ExecMessageControl, data)
}

// Close sends a close frame and closes the WebSocket connection.
// Close is idempotent and ignores benign errors from an already-closed connection.
func (s *ExecSession) Close() error {
	s.closeOnce.Do(func() {
		if s.conn == nil {
			return
		}
		_ = s.writeFrame(ExecMessageClose, nil)
		err := s.conn.Close()
		s.conn = nil
		if isBenignCloseError(err) {
			return
		}
		s.closeErr = err
	})
	return s.closeErr
}

func (s *ExecSession) writeFrame(msgType uint8, data []byte) error {
	if s.conn == nil {
		return NewInternalError("exec session is closed")
	}
	payload, err := encodeExecWSMessage(&execWSMessage{
		Type:             msgType,
		Data:             data,
		MicroserviceUUID: s.MicroserviceUUID,
		ExecID:           s.SessionID,
		SessionID:        s.SessionID,
		Timestamp:        time.Now().UnixMilli(),
	})
	if err != nil {
		return err
	}
	return s.conn.WriteMessage(ws.BinaryMessage, payload)
}

func (s *ExecSession) readActivation() (ExecSessionInfo, error) {
	frame, err := s.readRawFrame()
	if err != nil {
		return ExecSessionInfo{}, err
	}
	return parseExecSessionInfo(frame)
}

func (s *ExecSession) waitForAgentReady(onStatusLine func(string)) error {
	for {
		frame, err := s.readRawFrame()
		if err != nil {
			return err
		}
		if frame.Type != ExecMessageStderr {
			continue
		}
		line := string(frame.Data)
		if onStatusLine != nil && line != "" {
			onStatusLine(line)
		}
		if strings.Contains(line, execAgentReadySubstr) {
			return nil
		}
	}
}

func isBenignCloseError(err error) bool {
	if err == nil {
		return true
	}
	errStr := err.Error()
	return strings.Contains(errStr, "use of closed network connection") ||
		strings.Contains(errStr, "connection reset by peer") ||
		strings.Contains(errStr, "broken pipe")
}

func (s *ExecSession) readRawFrame() (*execWSMessage, error) {
	if s.conn == nil {
		return nil, NewInternalError("exec session is closed")
	}

	_, data, err := s.conn.ReadMessage()
	if err != nil {
		return nil, mapExecCloseError(err)
	}
	return decodeExecWSMessage(data)
}

func mapExecCloseError(err error) error {
	closeErr := &ws.CloseError{}
	if !errors.As(err, &closeErr) {
		return err
	}

	reason := closeErr.Text
	if closeErr.Code == ws.ClosePolicyViolation { // 1008
		switch {
		case strings.Contains(reason, "Maximum of 3 concurrent exec sessions"):
			return fmt.Errorf("%w: %s", ErrExecSessionQuotaExceeded, reason)
		case strings.Contains(reason, "Timeout waiting for agent connection"):
			return fmt.Errorf("%w: %s", ErrWsAgentTimeout, reason)
		case strings.Contains(reason, "not running"), strings.Contains(reason, "Not running"):
			return fmt.Errorf("%w: %s", ErrMicroserviceNotRunning, reason)
		}
	}
	if err := sharedWSCloseErr(closeErr.Code, reason); err != nil {
		return err
	}
	return fmt.Errorf("exec WebSocket closed (code %d): %s", closeErr.Code, reason)
}
