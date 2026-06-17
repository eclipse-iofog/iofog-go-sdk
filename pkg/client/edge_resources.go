package client

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	edgeResourceLoggedInErr = "Controller client must be logged in to perform Edge Resource requests"
)

func (clt *Client) IsEdgeResourceCapable() error {
	if _, err := clt.doRequest("HEAD", "/capabilities/edgeResources", nil); err != nil {
		notFoundError := &NotFoundError{}
		if errors.As(err, &notFoundError) {
			return NewNotSupportedError("Edge Resources")
		}
		return err
	}
	return nil
}

func (clt *Client) edgeResourcePreflight() error {
	if err := clt.IsEdgeResourceCapable(); err != nil {
		return err
	}

	if !clt.isLoggedIn() {
		return NewError(edgeResourceLoggedInErr)
	}

	return nil
}

// CreateHttpEdgeResource creates an Edge Resource using Controller REST API
func (clt *Client) CreateHTTPEdgeResource(request *EdgeResourceMetadata) error {
	if err := clt.edgeResourcePreflight(); err != nil {
		return err
	}

	if _, err := clt.doRequest("POST", "/edgeResource", request); err != nil {
		return err
	}

	return nil
}

// GetHttpEdgeResourceByName gets an Edge Resource using Controller REST API
func (clt *Client) GetHTTPEdgeResourceByName(name, version string) (EdgeResourceMetadata, error) {
	var response EdgeResourceMetadata
	if err := clt.edgeResourcePreflight(); err != nil {
		return response, err
	}

	body, err := clt.doRequest("GET", fmt.Sprintf("/edgeResource/%s/%s", name, version), nil)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// ListEdgeResources list all Edge Resources using Controller REST API
func (clt *Client) ListEdgeResources() (ListEdgeResourceResponse, error) {
	var response ListEdgeResourceResponse
	if err := clt.edgeResourcePreflight(); err != nil {
		return response, err
	}

	body, err := clt.doRequest("GET", "/edgeResources", nil)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// UpdateHttpEdgeResource updates an HTTP Based Edge Resources using Controller REST API
func (clt *Client) UpdateHTTPEdgeResource(name string, request *EdgeResourceMetadata) error {
	if err := clt.edgeResourcePreflight(); err != nil {
		return err
	}

	if _, err := clt.doRequest("PUT", fmt.Sprintf("/edgeResource/%s/%s", name, request.Version), request); err != nil {
		return err
	}

	return nil
}

// DeleteEdgeResource deletes an Edge Resource using Controller REST API
func (clt *Client) DeleteEdgeResource(name, version string) error {
	if err := clt.edgeResourcePreflight(); err != nil {
		return err
	}

	if _, err := clt.doRequest("DELETE", fmt.Sprintf("/edgeResource/%s/%s", name, version), nil); err != nil {
		return err
	}

	return nil
}

// LinkEdgeResource links an Edge Resource to an Agent using Controller REST API
func (clt *Client) LinkEdgeResource(request LinkEdgeResourceRequest) error {
	if err := clt.edgeResourcePreflight(); err != nil {
		return err
	}

	url := fmt.Sprintf("/edgeResource/%s/%s/link", request.EdgeResourceName, request.EdgeResourceVersion)
	if _, err := clt.doRequest("POST", url, request); err != nil {
		return err
	}

	return nil
}

// UnlinkEdgeResource unlinks an Edge Resource from an Agent using Controller REST API
func (clt *Client) UnlinkEdgeResource(request LinkEdgeResourceRequest) error {
	if err := clt.edgeResourcePreflight(); err != nil {
		return err
	}
	url := fmt.Sprintf("/edgeResource/%s/%s/link", request.EdgeResourceName, request.EdgeResourceVersion)
	if _, err := clt.doRequest("DELETE", url, request); err != nil {
		return err
	}

	return nil
}
