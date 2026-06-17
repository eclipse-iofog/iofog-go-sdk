package client

import (
	"encoding/json"
	"fmt"
)

// CreateRegistry creates a new registry using the Controller REST API
func (clt *Client) CreateRegistry(request *RegistryCreateRequest) (int, error) {
	response := RegistryCreateResponse{}
	body, err := clt.doRequest("POST", "/registries", request)
	if err != nil {
		return -1, err
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return -1, err
	}
	return response.ID, nil
}

// UpdateRegistry patches a registry using the Controller REST API
func (clt *Client) UpdateRegistry(request RegistryUpdateRequest) error {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/registries/%d", request.ID), request)
	if err != nil {
		return err
	}
	return nil
}

// GetRegistry retrieves a single registry by ID from the Controller REST API
func (clt *Client) GetRegistry(id int) (*RegistryInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/registries/%d", id), nil)
	if err != nil {
		return nil, err
	}
	registry := new(RegistryInfo)
	if err := json.Unmarshal(body, registry); err != nil {
		return nil, err
	}
	return registry, nil
}

// ListRegistries retrieve all registries information from the Controller REST API
func (clt *Client) ListRegistries() (response RegistryListResponse, err error) {
	body, err := clt.doRequest("GET", "/registries", nil)
	if err != nil {
		return
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	return response, nil
}

// DeleteRegistry deletes a registry using the Controller REST API
func (clt *Client) DeleteRegistry(id int) (err error) {
	_, err = clt.doRequest("DELETE", fmt.Sprintf("/registries/%d", id), nil)
	return
}
