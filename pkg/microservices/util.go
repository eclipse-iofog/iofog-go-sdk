package microservices

import (
	"errors"
	"fmt"
	"net"
	"time"

	ws "github.com/gorilla/websocket"
)

func setCustomPingHandler(conn *ws.Conn) {
	conn.SetPingHandler(func(message string) error {
		if message == fmt.Sprint(ws.PingMessage) {
			message = fmt.Sprint(ws.PongMessage)
		}
		err := conn.WriteControl(ws.PongMessage, []byte(message), time.Now().Add(time.Second))
		if errors.Is(err, ws.ErrCloseSent) {
			return nil
		}
		var e net.Error
		if errors.As(err, &e) && e.Timeout() {
			return nil
		}
		return err
	})
}
