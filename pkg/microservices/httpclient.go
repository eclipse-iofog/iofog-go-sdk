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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ioFogHttpClient struct {
	options       ClientOptions
	urlBaseREST   string
	urlGetConfig  string
	tokenProvider func(string) (string, error)
}

func newIoFogHttpClient(options ClientOptions) *ioFogHttpClient {
	client := ioFogHttpClient{options: options, tokenProvider: readBearerToken}
	client.urlBaseREST = options.restBaseURL()
	client.urlGetConfig = fmt.Sprint(client.urlBaseREST, URLGetConfigV3)
	return &client
}

func (client *ioFogHttpClient) getConfig() (map[string]interface{}, error) {
	resp, err := client.makeRequest(http.MethodGet, client.urlGetConfig, nil)
	if err != nil {
		return nil, err
	}
	payload, ok := resp["config"]
	if !ok {
		return nil, fmt.Errorf("missing config payload in response")
	}
	switch typed := payload.(type) {
	case map[string]interface{}:
		return typed, nil
	case string:
		config := make(map[string]interface{})
		if err := json.Unmarshal([]byte(typed), &config); err != nil {
			return nil, fmt.Errorf("failed to decode config string payload: %w", err)
		}
		return config, nil
	default:
		return nil, fmt.Errorf("unsupported config payload type: %T", payload)
	}
}

func (client *ioFogHttpClient) getConfigIntoStruct(config interface{}) error {
	configMap, err := client.getConfig()
	if err != nil {
		return err
	}
	configBytes, err := json.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("failed to marshal config payload: %w", err)
	}
	if err := json.Unmarshal(configBytes, config); err != nil {
		return fmt.Errorf("failed to decode config into target struct: %w", err)
	}
	return nil
}

func (client *ioFogHttpClient) makeRequest(method, url string, body io.Reader) (map[string]interface{}, error) {
	httpClient, err := buildHTTPClient(client.options)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP client: %w", err)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", ApplicationJSON)
	token, err := client.tokenProvider(client.options.TokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load bearer token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseV3Envelope(responseBody, resp.StatusCode)
}
