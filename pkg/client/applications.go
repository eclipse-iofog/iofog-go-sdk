package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
)

// GetApplicationByName retrieve application information using the Controller REST API
func (clt *Client) GetApplicationByName(name string) (*ApplicationInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/application/%s", name), nil)
	if err != nil {
		return nil, err
	}
	application := new(ApplicationInfo)
	if err = json.Unmarshal(body, application); err != nil {
		return nil, err
	}
	return application, nil
}

// GetSystemApplicationByName retrieve system application information using the Controller REST API
func (clt *Client) GetSystemApplicationByName(name string) (*ApplicationInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/application/system/%s", name), nil)
	if err != nil {
		return nil, err
	}
	application := new(ApplicationInfo)
	if err = json.Unmarshal(body, application); err != nil {
		return nil, err
	}
	return application, nil
}

// CreateApplicationFromYAML creates a new application using the Controller REST API
// It sends the yaml file to Controller REST API
func (clt *Client) CreateApplicationFromYAML(file io.Reader) (*ApplicationInfo, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, _ := writer.CreateFormFile("application", "application.yaml")
	_, err := io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	body, err := clt.doRequestWithHeaders("POST", "/application/yaml", requestBody, headers)
	if err != nil {
		return nil, err
	}
	response := FlowCreateResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return clt.GetApplicationByName(response.Name)
}

// UpdateApplicationFromYAML updates an application using the Controller REST API
// It sends the yaml file to Controller REST API
func (clt *Client) UpdateApplicationFromYAML(name string, file io.Reader) (*ApplicationInfo, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, _ := writer.CreateFormFile("application", "application.yaml")
	_, err := io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	_, err = clt.doRequestWithHeaders("PUT", fmt.Sprintf("/application/yaml/%s", name), requestBody, headers)
	if err != nil {
		return nil, err
	}
	return clt.GetApplicationByName(name)
}

// UpdateApplication patches an application using the Controller REST API
func (clt *Client) PatchApplication(name string, request *ApplicationPatchRequest) (*ApplicationInfo, error) {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/application/%s", name), *request)
	if err != nil {
		return nil, err
	}
	newName := name
	if request.Name != nil {
		newName = *request.Name
	}
	return clt.GetApplicationByName(newName)
}

// StartApplication set the application as active using the Controller REST API
func (clt *Client) StartApplication(name string) (*ApplicationInfo, error) {
	active := true
	return clt.PatchApplication(name, &ApplicationPatchRequest{IsActivated: &active})
}

// StopApplication set the application as inactive using the Controller REST API
func (clt *Client) StopApplication(name string) (*ApplicationInfo, error) {
	active := false
	return clt.PatchApplication(name, &ApplicationPatchRequest{IsActivated: &active})
}

// GetAllApplications retrieve all flows information from the Controller REST API
func (clt *Client) GetAllApplications() (*ApplicationListResponse, error) {
	body, err := clt.doRequest("GET", "/application", nil)
	if err != nil {
		return nil, err
	}
	response := new(ApplicationListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetAllSystemApplications retrieve all system applications information from the Controller REST API
func (clt *Client) GetAllSystemApplications() (*ApplicationListResponse, error) {
	body, err := clt.doRequest("GET", "/application/system", nil)
	if err != nil {
		return nil, err
	}
	response := new(ApplicationListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// DeleteApplication deletes an application using the Controller REST API
func (clt *Client) DeleteApplication(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/application/%s", name), nil)
	return err
}

// DeleteSystemApplication deletes an application using the Controller REST API
func (clt *Client) DeleteSystemApplication(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/application/system/%s", name), nil)
	return err
}
