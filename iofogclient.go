package container_sdk_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"os/exec"
	ws "github.com/gorilla/websocket"
	"time"
	"sync"
	"encoding/binary"
	"log"
)

var (
	//f, _ = os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE, 0777)
	logger = log.New(os.Stderr, "", log.LstdFlags)
)

type ioFogClient struct {
	id                          string
	ssl                         bool
	host                        string
	port                        int
	url_base_rest               string
	url_base_ws                 string
	url_get_config              string
	url_get_next_messages       string
	url_get_publishers_messages string
	url_post_message            string
	url_get_control_ws          string
	url_get_message_ws          string
	requestBodyId               []byte
	wsControl                   *ws.Conn
	wsMessage                   *ws.Conn
	wsControlAttempt            uint
	wsMessageAttempt            uint
	mutexMessageWs              sync.Mutex
}

func (client *ioFogClient) initClient() {
	protocol_rest := HTTP
	protocol_ws := WS
	if client.ssl {
		protocol_rest = HTTPS
		protocol_ws = WSS
	}
	client.url_base_rest = fmt.Sprintf("%s://%s:%d", protocol_rest, client.host, client.port)
	client.url_base_ws = fmt.Sprintf("%s://%s:%d", protocol_ws, client.host, client.port)
	client.url_get_config = fmt.Sprint(client.url_base_rest, URL_GET_CONFIG)
	client.url_get_next_messages = fmt.Sprint(client.url_base_rest, URL_GET_NEXT_MESSAGES)
	client.url_get_publishers_messages = fmt.Sprint(client.url_base_rest, URL_GET_PUBLISHERS_MESSAGES)
	client.url_post_message = fmt.Sprint(client.url_base_rest, URL_POST_MESSAGE)
	client.url_get_control_ws = fmt.Sprint(client.url_base_ws, URL_GET_CONTROL_WS, client.id)
	client.url_get_message_ws = fmt.Sprint(client.url_base_ws, URL_GET_MESSAGE_WS, client.id)
	client.requestBodyId, _ = json.Marshal(map[string]interface{}{
		ID : client.id,
	})
}

func NewIoFogClient(host string, port int, ssl bool, id string) *ioFogClient {
	if id == "" {
		logger.Print("Id is empty. IoFog client is not created")
		return nil
	}
	client := ioFogClient{id: id, ssl: ssl}
	client.initClient()
	return &client
}

func NewDefaultIoFogClient() *ioFogClient {
	selfname := os.Getenv(SELFNAME)
	if selfname == "" {
		logger.Println("Empty ", SELFNAME, " environment virable. IoFog client is not created")
		return nil
	}
	ssl, err := strconv.ParseBool(os.Getenv(SSL))
	if err != nil {
		logger.Println("Empty or malformed ", SSL, " environment variable. Using default value of ", SSL_DEFAULT)
		ssl = SSL_DEFAULT
	}

	host := IOFOG
	if cmd := exec.Command("ping", "-c 3", host); cmd.Run() != nil {
		logger.Println("Host ", host, " is unreachable. Switching to ", HOST_DEFAULT)
		host = HOST_DEFAULT
	}

	client := ioFogClient{id: selfname, ssl: ssl, host: host, port: PORT_IOFOG}
	client.initClient()
	return &client
}

func (client *ioFogClient) GetConfig() (map[string]interface{}, error) {
	resp, err := makePostRequest(client.url_get_config, APPLICATION_JSON, bytes.NewBuffer(client.requestBodyId))
	if err != nil {
		return nil, err
	}
	configResponse := new(getConfigResponse)
	config := make(map[string]interface{})
	if err := json.Unmarshal(resp, configResponse); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(configResponse.Config), &config); err != nil {
		return nil, err
	}
	return config, nil
}

func (client *ioFogClient) GetNextMessages(decodeBase64 bool) ([]IoMessage, error) {
	resp, err := makePostRequest(client.url_get_next_messages, APPLICATION_JSON, bytes.NewBuffer(client.requestBodyId))
	if err != nil {
		return nil, err
	}
	nextMessagesResponse := new(getNextMessagesResponse)
	if err := json.Unmarshal(resp, nextMessagesResponse); err != nil {
		return nil, err
	}
	if decodeBase64 {
		for i := range nextMessagesResponse.Messages {
			nextMessagesResponse.Messages[i].DecodeData()
		}
	}
	return nextMessagesResponse.Messages, nil
}

