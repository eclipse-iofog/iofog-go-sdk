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
	"strings"
)

// CreateAgent creates an ioFog Agent using Controller REST API
func (clt *Client) CreateAgent(request *CreateAgentRequest) (response CreateAgentResponse, err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform Create Agent request")
		return
	}

	// Send request
	body, err := clt.doRequest("POST", "/iofog", request)
	if err != nil {
		return
	}

	// TODO: Determine full type returned from this endpoint
	// Read uuid from response
	var respMap map[string]interface{}
	if err = json.Unmarshal(body, &respMap); err != nil {
		return
	}
	uuid, exists := respMap["uuid"].(string)
	if !exists {
		err = NewInternalError("Failed to get new Agent UUID from Controller")
		return
	}

	response.UUID = uuid
	return
}

// GetAgentProvisionKey get a provisioning key for an ioFog Agent using Controller REST API
func (clt *Client) GetAgentProvisionKey(uuid string) (response GetAgentProvisionKeyResponse, err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform Get Agent Provisioning Key request")
		return
	}

	// Send request
	body, err := clt.doRequest("GET", fmt.Sprintf("/iofog/%s/provisioning-key", uuid), nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	return
}

// ListAgents returns all ioFog Agents information using Controller REST API
func (clt *Client) ListAgents(request ListAgentsRequest) (response ListAgentsResponse, err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform List Agents request")
		return
	}

	// Send request
	body, err := clt.doRequest("GET", generateListAgentURL(request), nil)
	if err != nil {
		return
	}

	// Return body
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}

	return
}

// GetAgentByID returns an ioFog Agent information using Controller REST API
func (clt *Client) GetAgentByID(uuid string) (response *AgentInfo, err error) {
	if !clt.isLoggedIn() {
		err = NewError("Controller client must be logged into perform Get Agent request")
		return
	}

	// Send request
	body, err := clt.doRequest("GET", fmt.Sprintf("/iofog/%s", uuid), nil)
	if err != nil {
		return
	}

	// Return body
	response = new(AgentInfo)
	if err = json.Unmarshal(body, response); err != nil {
		return
	}

	return
}

// UpdateAgent patches an ioFog Agent using Controller REST API
func (clt *Client) UpdateAgent(request *AgentUpdateRequest) (*AgentInfo, error) {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/iofog/%s", request.UUID), request)
	if err != nil {
		return nil, err
	}
	return clt.GetAgentByID(request.UUID)
}

// RebootAgent reboots an ioFog Agent using Controller REST API
func (clt *Client) RebootAgent(uuid string) (err error) {
	_, err = clt.doRequest("POST", fmt.Sprintf("/iofog/%s/reboot", uuid), nil)
	return
}

// DeleteAgent removes an ioFog Agent from the Controller using Controller REST API
func (clt *Client) DeleteAgent(uuid string) error {
	if !clt.isLoggedIn() {
		return NewError("Controller client must be logged into perform Delete Agent request")
	}

	// Send request
	if _, err := clt.doRequest("DELETE", fmt.Sprintf("/iofog/%s", uuid), nil); err != nil {
		return err
	}

	return nil
}

// GetAgentByName retrieve the agent information by getting all agents then searching for the first occurance in the list
func (clt *Client) GetAgentByName(name string, system bool) (*AgentInfo, error) {
	list, err := clt.ListAgents(ListAgentsRequest{System: system})
	if err != nil {
		return nil, err
	}
	for idx := range list.Agents {
		if list.Agents[idx].Name == name {
			return &list.Agents[idx], nil
		}
	}
	return nil, NewNotFoundError(fmt.Sprintf("Could not find agent: %s", name))
}

// PruneAgent prunes an ioFog Agent using Controller REST API
func (clt *Client) PruneAgent(uuid string) (err error) {
	_, err = clt.doRequest("POST", fmt.Sprintf("/iofog/%s/prune", uuid), nil)
	return
}

func generateListAgentURL(request ListAgentsRequest) string {
	// Embed request options into URL as query params
	url := "/iofog-list?system=false"
	if request.System {
		url = strings.Replace(url, "false", "true", 1)
	}
	for idx, filter := range request.Filters {
		params := []string{
			fmt.Sprintf("&filters[%d][key]=%s", idx, filter.Key),
			fmt.Sprintf("&filters[%d][value]=%s", idx, filter.Value),
			fmt.Sprintf("&filters[%d][condition]=%s", idx, filter.Condition),
		}
		for _, param := range params {
			url = fmt.Sprintf("%s%s", url, param)
		}
	}
	return url
}

func (clt *Client) UpgradeAgent(name string) error {
	// Get Agent uuid
	agent, err := clt.GetAgentByName(name, false)
	if err != nil {
		return err
	}

	// Send request
	if _, err := clt.doRequest("POST", fmt.Sprintf("/iofog/%s/version/upgrade", agent.UUID), nil); err != nil {
		return err
	}

	return nil
}

func (clt *Client) RollbackAgent(name string) error {
	// Get Agent uuid
	agent, err := clt.GetAgentByName(name, false)
	if err != nil {
		return err
	}

	// Send request
	if _, err := clt.doRequest("POST", fmt.Sprintf("/iofog/%s/version/rollback", agent.UUID), nil); err != nil {
		return err
	}

	return nil
}
