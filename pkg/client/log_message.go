package client

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

// Log message types for Controller WebSocket log sessions (MessagePack).
const (
	LogMessageLine  uint8 = 6
	LogMessageStart uint8 = 7
	LogMessageStop  uint8 = 8
	LogMessageError uint8 = 9
)

type logWSMessage struct {
	Type             uint8  `msgpack:"type"`
	Data             []byte `msgpack:"data"`
	MicroserviceUUID string `msgpack:"microserviceUuid"`
	ExecID           string `msgpack:"execId"`
	SessionID        string `msgpack:"sessionId"`
	Timestamp        int64  `msgpack:"timestamp"`
}

// LogFrame is a decoded log WebSocket frame from the Controller.
type LogFrame struct {
	Type      uint8
	Data      []byte
	SessionID string
}

func decodeLogWSMessage(data []byte) (*logWSMessage, error) {
	var msg logWSMessage
	if err := msgpack.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("decode log MessagePack frame: %w", err)
	}
	return &msg, nil
}

func logFrameFromMessage(msg *logWSMessage) *LogFrame {
	return &LogFrame{
		Type:      msg.Type,
		Data:      append([]byte(nil), msg.Data...),
		SessionID: firstNonEmpty(msg.SessionID, msg.ExecID),
	}
}
