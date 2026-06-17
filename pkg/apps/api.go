package apps

import (
	"net/url"
)

func DeployApplicationTemplate(controller IofogController, controllerBaseURL *url.URL, template interface{}, name string) error {
	exe := newApplicationTemplateExecutor(controller, controllerBaseURL, template, name)
	return exe.execute()
}

func DeployApplication(controller IofogController, application interface{}, name string) error {
	exe := newApplicationExecutor(controller, application, name)
	return exe.execute()
}

func DeployMicroservice(controller IofogController, microservice interface{}, appName, name string) error {
	exe := newMicroserviceExecutor(controller, microservice, appName, name)
	return exe.execute()
}
