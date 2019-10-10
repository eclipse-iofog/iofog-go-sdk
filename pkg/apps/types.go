/*
 *  *******************************************************************************
 *  * Copyright (c) 2019 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package apps

// HeaderMetadata contains k8s metadata
// +k8s:deepcopy-gen=true
type HeaderMetadata struct {
	Name      string
	Namespace string
}

// Deployable can be deployed using iofogctl deploy
type Deployable interface {
	Deploy(namespace string) error
}

// Kind contains available types
type Kind string

// Available kind of deploy
const (
	ApplicationKind  Kind = "iofog-application"
	MicroserviceKind Kind = "iofog-microservice"
	ControlPlaneKind Kind = "iofog-controlplane"
	AgentKind        Kind = "iofog-agent"
	ConnectorKind    Kind = "iofog-connector"
	ControllerKind   Kind = "iofog-controller"
)

// Header contains k8s yaml header
// +k8s:deepcopy-gen=true
type Header struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       Kind           `yaml:"kind"`
	Metadata   HeaderMetadata `yaml:"metadata"`
	Spec       NestedMap
}

// CatalogItem contains information about a catalog item
// +k8s:deepcopy-gen=true
type CatalogItem struct {
	ID            int
	X86           string
	ARM           string
	Registry      string
	Name          string
	Description   string
	ConfigExample string
}

// MicroserviceImages contains information about the images for a microservice
// +k8s:deepcopy-gen=true
type MicroserviceImages struct {
	CatalogID int
	X86       string
	ARM       string
	Registry  string
}

// MicroserviceAgent contains information about required agent configuration for a microservice
// +k8s:deepcopy-gen=true
type MicroserviceAgent struct {
	Name   string
	Config AgentConfiguration
}

// Microservice contains information for configuring a microservice
// +k8s:deepcopy-gen=true
type Microservice struct {
	UUID           string `yaml:"-"`
	Name           string
	Agent          MicroserviceAgent
	Images         MicroserviceImages
	Config         NestedMap
	RootHostAccess bool
	Ports          []MicroservicePortMapping   `yaml:"ports"`
	Volumes        []MicroserviceVolumeMapping `yaml:"volumes"`
	Env            []MicroserviceEnvironment   `yaml:"env"`
	Routes         []string                           `yaml:"routes,omitempty"`
	Flow           *string                            `yaml:"application,omitempty"`
	Created        string                             `yaml:"created,omitempty"`
}

type NestedMap map[string]interface{}

func (j NestedMap) DeepCopy() NestedMap {
    copy := make(NestedMap)
    deepCopyNestedMap(j, copy)
    return copy
}

func deepCopyNestedMap(src NestedMap, dest NestedMap) {
    for key, value := range src {
        switch src[key].(type) {
        case NestedMap:
            dest[key] = NestedMap{}
            deepCopyNestedMap(src[key].(NestedMap), dest[key].(NestedMap))
        default:
            dest[key] = value
        }
    }
}

// +k8s:deepcopy-gen=true
type MicroservicePortMapping struct {
	Internal   int  `json:"internal"`
	External   int  `json:"external"`
	PublicMode bool `json:"publicMode"`
}

// +k8s:deepcopy-gen=true
type MicroserviceVolumeMapping struct {
	HostDestination      string `json:"hostDestination"`
	ContainerDestination string `json:"containerDestination"`
	AccessMode           string `json:"accessMode"`
}

// +k8s:deepcopy-gen=true
type MicroserviceEnvironment struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// +k8s:deepcopy-gen=true
type AgentConfiguration struct {
	DockerURL                 *string  `json:"dockerUrl,omitempty"`
	DiskLimit                 *int64   `json:"diskLimit,omitempty"`
	DiskDirectory             *string  `json:"diskDirectory,omitempty"`
	MemoryLimit               *int64   `json:"memoryLimit,omitempty"`
	CPULimit                  *int64   `json:"cpuLimit,omitempty"`
	LogLimit                  *int64   `json:"logLimit,omitempty"`
	LogDirectory              *string  `json:"logDirectory,omitempty"`
	LogFileCount              *int64   `json:"logFileCount,omitempty"`
	StatusFrequency           *float64 `json:"statusFrequency,omitempty"`
	ChangeFrequency           *float64 `json:"changeFrequency,omitempty"`
	DeviceScanFrequency       *float64 `json:"deviceScanFrequency,omitempty"`
	BluetoothEnabled          *bool    `json:"bluetoothEnabled,omitempty"`
	WatchdogEnabled           *bool    `json:"watchdogEnabled,omitempty"`
	AbstractedHardwareEnabled *bool    `json:"abstractedHardwareEnabled,omitempty"`
}

// Microservices is a list of Microservice
// +k8s:deepcopy-gen=true
type Microservices struct {
	Microservices []Microservices
}

// Route contains information about a route from one microservice to another
// +k8s:deepcopy-gen=true
type Route struct {
	From string
	To   string
}

// Application contains information for configuring an application
// +k8s:deepcopy-gen=true
type Application struct {
	Name          string
	Microservices []Microservice
	Routes        []Route
	ID            int
}

// Applications is a list of applications
// +k8s:deepcopy-gen=true
type Applications struct {
	Applications []Application
}

// IofogController contains informations needed to connect to the controller
// +k8s:deepcopy-gen=true
type IofogController struct {
	Email    string
	Password string
	Endpoint string
	Token    string
}
