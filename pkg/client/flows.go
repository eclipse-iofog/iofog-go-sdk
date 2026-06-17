package client

import (
	"encoding/json"
	"fmt"
)

// GetFlowByID retrieve flow information using the Controller REST API
func (clt *Client) GetFlowByID(id int) (*FlowInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/flow/%d", id), nil)
	if err != nil {
		return nil, err
	}
	flow := new(FlowInfo)
	if err = json.Unmarshal(body, flow); err != nil {
		return nil, err
	}
	return flow, nil
}

// CreateFlow creates a new flow using the Controller REST API
func (clt *Client) CreateFlow(name, description string) (*FlowInfo, error) {
	response := FlowCreateResponse{}
	body, err := clt.doRequest("POST", "/flow", FlowCreateRequest{Name: name, Description: description})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return clt.GetFlowByID(response.ID)
}

// UpdateFlow patches a flow using the Controller REST API
func (clt *Client) UpdateFlow(request *FlowUpdateRequest) (*FlowInfo, error) {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/flow/%d", request.ID), *request)
	if err != nil {
		return nil, err
	}
	return clt.GetFlowByID(request.ID)
}

// StartFlow set the flow as active using the Controller REST API
func (clt *Client) StartFlow(id int) (*FlowInfo, error) {
	active := true
	return clt.UpdateFlow(&FlowUpdateRequest{ID: id, IsActivated: &active})
}

// StopFlow set the flow as inactive using the Controller REST API
func (clt *Client) StopFlow(id int) (*FlowInfo, error) {
	active := false
	return clt.UpdateFlow(&FlowUpdateRequest{ID: id, IsActivated: &active})
}

// GetAllFlows retrieve all flows information from the Controller REST API
func (clt *Client) GetAllFlows() (*FlowListResponse, error) {
	body, err := clt.doRequest("GET", "/flow", nil)
	if err != nil {
		return nil, err
	}
	response := new(FlowListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetFlowByName retrieve the flow information by getting all flows then searching for the first occurrence in the list
func (clt *Client) GetFlowByName(name string) (*FlowInfo, error) {
	list, err := clt.GetAllFlows()
	if err != nil {
		return nil, err
	}
	for _, flow := range list.Flows {
		if flow.Name == name {
			return &flow, nil
		}
	}
	return nil, NewNotFoundError(fmt.Sprintf("Could not find flow: %s", name))
}

// DeleteFlow deletes a flow using the Controller REST API
func (clt *Client) DeleteFlow(id int) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/flow/%d", id), nil)
	return err
}
