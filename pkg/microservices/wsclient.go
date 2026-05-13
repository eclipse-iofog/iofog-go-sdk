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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	ws "github.com/gorilla/websocket"
)

type ioFogWsClient struct {
	options       ClientOptions
	urlBaseWS     string
	urlGetControl string
}

func newIoFogWsClient(options ClientOptions) *ioFogWsClient {
	client := ioFogWsClient{options: options}
	client.urlBaseWS = options.wsBaseURL()
	client.urlGetControl = fmt.Sprint(client.urlBaseWS, URLGetControlWSV3)
	return &client
}

func (client *ioFogWsClient) connectToControlWs(signalChannel chan<- byte) {
	attempt := 0
	for {
		conn, err := client.dialControlWS()
		if err != nil {
			attempt++
			logger.Println(err.Error(), "Reconnecting to v3 control ws...")
			sleepWithCap(attempt, client.options.WSReconnectBaseDelay, client.options.WSReconnectMaxDelay)
			continue
		}
		logger.Println("Control ws connection has been established.")
		attempt = 0
		if err := client.listenControl(conn, signalChannel); err != nil {
			logger.Println("Control ws reconnect after corruption:", err.Error())
		}
		_ = conn.Close()
	}
}

func (client *ioFogWsClient) dialControlWS() (*ws.Conn, error) {
	dialer := ws.Dialer{
		HandshakeTimeout: client.options.WSHandshakeTimeout,
	}
	if client.options.UseTLS {
		rootCAs, err := loadRootCAs(client.options.CAPath)
		if err != nil {
			return nil, err
		}
		dialer.TLSClientConfig = buildTLSConfig(rootCAs)
	}
	token, err := readBearerToken(client.options.TokenPath)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	conn, _, err := dialer.Dial(client.urlGetControl, header)
	if err != nil {
		return nil, err
	}
	setCustomPingHandler(conn)
	return conn, nil
}

func (client *ioFogWsClient) listenControl(conn *ws.Conn, signalChannel chan<- byte) error {
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if len(payload) == 0 {
			continue
		}
		if payload[0] != CODE_CONTROL_SIGNAL {
			continue
		}
		signalChannel <- payload[0]
		if err := conn.WriteMessage(ws.BinaryMessage, []byte{CODE_ACK}); err != nil {
			return err
		}
	}
}

func buildTLSConfig(rootCAs *x509.CertPool) *tls.Config {
	return &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
	}
}
