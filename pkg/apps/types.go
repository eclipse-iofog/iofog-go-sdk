package apps

import (
	"encoding/json"
	"errors"
	"fmt"

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
	ApplicationKind          Kind = "Application"
	ApplicationTemplateKind  Kind = "ApplicationTemplate"
	MicroserviceKind         Kind = "Microservice"
	MicroserviceTemplateKind Kind = "MicroserviceTemplate"
	ModelKind                Kind = "Model"
	RuntimeClassKind         Kind = "RuntimeClass"
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
	ID            int         `yaml:"id" json:"id"`
	Registry      RegistryRef `yaml:"registry" json:"registry"`
	ARM64         string      `yaml:"arm64" json:"arm64"`
	AMD64         string      `yaml:"amd64" json:"amd64"`
	RISCV64       string      `yaml:"riscv64,omitempty" json:"riscv64,omitempty"`
	ARM           string      `yaml:"arm,omitempty" json:"arm,omitempty"`
	Name          string      `yaml:"name" json:"name"`
	Description   string      `yaml:"description" json:"description"`
	ConfigExample string      `yaml:"configExample" json:"configExample"`
}

// MicroserviceImages contains information about the images for a microservice
// +k8s:deepcopy-gen=true
type MicroserviceImages struct {
	Registry  RegistryRef `yaml:"registry" json:"registry"`
	ARM64     string      `yaml:"arm64" json:"arm64"`
	AMD64     string      `yaml:"amd64" json:"amd64"`
	RISCV64   string      `yaml:"riscv64,omitempty" json:"riscv64,omitempty"`
	ARM       string      `yaml:"arm,omitempty" json:"arm,omitempty"`
	CatalogID int         `yaml:"catalogId,omitempty" json:"catalogId,omitempty"`
}

// RoleRef references a Role by kind and name.
// +k8s:deepcopy-gen=true
type RoleRef struct {
	Kind     string `yaml:"kind" json:"kind"`
	Name     string `yaml:"name" json:"name"`
	APIGroup string `yaml:"apiGroup,omitempty" json:"apiGroup,omitempty"`
}

// MicroserviceServiceAccountRef binds a microservice to an application Role.
// +k8s:deepcopy-gen=true
type MicroserviceServiceAccountRef struct {
	RoleRef RoleRef `yaml:"roleRef" json:"roleRef"`
}

// MicroserviceAgent names the edge agent that runs the microservice.
// +k8s:deepcopy-gen=true
type MicroserviceAgent struct {
	Name string `yaml:"name" json:"name"`
}

// MicroserviceContainer contains information for configuring a microservice container
// +k8s:deepcopy-gen=true
type MicroserviceContainer struct {
	HostNetworkMode        bool                          `yaml:"hostNetworkMode" json:"hostNetworkMode"`
	IsPrivileged           bool                          `yaml:"isPrivileged" json:"isPrivileged"`
	RunAsUser              string                        `yaml:"runAsUser,omitempty" json:"runAsUser,omitempty"`
	RunAsGroup             string                        `yaml:"runAsGroup,omitempty" json:"runAsGroup,omitempty"`
	ReadOnlyRootFilesystem bool                          `yaml:"readOnlyRootFilesystem" json:"readOnlyRootFilesystem"`
	IpcMode                string                        `yaml:"ipcMode,omitempty" json:"ipcMode,omitempty"`
	PidMode                string                        `yaml:"pidMode,omitempty" json:"pidMode,omitempty"`
	Platform               string                        `yaml:"platform,omitempty" json:"platform,omitempty"`
	Runtime                string                        `yaml:"runtime,omitempty" json:"runtime,omitempty"`
	CapAdd                 []string                      `yaml:"capAdd,omitempty" json:"capAdd,omitempty"`
	CapDrop                []string                      `yaml:"capDrop,omitempty" json:"capDrop,omitempty"`
	Annotations            ArbitraryJSON                 `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	Sysctls                map[string]string             `yaml:"sysctls,omitempty" json:"sysctls,omitempty"`
	Ulimits                map[string]MicroserviceUlimit `yaml:"ulimits,omitempty" json:"ulimits,omitempty"`
	CPUSetCpus             string                        `yaml:"cpuSetCpus,omitempty" json:"cpuSetCpus,omitempty"`
	CPUs                   *float64                      `yaml:"cpus,omitempty" json:"cpus,omitempty"`
	MemoryLimit            *int64                        `yaml:"memoryLimit,omitempty" json:"memoryLimit,omitempty"`
	MemoryReservation      *int64                        `yaml:"memoryReservation,omitempty" json:"memoryReservation,omitempty"`
	MemorySwap             *int64                        `yaml:"memorySwap,omitempty" json:"memorySwap,omitempty"`
	ShmSize                *int64                        `yaml:"shmSize,omitempty" json:"shmSize,omitempty"`
	CdiDevices             []string                      `yaml:"cdiDevices,omitempty" json:"cdiDevices,omitempty"`
	Devices                []MicroserviceDevice          `yaml:"devices,omitempty" json:"devices,omitempty"`
	Volumes                *[]MicroserviceVolumeMapping  `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	Tmpfs                  []MicroserviceTmpfs           `yaml:"tmpfs,omitempty" json:"tmpfs,omitempty"`
	ExtraHosts             *[]MicroserviceExtraHost      `yaml:"extraHosts,omitempty" json:"extraHosts,omitempty"`
	Env                    *[]MicroserviceEnvironment    `yaml:"env,omitempty" json:"env,omitempty"`
	Ports                  []MicroservicePortMapping     `yaml:"ports" json:"ports"`
	WorkingDir             string                        `yaml:"workingDir,omitempty" json:"workingDir,omitempty"`
	Entrypoint             []string                      `yaml:"entrypoint,omitempty" json:"entrypoint,omitempty"`
	Commands               []string                      `yaml:"commands,omitempty" json:"commands,omitempty"`
	HealthCheck            *MicroserviceHealthCheck      `yaml:"healthCheck,omitempty" json:"healthCheck,omitempty"`
}

