package container_sdk_go

import (
	"fmt"
	ws "github.com/gorilla/websocket"
	"time"
	"errors"
)

type ioFogWsClient struct {
	url_base_ws         string
	url_get_control_ws  string
	url_get_message_ws  string
	wsControl           *ws.Conn
	wsMessage           *ws.Conn
	wsControlAttempt    uint
	wsMessageAttempt    uint
	writeMessageChannel chan []byte
}

func newIoFogWsClient(id string, ssl bool, host string, port int) *ioFogWsClient {
	client := ioFogWsClient{}
	protocol_ws := WS
	if ssl {
		protocol_ws = WSS
	}
	client.url_base_ws = fmt.Sprintf("%s://%s:%d", protocol_ws, host, port)
	client.url_get_control_ws = fmt.Sprint(client.url_base_ws, URL_GET_CONTROL_WS, id)
	client.url_get_message_ws = fmt.Sprint(client.url_base_ws, URL_GET_MESSAGE_WS, id)
	return &client
}

// TODO return error if socket corruption
// todo check nil client?
func (client *ioFogWsClient) sendMessage(msg *IoMessage) (e error) {
	if client.wsMessage == nil {
		return errors.New("Socket is not initialized")
	}
	bytesToSend, err := prepareMessageForSendingViaSocket(msg)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			logger.Println(r)
			e = errors.New("Error while sending message")
		}
	}()
	client.writeMessageChannel <- bytesToSend
	return nil
}

// TODO loop or goroutine?
func (client *ioFogWsClient) connectToControlWs(signalChannel chan <- byte) {
	for {
		conn, _, err := ws.DefaultDialer.Dial(client.url_get_control_ws, nil)
		if conn == nil {
			logger.Println(err.Error(), "Reconnecting to control ws...")
			sleepTime := 1 << client.wsControlAttempt * wsConnectTimeout
			if client.wsControlAttempt < wsAttemptLimit {
				client.wsControlAttempt++
			}
			time.Sleep(sleepTime)
		} else {
			client.wsControlAttempt = 0
			client.wsControl = conn
			setCustomPingHandler(client.wsControl)
			errChanel := make(chan byte, 2)
			writeChannel := make(chan []byte)
			go client.listenControlSocket(errChanel, signalChannel, writeChannel)
			go client.writeControlSocket(errChanel, writeChannel)
			loop:
			for {
				select {
				case <-errChanel:
					logger.Println("Reconnecting after control ws corruption")
					client.wsControl.Close()
					break loop
				}
			}
		}
	}
}

// TODO loop or goroutine?
func (client *ioFogWsClient) connectToMessageWs(messageChannel chan <- *IoMessage, receiptChannel chan <- *PostMessageResponse) {
	for {
		conn, _, err := ws.DefaultDialer.Dial(client.url_get_message_ws, nil)
		if conn == nil {
			logger.Println(err.Error(), "Reconnecting to message ws...")
			sleepTime := 1 << client.wsMessageAttempt * wsConnectTimeout
			if client.wsMessageAttempt < wsAttemptLimit {
				client.wsMessageAttempt++
			}
			time.Sleep(sleepTime)
		} else {
			client.wsMessageAttempt = 0
			client.wsMessage = conn
			setCustomPingHandler(client.wsMessage)
			errChanel := make(chan byte, 2)
			writeChannel := make(chan []byte, 5)
			client.writeMessageChannel = writeChannel
			go client.listenMessageSocket(errChanel, messageChannel, receiptChannel, writeChannel)
			go client.writeMessageSocket(errChanel, writeChannel)
			loop:
			for {
				select {
				case <-errChanel:
					logger.Println("Reconnecting after message ws corruption")
					client.wsMessage.Close()
					break loop
				}
			}
		}
	}
}

func (client*ioFogWsClient) listenControlSocket(errChanel chan <- byte, signalChannel chan <- byte, writeChannel chan <- []byte) {
	for {
		_, p, err := client.wsControl.ReadMessage()
		if err != nil {
			logger.Println("Control ws read error:", err.Error())
			errChanel <- 0
			close(writeChannel)
			return
		}
		if p[0] == CODE_CONTROL_SIGNAL {
			signalChannel <- p[0]
			writeChannel <- []byte{CODE_ACK}
		}
	}
}

func (client*ioFogWsClient) writeControlSocket(errChanel chan <- byte, writeChannel <- chan []byte) {
	for data := range writeChannel {
		err := client.wsControl.WriteMessage(ws.BinaryMessage, data)
		if err != nil {
			logger.Println("Control ws write error:", err.Error())
			errChanel <- 0
			return
		}
	}
}

func (client*ioFogWsClient) listenMessageSocket(errChanel chan <- byte, messageChannel chan <- *IoMessage, receiptChannel chan <-  *PostMessageResponse, writeChannel chan <- []byte) {
	for {
		_, p, err := client.wsMessage.ReadMessage()
		if err != nil {
			logger.Println("Message ws read error:", err.Error())
			errChanel <- 0
			close(writeChannel)
			return
		}
		if p[0] == CODE_MSG {
			msg, err := getMessageReceivedViaSocket(p)
			if err != nil {
				logger.Println(err.Error())
			}
			messageChannel <- msg
			writeChannel <- []byte{CODE_ACK}
		} else if p[0] == CODE_RECEIPT {
			receiptResponse, err := getReceiptReceivedViaSocket(p)
			if err != nil {
				logger.Println(err.Error())
			}
			receiptChannel <- receiptResponse
			writeChannel <- []byte{CODE_ACK}
		}
	}
}

func (client*ioFogWsClient) writeMessageSocket(errChanel chan <- byte, writeChannel <- chan []byte) {
	for data := range writeChannel {
		err := client.wsMessage.WriteMessage(ws.BinaryMessage, data)
		if err != nil {
			logger.Println("Message ws write error:", err.Error())
			errChanel <- 0
			return
		}
	}
}