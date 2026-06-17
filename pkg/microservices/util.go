package microservices

import (
	"errors"
	"fmt"
	"math"
	"net"
	"time"

	ws "github.com/gorilla/websocket"
)

func intToBytesBE(num int) ([]byte, int) {
	if num == 0 {
		return []byte{0}, 1
	}
	numOfBits := int(math.Log2(float64(num))) + 1
	numOfBytes := int(math.Ceil(float64(numOfBits) / 8.0))
	b := make([]byte, numOfBytes)
	shift := uint(8 * (numOfBytes - 1))
	for i := 0; i < numOfBytes; i++ {
		b[i] = byte((uint(num) >> shift) & 0xFF)
		shift -= 8
	}
	return b, numOfBytes
}

func int64ToBytesBE(num int64) ([]byte, int) {
	if num == 0 {
		return []byte{0}, 1
	}
	numOfBits := int(math.Log2(float64(num))) + 1
	numOfBytes := int(math.Ceil(float64(numOfBits) / 8.0))
	b := make([]byte, numOfBytes)
	shift := uint(8 * (numOfBytes - 1))
	for i := 0; i < numOfBytes; i++ {
		b[i] = byte((num >> int64(shift)) & 0xFF) // #nosec G115 -- minimal big-endian encoding of non-negative protocol integers
		shift -= 8
	}
	return b, numOfBytes
}

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
		if errors.As(err, &e) && e.Temporary() {
			return nil
		}
		return err
	})
}
