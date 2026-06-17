package client

import "encoding/json"

func (clt *Client) GetStatus() (ControllerStatus, error) {
	var status ControllerStatus
	body, err := clt.doRequest("GET", "/status", nil)
	if err != nil {
		return status, err
	}

	if err = json.Unmarshal(body, &status); err != nil {
		return status, err
	}
	return status, nil
}
