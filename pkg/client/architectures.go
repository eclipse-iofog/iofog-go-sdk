package client

import (
	"encoding/json"
)

// GetArchitectures retrieves all architectures using Controller REST API.
func (clt *Client) GetArchitectures() (*ArchitecturesListResponse, error) {
	body, err := clt.doRequest("GET", "/architectures/", nil)
	if err != nil {
		return nil, err
	}

	response := new(ArchitecturesListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}
