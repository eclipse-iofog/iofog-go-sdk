package client

import (
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

const defaultLogHandshakeTimeout = 45 * time.Second

// LogTailOptions configures log WebSocket query parameters.
type LogTailOptions struct {
	Tail   int
	Follow bool
	Since  string
	Until  string
}

// LogSession is a WebSocket log streaming session from the Controller.
type LogSession struct {
	SessionID string
	conn      *ws.Conn
	closeOnce sync.Once
	closeErr  error
}

// DialMicroserviceLogs opens a WebSocket log stream for an application microservice.
func (clt *Client) DialMicroserviceLogs(microserviceUUID string, opts *LogTailOptions) (*LogSession, error) {
	return clt.dialLogs(fmt.Sprintf("/microservices/%s/logs", microserviceUUID), opts)
}

// DialSystemMicroserviceLogs opens a WebSocket log stream for a system microservice.
func (clt *Client) DialSystemMicroserviceLogs(microserviceUUID string, opts *LogTailOptions) (*LogSession, error) {
	return clt.dialLogs(fmt.Sprintf("/microservices/system/%s/logs", microserviceUUID), opts)
}

// DialFogLogs opens a WebSocket log stream for a fog node (Agent).
func (clt *Client) DialFogLogs(fogUUID string, opts *LogTailOptions) (*LogSession, error) {
	return clt.dialLogs(fmt.Sprintf("/iofog/%s/logs", fogUUID), opts)
}

func (clt *Client) dialLogs(requestPath string, opts *LogTailOptions) (*LogSession, error) {
	if !clt.isLoggedIn() {
		return nil, NewInputError("client is not logged in")
	}

	wsURL, err := clt.logWSURL(requestPath, opts)
	if err != nil {
		return nil, err
	}

	dialer := ws.Dialer{
		HandshakeTimeout: defaultLogHandshakeTimeout,
		TLSClientConfig:  clt.execDialTLSConfig(),
	}

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+clt.accessToken)

	conn, resp, err := dialer.Dial(wsURL, headers)
	if err != nil {
		if resp != nil && resp.StatusCode != http.StatusSwitchingProtocols {
			return nil, NewHTTPError(fmt.Sprintf("log WebSocket upgrade failed with status %d", resp.StatusCode), resp.StatusCode)
		}
		return nil, fmt.Errorf("log WebSocket dial: %w", err)
	}

	return &LogSession{conn: conn}, nil
}

func (clt *Client) logWSURL(requestPath string, opts *LogTailOptions) (string, error) {
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
	if query := logTailQuery(opts); query != "" {
		u.RawQuery = query
	}
	return u.String(), nil
}

func logTailQuery(opts *LogTailOptions) string {
	values := url.Values{}
	tail := 100
	follow := true
	if opts != nil {
		if opts.Tail > 0 {
			tail = opts.Tail
		}
		follow = opts.Follow
		if opts.Since != "" {
			values.Set("since", opts.Since)
		}
		if opts.Until != "" {
			values.Set("until", opts.Until)
		}
	}
	values.Set("tail", fmt.Sprintf("%d", tail))
	values.Set("follow", fmt.Sprintf("%t", follow))
	return values.Encode()
}

// Read reads the next log frame from the Controller.
func (s *LogSession) Read() (*LogFrame, error) {
	if s.conn == nil {
		return nil, NewInternalError("log session is closed")
	}

	_, data, err := s.conn.ReadMessage()
	if err != nil {
		return nil, mapLogCloseError(err)
	}

	msg, err := decodeLogWSMessage(data)
	if err != nil {
		return nil, err
	}

	frame := logFrameFromMessage(msg)
	if frame.SessionID != "" {
		s.SessionID = frame.SessionID
	}
	return frame, nil
}

// Close closes the log WebSocket connection.
func (s *LogSession) Close() error {
	s.closeOnce.Do(func() {
		if s.conn == nil {
			return
		}
		err := s.conn.Close()
		s.conn = nil
		if isBenignCloseError(err) {
			return
		}
		s.closeErr = err
	})
	return s.closeErr
}

func mapLogCloseError(err error) error {
	closeErr := &ws.CloseError{}
	if !errors.As(err, &closeErr) {
		return err
	}

	reason := closeErr.Text
	switch closeErr.Code {
	case ws.ClosePolicyViolation: // 1008
		switch {
		case strings.Contains(reason, "No available log session"):
			return fmt.Errorf("%w: %s", ErrLogSessionUnavailable, reason)
		case strings.Contains(reason, "Authentication failed"):
			return fmt.Errorf("%w: %s", ErrLogAuthenticationFailed, reason)
		case strings.Contains(reason, "Agent is not running"), strings.Contains(reason, "not running"):
			return fmt.Errorf("%w: %s", ErrAgentNotRunning, reason)
		case strings.Contains(reason, "Microservice is not running"), strings.Contains(reason, "Not running"):
			return fmt.Errorf("%w: %s", ErrMicroserviceNotRunning, reason)
		case strings.Contains(reason, "Insufficient permissions"):
			return fmt.Errorf("%w: %s", ErrLogInsufficientPermissions, reason)
		}
		return fmt.Errorf("%w: %s", ErrLogPolicyViolation, reason)
	case ws.CloseAbnormalClosure: // 1006
		return fmt.Errorf("%w: %s", ErrLogConnectionLost, reason)
	case ws.CloseMessageTooBig: // 1009
		return fmt.Errorf("%w: %s", ErrLogMessageTooLarge, reason)
	case ws.CloseInternalServerErr: // 1011
		return fmt.Errorf("%w: %s", ErrLogServerError, reason)
	}
	return fmt.Errorf("log WebSocket closed (code %d): %s", closeErr.Code, reason)
}
