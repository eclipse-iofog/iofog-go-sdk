package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
)

// GetMicroserviceByName retrieves a microservice information using Controller REST API
func (clt *Client) GetMicroserviceByName(appName, name string) (*MicroserviceInfo, error) {
	listMsvcs, err := clt.GetMicroservicesByApplication(appName)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(listMsvcs.Microservices); i++ {
		if listMsvcs.Microservices[i].Name == name {
			return &listMsvcs.Microservices[i], nil
		}
	}
	return nil, NewNotFoundError(fmt.Sprintf("Could not find a microservice named %s/%s", appName, name))
}

// GetSystemMicroserviceByName retrieves a system microservice information using Controller REST API
func (clt *Client) GetSystemMicroserviceByName(appName, name string) (*MicroserviceInfo, error) {
	listMsvcs, err := clt.GetSystemMicroservicesByApplication(appName)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(listMsvcs.Microservices); i++ {
		if listMsvcs.Microservices[i].Name == name {
			return &listMsvcs.Microservices[i], nil
		}
	}
	return nil, NewNotFoundError(fmt.Sprintf("Could not find a system microservice named %s/%s", appName, name))
}

// GetMicroserviceByID retrieves a microservice information using Controller REST API
func (clt *Client) GetMicroserviceByID(uuid string) (*MicroserviceInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/microservices/%s", uuid), nil)
	if err != nil {
		return nil, err
	}

	response := new(MicroserviceInfo)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetSystemMicroserviceByID retrieves a system microservice information using Controller REST API
func (clt *Client) GetSystemMicroserviceByID(uuid string) (*MicroserviceInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/microservices/system/%s", uuid), nil)
	if err != nil {
		return nil, err
	}

	response := new(MicroserviceInfo)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// CreateMicroserviceFromYAML creates a new microservice using the Controller REST API
// It sends the yaml file to Controller REST API
func (clt *Client) CreateMicroserviceFromYAML(file io.Reader) (*MicroserviceInfo, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, _ := writer.CreateFormFile("microservice", "microservice.yaml")
	_, err := io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	body, err := clt.doRequestWithHeaders("POST", "/microservices/yaml", requestBody, headers)
	if err != nil {
		return nil, err
	}
	response := MicroserviceCreateResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return clt.GetMicroserviceByID(response.UUID)
}

// GetMicroservicesByApplication returns a list of microservices in a specific application using Controller REST API
func (clt *Client) GetMicroservicesByApplication(application string) (*MicroserviceListResponse, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/microservices?application=%s", application), nil)
	if err != nil {
		return nil, err
	}
	response := new(MicroserviceListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetSystemMicroservicesByApplication returns a list of microservices in a specific application using Controller REST API
func (clt *Client) GetSystemMicroservicesByApplication(application string) (*MicroserviceListResponse, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/microservices/system?application=%s", application), nil)
	if err != nil {
		return nil, err
	}
	response := new(MicroserviceListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// getAllMicroservices returns all microservices on the Controller across all applications
func (clt *Client) getAllMicroservices() (*MicroserviceListResponse, error) {
	body, err := clt.doRequest("GET", "/microservices", nil)
	if err != nil {
		return nil, err
	}
	response := new(MicroserviceListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetAllSystemMicroservices returns all system microservices on the Controller across all system applications
func (clt *Client) GetAllSystemMicroservices() (*MicroserviceListResponse, error) {
	body, err := clt.doRequest("GET", "/microservices/system", nil)
	if err != nil {
		return nil, err
	}
	response := new(MicroserviceListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetAllMicroservices returns all microservices on the Controller across all applications.
func (clt *Client) GetAllMicroservices() (*MicroserviceListResponse, error) {
	return clt.getAllMicroservices()
}

// GetMicroservicePortMapping retrieves a microservice port mappings using Controller REST API
func (clt *Client) GetMicroservicePortMapping(uuid string) (*MicroservicePortMappingListResponse, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/microservices/%s/port-mapping", uuid), nil)
	if err != nil {
		return nil, err
	}

	response := new(MicroservicePortMappingListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// DeleteMicroservicePortMapping deletes a microservice port mapping using Controller REST API
func (clt *Client) DeleteMicroservicePortMapping(uuid string, portMapping *MicroservicePortMappingInfo) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/microservices/%s/port-mapping/%v", uuid, portMapping.Internal), nil)
	return err
}

// CreateMicroservicePortMapping creates a microservice port mapping using Controller REST API
func (clt *Client) CreateMicroservicePortMapping(uuid string, portMapping *MicroservicePortMappingInfo) error {
	_, err := clt.doRequest("POST", fmt.Sprintf("/microservices/%s/port-mapping", uuid), portMapping)
	return err
}

// CreateMicroserviceRoute creates a microservice route using Controller REST API.
// Deprecated: Controller no longer exposes /api/v3/microservices/:uuid/routes.
// Use NATS for messaging. CreateMicroserviceRoute returns ErrRoutesNotSupported.
func (clt *Client) CreateMicroserviceRoute(uuid, destUUID string) error {
	return ErrRoutesNotSupported
}

// DeleteMicroserviceRoute deletes a microservice route using Controller REST API.
// Deprecated: Controller no longer exposes /api/v3/microservices/:uuid/routes.
// Use NATS for messaging. DeleteMicroserviceRoute returns ErrRoutesNotSupported.
func (clt *Client) DeleteMicroserviceRoute(uuid, destUUID string) error {
	return ErrRoutesNotSupported
}

// UpdateMicroserviceRoutes syncs microservice routes to the given list.
// Deprecated: Controller no longer exposes /api/v3/microservices/:uuid/routes.
// Use NATS for messaging. UpdateMicroserviceRoutes returns ErrRoutesNotSupported.
func (clt *Client) UpdateMicroserviceRoutes(uuid string, currentRoutes, newRoutes []string) error {
	return ErrRoutesNotSupported
}

// UpdateMicroserviceFromYAML updates a microservice using the Controller REST API
// It sends the yaml file to Controller REST API
func (clt *Client) UpdateMicroserviceFromYAML(uuid string, file io.Reader) (*MicroserviceInfo, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, _ := writer.CreateFormFile("microservice", "microservice.yaml")
	_, err := io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	_, err = clt.doRequestWithHeaders("PATCH", fmt.Sprintf("/microservices/yaml/%s", uuid), requestBody, headers)
	if err != nil {
		return nil, err
	}
	return clt.GetMicroserviceByID(uuid)
}

// UpdateSystemMicroserviceFromYAML updates a system microservice using the Controller REST API
// It sends the yaml file to Controller REST API
func (clt *Client) UpdateSystemMicroserviceFromYAML(uuid string, file io.Reader) (*MicroserviceInfo, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, _ := writer.CreateFormFile("microservice", "microservice.yaml")
	_, err := io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	_, err = clt.doRequestWithHeaders("PATCH", fmt.Sprintf("/microservices/system/yaml/%s", uuid), requestBody, headers)
	if err != nil {
		return nil, err
	}
	return clt.GetSystemMicroserviceByID(uuid)
}

// DeleteMicroservice deletes a microservice using Controller REST API
func (clt *Client) DeleteMicroservice(uuid string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/microservices/%s", uuid), nil)
	return err
}

// RebuildsMicroservice rebuilds a microservice using Controller REST API
func (clt *Client) RebuildsMicroservice(uuid string) error {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/microservices/%s/rebuild", uuid), nil)
	return err
}

// RebuildsSystemMicroservice rebuilds a system microservice using Controller REST API
func (clt *Client) RebuildsSystemMicroservice(uuid string) error {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/microservices/system/%s/rebuild", uuid), nil)
	return err
}

// StartMicroservice starts a microservice using Controller REST API
func (clt *Client) StartMicroservice(uuid string) error {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/microservices/%s/start", uuid), nil)
	return err
}

// StopMicroservice stops a microservice using Controller REST API
func (clt *Client) StopMicroservice(uuid string) error {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/microservices/%s/stop", uuid), nil)
	return err
}
