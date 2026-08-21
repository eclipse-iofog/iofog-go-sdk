package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func authGroupPath(name string) string {
	return fmt.Sprintf("/groups/%s", url.PathEscape(name))
}

// ListAuthGroups retrieves all embedded auth groups using Controller REST API.
func (clt *Client) ListAuthGroups() ([]AuthGroupResponse, error) {
	body, err := clt.doRequest("GET", "/groups", nil)
	if err != nil {
		return nil, err
	}

	var groups []AuthGroupResponse
	if err = json.Unmarshal(body, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// CreateAuthGroup creates a custom embedded auth group using Controller REST API.
func (clt *Client) CreateAuthGroup(request AuthGroupCreateRequest) (AuthGroupResponse, error) {
	var response AuthGroupResponse
	body, err := clt.doRequest("POST", "/groups", request)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetAuthGroup retrieves an embedded auth group by name using Controller REST API.
func (clt *Client) GetAuthGroup(name string) (AuthGroupResponse, error) {
	var response AuthGroupResponse
	body, err := clt.doRequest("GET", authGroupPath(name), nil)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// UpdateAuthGroup updates an embedded auth group using Controller REST API.
// After a successful rename, use response.Name as the canonical group name.
func (clt *Client) UpdateAuthGroup(name string, request AuthGroupUpdateRequest) (AuthGroupResponse, error) {
	var response AuthGroupResponse
	body, err := clt.doRequest("PATCH", authGroupPath(name), request)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// DeleteAuthGroup deletes a custom embedded auth group using Controller REST API.
func (clt *Client) DeleteAuthGroup(name string) error {
	_, err := clt.doRequest("DELETE", authGroupPath(name), nil)
	return err
}
