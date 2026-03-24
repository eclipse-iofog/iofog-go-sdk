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
	"bytes"
	"encoding/json"
	"fmt"
)

type ioFogHttpClient struct {
	url_base_rest               string
	url_get_config              string
	url_get_next_messages       string
	url_get_publishers_messages string
	url_post_message            string
	requestBodyId               []byte
}

func newIoFogHttpClient(id string, ssl bool, host string, port int) *ioFogHttpClient {
	client := ioFogHttpClient{}
	protocol_rest := HTTP
	if ssl {
		protocol_rest = HTTPS
	}
	client.url_base_rest = fmt.Sprintf("%s://%s:%d", protocol_rest, host, port)
	client.url_get_config = fmt.Sprint(client.url_base_rest, URL_GET_CONFIG)
	client.url_get_next_messages = fmt.Sprint(client.url_base_rest, URL_GET_NEXT_MESSAGES)
	client.url_get_publishers_messages = fmt.Sprint(client.url_base_rest, URL_GET_PUBLISHERS_MESSAGES)
	client.url_post_message = fmt.Sprint(client.url_base_rest, URL_POST_MESSAGE)
	client.requestBodyId, _ = json.Marshal(map[string]interface{}{
		ID: id,
	})
	return &client
}

func (client *ioFogHttpClient) getConfig() (map[string]interface{}, error) {
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

func (client *ioFogHttpClient) getConfigIntoStruct(config interface{}) error {
	resp, err := makePostRequest(client.url_get_config, APPLICATION_JSON, bytes.NewBuffer(client.requestBodyId))
	if err != nil {
		return err
	}
	configResponse := new(getConfigResponse)
	if err := json.Unmarshal(resp, configResponse); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(configResponse.Config), config); err != nil {
		return err
	}
	return nil
}

func (client *ioFogHttpClient) getNextMessages() ([]IoMessageReadable, error) {
	resp, err := makePostRequest(client.url_get_next_messages, APPLICATION_JSON, bytes.NewBuffer(client.requestBodyId))
	if err != nil {
		return nil, err
	}
	nextMessagesResponse := new(getNextMessagesReadableResponse)
	if err := json.Unmarshal(resp, nextMessagesResponse); err != nil {
		return nil, err
	}

	// Decode Base64-encoded fields for each message
	var decodedMessages []IoMessageReadable
	for _, msg := range nextMessagesResponse.Messages {
		decodedMessages = append(decodedMessages, *decodeJson(&msg))
	}

	return decodedMessages, nil
}

func (client *ioFogHttpClient) postMessage(msg *IoMessage) (*PostMessageResponse, error) {
	encodedMsg := encodeJson(msg)

	requestBytes, _ := json.Marshal(encodedMsg)
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

func (client *ioFogHttpClient) getMessagesFromPublishersWithinTimeFrame(query *MessagesQueryParameters) (*TimeFrameReadableMessages, error) {
	requestBytes, _ := json.Marshal(query)
	resp, err := makePostRequest(client.url_get_publishers_messages, APPLICATION_JSON, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	nextMessagesResponse := new(getNextMessagesReadableResponse)
	if err := json.Unmarshal(resp, nextMessagesResponse); err != nil {
		return nil, err
	}

	// Decode Base64-encoded fields for each message
	var decodedMessages []IoMessageReadable
	for _, msg := range nextMessagesResponse.Messages {
		decodedMessages = append(decodedMessages, *decodeJson(&msg))
	}

	// Convert to TimeFrameReadableMessages
	readableResponse := &TimeFrameReadableMessages{
		TimeFrameStart: nextMessagesResponse.TimeFrameStart,
		TimeFrameEnd:   nextMessagesResponse.TimeFrameEnd,
		Messages:       decodedMessages,
	}

	return readableResponse, nil
}
