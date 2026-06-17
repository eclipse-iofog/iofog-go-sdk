package client

import (
	"encoding/json"
	"fmt"
	// "strings"
)

// CreateAgent creates an ioFog Agent using Controller REST API
func (clt *Client) CreateAgent(request *CreateAgentRequest) (CreateAgentResponse, error) {
	var response CreateAgentResponse
	if !clt.isLoggedIn() {
		return response, NewError("Controller client must be logged into perform Create Agent request")
	}

	body, err := clt.doRequest("POST", "/iofog", request)
	if err != nil {
		return response, err
	}

	var respMap map[string]any
	if err = json.Unmarshal(body, &respMap); err != nil {
		return response, err
	}
	uuid, exists := respMap["uuid"].(string)
	if !exists {
		return response, NewInternalError("Failed to get new Agent UUID from Controller")
	}

	response.UUID = uuid
	return response, nil
}

// GetAgentProvisionKey get a provisioning key for an ioFog Agent using Controller REST API
func (clt *Client) GetAgentProvisionKey(uuid string) (GetAgentProvisionKeyResponse, error) {
	var response GetAgentProvisionKeyResponse
	if !clt.isLoggedIn() {
		return response, NewError("Controller client must be logged into perform Get Agent Provisioning Key request")
	}

	body, err := clt.doRequest("GET", fmt.Sprintf("/iofog/%s/provisioning-key", uuid), nil)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// ListAgents returns all ioFog Agents information using Controller REST API
func (clt *Client) ListAgents(request ListAgentsRequest) (ListAgentsResponse, error) {
	var response ListAgentsResponse
	if !clt.isLoggedIn() {
		return response, NewError("Controller client must be logged into perform List Agents request")
	}

	body, err := clt.doRequest("GET", generateListAgentURL(request), nil)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}

	return response, nil
}

// GetAgentByID returns an ioFog Agent information using Controller REST API
func (clt *Client) GetAgentByID(uuid string) (*AgentInfo, error) {
	if !clt.isLoggedIn() {
		return nil, NewError("Controller client must be logged into perform Get Agent request")
	}

	body, err := clt.doRequest("GET", fmt.Sprintf("/iofog/%s", uuid), nil)
	if err != nil {
		return nil, err
	}

	response := new(AgentInfo)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
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
func (clt *Client) RebootAgent(uuid string) error {
	_, err := clt.doRequest("POST", fmt.Sprintf("/iofog/%s/reboot", uuid), nil)
	return err
}

// DeleteAgent removes an ioFog Agent from the Controller using Controller REST API
func (clt *Client) DeleteAgent(uuid string) error {
	if !clt.isLoggedIn() {
		return NewError("Controller client must be logged into perform Delete Agent request")
	}

	if _, err := clt.doRequest("DELETE", fmt.Sprintf("/iofog/%s", uuid), nil); err != nil {
		return err
	}

	return nil
}

// GetAgentByName retrieve the agent information by getting all agents then searching for the first occurrence in the list
func (clt *Client) GetAgentByName(name string) (*AgentInfo, error) {
	list, err := clt.ListAgents(ListAgentsRequest{})
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
func (clt *Client) PruneAgent(uuid string) error {
	_, err := clt.doRequest("POST", fmt.Sprintf("/iofog/%s/prune", uuid), nil)
	return err
}

func generateListAgentURL(request ListAgentsRequest) string {
	url := "/iofog-list"
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
	agent, err := clt.GetAgentByName(name)
	if err != nil {
		return err
	}

	if _, err := clt.doRequest("POST", fmt.Sprintf("/iofog/%s/version/upgrade", agent.UUID), nil); err != nil {
		return err
	}

	return nil
}

func (clt *Client) RollbackAgent(name string) error {
	agent, err := clt.GetAgentByName(name)
	if err != nil {
		return err
	}

	if _, err := clt.doRequest("POST", fmt.Sprintf("/iofog/%s/version/rollback", agent.UUID), nil); err != nil {
		return err
	}

	return nil
}
