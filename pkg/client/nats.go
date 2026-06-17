package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
)

const natsBasePath = "/nats"

// GetNatsOperator returns the NATS operator (JWT, publicKey) from the Controller.
func (clt *Client) GetNatsOperator() (NatsOperatorResponse, error) {
	var out NatsOperatorResponse
	body, err := clt.doRequest("GET", natsBasePath+"/operator", nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// RotateNatsOperator rotates the NATS operator and returns the new operator info.
func (clt *Client) RotateNatsOperator() (NatsOperatorResponse, error) {
	var out NatsOperatorResponse
	body, err := clt.doRequest("POST", natsBasePath+"/operator/rotate", nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// GetNatsBootstrap returns NATS bootstrap data (operator JWT/publicKey/seed and system user creds).
// Only available when the Controller runs on the Kubernetes control plane; returns 403 otherwise.
// Intended for use by the K8s operator to bootstrap NATS.
func (clt *Client) GetNatsBootstrap() (NatsBootstrapResponse, error) {
	var out NatsBootstrapResponse
	body, err := clt.doRequest("GET", natsBasePath+"/bootstrap", nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// GetNatsHub returns the default NATS hub configuration.
func (clt *Client) GetNatsHub() (NatsHubResponse, error) {
	var out NatsHubResponse
	body, err := clt.doRequest("GET", natsBasePath+"/hub", nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// UpsertNatsHub creates or updates the default NATS hub.
func (clt *Client) UpsertNatsHub(req *NatsHubRequest) (NatsHubResponse, error) {
	var out NatsHubResponse
	body, err := clt.doRequest("PUT", natsBasePath+"/hub", req)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// ListNatsAccounts returns all NATS accounts.
func (clt *Client) ListNatsAccounts() (NatsListAccountsResponse, error) {
	var out NatsListAccountsResponse
	body, err := clt.doRequest("GET", natsBasePath+"/accounts", nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// ListNatsUsers returns all NATS users across accounts.
func (clt *Client) ListNatsUsers() (NatsListUsersResponse, error) {
	var out NatsListUsersResponse
	body, err := clt.doRequest("GET", natsBasePath+"/users", nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// ListNatsAccountRules returns all NATS account rules.
func (clt *Client) ListNatsAccountRules() (NatsListAccountRulesResponse, error) {
	var out NatsListAccountRulesResponse
	body, err := clt.doRequest("GET", natsBasePath+"/account-rules", nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// CreateNatsAccountRule creates an account rule with a JSON body.
func (clt *Client) CreateNatsAccountRule(payload NatsAccountRulePayload) (NatsRuleInfo, error) {
	var out NatsRuleInfo
	body, err := clt.doRequest("POST", natsBasePath+"/account-rules", payload)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// CreateNatsAccountRuleFromYAML creates an account rule from a YAML file.
func (clt *Client) CreateNatsAccountRuleFromYAML(file io.Reader) (NatsRuleInfo, error) {
	var out NatsRuleInfo
	reqBody := &bytes.Buffer{}
	w := multipart.NewWriter(reqBody)
	part, err := w.CreateFormFile("natsAccountRule", "natsAccountRule.yaml")
	if err != nil {
		return out, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return out, err
	}
	if err = w.Close(); err != nil {
		return out, err
	}
	headers := map[string]string{"Content-Type": w.FormDataContentType()}
	body, err := clt.doRequestWithHeaders("POST", natsBasePath+"/account-rules/yaml", reqBody, headers)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// UpdateNatsAccountRule updates an account rule by name with a JSON body.
func (clt *Client) UpdateNatsAccountRule(ruleName string, payload NatsAccountRulePayload) (NatsRuleInfo, error) {
	var out NatsRuleInfo
	body, err := clt.doRequest("PATCH", natsBasePath+"/account-rules/"+ruleName, payload)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// UpdateNatsAccountRuleFromYAML updates an account rule by name from a YAML file.
func (clt *Client) UpdateNatsAccountRuleFromYAML(ruleName string, file io.Reader) (NatsRuleInfo, error) {
	var out NatsRuleInfo
	reqBody := &bytes.Buffer{}
	w := multipart.NewWriter(reqBody)
	part, err := w.CreateFormFile("natsAccountRule", "natsAccountRule.yaml")
	if err != nil {
		return out, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return out, err
	}
	if err = w.Close(); err != nil {
		return out, err
	}
	headers := map[string]string{"Content-Type": w.FormDataContentType()}
	body, err := clt.doRequestWithHeaders("PATCH", natsBasePath+"/account-rules/yaml/"+ruleName, reqBody, headers)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// DeleteNatsAccountRule deletes an account rule by name.
func (clt *Client) DeleteNatsAccountRule(ruleName string) error {
	_, err := clt.doRequest("DELETE", natsBasePath+"/account-rules/"+ruleName, nil)
	return err
}

// ListNatsUserRules returns all NATS user rules.
func (clt *Client) ListNatsUserRules() (NatsListUserRulesResponse, error) {
	var out NatsListUserRulesResponse
	body, err := clt.doRequest("GET", natsBasePath+"/user-rules", nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// CreateNatsUserRule creates a user rule with a JSON body.
func (clt *Client) CreateNatsUserRule(payload NatsUserRulePayload) (NatsRuleInfo, error) {
	var out NatsRuleInfo
	body, err := clt.doRequest("POST", natsBasePath+"/user-rules", payload)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// CreateNatsUserRuleFromYAML creates a user rule from a YAML file.
func (clt *Client) CreateNatsUserRuleFromYAML(file io.Reader) (NatsRuleInfo, error) {
	var out NatsRuleInfo
	reqBody := &bytes.Buffer{}
	w := multipart.NewWriter(reqBody)
	part, err := w.CreateFormFile("natsUserRule", "natsUserRule.yaml")
	if err != nil {
		return out, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return out, err
	}
	if err = w.Close(); err != nil {
		return out, err
	}
	headers := map[string]string{"Content-Type": w.FormDataContentType()}
	body, err := clt.doRequestWithHeaders("POST", natsBasePath+"/user-rules/yaml", reqBody, headers)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// UpdateNatsUserRule updates a user rule by name with a JSON body.
func (clt *Client) UpdateNatsUserRule(ruleName string, payload NatsUserRulePayload) (NatsRuleInfo, error) {
	var out NatsRuleInfo
	body, err := clt.doRequest("PATCH", natsBasePath+"/user-rules/"+ruleName, payload)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// UpdateNatsUserRuleFromYAML updates a user rule by name from a YAML file.
func (clt *Client) UpdateNatsUserRuleFromYAML(ruleName string, file io.Reader) (NatsRuleInfo, error) {
	var out NatsRuleInfo
	reqBody := &bytes.Buffer{}
	w := multipart.NewWriter(reqBody)
	part, err := w.CreateFormFile("natsUserRule", "natsUserRule.yaml")
	if err != nil {
		return out, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return out, err
	}
	if err = w.Close(); err != nil {
		return out, err
	}
	headers := map[string]string{"Content-Type": w.FormDataContentType()}
	body, err := clt.doRequestWithHeaders("PATCH", natsBasePath+"/user-rules/yaml/"+ruleName, reqBody, headers)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// DeleteNatsUserRule deletes a user rule by name.
func (clt *Client) DeleteNatsUserRule(ruleName string) error {
	_, err := clt.doRequest("DELETE", natsBasePath+"/user-rules/"+ruleName, nil)
	return err
}

// GetNatsAccount returns the NATS account for the given application name.
func (clt *Client) GetNatsAccount(appName string) (NatsAccountInfo, error) {
	var out NatsAccountInfo
	body, err := clt.doRequest("GET", natsBasePath+"/accounts/"+appName, nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// EnsureNatsAccount ensures a NATS account exists for the application (enable NATS for app).
func (clt *Client) EnsureNatsAccount(appName string, req *NatsEnsureAccountRequest) (NatsAccountInfo, error) {
	var out NatsAccountInfo
	body, err := clt.doRequest("POST", natsBasePath+"/accounts/"+appName, req)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// ListNatsAccountUsers returns NATS users for the given application's account.
func (clt *Client) ListNatsAccountUsers(appName string) (NatsListAccountUsersResponse, error) {
	var out NatsListAccountUsersResponse
	body, err := clt.doRequest("GET", natsBasePath+"/accounts/"+appName+"/users", nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// CreateNatsUser creates a NATS user for the given application's account.
func (clt *Client) CreateNatsUser(appName string, req *NatsCreateUserRequest) (NatsUserInfo, error) {
	var out NatsUserInfo
	body, err := clt.doRequest("POST", natsBasePath+"/accounts/"+appName+"/users", req)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// GetNatsUserCreds returns the credentials (creds file as base64) for a NATS user.
func (clt *Client) GetNatsUserCreds(appName, userName string) (NatsUserCredsResponse, error) {
	var out NatsUserCredsResponse
	path := fmt.Sprintf("%s/accounts/%s/users/%s/creds", natsBasePath, appName, userName)
	body, err := clt.doRequest("GET", path, nil)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// DeleteNatsUser deletes a NATS user for the given application.
func (clt *Client) DeleteNatsUser(appName, userName string) error {
	path := fmt.Sprintf("%s/accounts/%s/users/%s", natsBasePath, appName, userName)
	_, err := clt.doRequest("DELETE", path, nil)
	return err
}

// CreateNatsMqttBearer creates an MQTT bearer user for the given application.
func (clt *Client) CreateNatsMqttBearer(appName string, req *NatsCreateMqttBearerRequest) (NatsUserInfo, error) {
	var out NatsUserInfo
	body, err := clt.doRequest("POST", natsBasePath+"/accounts/"+appName+"/mqtt-bearer", req)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	return out, err
}

// DeleteNatsMqttBearer deletes an MQTT bearer user for the given application.
func (clt *Client) DeleteNatsMqttBearer(appName, userName string) error {
	path := fmt.Sprintf("%s/accounts/%s/mqtt-bearer/%s", natsBasePath, appName, userName)
	_, err := clt.doRequest("DELETE", path, nil)
	return err
}
