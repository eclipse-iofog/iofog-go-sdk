# Changelog

## [Unreleased]


## [v2.0.0-rc1] - 2020-04-29

* Added type to volume mapping

## [v2.0.0-beta3] - 2020-04-23

* Update WaitForLoadBalancer to get hostname too
* Fix error reporting in WaitForLoadBalancer func

## [v2.0.0-beta2] - 2020-04-06

### Features

* Update GetAgentByName to take system flag
* Update ListAgents and allow users to specify filters and system flag
* Remove iofogctl resource kinds
* Update AgentConfiguration for configuring docker frequency

### Bugs

* Stop passing ListAgentsRequest into body of request

## [v2.0.0-alpha] - 2020-03-10

### Features

* Add PutDefaultProxy function to client pkg
* Set retries as optional on new client
* Add retry policy to client
* Add omitempty to optional PATCH msvc args
* Add GetAllMicroservicePublicPorts function
* Move port validation to be run in application deployment too
* Add PublicLink to msvcPortMapping and update microservice update to detect public port mapping changes
* Add DefaultRouterName constant
* Update publicPort json key
* Add router fields to AgentConfiguration struct
* Update agent info to contain router info 
* Update routerConfig in agent config yaml
* Add networkRouter to AgentConfiguration
* Add isSystem in agent and applications
* Update AgentCreateRequest to allow for configuration
* Add Agent Prune API call
* Add PORT to apps.Microservice.Container
* Make CMD optional on microservice update
* Allows CMD in microservice creation and update

### Bugs

* Fix make gen to update file in $PWD/pkg/apps

## [v1.3.0]

* Add client package to the repo
* Re-organize the repo to maintain multiple packages
  
[Unreleased]: https://github.com/eclipse-iofog/iofog-go-sdk/compare/v2.0.0-beta3..HEAD
[v2.0.0-beta3]: https://github.com/eclipse-iofog/iofog-go-sdk/compare/v2.0.0-beta2..v2.0.0-beta3
[v2.0.0-beta]: https://github.com/eclipse-iofog/iofog-go-sdk/compare/v2.0.0-alpha..v2.0.0-beta2
[v2.0.0-alpha]: https://github.com/eclipse-iofog/iofog-go-sdk/compare/v1.3.0..v2.0.0-beta
[v1.3.0]: https://github.com/eclipse-iofog/iofog-go-sdk/tree/v1.3.0
