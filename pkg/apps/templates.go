package apps

import (
	"bytes"
	"errors"
	"net/url"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
	"gopkg.in/yaml.v2"
)

type applicationTemplateExecutor struct {
	controller IofogController
	baseURL    *url.URL
	template   any
	name       string
	apiVersion string
	client     *client.Client
}

func newApplicationTemplateExecutor(controller IofogController, controllerBaseURL *url.URL, template any, name string, opts ...DeployOption) *applicationTemplateExecutor {
	resolved := resolveDeployOptions(opts...)
	return &applicationTemplateExecutor{
		controller: controller,
		baseURL:    controllerBaseURL,
		name:       name,
		template:   template,
		apiVersion: resolved.apiVersion,
	}
}

func (exe *applicationTemplateExecutor) execute() error {
	// Init remote resources
	if err := exe.init(); err != nil {
		return err
	}

	// Deploy application
	return exe.deploy()
}

func (exe *applicationTemplateExecutor) init() (err error) {
	if exe.controller.Token != "" {
		exe.client, err = client.NewWithToken(client.Options{BaseURL: exe.baseURL}, exe.controller.Token)
	} else {
		exe.client, err = client.SessionLogin(client.Options{BaseURL: exe.baseURL}, exe.controller.RefreshToken, exe.controller.Email, exe.controller.Password)
	}

	return err
}

func (exe *applicationTemplateExecutor) deploy() error {
	file := IofogHeader{
		APIVersion: exe.apiVersion,
		Kind:       ApplicationTemplateKind,
		Metadata: HeaderMetadata{
			Name: exe.name,
		},
		Spec: exe.template,
	}
	yamlBytes, err := yaml.Marshal(file)
	if err != nil {
		return err
	}
	existingAppTemplate, err := exe.client.GetApplicationTemplate(exe.name)
	// If not notfound error, return error
	notFoundError := &client.NotFoundError{}
	if errors.As(err, &notFoundError) {
		return err
	}
	if existingAppTemplate == nil {
		if _, err := exe.client.CreateApplicationTemplateFromYAML(bytes.NewReader(yamlBytes)); err != nil {
			return err
		}
		return nil
	}
	if _, err := exe.client.UpdateApplicationTemplateFromYAML(exe.name, bytes.NewReader(yamlBytes)); err != nil {
		return err
	}
	return nil
}
