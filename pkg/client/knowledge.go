package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

func knowledgePath(name string) string {
	return fmt.Sprintf("/knowledge/%s", url.PathEscape(name))
}

func knowledgeYAMLPath(name string) string {
	return fmt.Sprintf("/knowledge/yaml/%s", url.PathEscape(name))
}

func knowledgeLinkPath(name string) string {
	return fmt.Sprintf("/knowledge/%s/link", url.PathEscape(name))
}

func (clt *Client) uploadKnowledgeYAML(method, requestPath string, file io.Reader) (*Knowledge, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("knowledge", "knowledge.yaml")
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
	knowledge := new(Knowledge)
	if err = json.Unmarshal(body, knowledge); err != nil {
		return nil, err
	}
	return knowledge, nil
}

// ListKnowledge retrieves fleet Knowledge using the Controller REST API.
func (clt *Client) ListKnowledge() (*KnowledgeListResponse, error) {
	body, err := clt.doRequest("GET", "/knowledge", nil)
	if err != nil {
		return nil, err
	}
	response := new(KnowledgeListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetKnowledge retrieves fleet Knowledge by name. Linked fog UUIDs are not included;
// use GetKnowledgeLink.
func (clt *Client) GetKnowledge(name string) (*Knowledge, error) {
	body, err := clt.doRequest("GET", knowledgePath(name), nil)
	if err != nil {
		return nil, err
	}
	knowledge := new(Knowledge)
	if err = json.Unmarshal(body, knowledge); err != nil {
		return nil, err
	}
	return knowledge, nil
}

// CreateKnowledge creates fleet Knowledge using the Controller REST API.
func (clt *Client) CreateKnowledge(request *KnowledgeCreateRequest) (*Knowledge, error) {
	body, err := clt.doRequest("POST", "/knowledge", request)
	if err != nil {
		return nil, err
	}
	knowledge := new(Knowledge)
	if err = json.Unmarshal(body, knowledge); err != nil {
		return nil, err
	}
	return knowledge, nil
}

// UpdateKnowledge patches fleet Knowledge. Name is immutable and is taken from the path.
func (clt *Client) UpdateKnowledge(name string, request *KnowledgeUpdateRequest) (*Knowledge, error) {
	body, err := clt.doRequest("PATCH", knowledgePath(name), request)
	if err != nil {
		return nil, err
	}
	knowledge := new(Knowledge)
	if err = json.Unmarshal(body, knowledge); err != nil {
		return nil, err
	}
	return knowledge, nil
}

// DeleteKnowledge deletes fleet Knowledge. Success is HTTP 202. Returns HTTP 409
// when a microservice still binds the Knowledge name.
func (clt *Client) DeleteKnowledge(name string) error {
	_, err := clt.doRequest("DELETE", knowledgePath(name), nil)
	return err
}

// CreateKnowledgeFromYAML creates fleet Knowledge from a YAML document (multipart field "knowledge").
func (clt *Client) CreateKnowledgeFromYAML(file io.Reader) (*Knowledge, error) {
	return clt.uploadKnowledgeYAML("POST", "/knowledge/yaml", file)
}

// UpsertKnowledgeFromYAML creates or updates fleet Knowledge from a YAML document
// (multipart field "knowledge").
func (clt *Client) UpsertKnowledgeFromYAML(name string, file io.Reader) (*Knowledge, error) {
	return clt.uploadKnowledgeYAML("PUT", knowledgeYAMLPath(name), file)
}

// GetKnowledgeLink retrieves fog UUIDs linked to fleet Knowledge.
func (clt *Client) GetKnowledgeLink(name string) (*FogLinkSet, error) {
	body, err := clt.doRequest("GET", knowledgeLinkPath(name), nil)
	if err != nil {
		return nil, err
	}
	response := new(FogLinkSet)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// LinkKnowledge links fleet Knowledge to fog nodes. Name is in the URL path only.
func (clt *Client) LinkKnowledge(name string, request *FogLinkRequest) error {
	_, err := clt.doRequest("POST", knowledgeLinkPath(name), request)
	return err
}

// UnlinkKnowledge unlinks fleet Knowledge from fog nodes. Name is in the URL path
// only. Returns HTTP 409 when a microservice on that fog still binds the Knowledge name.
func (clt *Client) UnlinkKnowledge(name string, request *FogLinkRequest) error {
	_, err := clt.doRequest("DELETE", knowledgeLinkPath(name), request)
	return err
}
