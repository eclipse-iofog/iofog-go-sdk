package client

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
	"github.com/vmihailenco/msgpack/v5"
)

// Exec message types for Controller WebSocket exec sessions (MessagePack).
const (
	ExecMessageStdin      uint8 = 0
	ExecMessageStdout     uint8 = 1
	ExecMessageStderr     uint8 = 2
	ExecMessageControl    uint8 = 3
	ExecMessageClose      uint8 = 4
	ExecMessageActivation uint8 = 5
)

const execAgentReadySubstr = "Agent connected"

type execWSMessage struct {
	Type             uint8  `msgpack:"type"`
	Data             []byte `msgpack:"data"`
	MicroserviceUUID string `msgpack:"microserviceUuid"`
	ExecID           string `msgpack:"execId"`
	SessionID        string `msgpack:"sessionId"`
	Timestamp        int64  `msgpack:"timestamp"`
}

// ExecFrame is a decoded exec WebSocket frame from the Controller.
type ExecFrame struct {
	Type             uint8
	Data             []byte
	MicroserviceUUID string
	SessionID        string
}

func decodeExecWSMessage(data []byte) (*execWSMessage, error) {
	var msg execWSMessage
	if err := msgpack.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("decode exec MessagePack frame: %w", err)
	}
	return &msg, nil
}

func encodeExecWSMessage(msg *execWSMessage) ([]byte, error) {
	return msgpack.Marshal(msg)
}

func parseExecSessionInfo(msg *execWSMessage) (ExecSessionInfo, error) {
	if msg.Type != ExecMessageActivation {
		return ExecSessionInfo{}, fmt.Errorf("expected ACTIVATION frame, got type %d", msg.Type)
	}

	info := ExecSessionInfo{
		SessionID:        firstNonEmpty(msg.SessionID, msg.ExecID),
		MicroserviceUUID: msg.MicroserviceUUID,
	}

	if len(msg.Data) > 0 {
		var fromData ExecSessionInfo
		if err := json.Unmarshal(msg.Data, &fromData); err != nil {
			return ExecSessionInfo{}, fmt.Errorf("decode ACTIVATION data JSON: %w", err)
		}
		if fromData.SessionID != "" {
			info.SessionID = fromData.SessionID
		}
		if fromData.MicroserviceUUID != "" {
			info.MicroserviceUUID = fromData.MicroserviceUUID
		}
	}

	if info.SessionID == "" {
		return ExecSessionInfo{}, errors.New("ACTIVATION frame missing sessionId")
	}
	if info.MicroserviceUUID == "" {
		return ExecSessionInfo{}, errors.New("ACTIVATION frame missing microserviceUuid")
	}
	return info, nil
}

func execFrameFromMessage(msg *execWSMessage) *ExecFrame {
	sessionID := firstNonEmpty(msg.SessionID, msg.ExecID)
	return &ExecFrame{
		Type:             msg.Type,
		Data:             append([]byte(nil), msg.Data...),
		MicroserviceUUID: msg.MicroserviceUUID,
		SessionID:        sessionID,
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
