package container_sdk_go

import (
	"fmt"
	ws "github.com/gorilla/websocket"
	"time"
	"sync"
	"encoding/binary"
	"errors"
)

type ioFogWsClient struct {
	ssl                         bool
	host                        string
	port                        int
	url_base_ws                 string
	url_get_control_ws          string
	url_get_message_ws          string
	wsControl                   *ws.Conn
	wsMessage                   *ws.Conn
	wsControlAttempt            uint
	wsMessageAttempt            uint
	mutexMessageWs              sync.Mutex
}

func newIoFogWsClient(id string, ssl bool, host string, port int) *ioFogWsClient {
	client := ioFogWsClient {
		host:host,
		port:port,
		ssl:ssl,
	}
	protocol_ws := WS
	if client.ssl {
		protocol_ws = WSS
	}
	client.url_base_ws = fmt.Sprintf("%s://%s:%d", protocol_ws, client.host, client.port)
	client.url_get_control_ws = fmt.Sprint(client.url_base_ws, URL_GET_CONTROL_WS, id)
	client.url_get_message_ws = fmt.Sprint(client.url_base_ws, URL_GET_MESSAGE_WS, id)
	return &client
}

func (client *ioFogWsClient) sendMessage(msg *IoMessage) error {
	if client.wsMessage == nil {
		return errors.New("Socket is not initialized")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Println("Recovered after", r)
		}
	}()
	msgBytes, err := msg.EncodeBinary()
	if err != nil {
		return err
	}
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(msgBytes)))
	bytesToSend := make([]byte, 0, len(msgBytes) + 5)
	bytesToSend = append(bytesToSend, CODE_MSG)
	bytesToSend = append(bytesToSend, lengthBytes...)
	bytesToSend = append(bytesToSend, msgBytes...)
	defer client.mutexMessageWs.Unlock()
	client.mutexMessageWs.Lock()
	if err = client.wsMessage.WriteMessage(ws.BinaryMessage, bytesToSend); err != nil {
		logger.Println("Error while sending message:", err.Error())
	}
	return err
}

func (client *ioFogWsClient) establishControlWsConnection(signalChannel chan <- int) {
	for {
		conn, _, err := ws.DefaultDialer.Dial(client.url_get_control_ws, nil)
		if conn == nil {
			logger.Println(err.Error(), "Reconnecting...")
			sleepTime := 1 << client.wsControlAttempt * wsConnectTimeout
			if client.wsControlAttempt < wsAttemptLimit {
				client.wsControlAttempt++
			}
			time.Sleep(sleepTime)
		} else {
			client.wsControlAttempt = 0
			client.wsControl = conn
			setCustomPingHandler(client.wsControl)
			errChanel := make(chan int);
			go client.wsControlLoop(errChanel, signalChannel)
			loop:
			for {
				select {
				case <-errChanel:
					logger.Println("Reconnecting after control ws corruption")
					break loop
				}
			}
		}
	}
}

func (client *ioFogWsClient) establishMessageWsConnection(messageChannel chan <- *IoMessage, receiptChannel chan <- *PostMessageResponse) {
	for {
		conn, _, err := ws.DefaultDialer.Dial(client.url_get_message_ws, nil)
		if conn == nil {
			logger.Println(err.Error(), "Reconnecting...")
			sleepTime := 1 << client.wsMessageAttempt * wsConnectTimeout
			if client.wsMessageAttempt < wsAttemptLimit {
				client.wsMessageAttempt++
			}
			time.Sleep(sleepTime)
		} else {
			client.wsMessageAttempt = 0
			client.wsMessage = conn
			setCustomPingHandler(client.wsMessage)
			errChanel := make(chan int);
			go client.wsMessageLoop(errChanel, messageChannel, receiptChannel)
			loop:
			for {
				select {
				case <-errChanel:
					logger.Println("Reconnecting after message ws corruption")
					break loop
				}
			}
		}
	}
}

func (client*ioFogWsClient) wsControlLoop(errChanel chan <- int, dataChannel chan <- int) {
	defer func() {
		client.wsControl.Close()
		client.wsControl = nil
	}()
	for {
		_, p, err := client.wsControl.ReadMessage()
		if err != nil {
			logger.Println("Control ws error:", err.Error())
			errChanel <- 0
			return
		}
		if p[0] == CODE_CONTROL_SIGNAL {
			dataChannel <- int(p[0])
			err := client.wsControl.WriteMessage(ws.BinaryMessage, []byte{CODE_ACK})
			if err != nil {
				logger.Println("Error while sending acknowledgement:", err.Error())
			}
		}
	}
}

func (client*ioFogWsClient) wsMessageLoop(errChanel chan <- int, messageChannel chan <- *IoMessage, receiptChannel chan <-  *PostMessageResponse) {
	defer func() {
		client.wsMessage.Close()
		client.wsMessage = nil
	}()
	for {
		_, p, err := client.wsMessage.ReadMessage()
		if err != nil {
			logger.Println("Message ws error:", err.Error())
			errChanel <- 0
			return
		}
		if p[0] == CODE_MSG {
			msgLen := binary.BigEndian.Uint32(p[1:5])
			msgBytes := p[5: 5 + msgLen]
			msg := new(IoMessage)
			err = msg.DecodeBinary(msgBytes)
			if err != nil {
				logger.Println(err.Error())
			}
			messageChannel <- msg
			client.mutexMessageWs.Lock()
			err = client.wsMessage.WriteMessage(ws.BinaryMessage, []byte{CODE_ACK})
			client.mutexMessageWs.Unlock()
			if err != nil {
				logger.Println("Error while sending acknowledgement:", err.Error())
			}
		} else if p[0] == CODE_RECEIPT {
			idLen := int(p[1])
			tsLen := int(p[2])
			receiptResponse := new(PostMessageResponse)
			dataPos := 3
			if idLen != 0 {
				receiptResponse.ID = string(p[dataPos: dataPos + idLen])
				dataPos += idLen
			}
			if tsLen != 0 {
				receiptResponse.Timestamp = int(binary.BigEndian.Uint32(p[dataPos: dataPos + tsLen]))
			}
			receiptChannel <- receiptResponse
			client.mutexMessageWs.Lock()
			err = client.wsMessage.WriteMessage(ws.BinaryMessage, []byte{CODE_ACK})
			client.mutexMessageWs.Unlock()
			if err != nil {
				logger.Println("Error while sending acknowledgement:", err.Error())
			}
		}
	}
}