// MicroserviceUlimit is a soft/hard resource limit. -1 means unlimited.
// +k8s:deepcopy-gen=true
type MicroserviceUlimit struct {
	Soft int `yaml:"soft" json:"soft"`
	Hard int `yaml:"hard" json:"hard"`
}

// MicroserviceDevice maps a host /dev node into the container.
// +k8s:deepcopy-gen=true
type MicroserviceDevice struct {
	HostPath      string `yaml:"hostPath" json:"hostPath"`
	ContainerPath string `yaml:"containerPath" json:"containerPath"`
	Permissions   string `yaml:"permissions,omitempty" json:"permissions,omitempty"`
}

// MicroserviceTmpfs is an in-memory mount. Size is MiB; Mode is an optional octal string.
// +k8s:deepcopy-gen=true
type MicroserviceTmpfs struct {
	ContainerPath string `yaml:"containerPath" json:"containerPath"`
	Size          *int64 `yaml:"size,omitempty" json:"size,omitempty"`
	Mode          string `yaml:"mode,omitempty" json:"mode,omitempty"`
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
	PodID             string   `yaml:"podId,omitempty" json:"podId,omitempty"`
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
	UUID           string                         `yaml:"uuid" json:"uuid"`
	Application    string                         `yaml:"application,omitempty" json:"application,omitempty"`
	Name           string                         `yaml:"name" json:"name"`
	Agent          MicroserviceAgent              `yaml:"agent" json:"agent"`
	Images         *MicroserviceImages            `yaml:"images,omitempty" json:"images,omitempty"`
	NatsConfig     *MicroserviceNatsConfig        `yaml:"natsConfig,omitempty" json:"natsConfig,omitempty"`
	Models         *MicroserviceCatalog           `yaml:"models,omitempty" json:"models,omitempty"`
	Container      MicroserviceContainer          `yaml:"container,omitempty" json:"container,omitempty"`
	Schedule       int                            `yaml:"schedule" json:"schedule"`
	Config         ArbitraryJSON                  `yaml:"config" json:"config"`
	ServiceAccount *MicroserviceServiceAccountRef `yaml:"serviceAccount,omitempty" json:"serviceAccount,omitempty"`
	Template       *MicroserviceTemplateRef       `yaml:"template,omitempty" json:"template,omitempty"`
	Created        string                         `yaml:"created,omitempty" json:"created,omitempty"`
	Rebuild        bool                           `yaml:"rebuild,omitempty" json:"rebuild,omitempty"`
	Status         MicroserviceStatusInfo         `yaml:"status,omitempty" json:"status,omitempty"`
	ExecStatus     MicroserviceExecStatusInfo     `yaml:"execStatus,omitempty" json:"execStatus,omitempty"`
}

// MicroserviceCatalogItem names a fleet model bound into the container.
// +k8s:deepcopy-gen=true
type MicroserviceCatalogItem struct {
	// Name is the fleet model metadata name; never a uuid or host path.
	Name string `yaml:"name" json:"name"`
}

