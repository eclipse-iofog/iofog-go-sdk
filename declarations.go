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
	"log"
	"os"
	"time"
)

const (
	IOFOG = "iofog"
	ID    = "id"

	PORT_IOFOG   = 54321
	SELFNAME     = "SELFNAME"
	SSL          = "SSL"
	SSL_DEFAULT  = false
	HOST_DEFAULT = "127.0.0.1"

	URL_GET_CONFIG              = "/v2/config/get"
	URL_GET_NEXT_MESSAGES       = "/v2/messages/next"
	URL_GET_PUBLISHERS_MESSAGES = "/v2/messages/query"
	URL_POST_MESSAGE            = "/v2/messages/new"
	URL_GET_CONTROL_WS          = "/v2/control/socket/id/"
	URL_GET_MESSAGE_WS          = "/v2/message/socket/id/"

	APPLICATION_JSON = "application/json"
	HTTP             = "http"
	HTTPS            = "https"
	WS               = "ws"
	WSS              = "wss"

	CODE_ACK            = 0xB
	CODE_CONTROL_SIGNAL = 0xC
	CODE_MSG            = 0xD
	CODE_RECEIPT        = 0xE

	WS_ATTEMPT_LIMIT   = 10
	WS_CONNECT_TIMEOUT = time.Second

	DEFAULT_SIGNAL_BUFFER_SIZE  = 5
	DEFAULT_MESSAGE_BUFFER_SIZE = 200
	DEFAULT_RECEIPT_BUFFER_SIZE = 200
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
)

type getConfigResponse struct {
	Config string `json:"config"`
}

type PostMessageResponse struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type MessagesQueryParameters struct {
	ID             string   `json:"id"`
	TimeFrameStart int64    `json:"timeframestart"`
	TimeFrameEnd   int64    `json:"timeframeend"`
	Publishers     []string `json:"publishers"`
}

type getNextMessagesResponse struct {
	TimeFrameStart int64       `json:"timeframestart"`
	TimeFrameEnd   int64       `json:"timeframeend"`
	Messages       []IoMessage `json:"messages"`
}

type TimeFrameMessages getNextMessagesResponse
