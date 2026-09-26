package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

// Roles

// ListRoles retrieves all roles from the Controller REST API
func (clt *Client) ListRoles() (*RoleListResponse, error) {
	body, err := clt.doRequest("GET", "/roles", nil)
	if err != nil {
		return nil, err
	}
	response := new(RoleListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetRole retrieves a single role by name from the Controller REST API
func (clt *Client) GetRole(name string) (*RoleInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/roles/%s", name), nil)
	if err != nil {
		return nil, err
	}
	response := new(RoleResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return &response.Role, nil
}

// CreateRole creates a new role using the Controller REST API
func (clt *Client) CreateRole(request *RoleCreateRequest) (*RoleInfo, error) {
	body, err := clt.doRequest("POST", "/roles", request)
	if err != nil {
		return nil, err
	}
	response := new(RoleResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return &response.Role, nil
}

// CreateRoleFromYaml creates a role from a YAML file using the Controller REST API
func (clt *Client) CreateRoleFromYaml(file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("role", "role.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("POST", "/roles/yaml", requestBody, headers)
	return err
}

// UpdateRole updates a role by name using the Controller REST API
func (clt *Client) UpdateRole(name string, request *RoleUpdateRequest) (*RoleInfo, error) {
	body, err := clt.doRequest("PATCH", fmt.Sprintf("/roles/%s", name), request)
	if err != nil {
		return nil, err
	}
	response := new(RoleResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return &response.Role, nil
}

// UpdateRoleFromYaml updates a role from a YAML file using the Controller REST API
func (clt *Client) UpdateRoleFromYaml(name string, file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("role", "role.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("PATCH", fmt.Sprintf("/roles/yaml/%s", name), requestBody, headers)
	return err
}

// DeleteRole deletes a role by name using the Controller REST API
func (clt *Client) DeleteRole(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/roles/%s", name), nil)
	return err
}

// RoleBindings

// ListRoleBindings retrieves all role bindings from the Controller REST API
func (clt *Client) ListRoleBindings() (*RoleBindingListResponse, error) {
	body, err := clt.doRequest("GET", "/rolebindings", nil)
	if err != nil {
		return nil, err
	}
	response := new(RoleBindingListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetRoleBinding retrieves a single role binding by name from the Controller REST API
func (clt *Client) GetRoleBinding(name string) (*RoleBindingInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/rolebindings/%s", name), nil)
	if err != nil {
		return nil, err
	}
	response := new(RoleBindingResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return &response.Binding, nil
}

// CreateRoleBinding creates a new role binding using the Controller REST API
func (clt *Client) CreateRoleBinding(request *RoleBindingCreateRequest) (*RoleBindingInfo, error) {
	body, err := clt.doRequest("POST", "/rolebindings", request)
	if err != nil {
		return nil, err
	}
	response := new(RoleBindingResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return &response.Binding, nil
}

// CreateRoleBindingFromYaml creates a role binding from a YAML file using the Controller REST API
func (clt *Client) CreateRoleBindingFromYaml(file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("rolebinding", "rolebinding.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("POST", "/rolebindings/yaml", requestBody, headers)
	return err
}

// UpdateRoleBinding updates a role binding by name using the Controller REST API
func (clt *Client) UpdateRoleBinding(name string, request *RoleBindingUpdateRequest) (*RoleBindingInfo, error) {
	body, err := clt.doRequest("PATCH", fmt.Sprintf("/rolebindings/%s", name), request)
	if err != nil {
		return nil, err
	}
	response := new(RoleBindingResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return &response.Binding, nil
}

// UpdateRoleBindingFromYaml updates a role binding from a YAML file using the Controller REST API
func (clt *Client) UpdateRoleBindingFromYaml(name string, file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("rolebinding", "rolebinding.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("PATCH", fmt.Sprintf("/rolebindings/yaml/%s", name), requestBody, headers)
	return err
}

// DeleteRoleBinding deletes a role binding by name using the Controller REST API
func (clt *Client) DeleteRoleBinding(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/rolebindings/%s", name), nil)
	return err
}

// ServiceAccounts
// Service accounts are application-scoped: identify by (applicationName, name).
// Controller routes: GET/POST /serviceaccounts, GET/PATCH/DELETE /serviceaccounts/:appName/:name.

// ListServiceAccounts retrieves all service accounts from the Controller REST API.
// If applicationName is non-empty, only service accounts in that application are returned.
func (clt *Client) ListServiceAccounts(applicationName string) (*ServiceAccountListResponse, error) {
	path := "/serviceaccounts"
	if applicationName != "" {
		path = "/serviceaccounts?applicationName=" + url.QueryEscape(applicationName)
	}
	body, err := clt.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	response := new(ServiceAccountListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetServiceAccount retrieves a single service account by application name and name from the Controller REST API.
func (clt *Client) GetServiceAccount(appName, name string) (*ServiceAccountInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/serviceaccounts/%s/%s", appName, name), nil)
	if err != nil {
		return nil, err
	}
	response := new(ServiceAccountResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return &response.ServiceAccount, nil
}

// CreateServiceAccount creates a new service account using the Controller REST API
func (clt *Client) CreateServiceAccount(request *ServiceAccountCreateRequest) (*ServiceAccountInfo, error) {
	body, err := clt.doRequest("POST", "/serviceaccounts", request)
	if err != nil {
		return nil, err
	}
	response := new(ServiceAccountResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return &response.ServiceAccount, nil
}

// CreateServiceAccountFromYaml creates a service account from a YAML file using the Controller REST API
func (clt *Client) CreateServiceAccountFromYaml(file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("serviceaccount", "serviceaccount.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("POST", "/serviceaccounts/yaml", requestBody, headers)
	return err
}

// UpdateServiceAccount updates a service account by application name and name using the Controller REST API.
func (clt *Client) UpdateServiceAccount(appName, name string, request *ServiceAccountUpdateRequest) (*ServiceAccountInfo, error) {
	body, err := clt.doRequest("PATCH", fmt.Sprintf("/serviceaccounts/%s/%s", appName, name), request)
	if err != nil {
		return nil, err
	}
	response := new(ServiceAccountResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return &response.ServiceAccount, nil
}

// UpdateServiceAccountFromYaml updates a service account from a YAML file using the Controller REST API.
func (clt *Client) UpdateServiceAccountFromYaml(appName, name string, file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("serviceaccount", "serviceaccount.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	_ = writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("PATCH", fmt.Sprintf("/serviceaccounts/yaml/%s/%s", appName, name), requestBody, headers)
	return err
}

// DeleteServiceAccount deletes a service account by application name and name using the Controller REST API.
func (clt *Client) DeleteServiceAccount(appName, name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/serviceaccounts/%s/%s", appName, name), nil)
	return err
}
