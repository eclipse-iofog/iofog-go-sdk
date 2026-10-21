package apps

import (
	"bytes"
	"net/url"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
	"gopkg.in/yaml.v2"
)

type microserviceTemplateExecutor struct {
	controller IofogController
	baseURL    *url.URL
	template   any
	name       string
	apiVersion string
	client     *client.Client
}

func newMicroserviceTemplateExecutor(controller IofogController, controllerBaseURL *url.URL, template any, name string, opts ...DeployOption) *microserviceTemplateExecutor {
	resolved := resolveDeployOptions(opts...)
	return &microserviceTemplateExecutor{
		controller: controller,
		baseURL:    controllerBaseURL,
		name:       name,
		template:   template,
		apiVersion: resolved.apiVersion,
	}
}

func (exe *microserviceTemplateExecutor) execute() error {
	// Init remote resources
	if err := exe.init(); err != nil {
		return err
	}

	// Deploy microservice template
	return exe.deploy()
}

func (exe *microserviceTemplateExecutor) init() (err error) {
	if exe.controller.Token != "" {
		exe.client, err = client.NewWithToken(client.Options{BaseURL: exe.baseURL}, exe.controller.Token)
	} else {
		exe.client, err = client.SessionLogin(client.Options{BaseURL: exe.baseURL}, exe.controller.RefreshToken, exe.controller.Email, exe.controller.Password)
	}

	return err
}

func (exe *microserviceTemplateExecutor) marshalYAML() ([]byte, error) {
	file := IofogHeader{
		APIVersion: exe.apiVersion,
		Kind:       MicroserviceTemplateKind,
		Metadata: HeaderMetadata{
			Name: exe.name,
		},
		Spec: exe.template,
	}
	return yaml.Marshal(file)
}

func (exe *microserviceTemplateExecutor) deploy() error {
	yamlBytes, err := exe.marshalYAML()
	if err != nil {
		return err
	}
	existing, err := exe.client.GetMicroserviceTemplate(exe.name)
	if err = lookupErrorAllowNotFound(err); err != nil {
		return err
	}
	if existing == nil {
		if _, err := exe.client.CreateMicroserviceTemplateFromYAML(bytes.NewReader(yamlBytes)); err != nil {
			return err
		}
		return nil
	}
	if _, err := exe.client.UpdateMicroserviceTemplateFromYAML(exe.name, bytes.NewReader(yamlBytes)); err != nil {
		return err
	}
	return nil
}