// MicroserviceCatalog binds Ready model content into the container.
// +k8s:deepcopy-gen=true
type MicroserviceCatalog struct {
	BindPath    string                    `yaml:"bindPath,omitempty" json:"bindPath,omitempty"`
	Permissions string                    `yaml:"permissions,omitempty" json:"permissions,omitempty"` // ro | rw
	Items       []MicroserviceCatalogItem `yaml:"items,omitempty" json:"items,omitempty"`
}

// MicroserviceTemplateVariables is a map of template instance variable values.
// YAML and JSON unmarshal accept a mapping or a list of {key, value} pairs; marshal emits a mapping.
// +k8s:deepcopy-gen=ignore
type MicroserviceTemplateVariables map[string]any

// DeepCopyInto copies the receiver into out.
func (in *MicroserviceTemplateVariables) DeepCopyInto(out *MicroserviceTemplateVariables) {
	if in == nil {
		return
	}
	*out = in.DeepCopy()
}

// DeepCopy returns a shallow copy of the variable map.
func (in MicroserviceTemplateVariables) DeepCopy() MicroserviceTemplateVariables {
	if in == nil {
		return nil
	}
	out := make(MicroserviceTemplateVariables, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// UnmarshalJSON accepts a JSON object or a list of {key, value} objects.
func (in *MicroserviceTemplateVariables) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*in = nil
		return nil
	}
	switch data[0] {
	case '{':
		var asMap map[string]any
		if err := json.Unmarshal(data, &asMap); err != nil {
			return err
		}
		*in = asMap
		return nil
	case '[':
		var asList []struct {
			Key   string `json:"key"`
			Value any    `json:"value"`
		}
		if err := json.Unmarshal(data, &asList); err != nil {
			return err
		}
		out := make(MicroserviceTemplateVariables, len(asList))
		for _, item := range asList {
			if item.Key == "" {
				return errors.New("template variables list entries require a key")
			}
			out[item.Key] = item.Value
		}
		*in = out
		return nil
	default:
		return errors.New("template variables must be a mapping or a list of key/value pairs")
	}
}

// UnmarshalYAML accepts a YAML mapping or a list of {key, value} objects.
func (in *MicroserviceTemplateVariables) UnmarshalYAML(unmarshal func(any) error) error {
	var raw any
	if err := unmarshal(&raw); err != nil {
		return err
	}
	parsed, err := parseTemplateVariables(raw)
	if err != nil {
		return err
	}
	*in = parsed
	return nil
}

func parseTemplateVariables(raw any) (MicroserviceTemplateVariables, error) {
	if raw == nil {
		return nil, nil
	}
	if m, ok := stringifyYAMLMap(raw); ok {
		return MicroserviceTemplateVariables(m), nil
	}
	list, ok := raw.([]any)
	if !ok {
		return nil, errors.New("template variables must be a mapping or a list of key/value pairs")
	}
	out := make(MicroserviceTemplateVariables, len(list))
	for _, item := range list {
		pair, ok := stringifyYAMLMap(item)
		if !ok {
			return nil, errors.New("template variables list entries must be key/value objects")
		}
		key, ok := pair["key"].(string)
		if !ok || key == "" {
			return nil, errors.New("template variables list entries require a key")
		}
		out[key] = pair["value"]
	}
	return out, nil
}

func stringifyYAMLMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return m, true
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, val := range m {
			ks, ok := k.(string)
			if !ok {
				return nil, false
			}
			out[ks] = val
		}
		return out, true
	default:
		return nil, false
	}
}

// MicroserviceTemplateRef selects a microservice template and optional variable values.
// +k8s:deepcopy-gen=ignore
type MicroserviceTemplateRef struct {
	Name      string                        `yaml:"name" json:"name"`
	Variables MicroserviceTemplateVariables `yaml:"variables,omitempty" json:"variables,omitempty"`
}

// DeepCopyInto copies the receiver into out.
func (in *MicroserviceTemplateRef) DeepCopyInto(out *MicroserviceTemplateRef) {
	*out = *in
	if in.Variables != nil {
		out.Variables = in.Variables.DeepCopy()
	}
}

