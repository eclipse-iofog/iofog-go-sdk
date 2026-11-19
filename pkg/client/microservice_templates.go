package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

func microserviceTemplatePath(name string) string {
	return fmt.Sprintf("/microserviceTemplates/%s", url.PathEscape(name))
}

func microserviceTemplateYAMLPath(name string) string {
	return fmt.Sprintf("/microserviceTemplates/yaml/%s", url.PathEscape(name))
}

func (clt *Client) uploadMicroserviceTemplateYAML(method, requestPath string, file io.Reader) (*MicroserviceTemplate, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("template", "template.yaml")
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
	template := new(MicroserviceTemplate)
	if err = json.Unmarshal(body, template); err != nil {
		return nil, err
	}
	return template, nil
}

// ListMicroserviceTemplates retrieves microservice templates using the Controller REST API.
func (clt *Client) ListMicroserviceTemplates() (*MicroserviceTemplateListResponse, error) {
	body, err := clt.doRequest("GET", "/microserviceTemplates", nil)
	if err != nil {
		return nil, err
	}
	response := new(MicroserviceTemplateListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetMicroserviceTemplate retrieves a microservice template by name.
func (clt *Client) GetMicroserviceTemplate(name string) (*MicroserviceTemplate, error) {
	body, err := clt.doRequest("GET", microserviceTemplatePath(name), nil)
	if err != nil {
		return nil, err
	}
	template := new(MicroserviceTemplate)
	if err = json.Unmarshal(body, template); err != nil {
		return nil, err
	}
	return template, nil
}

// CreateMicroserviceTemplate creates a microservice template using the Controller REST API.
func (clt *Client) CreateMicroserviceTemplate(request *MicroserviceTemplateCreateRequest) (*MicroserviceTemplate, error) {
	body, err := clt.doRequest("POST", "/microserviceTemplates", request)
	if err != nil {
		return nil, err
	}
	template := new(MicroserviceTemplate)
	if err = json.Unmarshal(body, template); err != nil {
		return nil, err
	}
	return template, nil
}

// UpdateMicroserviceTemplate patches a microservice template. Name is immutable and is taken from the path.
func (clt *Client) UpdateMicroserviceTemplate(name string, request *MicroserviceTemplateUpdateRequest) (*MicroserviceTemplate, error) {
	body, err := clt.doRequest("PATCH", microserviceTemplatePath(name), request)
	if err != nil {
		return nil, err
	}
	template := new(MicroserviceTemplate)
	if err = json.Unmarshal(body, template); err != nil {
		return nil, err
	}
	return template, nil
}

// DeleteMicroserviceTemplate deletes a microservice template.
func (clt *Client) DeleteMicroserviceTemplate(name string) error {
	_, err := clt.doRequest("DELETE", microserviceTemplatePath(name), nil)
	return err
}

// CreateMicroserviceTemplateFromYAML creates a microservice template from a YAML document
// (multipart field "template").
func (clt *Client) CreateMicroserviceTemplateFromYAML(file io.Reader) (*MicroserviceTemplate, error) {
	return clt.uploadMicroserviceTemplateYAML("POST", "/microserviceTemplates/yaml", file)
}

// UpdateMicroserviceTemplateFromYAML creates or updates a microservice template from a YAML
// document (multipart field "template").
func (clt *Client) UpdateMicroserviceTemplateFromYAML(name string, file io.Reader) (*MicroserviceTemplate, error) {
	return clt.uploadMicroserviceTemplateYAML("PUT", microserviceTemplateYAMLPath(name), file)
}
