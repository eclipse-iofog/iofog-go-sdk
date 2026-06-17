package apps

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
	"gopkg.in/yaml.v2"
)

type applicationExecutor struct {
	controller      IofogController
	app             any
	name            string
	apiVersion      string
	applicationInfo *client.ApplicationInfo
	client          *client.Client
}

func newApplicationExecutor(controller IofogController, app any, name string, opts ...DeployOption) *applicationExecutor {
	resolved := resolveDeployOptions(opts...)
	return &applicationExecutor{
		controller: controller,
		app:        app,
		name:       name,
		apiVersion: resolved.apiVersion,
	}
}

func (exe *applicationExecutor) execute() (err error) {
	// Init remote resources
	if err := exe.init(); err != nil {
		return err
	}

	// Try application API
	// Look for exisiting application
	exe.applicationInfo, err = exe.client.GetApplicationByName(exe.name)

	// If not notfound error, return error
	notFoundError := &client.NotFoundError{}
	if errors.As(err, &notFoundError) {
		return err
	}

	// Deploy application
	return exe.deploy()
}

func (exe *applicationExecutor) init() (err error) {
	baseURL, err := url.Parse(exe.controller.Endpoint)
	if err != nil {
		return fmt.Errorf(errParseControllerURL, err.Error())
	}
	if exe.controller.Token != "" {
		exe.client, err = client.NewWithToken(client.Options{BaseURL: baseURL}, exe.controller.Token)
	} else {
		exe.client, err = client.SessionLogin(client.Options{BaseURL: baseURL}, exe.controller.RefreshToken, exe.controller.Email, exe.controller.Password)
	}
	return err
}

func (exe *applicationExecutor) create() (err error) {
	file := IofogHeader{
		APIVersion: exe.apiVersion,
		Kind:       ApplicationKind,
		Metadata: HeaderMetadata{
			Name: exe.name,
		},
		Spec: exe.app,
	}
	yamlBytes, err := yaml.Marshal(file)
	if err != nil {
		return err
	}
	if _, err = exe.client.CreateApplicationFromYAML(bytes.NewReader(yamlBytes)); err != nil {
		return err
	}
	return nil
}

func (exe *applicationExecutor) update() (err error) {
	file := IofogHeader{
		APIVersion: exe.apiVersion,
		Kind:       ApplicationKind,
		Metadata: HeaderMetadata{
			Name: exe.name,
		},
		Spec: exe.app,
	}
	yamlBytes, err := yaml.Marshal(file)
	if err != nil {
		return err
	}

	if _, err = exe.client.UpdateApplicationFromYAML(exe.name, bytes.NewReader(yamlBytes)); err != nil {
		return err
	}
	return nil
}

func (exe *applicationExecutor) deploy() (err error) {
	// Existing app info retrieved in init
	if exe.applicationInfo == nil {
		if err := exe.create(); err != nil {
			return err
		}
	} else {
		if err := exe.update(); err != nil {
			return err
		}
	}

	// Start application
	if _, err = exe.client.StartApplication(exe.name); err != nil {
		return err
	}
	return nil
}