func (client *ioFogClient) PostMessage(msg *IoMessage) (*PostMessageResponse, error) {
	msg.Publisher = client.id
	msg.Version = IOMESSAGE_VERSION
	requestBytes, _ := json.Marshal(msg)
	resp, err := makePostRequest(client.url_post_message, APPLICATION_JSON, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	postMessageResponse := new(PostMessageResponse)
	if err := json.Unmarshal(resp, postMessageResponse); err != nil {
		return nil, err
	}
	return postMessageResponse, nil
}

func (client *ioFogClient) GetMessagesFromPublishersWithinTimeFrame(query *MessagesQueryParameters, decodeBase64 bool) (*TimeFrameMessages, error) {
	query.ID = client.id
	requestBytes, _ := json.Marshal(query)
	resp, err := makePostRequest(client.url_get_publishers_messages, APPLICATION_JSON, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	nextMessagesResponse := new(getNextMessagesResponse)
	if err := json.Unmarshal(resp, nextMessagesResponse); err != nil {
		return nil, err
	}
	if decodeBase64 {
		for i := range nextMessagesResponse.Messages {
			nextMessagesResponse.Messages[i].DecodeData()
		}
	}
	return (*TimeFrameMessages)(nextMessagesResponse), nil
}

func (client *ioFogClient) EstablishControlWsConnection(signalChannel chan <- int) {
	go client.establishControlWsConnection(signalChannel)
}

func (client *ioFogClient) EstablishMessageWsConnection(messageChannel chan <- *IoMessage, receiptChannel chan <- *PostMessageResponse) {
	go client.establishMessageWsConnection(messageChannel, receiptChannel)
}

func (client *ioFogClient) SendMessageViaSocket(msg *IoMessage) error {
	defer func() {
		client.mutexMessageWs.Unlock()
		if r := recover(); r != nil {
			logger.Println("Recovered after", r)
		}
	}()
	msg.ID = "";
	msg.Timestamp = 0
	msg.Publisher = client.id
	msgBytes, err := msg.EncodeBinary()
	if err != nil {
		logger.Println("Error while encoding IoMessage to bytes:", err.Error())
		return err
	}
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(msgBytes)))
	bytesToSend := make([]byte, 0, len(msgBytes) + 5)
	bytesToSend = append(bytesToSend, CODE_MSG)
	bytesToSend = append(bytesToSend, lengthBytes...)
	bytesToSend = append(bytesToSend, msgBytes...)
	client.mutexMessageWs.Lock()
	if err = client.wsMessage.WriteMessage(ws.BinaryMessage, bytesToSend); err != nil {
		logger.Println("Error while sending message:", err.Error())
	}
	return err
}

func (client *ioFogClient) establishControlWsConnection(signalChannel chan <- int) {
	conn, _, err := ws.DefaultDialer.Dial(client.url_get_control_ws, nil)
	if conn == nil {
		logger.Println(err.Error(), "Reconnecting...")
		sleepTime := 1 << client.wsControlAttempt * wsConnectTimeout
		if client.wsControlAttempt < wsAttemptLimit {
			client.wsControlAttempt++
		}
		time.Sleep(sleepTime)
		client.establishControlWsConnection(signalChannel)
	} else {
		client.wsControlAttempt = 0
		client.wsControl = conn
		setCustomPingHandler(client.wsControl, "control")
		errChanel := make(chan int);
		go client.wsControlLoop(errChanel, signalChannel)
		for {
			select {
			case <-errChanel:
				logger.Println("Reconnecting after control ws corruption")
				client.establishControlWsConnection(signalChannel)
				return
			}
		}
	}
}

func (client *ioFogClient) establishMessageWsConnection(messageChannel chan <- *IoMessage, receiptChannel chan <- *PostMessageResponse) {
	conn, _, err := ws.DefaultDialer.Dial(client.url_get_message_ws, nil)
	if conn == nil {
		logger.Println(err.Error(), "Reconnecting...")
		sleepTime := 1 << client.wsMessageAttempt * wsConnectTimeout
		if client.wsMessageAttempt < wsAttemptLimit {
			client.wsMessageAttempt++
		}
		time.Sleep(sleepTime)
		client.establishMessageWsConnection(messageChannel, receiptChannel)
	} else {
		client.wsMessageAttempt = 0
		client.wsMessage = conn
		setCustomPingHandler(client.wsMessage, "message")
		errChanel := make(chan int);
		go client.wsMessageLoop(errChanel, messageChannel, receiptChannel)
		for {
			select {
			case <-errChanel:
				logger.Println("Reconnecting after message ws corruption")
				client.establishMessageWsConnection(messageChannel, receiptChannel)
				return
			}
		}
	}
}

func (client*ioFogClient) wsControlLoop(errChanel chan <- int, dataChannel chan <- int) {
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
			logger.Println("IoFog control signal received. Sending acknowledgement")
			err := client.wsControl.WriteMessage(ws.BinaryMessage, []byte{CODE_ACK})
			if err != nil {
				logger.Println("Error while sending acknowledgement:", err.Error())
			}
		}
	}
}

func (client*ioFogClient) wsMessageLoop(errChanel chan <- int, messageChannel chan <- *IoMessage, receiptChannel chan <-  *PostMessageResponse) {
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
			logger.Println("IoFog message received. Sending acknowledgement")
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
			logger.Println("IoFog receipt received. Sending acknowledgement")
			client.mutexMessageWs.Lock()
			err = client.wsMessage.WriteMessage(ws.BinaryMessage, []byte{CODE_ACK})
			client.mutexMessageWs.Unlock()
			if err != nil {
				logger.Println("Error while sending acknowledgement:", err.Error())
			}
		}
	}
}