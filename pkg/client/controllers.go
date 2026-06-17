package client

import "encoding/json"

func (clt *Client) GetStatus() (status ControllerStatus, err error) {
	// Prepare request
	body, err := clt.doRequest("GET", "/status", nil)
	if err != nil {
		return
	}

	// Return body
	if err = json.Unmarshal(body, &status); err != nil {
		return
	}
	return
}
