package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

func modelPath(name string) string {
	return fmt.Sprintf("/models/%s", url.PathEscape(name))
}

func modelYAMLPath(name string) string {
	return fmt.Sprintf("/models/yaml/%s", url.PathEscape(name))
}

func modelLinkPath(name string) string {
	return fmt.Sprintf("/models/%s/link", url.PathEscape(name))
}

func (clt *Client) uploadModelYAML(method, requestPath string, file io.Reader) (*Model, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("model", "model.yaml")
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	body, err := clt.doRequestWithHeaders(method, requestPath, requestBody, headers)
	if err != nil {
		return nil, err
	}
	model := new(Model)
	if err = json.Unmarshal(body, model); err != nil {
		return nil, err
	}
	return model, nil
}

// ListModels retrieves fleet models using the Controller REST API.
func (clt *Client) ListModels() (*ModelListResponse, error) {
	body, err := clt.doRequest("GET", "/models", nil)
	if err != nil {
		return nil, err
	}
	response := new(ModelListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetModel retrieves a fleet model by name. Linked fog UUIDs are not included;
// use GetModelLink.
func (clt *Client) GetModel(name string) (*Model, error) {
	body, err := clt.doRequest("GET", modelPath(name), nil)
	if err != nil {
		return nil, err
	}
	model := new(Model)
	if err = json.Unmarshal(body, model); err != nil {
		return nil, err
	}
	return model, nil
}

// CreateModel creates a fleet model using the Controller REST API.
func (clt *Client) CreateModel(request *ModelCreateRequest) (*Model, error) {
	body, err := clt.doRequest("POST", "/models", request)
	if err != nil {
		return nil, err
	}
	model := new(Model)
	if err = json.Unmarshal(body, model); err != nil {
		return nil, err
	}
	return model, nil
}

// UpdateModel patches a fleet model. Name is immutable and is taken from the path.
func (clt *Client) UpdateModel(name string, request *ModelUpdateRequest) (*Model, error) {
	body, err := clt.doRequest("PATCH", modelPath(name), request)
	if err != nil {
		return nil, err
	}
	model := new(Model)
	if err = json.Unmarshal(body, model); err != nil {
		return nil, err
	}
	return model, nil
}

// DeleteModel deletes a fleet model. Returns HTTP 409 when a microservice still
// binds the model name.
func (clt *Client) DeleteModel(name string) error {
	_, err := clt.doRequest("DELETE", modelPath(name), nil)
	return err
}

// CreateModelFromYAML creates a fleet model from a YAML document (multipart field "model").
func (clt *Client) CreateModelFromYAML(file io.Reader) (*Model, error) {
	return clt.uploadModelYAML("POST", "/models/yaml", file)
}

// UpsertModelFromYAML creates or updates a fleet model from a YAML document
// (multipart field "model").
func (clt *Client) UpsertModelFromYAML(name string, file io.Reader) (*Model, error) {
	return clt.uploadModelYAML("PUT", modelYAMLPath(name), file)
}

// GetModelLink retrieves fog UUIDs linked to a fleet model.
func (clt *Client) GetModelLink(name string) (*FogLinkSet, error) {
	body, err := clt.doRequest("GET", modelLinkPath(name), nil)
	if err != nil {
		return nil, err
	}
	response := new(FogLinkSet)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// LinkModel links a fleet model to fog nodes. Name is in the URL path only.
func (clt *Client) LinkModel(name string, request *FogLinkRequest) error {
	_, err := clt.doRequest("POST", modelLinkPath(name), request)
	return err
}

// UnlinkModel unlinks a fleet model from fog nodes. Name is in the URL path
// only. Returns HTTP 409 when a microservice on that fog still binds the model name.
func (clt *Client) UnlinkModel(name string, request *FogLinkRequest) error {
	_, err := clt.doRequest("DELETE", modelLinkPath(name), request)
	return err
}
