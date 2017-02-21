package container_sdk_go

import (
	"math"
	"io"
	"net/http"
	"errors"
	ws "github.com/gorilla/websocket"
	"time"
	"net"
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
		b[i] = byte(num >> shift)
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
		b[i] = byte(num >> shift)
		shift -= 8
	}

	return b, numOfBytes
}

func makePostRequest(url, bodyType string, body io.Reader) ([]byte, error) {
	resp, err := http.Post(url, bodyType, body)
	if err != nil {
		return nil, err
	}
	respBodyBytes := make([]byte, resp.ContentLength)
	resp.Body.Read(respBodyBytes)
	resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New(string(respBodyBytes))
	}
	return respBodyBytes, nil
}

func setCustomPingHandler(conn* ws.Conn) {
	conn.SetPingHandler(func(message string) error {
		if message == string(ws.PingMessage) {
			message = string(ws.PongMessage)
		}
		err := conn.WriteControl(ws.PongMessage, []byte(message), time.Now().Add(time.Second))
		if err == ws.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})
}