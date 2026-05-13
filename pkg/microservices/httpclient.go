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
	tokenProvider func(string) (string, error)
}

func newIoFogHttpClient(options ClientOptions) *ioFogHttpClient {
	return &ioFogHttpClient{options: options, tokenProvider: readBearerToken}
}

func (client *ioFogHttpClient) getConfig() (map[string]interface{}, error) {
	resp, err := client.makeRequest(http.MethodGet, URLGetConfigV3, nil)
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

func (client *ioFogHttpClient) makeRequest(method, path string, body io.Reader) (map[string]interface{}, error) {
	httpClient, err := buildHTTPClient(client.options)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP client: %w", err)
	}
	token, err := client.tokenProvider(client.options.TokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load bearer token: %w", err)
	}

	var lastErr error
	for idx, host := range client.options.requestHosts() {
		endpoint := client.options.restBaseURLForHost(host) + path
		req, reqErr := http.NewRequest(method, endpoint, body)
		if reqErr != nil {
			return nil, reqErr
		}
		req.Header.Set("Accept", ApplicationJSON)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, doErr := httpClient.Do(req)
		if doErr != nil {
			lastErr = doErr
			if idx < len(client.options.requestHosts())-1 && isRetriableHostError(doErr) {
				logger.Println("WARN: failed to reach host", host, "for", path, "trying fallback host:", doErr)
				continue
			}
			return nil, doErr
		}
		responseBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		return parseV3Envelope(responseBody, resp.StatusCode)
	}
	return nil, lastErr
}