// DeepCopy returns a deep copy of MicroserviceTemplateRef.
func (in *MicroserviceTemplateRef) DeepCopy() *MicroserviceTemplateRef {
	if in == nil {
		return nil
	}
	out := new(MicroserviceTemplateRef)
	in.DeepCopyInto(out)
	return out
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

func normalizeYAMLValue(v any) (any, error) {
	switch x := v.(type) {
	case map[any]any:
		out := make(map[string]any, len(x))
		for k, val := range x {
			ks, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("non-string yaml map key %T", k)
			}
			normalized, err := normalizeYAMLValue(val)
			if err != nil {
				return nil, err
			}
			out[ks] = normalized
		}
		return out, nil
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, val := range x {
			normalized, err := normalizeYAMLValue(val)
			if err != nil {
				return nil, err
			}
			out[k] = normalized
		}
		return out, nil
	case []any:
		out := make([]any, len(x))
		for i, val := range x {
			normalized, err := normalizeYAMLValue(val)
			if err != nil {
				return nil, err
			}
			out[i] = normalized
		}
		return out, nil
	default:
		return v, nil
	}
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
	normalized, err := normalizeYAMLValue(v)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(normalized)
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
	ContainerEngineURL *string   `yaml:"containerEngineUrl,omitempty" json:"containerEngineUrl,omitempty"`
	ContainerEngine    *string   `yaml:"containerEngine,omitempty" json:"containerEngine,omitempty"`
	DeploymentType     *string   `yaml:"deploymentType,omitempty" json:"deploymentType,omitempty"`
	DiskLimit          *int64    `yaml:"diskLimit,omitempty" json:"diskLimit,omitempty"`
	DiskDirectory      *string   `yaml:"diskDirectory,omitempty" json:"diskDirectory,omitempty"`
	MemoryLimit        *int64    `yaml:"memoryLimit,omitempty" json:"memoryLimit,omitempty"`
	CPULimit           *int64    `yaml:"cpuLimit,omitempty" json:"cpuLimit,omitempty"`
	LogLimit           *int64    `yaml:"logLimit,omitempty" json:"logLimit,omitempty"`
	LogDirectory       *string   `yaml:"logDirectory,omitempty" json:"logDirectory,omitempty"`
	LogFileCount       *int64    `yaml:"logFileCount,omitempty" json:"logFileCount,omitempty"`
	StatusFrequency    *float64  `yaml:"statusFrequency,omitempty" json:"statusFrequency,omitempty"`
	ChangeFrequency    *float64  `yaml:"changeFrequency,omitempty" json:"changeFrequency,omitempty"`
	GpsMode            *string   `yaml:"gpsMode,omitempty" json:"gpsMode,omitempty"`
	GpsScanFrequency   *float64  `yaml:"gpsScanFrequency,omitempty" json:"gpsScanFrequency,omitempty"`
	GpsDevice          *string   `yaml:"gpsDevice,omitempty" json:"gpsDevice,omitempty"`
	EdgeGuardFrequency *float64  `yaml:"edgeGuardFrequency,omitempty" json:"edgeGuardFrequency,omitempty"`
	WatchdogEnabled    *bool     `yaml:"watchdogEnabled,omitempty" json:"watchdogEnabled,omitempty"`
	RouterMode         *string   `yaml:"routerMode,omitempty" json:"routerMode,omitempty"`           // [edge, interior, none], default: edge
	RouterPort         *int      `yaml:"routerPort,omitempty" json:"routerPort,omitempty"`           // default: 5671
	UpstreamRouters    *[]string `yaml:"upstreamRouters,omitempty" json:"upstreamRouters,omitempty"` // ignored if routerMode: none
	NetworkRouter      *string   `yaml:"networkRouter,omitempty" json:"networkRouter,omitempty"`     // required if routerMone: none
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

// MicroserviceTemplate is the spec for kind: MicroserviceTemplate YAML.
// +k8s:deepcopy-gen=true
type MicroserviceTemplate struct {
	Name         string             `yaml:"name,omitempty" json:"name,omitempty"`
	Description  string             `yaml:"description,omitempty" json:"description,omitempty"`
	Variables    []TemplateVariable `yaml:"variables,omitempty" json:"variables,omitempty"`
	Microservice *Microservice      `yaml:"microservice,omitempty" json:"microservice,omitempty"`
}

// Model is the spec for kind: Model YAML.
// +k8s:deepcopy-gen=true
type Model struct {
	Repo       string   `yaml:"repo" json:"repo"`
	Revision   string   `yaml:"revision,omitempty" json:"revision,omitempty"`
	RegistryID int      `yaml:"registryId" json:"registryId"`
	Files      []string `yaml:"files,omitempty" json:"files,omitempty"`
	Format     string   `yaml:"format,omitempty" json:"format,omitempty"`
}

// RuntimeClass is the YAML document for kind: RuntimeClass.
// Handler is a top-level field (not under spec), matching Controller wire shape.
// +k8s:deepcopy-gen=true
type RuntimeClass struct {
	APIVersion string         `yaml:"apiVersion" json:"apiVersion"`
	Kind       Kind           `yaml:"kind" json:"kind"`
	Metadata   HeaderMetadata `yaml:"metadata" json:"metadata"`
	Handler    string         `yaml:"handler" json:"handler"`
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
