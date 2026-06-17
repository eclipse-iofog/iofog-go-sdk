package apps

import (
	"net/url"
)

func DeployApplicationTemplate(controller IofogController, controllerBaseURL *url.URL, template any, name string, opts ...DeployOption) error {
	exe := newApplicationTemplateExecutor(controller, controllerBaseURL, template, name, opts...)
	return exe.execute()
}

func DeployApplication(controller IofogController, application any, name string, opts ...DeployOption) error {
	exe := newApplicationExecutor(controller, application, name, opts...)
	return exe.execute()
}

func DeployMicroservice(controller IofogController, microservice any, appName, name string, opts ...DeployOption) error {
	exe := newMicroserviceExecutor(controller, microservice, appName, name, opts...)
	return exe.execute()
}
