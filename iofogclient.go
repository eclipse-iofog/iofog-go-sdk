/*
 *******************************************************************************
 * Copyright (c) 2018 Edgeworx, Inc.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License v. 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0
 *
 * SPDX-License-Identifier: EPL-2.0
 *******************************************************************************
*/

package iofog_sdk_go

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
)

type IoFogClient struct {
	id         string
	httpClient *ioFogHttpClient
	wsClient   *ioFogWsClient
}

func (client *IoFogClient) initClient(host string, port int, ssl bool) {
	client.httpClient = newIoFogHttpClient(client.id, ssl, host, port)
	client.wsClient = newIoFogWsClient(client.id, ssl, host, port)
}

func NewIoFogClient(id string, ssl bool, host string, port int) (*IoFogClient, error) {
	if id == "" {
		return nil, errors.New("Cannot create client with empty id")
	}
	client := IoFogClient{id: id}
	client.initClient(host, port, ssl)
	return &client, nil
}

func NewDefaultIoFogClient() (*IoFogClient, error) {
	selfname := os.Getenv(SELFNAME)
	if selfname == "" {
		return nil, errors.New("Cannot create client with empty id: " + SELFNAME + " environment variable is not set")
	}
	ssl, err := strconv.ParseBool(os.Getenv(SSL))
	if err != nil {
		logger.Println("Empty or malformed", SSL, "environment variable. Using default value of", SSL_DEFAULT)
		ssl = SSL_DEFAULT
	}

	host := IOFOG
	if cmd := exec.Command("ping", "-c", "3", host); cmd.Run() != nil {
		logger.Println("Host", host, "is unreachable. Switching to", HOST_DEFAULT)
		host = HOST_DEFAULT
	}

	client := IoFogClient{id: selfname}
	client.initClient(host, PORT_IOFOG, ssl)
	return &client, nil
}

func (client *IoFogClient) GetConfig() (map[string]interface{}, error) {
	return client.httpClient.getConfig()
}

func (client *IoFogClient) GetConfigIntoStruct(config interface{}) error {
	return client.httpClient.getConfigIntoStruct(config)
}

func (client *IoFogClient) GetNextMessages() ([]IoMessage, error) {
	return client.httpClient.getNextMessages()
}

func (client *IoFogClient) PostMessage(msg *IoMessage) (*PostMessageResponse, error) {
	msg.Publisher = client.id
	if msg.Version == 0 {
		msg.Version = IOMESSAGE_VERSION
	}
	return client.httpClient.postMessage(msg)
}

func (client *IoFogClient) GetMessagesFromPublishersWithinTimeFrame(query *MessagesQueryParameters) (*TimeFrameMessages, error) {
	query.ID = client.id
	return client.httpClient.getMessagesFromPublishersWithinTimeFrame(query)
}

func (client *IoFogClient) EstablishControlWsConnection(signalBufSize int) <-chan byte {
	if signalBufSize == 0 {
		signalBufSize = DEFAULT_SIGNAL_BUFFER_SIZE
	}
	signalChannel := make(chan byte, signalBufSize)
	go client.wsClient.connectToControlWs(signalChannel)
	return signalChannel
}

func (client *IoFogClient) EstablishMessageWsConnection(msgBufSize, receiptBufSize int) (<-chan *IoMessage, <-chan *PostMessageResponse) {
	if msgBufSize == 0 {
		msgBufSize = DEFAULT_MESSAGE_BUFFER_SIZE
	}
	if receiptBufSize == 0 {
		receiptBufSize = DEFAULT_RECEIPT_BUFFER_SIZE
	}
	messageChannel := make(chan *IoMessage, msgBufSize)
	receiptChannel := make(chan *PostMessageResponse, receiptBufSize)
	go client.wsClient.connectToMessageWs(messageChannel, receiptChannel)
	return messageChannel, receiptChannel
}

func (client *IoFogClient) SendMessageViaSocket(msg *IoMessage) error {
	msg.ID = ""
	msg.Timestamp = 0
	if msg.Version == 0 {
		msg.Version = IOMESSAGE_VERSION
	}
	msg.Publisher = client.id
	return client.wsClient.sendMessage(msg)
}
