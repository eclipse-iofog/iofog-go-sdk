package client

import "encoding/json"

func (clt *Client) PutDefaultRouter(router Router) error {
	_, err := clt.doRequest("PUT", "/router", router)
	return err
}

func (clt *Client) GetDefaultRouter() (Router, error) {
	var router Router
	body, err := clt.doRequest("GET", "/router", nil)
	if err != nil {
		return router, err
	}

	if err = json.Unmarshal(body, &router); err != nil {
		return router, err
	}
	return router, nil
}
