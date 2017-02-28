package container_sdk_go

import (
	"os"
	"strconv"
	"os/exec"
	"errors"
)

type ioFogClient struct {
	id         string
	httpClient *ioFogHttpClient
	wsClient   *ioFogWsClient
}

func (client *ioFogClient) initClient(host string, port int, ssl bool) {
	client.httpClient = newIoFogHttpClient(client.id, ssl, host, port)
	client.wsClient = newIoFogWsClient(client.id, ssl, host, port)
}

func NewIoFogClient(id string, ssl bool, host string, port int) (*ioFogClient, error) {
	if id == "" {
		return nil, errors.New("Cannot create client with empty id")
	}
	client := ioFogClient{id: id}
	client.initClient(host, port, ssl)
	return &client, nil
}

func NewDefaultIoFogClient() (*ioFogClient, error) {
	selfname := os.Getenv(SELFNAME)
	if selfname == "" {
		return nil, errors.New("Cannot create client with empty id: " + SELFNAME + " environment virable is not set")
	}
	ssl, err := strconv.ParseBool(os.Getenv(SSL))
	if err != nil {
		logger.Println("Empty or malformed", SSL, "environment variable. Using default value of", SSL_DEFAULT)
		ssl = SSL_DEFAULT
	}

	host := IOFOG
	if cmd := exec.Command("ping", "-c 3", host); cmd.Run() != nil {
		logger.Println("Host", host, "is unreachable. Switching to", HOST_DEFAULT)
		host = HOST_DEFAULT
	}

	client := ioFogClient{id: selfname}
	client.initClient(host, PORT_IOFOG, ssl)
	return &client, nil
}

func (client *ioFogClient) GetConfig() (map[string]interface{}, error) {
	return client.httpClient.getConfig()
}

func (client *ioFogClient) GetNextMessages() ([]IoMessage, error) {
	return client.httpClient.getNextMessages()
}

func (client *ioFogClient) PostMessage(msg *IoMessage) (*PostMessageResponse, error) {
	msg.Publisher = client.id
	if msg.Version == 0 {
		msg.Version = IOMESSAGE_VERSION
	}
	return client.httpClient.postMessage(msg)
}

func (client *ioFogClient) GetMessagesFromPublishersWithinTimeFrame(query *MessagesQueryParameters) (*TimeFrameMessages, error) {
	query.ID = client.id
	return client.httpClient.getMessagesFromPublishersWithinTimeFrame(query)
}

func (client *ioFogClient) EstablishControlWsConnection() <- chan byte {
	signalChannel := make(chan byte, 5)
	go client.wsClient.connectToControlWs(signalChannel)
	return signalChannel
}

func (client *ioFogClient) EstablishMessageWsConnection() (<- chan *IoMessage, <- chan *PostMessageResponse) {
	messageChannel := make(chan *IoMessage, 20)
	receiptChannel := make(chan *PostMessageResponse, 20)
	go client.wsClient.connectToMessageWs(messageChannel, receiptChannel)
	return messageChannel, receiptChannel
}

func (client *ioFogClient) SendMessageViaSocket(msg *IoMessage) error {
	msg.ID = "";
	msg.Timestamp = 0
	if msg.Version == 0 {
		msg.Version = IOMESSAGE_VERSION
	}
	msg.Publisher = client.id
	return client.wsClient.sendMessage(msg)
}