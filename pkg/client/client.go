/*
 *  *******************************************************************************
 *  * Copyright (c) 2019 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package client

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type controllerStatus struct {
	version         string
	versionNoSuffix string
	versionNums     []string
}

type Client struct {
	endpoint    string
	baseURL     string
	accessToken string
	retries     Retries
	status      controllerStatus
}

type Options struct {
	Endpoint string
	Retries  *Retries
}

var apiPrefix = "/api/v3"

func New(opt Options) *Client {
	// remember if we are using https
	var protocol string
	if strings.HasPrefix(opt.Endpoint, "https://") {
		protocol = "https"
	} else {
		protocol = "http"
	}

	// Remove prefix
	regex := regexp.MustCompile("https?://")
	endpoint := regex.ReplaceAllString(opt.Endpoint, "")

	// Add default port if none specified
	if !strings.Contains(endpoint, ":") {
		endpoint = endpoint + ":" + ControllerPortString
	}

	retries := GlobalRetriesPolicy
	if opt.Retries != nil {
		retries = *opt.Retries
	}
	client := &Client{
		endpoint: endpoint,
		retries:  retries,
		baseURL:  fmt.Sprintf("%s://%s%s", protocol, endpoint, apiPrefix),
	}
	// Get Controller version
	if status, err := client.GetStatus(); err == nil {
		versionNoSuffix := before(status.Versions.Controller, "-")
		versionNums := strings.Split(versionNoSuffix, ".")
		client.status = controllerStatus{
			version:         status.Versions.Controller,
			versionNoSuffix: versionNoSuffix,
			versionNums:     versionNums,
		}
	}
	return client
}

func NewAndLogin(opt Options, email, password string) (clt *Client, err error) {
	clt = New(opt)
	if err = clt.Login(LoginRequest{Email: email, Password: password}); err != nil {
		return
	}
	return
}

func NewWithToken(opt Options, token string) (clt *Client, err error) {
	clt = New(opt)
	clt.SetAccessToken(token)
	return
}

func (clt *Client) GetEndpoint() string {
	return clt.endpoint
}

func (clt *Client) GetRetries() Retries {
	return clt.retries
}

func (clt *Client) SetRetries(retries Retries) {
	clt.retries = retries
}

func (clt *Client) GetAccessToken() string {
	return clt.accessToken
}

func (clt *Client) SetAccessToken(token string) {
	clt.accessToken = token
}

func (clt *Client) makeRequestURL(url string) string {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	return clt.baseURL + url
}

func (clt *Client) doRequestWithRetries(currentRetries Retries, method, requestURL string, headers map[string]string, request interface{}) ([]byte, error) {
	// Send request
	bytes, err := httpDo(method, requestURL, headers, request)
	if err != nil {
		httpErr, ok := err.(*HTTPError)
		// If HTTP Error
		if ok {
			if httpErr.Code == 408 { // HTTP Timeout
				if currentRetries.Timeout < clt.retries.Timeout {
					currentRetries.Timeout++
					time.Sleep(time.Duration(currentRetries.Timeout) * time.Second)
					return clt.doRequestWithRetries(currentRetries, method, requestURL, headers, request)
				}
				return bytes, err
			}
		}
		// If custom retries defined
		if clt.retries.CustomMessage != nil {
			for message, allowedRetries := range clt.retries.CustomMessage {
				if strings.Contains(err.Error(), message) {
					if currentRetries.CustomMessage[message] < allowedRetries {
						currentRetries.CustomMessage[message]++
						time.Sleep(time.Duration(currentRetries.CustomMessage[message]) * time.Second)
						return clt.doRequestWithRetries(currentRetries, method, requestURL, headers, request)
					}
					return bytes, err
				}
			}
		}
	}
	return bytes, err
}

func (clt *Client) doRequest(method, url string, request interface{}) ([]byte, error) {
	// Prepare request
	requestURL := clt.makeRequestURL(url)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": clt.accessToken,
	}

	currentRetries := Retries{CustomMessage: make(map[string]int)}
	if clt.retries.CustomMessage != nil {
		for message := range clt.retries.CustomMessage {
			currentRetries.CustomMessage[message] = 0
		}
	}

	return clt.doRequestWithRetries(currentRetries, method, requestURL, headers, request)
}

func (clt *Client) isLoggedIn() bool {
	return clt.accessToken != ""
}
