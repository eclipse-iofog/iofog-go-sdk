package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

func runtimeClassPath(name string) string {
	return fmt.Sprintf("/runtimeClasses/%s", url.PathEscape(name))
}

func runtimeClassYAMLPath(name string) string {
	return fmt.Sprintf("/runtimeClasses/yaml/%s", url.PathEscape(name))
}

func runtimeClassLinkPath(name string) string {
	return fmt.Sprintf("/runtimeClasses/%s/link", url.PathEscape(name))
}

func (clt *Client) uploadRuntimeClassYAML(method, requestPath string, file io.Reader) (*RuntimeClass, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("runtimeClass", "runtimeClass.yaml")
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
	runtimeClass := new(RuntimeClass)
	if err = json.Unmarshal(body, runtimeClass); err != nil {
		return nil, err
	}
	return runtimeClass, nil
}

// ListRuntimeClasses retrieves RuntimeClasses using the Controller REST API.
func (clt *Client) ListRuntimeClasses() (*RuntimeClassListResponse, error) {
	body, err := clt.doRequest("GET", "/runtimeClasses", nil)
	if err != nil {
		return nil, err
	}
	response := new(RuntimeClassListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetRuntimeClass retrieves a RuntimeClass by name. Linked fog UUIDs are not
// included; use GetRuntimeClassLink.
func (clt *Client) GetRuntimeClass(name string) (*RuntimeClass, error) {
	body, err := clt.doRequest("GET", runtimeClassPath(name), nil)
	if err != nil {
		return nil, err
	}
	runtimeClass := new(RuntimeClass)
	if err = json.Unmarshal(body, runtimeClass); err != nil {
		return nil, err
	}
	return runtimeClass, nil
}

// CreateRuntimeClass creates a RuntimeClass using the Controller REST API.
func (clt *Client) CreateRuntimeClass(request *RuntimeClassCreateRequest) (*RuntimeClass, error) {
	body, err := clt.doRequest("POST", "/runtimeClasses", request)
	if err != nil {
		return nil, err
	}
	runtimeClass := new(RuntimeClass)
	if err = json.Unmarshal(body, runtimeClass); err != nil {
		return nil, err
	}
	return runtimeClass, nil
}

// UpdateRuntimeClass patches a RuntimeClass. Name is immutable and is taken from the path.
func (clt *Client) UpdateRuntimeClass(name string, request *RuntimeClassUpdateRequest) (*RuntimeClass, error) {
	body, err := clt.doRequest("PATCH", runtimeClassPath(name), request)
	if err != nil {
		return nil, err
	}
	runtimeClass := new(RuntimeClass)
	if err = json.Unmarshal(body, runtimeClass); err != nil {
		return nil, err
	}
	return runtimeClass, nil
}

// DeleteRuntimeClass deletes a RuntimeClass. Returns HTTP 409 when a
// microservice still pins runtime to this name.
func (clt *Client) DeleteRuntimeClass(name string) error {
	_, err := clt.doRequest("DELETE", runtimeClassPath(name), nil)
	return err
}

// CreateRuntimeClassFromYAML creates a RuntimeClass from a YAML document
// (multipart field "runtimeClass").
func (clt *Client) CreateRuntimeClassFromYAML(file io.Reader) (*RuntimeClass, error) {
	return clt.uploadRuntimeClassYAML("POST", "/runtimeClasses/yaml", file)
}

// UpsertRuntimeClassFromYAML creates or updates a RuntimeClass from a YAML
// document (multipart field "runtimeClass").
func (clt *Client) UpsertRuntimeClassFromYAML(name string, file io.Reader) (*RuntimeClass, error) {
	return clt.uploadRuntimeClassYAML("PUT", runtimeClassYAMLPath(name), file)
}

// GetRuntimeClassLink retrieves fog UUIDs linked to a RuntimeClass.
func (clt *Client) GetRuntimeClassLink(name string) (*FogLinkSet, error) {
	body, err := clt.doRequest("GET", runtimeClassLinkPath(name), nil)
	if err != nil {
		return nil, err
	}
	response := new(FogLinkSet)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// LinkRuntimeClass links a RuntimeClass to fog nodes. Name is in the URL path
// only. Linking a fog whose container engine is not edgelet returns HTTP 400.
func (clt *Client) LinkRuntimeClass(name string, request *FogLinkRequest) error {
	_, err := clt.doRequest("POST", runtimeClassLinkPath(name), request)
	return err
}

// UnlinkRuntimeClass unlinks a RuntimeClass from fog nodes. Name is in the URL
// path only. Returns HTTP 409 when a microservice on that fog still pins
// runtime to this name.
func (clt *Client) UnlinkRuntimeClass(name string, request *FogLinkRequest) error {
	_, err := clt.doRequest("DELETE", runtimeClassLinkPath(name), request)
	return err
}
