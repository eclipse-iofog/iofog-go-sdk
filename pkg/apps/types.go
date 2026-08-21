package apps

import (
	"encoding/json"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apimachinery/pkg/runtime"
)

// HeaderMetadata contains k8s metadata
// +k8s:deepcopy-gen=true
type HeaderMetadata struct {
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`
}

// Kind contains available types
type Kind string

// IofogHeader represent the file structure
type IofogHeader Header

// Available kind of deploy
const (
	ApplicationKind         Kind = "Application"
	ApplicationTemplateKind Kind = "ApplicationTemplate"
	MicroserviceKind        Kind = "Microservice"
	// RouteKind is deprecated. Controller no longer exposes routing; use NATS for messaging.
	RouteKind Kind = "Route"
)

// Header contains k8s yaml header
type Header struct {
	APIVersion string         `yaml:"apiVersion" json:"apiVersion"`
	Kind       Kind           `yaml:"kind" json:"kind"`
	Metadata   HeaderMetadata `yaml:"metadata" json:"metadata"`
	Spec       any            `yaml:"spec" json:"spec"`
}

// CatalogItem contains information about a catalog item
// +k8s:deepcopy-gen=true
type CatalogItem struct {
	ID            int    `yaml:"id" json:"id"`
	AMD64         string `yaml:"amd64" json:"amd64"`
	ARM64         string `yaml:"arm64" json:"arm64"`
	RISCV64       string `yaml:"riscv64" json:"riscv64"`
	ARM           string `yaml:"arm" json:"arm"`
	Registry      string `yaml:"registry" json:"registry"`
	Name          string `yaml:"name" json:"name"`
	Description   string `yaml:"description" json:"description"`
	ConfigExample string `yaml:"configExample" json:"configExample"`
}

// MicroserviceImages contains information about the images for a microservice
// +k8s:deepcopy-gen=true
type MicroserviceImages struct {
	CatalogID int    `yaml:"catalogId" json:"catalogId"`
	AMD64     string `yaml:"amd64" json:"amd64"`
	ARM64     string `yaml:"arm64" json:"arm64"`
	RISCV64   string `yaml:"riscv64" json:"riscv64"`
	ARM       string `yaml:"arm" json:"arm"`
	Registry  string `yaml:"registry" json:"registry"`
}

// MicroserviceAgent contains information about required agent configuration for a microservice
// +k8s:deepcopy-gen=true
type MicroserviceAgent struct {
	Name   string             `yaml:"name" json:"name"`
	Config AgentConfiguration `yaml:"config" json:"config"`
}

// MicroserviceContainer contains information for configuring a microservice container
// +k8s:deepcopy-gen=true
type MicroserviceContainer struct {
	Commands        []string                     `yaml:"commands,omitempty" json:"commands,omitempty"`
	Volumes         *[]MicroserviceVolumeMapping `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	Env             *[]MicroserviceEnvironment   `yaml:"env,omitempty" json:"env,omitempty"`
	ExtraHosts      *[]MicroserviceExtraHost     `yaml:"extraHosts,omitempty" json:"extraHosts,omitempty"`
	Ports           []MicroservicePortMapping    `yaml:"ports" json:"ports"`
	HostNetworkMode bool                         `yaml:"hostNetworkMode" json:"hostNetworkMode"`
	IsPrivileged    bool                         `yaml:"isPrivileged" json:"isPrivileged"`
	PidMode         string                       `yaml:"pidMode,omitempty" json:"pidMode,omitempty"`
	IpcMode         string                       `yaml:"ipcMode,omitempty" json:"ipcMode,omitempty"`
	Runtime         string                       `yaml:"runtime,omitempty" json:"runtime,omitempty"`
	Platform        string                       `yaml:"platform,omitempty" json:"platform,omitempty"`
	RunAsUser       string                       `yaml:"runAsUser,omitempty" json:"runAsUser,omitempty"`
	CdiDevices      []string                     `yaml:"cdiDevices,omitempty" json:"cdiDevices,omitempty"`
	CapAdd          []string                     `yaml:"capAdd,omitempty" json:"capAdd,omitempty"`
	CapDrop         []string                     `yaml:"capDrop,omitempty" json:"capDrop,omitempty"`
	Annotations     ArbitraryJSON                `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	CPUSetCpus      string                       `yaml:"cpuSetCpus,omitempty" json:"cpuSetCpus,omitempty"`
	MemoryLimit     *int64                       `yaml:"memoryLimit,omitempty" json:"memoryLimit,omitempty"`
	HealthCheck     *MicroserviceHealthCheck     `yaml:"healthCheck,omitempty" json:"healthCheck,omitempty"`
}

// MicroserviceHealthCheck contains information about the health check of a microservice
// +k8s:deepcopy-gen=true
type MicroserviceHealthCheck struct {
	Test          []string `yaml:"test" json:"test"`
	Interval      *int64   `yaml:"interval,omitempty" json:"interval,omitempty"`
	Timeout       *int64   `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Retries       *int     `yaml:"retries,omitempty" json:"retries,omitempty"`
	StartPeriod   *int64   `yaml:"startPeriod,omitempty" json:"startPeriod,omitempty"`
	StartInterval *int64   `yaml:"startInterval,omitempty" json:"startInterval,omitempty"`
}

