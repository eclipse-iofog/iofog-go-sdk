package client

import "encoding/json"

func (clt *Client) PutDefaultRouter(router Router) (err error) {
	// Send request
	_, err = clt.doRequest("PUT", "/router", router)
	return err
}

func (clt *Client) GetDefaultRouter() (router Router, err error) {
	// Send request
	body, err := clt.doRequest("GET", "/router", nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &router); err != nil {
		return
	}
	return
}
