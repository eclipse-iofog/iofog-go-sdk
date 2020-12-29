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
	"encoding/json"
	"fmt"
)

const (
	edgeResourceLoggedInErr = "Controller client must be logged in to perform Edge Resource requests"
)

// CreateHttpEdgeResource creates an Edge Resource using Controller REST API
func (clt *Client) CreateHTTPEdgeResource(request *EdgeResourceMetadata) error {
	if !clt.isLoggedIn() {
		return NewError(edgeResourceLoggedInErr)
	}

	// Send request
	if _, err := clt.doRequest("POST", "/edgeResource", request); err != nil {
		return err
	}

	return nil
}

// GetHttpEdgeResourceByName gets an Edge Resource using Controller REST API
func (clt *Client) GetHTTPEdgeResourceByName(name, version string) (response EdgeResourceMetadata, err error) {
	if !clt.isLoggedIn() {
		err = NewError(edgeResourceLoggedInErr)
		return
	}

	// Send request
	body, err := clt.doRequest("GET", fmt.Sprintf("/edgeResource/%s/%s", name, version), nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	return
}

// ListEdgeResources list all Edge Resources using Controller REST API
func (clt *Client) ListEdgeResources() (response ListEdgeResourceResponse, err error) {
	if !clt.isLoggedIn() {
		err = NewError(edgeResourceLoggedInErr)
		return
	}

	// Send request
	body, err := clt.doRequest("GET", "/edgeResources", nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	return
}

// UpdateHttpEdgeResource updates an HTTP Based Edge Resources using Controller REST API
func (clt *Client) UpdateHTTPEdgeResource(name string, request *EdgeResourceMetadata) error {
	if !clt.isLoggedIn() {
		return NewError(edgeResourceLoggedInErr)
	}

	// Send request
	if _, err := clt.doRequest("PUT", fmt.Sprintf("/edgeResource/%s/%s", name, request.Version), request); err != nil {
		return err
	}

	return nil
}

// ListEdgeResources list all Edge Resources using Controller REST API
func (clt *Client) DeleteEdgeResource(name, version string) error {
	if !clt.isLoggedIn() {
		return NewError(edgeResourceLoggedInErr)
	}

	// Send request
	if _, err := clt.doRequest("DELETE", fmt.Sprintf("/edgeResource/%s/%s", name, version), nil); err != nil {
		return err
	}

	return nil
}

// LinkEdgeResource links an Edge Resource to an Agent using Controller REST API
func (clt *Client) LinkEdgeResource(request LinkEdgeResourceRequest) error {
	if !clt.isLoggedIn() {
		return NewError(edgeResourceLoggedInErr)
	}

	// Send request
	url := fmt.Sprintf("/edgeResource/%s/%s/link", request.EdgeResourceName, request.EdgeResourceVersion)
	if _, err := clt.doRequest("POST", url, request); err != nil {
		return err
	}

	return nil
}

// UnlinkEdgeResource unlinks an Edge Resource from an Agent using Controller REST API
func (clt *Client) UnlinkEdgeResource(request LinkEdgeResourceRequest) error {
	if !clt.isLoggedIn() {
		return NewError(edgeResourceLoggedInErr)
	}

	// Send request
	url := fmt.Sprintf("/edgeResource/%s/%s/link", request.EdgeResourceName, request.EdgeResourceVersion)
	if _, err := clt.doRequest("DELETE", url, request); err != nil {
		return err
	}

	return nil
}