// MicroserviceStatusInfo contains information about the status of a microservice
// +k8s:deepcopy-gen=true
type MicroserviceStatusInfo struct {
	Status            string   `yaml:"status" json:"status"`
	StartTime         int64    `yaml:"startTime" json:"startTime"`
	OperatingDuration int64    `yaml:"operatingDuration" json:"operatingDuration"`
	MemoryUsage       float64  `yaml:"memoryUsage" json:"memoryUsage"`
	CPUUsage          float64  `yaml:"cpuUsage" json:"cpuUsage"`
	ContainerID       string   `yaml:"containerId" json:"containerId"`
	Percentage        float64  `yaml:"percentage" json:"percentage"`
	IPAddress         string   `yaml:"ipAddress" json:"ipAddress"`
	ErrorMessage      string   `yaml:"errorMessage" json:"errorMessage"`
	ExecSessionIDs    []string `yaml:"execSessionIds" json:"execSessionIds"`
	HealthStatus      string   `yaml:"healthStatus" json:"healthStatus"`
}

// MicroserviceExecStatusInfo contains information about the exec status of a microservice
// +k8s:deepcopy-gen=true
type MicroserviceExecStatusInfo struct {
	Status        string `yaml:"status" json:"status"`
	ExecSessionID string `yaml:"execSessionId" json:"execSessionId"`
}

// Microservice contains information for configuring a microservice
// +k8s:deepcopy-gen=true
type Microservice struct {
	UUID        string                     `yaml:"uuid" json:"uuid"`
	Name        string                     `yaml:"name" json:"name"`
	Agent       MicroserviceAgent          `yaml:"agent" json:"agent"`
	Images      *MicroserviceImages        `yaml:"images,omitempty" json:"images,omitempty"`
	Container   MicroserviceContainer      `yaml:"container,omitempty" json:"container,omitempty"`
	NatsConfig  *MicroserviceNatsConfig    `yaml:"natsConfig,omitempty" json:"natsConfig,omitempty"`
	Schedule    int                        `yaml:"schedule" json:"schedule"`
	Config      ArbitraryJSON              `yaml:"config" json:"config"`
	Application string                     `yaml:"application,omitempty" json:"application,omitempty"`
	Created     string                     `yaml:"created,omitempty" json:"created,omitempty"`
	Rebuild     bool                       `yaml:"rebuild,omitempty" json:"rebuild,omitempty"`
	Status      MicroserviceStatusInfo     `yaml:"status,omitempty" json:"status,omitempty"`
	ExecStatus  MicroserviceExecStatusInfo `yaml:"execStatus,omitempty" json:"execStatus,omitempty"`
}

// MicroserviceNatsConfig holds NATS configuration for a microservice (natsAccess, natsRule).
// +k8s:deepcopy-gen=true
type MicroserviceNatsConfig struct {
	NatsAccess bool   `yaml:"natsAccess" json:"natsAccess"`
	NatsRule   string `yaml:"natsRule,omitempty" json:"natsRule,omitempty"`
}

