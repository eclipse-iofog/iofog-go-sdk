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

// CreateEdgeResource creates an Edge Resource using Controller REST API
func (clt *Client) CreateHTTPBasedEdgeResource(request HTTPBasedEdgeResourceCreateRequest) (response GetEdgeResourceResponse, err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform request")
		return
	}

	// Send request
	body, err := clt.doRequest("POST", "/edgeResource", request)
	if err != nil {
		return
	}

	// TODO: Determine full type returned from this endpoint
	// Read uuid from response
	var respMap map[string]interface{}
	if err = json.Unmarshal(body, &respMap); err != nil {
		return
	}
	name, existsName := respMap["name"].(string)
	version, existsVersion := respMap["version"].(string)
	if !existsVersion || !existsName {
		err = NewInternalError("Failed to get new Edge Resource name/version from Controller")
		return
	}

	return clt.GetEdgeResourceByName(name, version)
}

// GetEdgeResourceByName gets an Edge Resource using Controller REST API
func (clt *Client) GetEdgeResourceByName(name, version string) (response GetEdgeResourceResponse, err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform request")
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
		err = NewError("Controller client must be logged into perform request")
		return
	}

	// Send request
	body, err := clt.doRequest("GET", fmt.Sprintf("/edgeResources"), nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	return
}

// UpdateHTTPBasedEdgeResource updates an HTTP Based Edge Resources using Controller REST API
func (clt *Client) UpdateHTTPBasedEdgeResource(name string, request HTTPBasedEdgeResourceCreateRequest) (err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform request")
		return
	}

	// Send request
	_, err = clt.doRequest("POST", fmt.Sprintf("/edgeResource/%s/%s", name, request.Version), request)
	return
}

// ListEdgeResources list all Edge Resources using Controller REST API
func (clt *Client) DeleteEdgeResource(name, version string) (err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform request")
		return
	}

	// Send request
	_, err = clt.doRequest("DELETE", fmt.Sprintf("/edgeResource/%s/%s", name, version), nil)
	return
}

// LinkEdgeResource links an Edge Resource to an Agent using Controller REST API
func (clt *Client) LinkEdgeResource(request LinkEdgeResourceRequest) (err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform request")
		return
	}

	// Send request
	_, err = clt.doRequest("POST", fmt.Sprintf("/edgeResource/%s/%s/link", request.EdgeResourceName, request.EdgeResourceVersion), request)
	return
}

// UnlinkEdgeResource unlinks an Edge Resource from an Agent using Controller REST API
func (clt *Client) UnlinkEdgeResource(request LinkEdgeResourceRequest) (err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform request")
		return
	}

	// Send request
	_, err = clt.doRequest("DELETE", fmt.Sprintf("/edgeResource/%s/%s/link", request.EdgeResourceName, request.EdgeResourceVersion), request)
	return
}
