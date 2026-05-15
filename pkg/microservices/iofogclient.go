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

package microservices

import (
	"errors"
	"os"
	"strconv"
)

// IoFogClient talks to ioFog agent LocalAPI v3.
type IoFogClient struct {
	id         string
	options    ClientOptions
	httpClient *ioFogHttpClient
	wsClient   *ioFogWsClient
}

func (client *IoFogClient) initClient(options ClientOptions) {
	client.options = options
	client.httpClient = newIoFogHttpClient(options)
	client.wsClient = newIoFogWsClient(options)
}

// NewIoFogClient keeps legacy constructor shape while targeting LocalAPI v3 transport.
func NewIoFogClient(id string, ssl bool, host string, port int) (*IoFogClient, error) {
	if id == "" {
		return nil, errors.New("Cannot create client with empty id")
	}
	client := IoFogClient{id: id}
	client.initClient(applyClientOptions(defaultClientOptions(),
		WithTLS(ssl),
		WithHost(host),
		WithPort(port),
	))
	return &client, nil
}

// NewIoFogClientV3 creates a LocalAPI v3 client with explicit options.
func NewIoFogClientV3(id string, opts ...ClientOption) (*IoFogClient, error) {
	if id == "" {
		return nil, errors.New("Cannot create client with empty id")
	}
	client := IoFogClient{id: id}
	client.initClient(applyClientOptions(defaultClientOptions(), opts...))
	return &client, nil
}

// NewDefaultIoFogClient creates a LocalAPI v3 client using env defaults.
func NewDefaultIoFogClient() (*IoFogClient, error) {
	return NewDefaultIoFogClientV3()
}

// NewDefaultIoFogClientV3 creates a LocalAPI v3 client using mounted token/CA defaults.
func NewDefaultIoFogClientV3() (*IoFogClient, error) {
	microserviceUID := os.Getenv(MicroserviceUID)
	if microserviceUID == "" {
		return nil, errors.New("Cannot create client with empty id: " + MicroserviceUID + " environment variable is not set")
	}
	ssl, err := strconv.ParseBool(os.Getenv(SSL))
	if err != nil {
		logger.Println("Empty or malformed", SSL, "environment variable. Using default value of", SSLDefault)
		ssl = SSLDefault
	}
	return NewIoFogClientV3(microserviceUID,
		WithTLS(ssl),
		WithHost(HostDefault),
		WithPort(PortIoFog),
	)
}

func (client *IoFogClient) GetConfig() (map[string]interface{}, error) {
	return client.httpClient.getConfig()
}

func (client *IoFogClient) GetConfigIntoStruct(config interface{}) error {
	return client.httpClient.getConfigIntoStruct(config)
}

func (client *IoFogClient) EstablishControlWsConnection(signalBufSize int) <-chan byte {
	if signalBufSize == 0 {
		signalBufSize = DefaultSignalBufferSize
	}
	signalChannel := make(chan byte, signalBufSize)
	go client.wsClient.connectToControlWs(signalChannel)
	return signalChannel
}