// ArbitraryJSON holds arbitrary JSON (e.g. config, annotations) as raw bytes.
// It is compatible with controller-gen and preserves round-trip for YAML/JSON.
// Use ToMap/FromMap for a map[string]any API.
// +k8s:deepcopy-gen=ignore
type ArbitraryJSON struct {
	//revive:disable-next-line:struct-tag k8s yaml inline convention
	runtime.RawExtension `yaml:",inline" json:",inline"`
}

// MarshalJSON implements json.Marshaler. Emits Raw or {} if empty.
func (a ArbitraryJSON) MarshalJSON() ([]byte, error) {
	if len(a.Raw) == 0 {
		return []byte("{}"), nil
	}
	return a.Raw, nil
}

// UnmarshalJSON implements json.Unmarshaler. Stores the value as raw bytes.
func (a *ArbitraryJSON) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		a.Raw = nil
		return nil
	}
	a.Raw = make([]byte, len(data))
	copy(a.Raw, data)
	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler so YAML objects (e.g. config: { key: value }) decode correctly.
func (a *ArbitraryJSON) UnmarshalYAML(unmarshal func(any) error) error {
	var v any
	if err := unmarshal(&v); err != nil {
		return err
	}
	if v == nil {
		a.Raw = []byte("{}")
		return nil
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	a.Raw = raw
	return nil
}

// MarshalYAML implements yaml.Marshaler so values emit as YAML objects.
func (a ArbitraryJSON) MarshalYAML() (any, error) {
	if len(a.Raw) == 0 {
		return map[string]any{}, nil
	}
	var v any
	if err := json.Unmarshal(a.Raw, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// DeepCopyInto copies the receiver into out, deep-copying Raw bytes.
func (a *ArbitraryJSON) DeepCopyInto(out *ArbitraryJSON) {
	if a == nil || len(a.Raw) == 0 {
		out.Raw = nil
		return
	}
	out.Raw = make([]byte, len(a.Raw))
	copy(out.Raw, a.Raw)
}

// DeepCopy returns a deep copy of the receiver.
func (a ArbitraryJSON) DeepCopy() ArbitraryJSON {
	if len(a.Raw) == 0 {
		return ArbitraryJSON{}
	}
	out := make([]byte, len(a.Raw))
	copy(out, a.Raw)
	return ArbitraryJSON{RawExtension: runtime.RawExtension{Raw: out}}
}

// ToMap unmarshals Raw into a map for code that needs a map API. Returns nil map and nil error if empty.
func (a ArbitraryJSON) ToMap() (map[string]any, error) {
	if len(a.Raw) == 0 {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(a.Raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// FromMap sets Raw from the given map (for backward compatibility in code that builds config/annotations).
func (a *ArbitraryJSON) FromMap(m map[string]any) error {
	if m == nil {
		a.Raw = []byte("{}")
		return nil
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return err
	}
	a.Raw = raw
	return nil
}

// +k8s:deepcopy-gen=true
type MicroservicePortMapping struct {
	Internal int64  `json:"internal"`
	External int64  `json:"external"`
	Protocol string `json:"protocol,omitempty"`
}

// MicroserviceVolumeMapping maps a host path or service account into the container.
// Type is the mapping kind (e.g. bind, serviceAccount).
// +k8s:deepcopy-gen=true
type MicroserviceVolumeMapping struct {
	HostDestination      string `yaml:"hostDestination" json:"hostDestination"`
	ContainerDestination string `yaml:"containerDestination" json:"containerDestination"`
	AccessMode           string `yaml:"accessMode" json:"accessMode"`
	Type                 string `yaml:"type,omitempty" json:"type,omitempty"`
}

// +k8s:deepcopy-gen=true
type MicroserviceEnvironment struct {
	Key                string `yaml:"key" json:"key"`
	Value              string `yaml:"value,omitempty" json:"value,omitempty"`
	ValueFromSecret    string `yaml:"valueFromSecret,omitempty" json:"valueFromSecret,omitempty"`
	ValueFromConfigMap string `yaml:"valueFromConfigMap,omitempty" json:"valueFromConfigMap,omitempty"`
}

// +k8s:deepcopy-gen=true
type MicroserviceExtraHost struct {
	Name    string `yaml:"name" json:"name,omitempty"`
	Address string `yaml:"address" json:"address,omitempty"`
	Value   string `yaml:"value" json:"value,omitempty"`
}

// +k8s:deepcopy-gen=true
type AgentConfiguration struct {
	ContainerEngineURL        *string   `yaml:"containerEngineUrl,omitempty" json:"containerEngineUrl,omitempty"`
	ContainerEngine           *string   `yaml:"containerEngine,omitempty" json:"containerEngine,omitempty"`
	DeploymentType            *string   `yaml:"deploymentType,omitempty" json:"deploymentType,omitempty"`
	DiskLimit                 *int64    `yaml:"diskLimit,omitempty" json:"diskLimit,omitempty"`
	DiskDirectory             *string   `yaml:"diskDirectory,omitempty" json:"diskDirectory,omitempty"`
	MemoryLimit               *int64    `yaml:"memoryLimit,omitempty" json:"memoryLimit,omitempty"`
	CPULimit                  *int64    `yaml:"cpuLimit,omitempty" json:"cpuLimit,omitempty"`
	LogLimit                  *int64    `yaml:"logLimit,omitempty" json:"logLimit,omitempty"`
	LogDirectory              *string   `yaml:"logDirectory,omitempty" json:"logDirectory,omitempty"`
	LogFileCount              *int64    `yaml:"logFileCount,omitempty" json:"logFileCount,omitempty"`
	StatusFrequency           *float64  `yaml:"statusFrequency,omitempty" json:"statusFrequency,omitempty"`
	ChangeFrequency           *float64  `yaml:"changeFrequency,omitempty" json:"changeFrequency,omitempty"`
	DeviceScanFrequency       *float64  `yaml:"deviceScanFrequency,omitempty" json:"deviceScanFrequency,omitempty"`
	GpsMode                   *string   `yaml:"gpsMode,omitempty" json:"gpsMode,omitempty"`
	GpsScanFrequency          *float64  `yaml:"gpsScanFrequency,omitempty" json:"gpsScanFrequency,omitempty"`
	GpsDevice                 *string   `yaml:"gpsDevice,omitempty" json:"gpsDevice,omitempty"`
	EdgeGuardFrequency        *float64  `yaml:"edgeGuardFrequency,omitempty" json:"edgeGuardFrequency,omitempty"`
	BluetoothEnabled          *bool     `yaml:"bluetoothEnabled,omitempty" json:"bluetoothEnabled,omitempty"`
	WatchdogEnabled           *bool     `yaml:"watchdogEnabled,omitempty" json:"watchdogEnabled,omitempty"`
	AbstractedHardwareEnabled *bool     `yaml:"abstractedHardwareEnabled,omitempty" json:"abstractedHardwareEnabled,omitempty"`
	RouterMode                *string   `yaml:"routerMode,omitempty" json:"routerMode,omitempty"`           // [edge, interior, none], default: edge
	RouterPort                *int      `yaml:"routerPort,omitempty" json:"routerPort,omitempty"`           // default: 5671
	UpstreamRouters           *[]string `yaml:"upstreamRouters,omitempty" json:"upstreamRouters,omitempty"` // ignored if routerMode: none
	NetworkRouter             *string   `yaml:"networkRouter,omitempty" json:"networkRouter,omitempty"`     // required if routerMone: none
	// NATS-related fields (Controller iofog schema)
	NatsMode            *string   `yaml:"natsMode,omitempty" json:"natsMode,omitempty"` // none, leaf, server
	NatsServerPort      *int      `yaml:"natsServerPort,omitempty" json:"natsServerPort,omitempty"`
	NatsLeafPort        *int      `yaml:"natsLeafPort,omitempty" json:"natsLeafPort,omitempty"`
	NatsClusterPort     *int      `yaml:"natsClusterPort,omitempty" json:"natsClusterPort,omitempty"`
	NatsMqttPort        *int      `yaml:"natsMqttPort,omitempty" json:"natsMqttPort,omitempty"`
	NatsHTTPPort        *int      `yaml:"natsHttpPort,omitempty" json:"natsHttpPort,omitempty"`
	UpstreamNatsServers *[]string `yaml:"upstreamNatsServers,omitempty" json:"upstreamNatsServers,omitempty"`
	JsStorageSize       *string   `yaml:"jsStorageSize,omitempty" json:"jsStorageSize,omitempty"`
	JsMemoryStoreSize   *string   `yaml:"jsMemoryStoreSize,omitempty" json:"jsMemoryStoreSize,omitempty"`
	PruningFrequency    *float64  `yaml:"pruningFrequency,omitempty" json:"pruningFrequency,omitempty"`
	ArchID              *int64    `yaml:"archId,omitempty" json:"archId,omitempty"`
}

// Microservices is a list of Microservice
// +k8s:deepcopy-gen=true
type Microservices struct {
	Microservices []Microservice `yaml:"microservices" json:"microservices"`
}

// Route contains information about a route from one microservice to another.
// Deprecated: Controller no longer exposes routing; use NATS for messaging.
// +k8s:deepcopy-gen=true
type Route struct {
	Name string `yaml:"name" json:"name"`
	From string `yaml:"from" json:"from"`
	To   string `yaml:"to" json:"to"`
}

// ApplicationNatsConfig holds NATS configuration for an application (natsAccess, natsRule).
// +k8s:deepcopy-gen=true
type ApplicationNatsConfig struct {
	NatsAccess bool   `yaml:"natsAccess" json:"natsAccess"`
	NatsRule   string `yaml:"natsRule,omitempty" json:"natsRule,omitempty"`
}

// Application contains information for configuring an application
// +k8s:deepcopy-gen=true
type Application struct {
	Name          string                 `yaml:"name" json:"name"`
	Microservices []Microservice         `yaml:"microservices,omitempty" json:"microservices,omitempty"`
	NatsConfig    *ApplicationNatsConfig `yaml:"natsConfig,omitempty" json:"natsConfig,omitempty"`
	ID            int                    `yaml:"id,omitempty" json:"id,omitempty"`
	Template      *ApplicationTemplate   `yaml:"template,omitempty" json:"template,omitempty"`
}

// ApplicationTemplate contains information for configuring an application template
// +k8s:deepcopy-gen=true
type ApplicationTemplate struct {
	Name        string                   `yaml:"name,omitempty"`
	Description string                   `yaml:"description,omitempty"`
	Variables   []TemplateVariable       `yaml:"variables,omitempty"`
	Application *ApplicationTemplateInfo `yaml:"application,omitempty"`
}

// TemplateVariable contains a key-value pair.
// +k8s:deepcopy-gen=ignore
// (apiextensions.JSON is not supported by deepcopy-gen; use shallow copy if needed.)
type TemplateVariable struct {
	Key          string              `yaml:"key"`
	Description  string              `yaml:"description"`
	DefaultValue *apiextensions.JSON `yaml:"defaultValue,omitempty"`
	Value        *apiextensions.JSON `yaml:"value,omitempty"`
}

// DeepCopyInto implements deepcopy-gen's requirement for types referenced by generated code.
func (in *TemplateVariable) DeepCopyInto(out *TemplateVariable) {
	*out = *in
}

// DeepCopy returns a shallow copy of TemplateVariable.
func (in *TemplateVariable) DeepCopy() *TemplateVariable {
	if in == nil {
		return nil
	}
	out := new(TemplateVariable)
	in.DeepCopyInto(out)
	return out
}

// ApplicationTemplateInfo contains microservice and NATS config details for template
// +k8s:deepcopy-gen=true
type ApplicationTemplateInfo struct {
	Microservices []Microservice         `yaml:"microservices"`
	NatsConfig    *ApplicationNatsConfig `yaml:"natsConfig,omitempty"`
}

// Applications is a list of applications
// +k8s:deepcopy-gen=true
type Applications struct {
	Applications []Application `yaml:"applications" json:"applications"`
}

// IofogController contains informations needed to connect to the controller
// +k8s:deepcopy-gen=true
type IofogController struct {
	Email        string `yaml:"email" json:"email"`
	Password     string `yaml:"password" json:"password"`
	Endpoint     string `yaml:"endpoint" json:"endpoint"`
	Token        string `yaml:"token" json:"token"`
	RefreshToken string `yaml:"refreshToken" json:"refreshToken"`
}
