package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
)

// Secrets
func (clt *Client) CreateSecret(request *SecretCreateRequest) error {
	_, err := clt.doRequest("POST", "/secrets", request)
	return err
}

func (clt *Client) CreateSecretFromYaml(file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("secret", "secret.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("POST", "/secrets/yaml", requestBody, headers)
	return err
}

func (clt *Client) UpdateSecret(name string, request *SecretUpdateRequest) error {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/secrets/%s", name), request)
	return err
}

func (clt *Client) UpdateSecretFromYaml(name string, file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("secret", "secret.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("PATCH", fmt.Sprintf("/secrets/yaml/%s", name), requestBody, headers)
	return err
}

func (clt *Client) GetSecret(name string) (*SecretInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/secrets/%s", name), nil)
	if err != nil {
		return nil, err
	}
	secret := new(SecretInfo)
	if err = json.Unmarshal(body, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func (clt *Client) ListSecrets() (*SecretListResponse, error) {
	body, err := clt.doRequest("GET", "/secrets", nil)
	if err != nil {
		return nil, err
	}

	// First try to unmarshal as an object with secrets field
	response := &SecretListResponse{}
	if err = json.Unmarshal(body, response); err != nil {
		// If that fails, try to unmarshal as an array
		var secrets []SecretInfo
		if err = json.Unmarshal(body, &secrets); err != nil {
			return nil, err
		}
		response = &SecretListResponse{
			Secrets: secrets,
		}
	}
	return response, nil
}

func (clt *Client) DeleteSecret(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/secrets/%s", name), nil)
	return err
}

// Services
func (clt *Client) CreateService(request *ServiceCreateRequest) error {
	_, err := clt.doRequest("POST", "/services", request)
	return err
}

func (clt *Client) CreateServiceFromYaml(file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("service", "service.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("POST", "/services/yaml", requestBody, headers)
	return err
}

func (clt *Client) UpdateService(name string, request *ServiceUpdateRequest) error {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/services/%s", name), request)
	return err
}

func (clt *Client) UpdateServiceFromYaml(name string, file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("service", "service.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("PATCH", fmt.Sprintf("/services/yaml/%s", name), requestBody, headers)
	return err
}

func (clt *Client) GetService(name string) (*ServiceInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/services/%s", name), nil)
	if err != nil {
		return nil, err
	}
	service := new(ServiceInfo)
	if err = json.Unmarshal(body, service); err != nil {
		return nil, err
	}
	return service, nil
}

func (clt *Client) ListServices() (*ServiceListResponse, error) {
	body, err := clt.doRequest("GET", "/services", nil)
	if err != nil {
		return nil, err
	}

	// First try to unmarshal as an object with services field
	response := &ServiceListResponse{}
	if err = json.Unmarshal(body, response); err != nil {
		// If that fails, try to unmarshal as an array
		var services []ServiceInfo
		if err = json.Unmarshal(body, &services); err != nil {
			return nil, err
		}
		response = &ServiceListResponse{
			Services: services,
		}
	}
	return response, nil
}

func (clt *Client) DeleteService(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/services/%s", name), nil)
	return err
}

// ConfigMaps
func (clt *Client) CreateConfigMap(request *ConfigMapCreateRequest) error {
	_, err := clt.doRequest("POST", "/configmaps", request)
	return err
}

func (clt *Client) CreateConfigMapFromYaml(file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("configMap", "configMap.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("POST", "/configmaps/yaml", requestBody, headers)
	return err
}

func (clt *Client) UpdateConfigMap(name string, request *ConfigMapUpdateRequest) error {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/configmaps/%s", name), request)
	return err
}

func (clt *Client) UpdateConfigMapFromYaml(name string, file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("configMap", "configMap.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("PATCH", fmt.Sprintf("/configmaps/yaml/%s", name), requestBody, headers)
	return err
}

func (clt *Client) GetConfigMap(name string) (*ConfigMapInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/configmaps/%s", name), nil)
	if err != nil {
		return nil, err
	}
	configMap := new(ConfigMapInfo)
	if err = json.Unmarshal(body, configMap); err != nil {
		return nil, err
	}
	return configMap, nil
}

func (clt *Client) ListConfigMaps() (*ConfigMapListResponse, error) {
	body, err := clt.doRequest("GET", "/configmaps", nil)
	if err != nil {
		return nil, err
	}

	// First try to unmarshal as an object with configMaps field
	response := &ConfigMapListResponse{}
	if err = json.Unmarshal(body, response); err != nil {
		// If that fails, try to unmarshal as an array
		var configMaps []ConfigMapInfo
		if err = json.Unmarshal(body, &configMaps); err != nil {
			return nil, err
		}
		response = &ConfigMapListResponse{
			ConfigMaps: configMaps,
		}
	}
	return response, nil
}

func (clt *Client) DeleteConfigMap(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/configmaps/%s", name), nil)
	return err
}

// VolumeMounts
func (clt *Client) CreateVolumeMount(request *VolumeMountCreateRequest) error {
	_, err := clt.doRequest("POST", "/volumeMounts", request)
	return err
}

func (clt *Client) CreateVolumeMountFromYaml(file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("volumeMount", "volumeMount.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("POST", "/volumeMounts/yaml", requestBody, headers)
	return err
}

func (clt *Client) UpdateVolumeMount(name string, request *VolumeMountUpdateRequest) error {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/volumeMounts/%s", name), request)
	return err
}

func (clt *Client) UpdateVolumeMountFromYaml(name string, file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("volumeMount", "volumeMount.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("PATCH", fmt.Sprintf("/volumeMounts/yaml/%s", name), requestBody, headers)
	return err
}

func (clt *Client) GetVolumeMount(name string) (*VolumeMountInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/volumeMounts/%s", name), nil)
	if err != nil {
		return nil, err
	}
	volumeMount := new(VolumeMountInfo)
	if err = json.Unmarshal(body, volumeMount); err != nil {
		return nil, err
	}
	return volumeMount, nil
}

func (clt *Client) ListVolumeMounts() (*VolumeMountListResponse, error) {
	body, err := clt.doRequest("GET", "/volumeMounts", nil)
	if err != nil {
		return nil, err
	}

	// First try to unmarshal as an object with volumeMounts field
	response := &VolumeMountListResponse{}
	if err = json.Unmarshal(body, response); err != nil {
		// If that fails, try to unmarshal as an array
		var volumeMounts []VolumeMountInfo
		if err = json.Unmarshal(body, &volumeMounts); err != nil {
			return nil, err
		}
		response = &VolumeMountListResponse{
			VolumeMounts: volumeMounts,
		}
	}
	return response, nil
}

func (clt *Client) DeleteVolumeMount(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/volumeMounts/%s", name), nil)
	return err
}

func (clt *Client) LinkVolumeMount(request *VolumeMountLinkRequest) error {
	_, err := clt.doRequest("POST", fmt.Sprintf("/volumeMounts/%s/link", request.Name), request)
	return err
}

func (clt *Client) UnlinkVolumeMount(request *VolumeMountUnlinkRequest) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/volumeMounts/%s/link", request.Name), request)
	return err
}

// Certificates
func (clt *Client) CreateCA(request *CACreateRequest) error {
	_, err := clt.doRequest("POST", "/certificates/ca", request)
	return err
}

func (clt *Client) GetCA(name string) (*CAInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/certificates/ca/%s", name), nil)
	if err != nil {
		return nil, err
	}
	ca := new(CAInfo)
	if err = json.Unmarshal(body, ca); err != nil {
		return nil, err
	}
	return ca, nil
}

func (clt *Client) ListCAs() (*CAListResponse, error) {
	body, err := clt.doRequest("GET", "/certificates/ca", nil)
	if err != nil {
		return nil, err
	}

	// First try to unmarshal as an object with cas field
	response := &CAListResponse{}
	if err = json.Unmarshal(body, response); err != nil {
		// If that fails, try to unmarshal as an array
		var cas []CAInfo
		if err = json.Unmarshal(body, &cas); err != nil {
			return nil, err
		}
		response = &CAListResponse{
			CAs: cas,
		}
	}
	return response, nil
}

func (clt *Client) DeleteCA(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/certificates/ca/%s", name), nil)
	return err
}

func (clt *Client) CreateCertificate(request *CertificateCreateRequest) error {
	_, err := clt.doRequest("POST", "/certificates", request)
	return err
}

func (clt *Client) CreateCertificateFromYaml(file io.Reader) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("certificate", "certificate.yaml")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	_, err = clt.doRequestWithHeaders("POST", "/certificates/yaml", requestBody, headers)
	return err
}

func (clt *Client) GetCertificate(name string) (*CertificateInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/certificates/%s", name), nil)
	if err != nil {
		return nil, err
	}
	cert := new(CertificateInfo)
	if err = json.Unmarshal(body, cert); err != nil {
		return nil, err
	}
	return cert, nil
}

func (clt *Client) ListCertificates() (*CertificateListResponse, error) {
	body, err := clt.doRequest("GET", "/certificates", nil)
	if err != nil {
		return nil, err
	}

	// First try to unmarshal as an object with certificates field
	response := &CertificateListResponse{}
	if err = json.Unmarshal(body, response); err != nil {
		// If that fails, try to unmarshal as an array
		var certificates []CertificateInfo
		if err = json.Unmarshal(body, &certificates); err != nil {
			return nil, err
		}
		response = &CertificateListResponse{
			Certificates: certificates,
		}
	}
	return response, nil
}

func (clt *Client) ListExpiringCertificates() (*CertificateListResponse, error) {
	body, err := clt.doRequest("GET", "/certificates/expiring", nil)
	if err != nil {
		return nil, err
	}
	response := new(CertificateListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (clt *Client) DeleteCertificate(name string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/certificates/%s", name), nil)
	return err
}

func (clt *Client) RenewCertificate(name string) error {
	_, err := clt.doRequest("POST", fmt.Sprintf("/certificates/%s/renew", name), nil)
	return err
}

func (clt *Client) AttachExecMicroservice(request *AttachExecMicroserviceRequest) error {
	_, err := clt.doRequest("POST", fmt.Sprintf("/microservices/%s/exec", request.UUID), request)
	return err
}

func (clt *Client) DetachExecMicroservice(request *DetachExecMicroserviceRequest) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/microservices/%s/exec", request.UUID), request)
	return err
}

func (clt *Client) AttachExecSystemMicroservice(request *AttachExecMicroserviceRequest) error {
	_, err := clt.doRequest("POST", fmt.Sprintf("/microservices/system/%s/exec", request.UUID), request)
	return err
}

func (clt *Client) DetachExecSystemMicroservice(request *DetachExecMicroserviceRequest) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/microservices/system/%s/exec", request.UUID), request)
	return err
}

func (clt *Client) AttachExecToAgent(request *AttachExecToAgentRequest) error {
	_, err := clt.doRequest("POST", fmt.Sprintf("/iofog/%s/exec", request.UUID), request)
	return err
}

func (clt *Client) DetachExecFromAgent(request *DetachExecFromAgentRequest) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/iofog/%s/exec", request.UUID), request)
	return err
}
