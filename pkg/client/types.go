package client

import (
	"encoding/json"
	"strconv"
	"time"
)

// Flows - Keep for legacy
type FlowInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActivated bool   `json:"isActivated"`
	IsSystem    bool   `json:"isSystem"`
	UserID      int    `json:"userId"`
	ID          int    `json:"id"`
}

type FlowCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type FlowCreateResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type FlowUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActivated *bool   `json:"isActivated,omitempty"`
	IsSystem    *bool   `json:"isSystem,omitempty"`
	ID          int     `json:"-"`
}

type FlowListResponse struct {
	Flows []FlowInfo `json:"flows"`
}

// ApplicationNatsConfig holds NATS configuration for an application (Controller applicationNatsConfig).
type ApplicationNatsConfig struct {
	NatsAccess bool   `json:"natsAccess"`
	NatsRule   string `json:"natsRule,omitempty"`
}

// Applications
type ApplicationInfo struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	IsActivated   bool                   `json:"isActivated"`
	IsSystem      bool                   `json:"isSystem"`
	UserID        int                    `json:"userId"`
	ID            int                    `json:"id"`
	Microservices []MicroserviceInfo     `json:"microservices"`
	NatsConfig    *ApplicationNatsConfig `json:"natsConfig,omitempty"`
}

type ApplicationCreateResponse struct {
	ID int `json:"id"`
}

type ApplicationPatchRequest struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	IsActivated *bool                  `json:"isActivated,omitempty"`
	IsSystem    *bool                  `json:"isSystem,omitempty"`
	NatsConfig  *ApplicationNatsConfig `json:"natsConfig,omitempty"`
}

type ApplicationListResponse struct {
	Applications []ApplicationInfo `json:"applications"`
}

// Application Templates
type ApplicationTemplate struct {
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	Variables   []TemplateVariable       `json:"variables,omitempty"`
	Application *ApplicationTemplateInfo `json:"application,omitempty"`
}

type ApplicationTemplateCreateRequest = ApplicationTemplate

type TemplateVariable struct {
	Key          string `json:"key" yaml:"key,omitempty"`
	Description  string `json:"description" yaml:"description,omitempty"`
	DefaultValue any    `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`
	Value        any    `json:"value,omitempty" yaml:"value,omitempty"`
}

type ApplicationTemplateInfo struct {
	Microservices []any                  `json:"microservices"`
	NatsConfig    *ApplicationNatsConfig `json:"natsConfig,omitempty"`
}

type ApplicationTemplateCreateResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type ApplicationTemplateUpdateResponse = ApplicationTemplateCreateResponse

type ApplicationTemplateMetadataUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type ApplicationTemplateListResponse struct {
	ApplicationTemplates []ApplicationTemplate
}

// Registries
type RegistryInfo struct {
	ID           int    `json:"id"`
	URL          string `json:"url"`
	IsPublic     bool   `json:"isPublic"`
	IsSecure     bool   `json:"isSecure"`
	Certificate  string `json:"certificate"`
	RequiresCert bool   `json:"requiresCert"`
	Username     string `json:"username"`
	Email        string `json:"userEmail"`
	Password     string `json:"password"`
}

type RegistryCreateRequest struct {
	URL          string `json:"url"`
	IsPublic     bool   `json:"isPublic"`
	Certificate  string `json:"certificate"`
	RequiresCert bool   `json:"requiresCert"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

type RegistryCreateResponse struct {
	ID int `json:"id"`
}

type RegistryUpdateRequest struct {
	URL          *string `json:"url,omitempty"`
	IsPublic     *bool   `json:"isPublic,omitempty"`
	Certificate  *string `json:"certificate,omitempty"`
	RequiresCert *bool   `json:"requiresCert,omitempty"`
	Username     *string `json:"username,omitempty"`
	Email        *string `json:"email,omitempty"`
	Password     *string `json:"password,omitempty"`
	ID           int     `json:"-"`
}

type RegistryListResponse struct {
	Registries []RegistryInfo `json:"registries"`
}

// RBAC

// RBACRule is a single rule in a Role (apiGroups, resources, verbs, optional resourceNames)
type RBACRule struct {
	APIGroups     []string `json:"apiGroups"`
	Resources     []string `json:"resources"`
	Verbs         []string `json:"verbs"`
	ResourceNames []string `json:"resourceNames,omitempty"`
}

// RoleRef references a Role by kind and name
type RoleRef struct {
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	APIGroup string `json:"apiGroup,omitempty"`
}

// Subject is a user, group, or service account in a RoleBinding
type Subject struct {
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	APIGroup string `json:"apiGroup,omitempty"`
}

// RoleInfo is a Role as returned by the API
type RoleInfo struct {
	Name  string     `json:"name"`
	Kind  string     `json:"kind,omitempty"`
	Rules []RBACRule `json:"rules"`
}

// RoleCreateRequest is the request body for creating a Role
type RoleCreateRequest struct {
	Name  string     `json:"name"`
	Kind  string     `json:"kind,omitempty"`
	Rules []RBACRule `json:"rules"`
}

// RoleUpdateRequest is the request body for updating a Role
type RoleUpdateRequest struct {
	Name  *string    `json:"name,omitempty"`
	Kind  string     `json:"kind,omitempty"`
	Rules []RBACRule `json:"rules"`
}

// RoleListResponse is the response for listing roles
type RoleListResponse struct {
	Roles []RoleInfo `json:"roles"`
}

// RoleResponse is the response for get/create/update role (single role wrapped)
type RoleResponse struct {
	Role RoleInfo `json:"role"`
}

// RoleBindingInfo is a RoleBinding as returned by the API
type RoleBindingInfo struct {
	Name     string    `json:"name"`
	Kind     string    `json:"kind,omitempty"`
	RoleRef  RoleRef   `json:"roleRef"`
	Subjects []Subject `json:"subjects"`
}

// RoleBindingCreateRequest is the request body for creating a RoleBinding
type RoleBindingCreateRequest struct {
	Name     string    `json:"name"`
	Kind     string    `json:"kind,omitempty"`
	RoleRef  RoleRef   `json:"roleRef"`
	Subjects []Subject `json:"subjects"`
}

// RoleBindingUpdateRequest is the request body for updating a RoleBinding
type RoleBindingUpdateRequest struct {
	Name     *string   `json:"name,omitempty"`
	Kind     string    `json:"kind,omitempty"`
	RoleRef  RoleRef   `json:"roleRef"`
	Subjects []Subject `json:"subjects"`
}

// RoleBindingListResponse is the response for listing role bindings
type RoleBindingListResponse struct {
	Bindings []RoleBindingInfo `json:"bindings"`
}

// RoleBindingResponse is the response for get/create/update role binding
type RoleBindingResponse struct {
	Binding RoleBindingInfo `json:"binding"`
}

// ServiceAccountInfo is a ServiceAccount as returned by the API (application-scoped).
type ServiceAccountInfo struct {
	ID              int     `json:"id,omitempty"`
	Name            string  `json:"name"`
	ApplicationName string  `json:"applicationName,omitempty"`
	RoleRef         RoleRef `json:"roleRef"`
}

// ServiceAccountCreateRequest is the request body for creating a ServiceAccount.
// ApplicationName is required; the service account is created in that application.
type ServiceAccountCreateRequest struct {
	Name            string  `json:"name"`
	ApplicationName string  `json:"applicationName"`
	RoleRef         RoleRef `json:"roleRef"`
}

// ServiceAccountUpdateRequest is the request body for updating a ServiceAccount.
// AppName and name are passed in the URL path; body contains name and roleRef (both required by Controller).
type ServiceAccountUpdateRequest struct {
	Name    string  `json:"name"`
	RoleRef RoleRef `json:"roleRef"`
}

// ServiceAccountListResponse is the response for listing service accounts
type ServiceAccountListResponse struct {
	ServiceAccounts []ServiceAccountInfo `json:"serviceAccounts"`
}

// ServiceAccountResponse is the response for get/create/update service account
type ServiceAccountResponse struct {
	ServiceAccount ServiceAccountInfo `json:"serviceAccount"`
}

// Catalog (Keeping it basic, because it will be reworked soon)

type CatalogImage struct {
	ContainerImage string `json:"containerImage"`
	AgentTypeID    int    `json:"fogTypeId"`
}

type CatalogItemInfo struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      []CatalogImage `json:"images"`
	RegistryID  int            `json:"registryId"`
	Category    string         `json:"category"`
}

type CatalogItemCreateRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      []CatalogImage `json:"images"`
	RegistryID  int            `json:"registryId"`
}

type CatalogItemCreateResponse struct {
	ID int `json:"id"`
}

type CatalogItemUpdateRequest struct {
	ID          int
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Images      []CatalogImage `json:"images,omitempty"`
	RegistryID  int            `json:"registryId,omitempty"`
}

type CatalogListResponse struct {
	CatalogItems []CatalogItemInfo `json:"catalogItems"`
}

// Microservices

type MicroservicePortMappingInfo struct {
	Internal int64  `json:"internal"`
	External int64  `json:"external"`
	Protocol string `json:"protocol,omitempty"`
}

type MicroserviceVolumeMappingInfo struct {
	HostDestination      string `json:"hostDestination"`
	ContainerDestination string `json:"containerDestination"`
	AccessMode           string `json:"accessMode"`
	Type                 string `json:"type,omitempty"`
}

type MicroserviceEnvironmentInfo struct {
	Key                string `json:"key"`
	Value              string `json:"value,omitempty"`
	ValueFromSecret    string `json:"valueFromSecret,omitempty"`
	ValueFromConfigMap string `json:"valueFromConfigMap,omitempty"`
}

type MicroserviceStatusInfo struct {
	Status            string   `json:"status"`
	StartTime         int64    `json:"startTime"`
	OperatingDuration int64    `json:"operatingDuration"`
	MemoryUsage       float64  `json:"memoryUsage"`
	CPUUsage          float64  `json:"cpuUsage"`
	ContainerID       string   `json:"containerId"`
	Percentage        float64  `json:"percentage"`
	IPAddress         string   `json:"ipAddress"`
	ErrorMessage      string   `json:"errorMessage"`
	ExecSessionIDs    []string `json:"execSessionIds"`
	HealthStatus      string   `json:"healthStatus"`
}

type MicroserviceExecStatusInfo struct {
	Status        string `json:"status"`
	ExecSessionID string `json:"execSessionId"`
}

type MicroserviceInfo struct {
	UUID              string                          `json:"uuid"`
	Config            string                          `json:"config"`
	Name              string                          `json:"name"`
	HostNetworkMode   bool                            `json:"hostNetworkMode"`
	IsPrivileged      bool                            `json:"isPrivileged"`
	Schedule          int                             `json:"schedule"`
	PidMode           string                          `json:"pidMode,omitempty"`
	IpcMode           string                          `json:"ipcMode,omitempty"`
	Runtime           string                          `json:"runtime,omitempty"`
	Platform          string                          `json:"platform,omitempty"`
	RunAsUser         string                          `json:"runAsUser,omitempty"`
	CdiDevices        []string                        `json:"cdiDevices,omitempty"`
	CapAdd            []string                        `json:"capAdd,omitempty"`
	CapDrop           []string                        `json:"capDrop,omitempty"`
	LogSize           int                             `json:"logSize"`
	Delete            bool                            `json:"delete"`
	DeleteWithCleanup bool                            `json:"deleteWithCleanup"`
	FlowID            int                             `json:"flowId"`
	ApplicationID     int                             `json:"applicationID"`
	Application       string                          `json:"application"`
	CatalogItemID     int                             `json:"catalogItemId"`
	AgentUUID         string                          `json:"iofogUuid"`
	UserID            int                             `json:"userId"`
	RegistryID        int                             `json:"registryId"`
	Ports             []MicroservicePortMappingInfo   `json:"ports"`
	Volumes           []MicroserviceVolumeMappingInfo `json:"volumeMappings"`
	Commands          []string                        `json:"cmd"`
	Env               []MicroserviceEnvironmentInfo   `json:"env"`
	ExtraHosts        []MicroserviceExtraHost         `json:"extraHosts"`
	Status            MicroserviceStatusInfo          `json:"status"`
	ExecStatus        MicroserviceExecStatusInfo      `json:"execStatus"`
	Images            []CatalogImage                  `json:"images"`
	PubTags           []string                        `json:"pubTags"`
	SubTags           []string                        `json:"subTags"`
	Annotations       string                          `json:"annotations"`
	CPUSetCpus        string                          `json:"cpuSetCpus,omitempty"`
	MemoryLimit       int64                           `json:"memoryLimit,omitempty"`
	HealthCheck       MicroserviceHealthCheck         `json:"healthCheck,omitempty"`
	NatsConfig        *MicroserviceNatsConfig         `json:"natsConfig,omitempty"`
	ServiceAccount    *MicroserviceServiceAccountRef  `json:"serviceAccount,omitempty"`
}

// MicroserviceServiceAccountRef is the optional serviceAccount field in a microservice spec (YAML or API).
// When set, the microservice uses a service account with the given roleRef (e.g. roleRef.name).
type MicroserviceServiceAccountRef struct {
	RoleRef RoleRef `json:"roleRef"`
}

type MicroserviceHealthCheck struct {
	Test          []string `json:"test"`
	Interval      *int64   `json:"interval,omitempty"`
	Timeout       *int64   `json:"timeout,omitempty"`
	Retries       *int     `json:"retries,omitempty"`
	StartPeriod   *int64   `json:"startPeriod,omitempty"`
	StartInterval *int64   `json:"startInterval,omitempty"`
}

// MicroserviceNatsConfig holds NATS configuration for a microservice (Controller microserviceNatsConfig).
type MicroserviceNatsConfig struct {
	NatsAccess bool   `json:"natsAccess"`
	NatsRule   string `json:"natsRule,omitempty"`
}

type MicroserviceExtraHost struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Value   string `json:"value,omitempty"`
}

type MicroserviceCreateResponse struct {
	UUID string `json:"uuid"`
}

type MicroserviceListResponse struct {
	Microservices []MicroserviceInfo
}

type MicroservicePortMappingListResponse struct {
	PortMappings []MicroservicePortMappingInfo `json:"ports"`
}

// Users

type User struct {
	Name            string `json:"firstName"`
	Surname         string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	SubscriptionKey string `json:"subscriptionKey"`
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
}

type UserResponse struct {
	Name            string `json:"firstName"`
	Surname         string `json:"lastName"`
	Email           string `json:"email"`
	SubscriptionKey string `json:"subscriptionKey"`
}

type ControllerVersions struct {
	Controller string `json:"controller"`
	EcnViewer  string `json:"ecnViewer"`
}

type ControllerStatus struct {
	Status        string             `json:"status"`
	UptimeSeconds float64            `json:"uptimeSec"`
	Versions      ControllerVersions `json:"versions"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Totp     string `json:"totp"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type WithTokenRequest struct {
	AccessToken string `json:"accessToken"`
}

type UpdateUserPasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ListAgentsRequest struct {
	System  bool              `json:"system"`
	Filters []AgentListFilter `json:"filters"`
}

type CreateAgentRequest struct {
	AgentUpdateRequest `json:",inline"` //nolint:revive // inline embedding matches Controller API JSON
}

type CreateAgentResponse struct {
	UUID string
}

type GetAgentProvisionKeyResponse struct {
	Key             string `json:"key"`
	CaCert          string `json:"caCert"`
	ExpireTimeMsUTC int64  `json:"expirationTime"`
}

type AgentInfo struct {
	UUID                      string            `json:"uuid" yaml:"uuid"`
	Name                      string            `json:"name" yaml:"name"`
	Host                      string            `json:"host" yaml:"host"`
	Location                  string            `json:"location" yaml:"location"`
	Latitude                  float64           `json:"latitude" yaml:"latitude"`
	Longitude                 float64           `json:"longitude" yaml:"longitude"`
	Description               string            `json:"description" yaml:"description"`
	DockerURL                 string            `json:"dockerUrl" yaml:"dockerUrl"`
	ContainerEngine           string            `json:"containerEngine" yaml:"containerEngine"`
	DeploymentType            string            `json:"deploymentType" yaml:"deploymentType"`
	DiskLimit                 int64             `json:"diskLimit" yaml:"diskLimit"`
	DiskDirectory             string            `json:"diskDirectory" yaml:"diskDirectory"`
	MemoryLimit               int64             `json:"memoryLimit" yaml:"memoryLimit"`
	CPULimit                  int64             `json:"cpuLimit" yaml:"cpuLimit"`
	LogLimit                  int64             `json:"logLimit" yaml:"logLimit"`
	LogDirectory              string            `json:"logDirectory" yaml:"logDirectory"`
	LogFileCount              int64             `json:"logFileCount" yaml:"logFileCount"`
	StatusFrequency           float64           `json:"statusFrequency" yaml:"statusFrequency"`
	ChangeFrequency           float64           `json:"changeFrequency" yaml:"changeFrequency"`
	DeviceScanFrequency       float64           `json:"deviceScanFrequency" yaml:"deviceScanFrequency"`
	BluetoothEnabled          bool              `json:"bluetoothEnabled" yaml:"bluetoothEnabled"`
	WatchdogEnabled           bool              `json:"watchdogEnabled" yaml:"watchdogEnabled"`
	GpsMode                   string            `json:"gpsMode" yaml:"gpsMode"`
	GpsScanFrequency          float64           `json:"gpsScanFrequency" yaml:"gpsScanFrequency"`
	GpsDevice                 string            `json:"gpsDevice" yaml:"gpsDevice"`
	EdgeGuardFrequency        float64           `json:"edgeGuardFrequency" yaml:"edgeGuardFrequency"`
	AbstractedHardwareEnabled bool              `json:"abstractedHardwareEnabled" yaml:"abstractedHardwareEnabled"`
	CreatedTimeRFC3339        string            `json:"createdAt" yaml:"created"`
	UpdatedTimeRFC3339        string            `json:"updatedAt" yaml:"updated"`
	LastActive                int64             `json:"lastActive" yaml:"lastActive"`
	DaemonStatus              string            `json:"daemonStatus" yaml:"daemonStatus"`
	UptimeMs                  int64             `json:"daemonOperatingDuration" yaml:"uptime"`
	MemoryUsage               float64           `json:"memoryUsage" yaml:"memoryUsage"`
	DiskUsage                 float64           `json:"diskUsage" yaml:"diskUsage"`
	CPUUsage                  float64           `json:"cpuUsage" yaml:"cpuUsage"`
	SystemAvailableMemory     float64           `json:"systemAvailableMemory" yaml:"systemAvailableMemory"`
	SystemAvailableDisk       float64           `json:"systemAvailableDisk" yaml:"systemAvailableDisk"`
	SystemTotalCPU            float64           `json:"systemTotalCPU" yaml:"systemTotalCPU"`
	MemoryViolation           string            `json:"memoryViolation" yaml:"memoryViolation"`
	DiskViolation             string            `json:"diskViolation" yaml:"diskViolation"`
	CPUViolation              string            `json:"cpuViolation" yaml:"cpuViolation"`
	MicroserviceStatus        string            `json:"microserviceStatus" yaml:"microserviceStatus"`
	RepositoryCount           int64             `json:"repositoryCount" yaml:"repositoryCount"`
	RepositoryStatus          string            `json:"repositoryStatus" yaml:"repositoryStatus"`
	LastStatusTimeMsUTC       int64             `json:"lastStatusTime" yaml:"lastStatusTime"`
	IPAddress                 string            `json:"ipAddress" yaml:"ipAddress"`
	IPAddressExternal         string            `json:"ipAddressExternal" yaml:"ipAddressExternal"`
	ProcessedMessaged         int64             `json:"processedMessages" yaml:"ProcessedMessages"`
	MicroserviceMessageCount  int64             `json:"microserviceMessageCounts" yaml:"microserviceMessageCount"`
	MessageSpeed              float64           `json:"messageSpeed" yaml:"messageSpeed"`
	LastCommandTimeMsUTC      int64             `json:"lastCommandTime" yaml:"lastCommandTime"`
	NetworkInterface          string            `json:"networkInterface" yaml:"networkInterface"`
	Version                   string            `json:"version" yaml:"version"`
	IsReadyToUpgrade          bool              `json:"isReadyToUpgrade" yaml:"isReadyToUpgrade"`
	IsReadyToRollback         bool              `json:"isReadyToRollback" yaml:"isReadyToRollback"`
	Tunnel                    string            `json:"tunnel" yaml:"tunnel"`
	FogType                   int               `json:"fogTypeId" yaml:"fogTypeId"`
	RouterMode                string            `json:"routerMode" yaml:"routerMode"`
	NetworkRouter             *string           `json:"networkRouter,omitempty" yaml:"networkRouter,omitempty"`
	UpstreamRouters           *[]string         `json:"upstreamRouters,omitempty" yaml:"upstreamRouters,omitempty"`
	MessagingPort             *int              `json:"messagingPort,omitempty" yaml:"messagingPort,omitempty"`
	EdgeRouterPort            *int              `json:"edgeRouterPort,omitempty" yaml:"edgeRouterPort,omitempty"`
	InterRouterPort           *int              `json:"interRouterPort,omitempty" yaml:"interRouterPort,omitempty"`
	LogLevel                  *string           `json:"logLevel" yaml:"logLevel"`
	DockerPruningFrequency    *float64          `json:"dockerPruningFrequency" yaml:"dockerPruningFrequency"`
	AvailableDiskThreshold    *float64          `json:"availableDiskThreshold" yaml:"availableDiskThreshold"`
	Tags                      *[]string         `json:"tags,omitempty" yaml:"tags,omitempty"`
	TimeZone                  string            `json:"timeZone" yaml:"timeZone"`
	IsSystem                  bool              `json:"isSystem" yaml:"-"`
	VolumeMounts              []VolumeMountInfo `json:"volumeMounts" yaml:"volumeMounts"`
	SecurityStatus            string            `json:"securityStatus" yaml:"securityStatus"`
	SecurityViolationInfo     string            `json:"securityViolationInfo" yaml:"securityViolationInfo"`
	WarningMessage            string            `json:"warningMessage" yaml:"warningMessage"`
	GpsStatus                 string            `json:"gpsStatus" yaml:"gpsStatus"`
	// NATS-related fields (Controller iofog schema)
	NatsMode            *string   `json:"natsMode,omitempty" yaml:"natsMode,omitempty"` // none, leaf, server
	NatsServerPort      *int      `json:"natsServerPort,omitempty" yaml:"natsServerPort,omitempty"`
	NatsLeafPort        *int      `json:"natsLeafPort,omitempty" yaml:"natsLeafPort,omitempty"`
	NatsClusterPort     *int      `json:"natsClusterPort,omitempty" yaml:"natsClusterPort,omitempty"`
	NatsMqttPort        *int      `json:"natsMqttPort,omitempty" yaml:"natsMqttPort,omitempty"`
	NatsHTTPPort        *int      `json:"natsHttpPort,omitempty" yaml:"natsHttpPort,omitempty"`
	UpstreamNatsServers *[]string `json:"upstreamNatsServers,omitempty" yaml:"upstreamNatsServers,omitempty"`
	JsStorageSize       *string   `json:"jsStorageSize,omitempty" yaml:"jsStorageSize,omitempty"`
	JsMemoryStoreSize   *string   `json:"jsMemoryStoreSize,omitempty" yaml:"jsMemoryStoreSize,omitempty"`
}

type RouterConfig struct {
	RouterMode      *string `json:"routerMode,omitempty" yaml:"routerMode,omitempty"`
	MessagingPort   *int    `json:"messagingPort,omitempty" yaml:"messagingPort,omitempty"`
	EdgeRouterPort  *int    `json:"edgeRouterPort,omitempty" yaml:"edgeRouterPort,omitempty"`
	InterRouterPort *int    `json:"interRouterPort,omitempty" yaml:"interRouterPort,omitempty"`
}

type NatsConfig struct {
	NatsMode          *string `json:"natsMode,omitempty" yaml:"natsMode,omitempty"` // none, leaf, server
	NatsServerPort    *int    `json:"natsServerPort,omitempty" yaml:"natsServerPort,omitempty"`
	NatsLeafPort      *int    `json:"natsLeafPort,omitempty" yaml:"natsLeafPort,omitempty"`
	NatsClusterPort   *int    `json:"natsClusterPort,omitempty" yaml:"natsClusterPort,omitempty"`
	NatsMqttPort      *int    `json:"natsMqttPort,omitempty" yaml:"natsMqttPort,omitempty"`
	NatsHTTPPort      *int    `json:"natsHttpPort,omitempty" yaml:"natsHttpPort,omitempty"`
	JsStorageSize     *string `json:"jsStorageSize,omitempty" yaml:"jsStorageSize,omitempty"`
	JsMemoryStoreSize *string `json:"jsMemoryStoreSize,omitempty" yaml:"jsMemoryStoreSize,omitempty"`
}

type AgentConfiguration struct {
	NetworkInterface          *string   `json:"networkInterface,omitempty" yaml:"networkInterface"`
	DockerURL                 *string   `json:"dockerUrl,omitempty" yaml:"dockerUrl"`
	ContainerEngine           *string   `json:"containerEngine,omitempty" yaml:"containerEngine"`
	DeploymentType            *string   `json:"deploymentType,omitempty" yaml:"deploymentType"`
	DiskLimit                 *int64    `json:"diskLimit,omitempty" yaml:"diskLimit"`
	DiskDirectory             *string   `json:"diskDirectory,omitempty" yaml:"diskDirectory"`
	MemoryLimit               *int64    `json:"memoryLimit,omitempty" yaml:"memoryLimit"`
	CPULimit                  *int64    `json:"cpuLimit,omitempty" yaml:"cpuLimit"`
	LogLimit                  *int64    `json:"logLimit,omitempty" yaml:"logLimit"`
	LogDirectory              *string   `json:"logDirectory,omitempty" yaml:"logDirectory"`
	LogFileCount              *int64    `json:"logFileCount,omitempty" yaml:"logFileCount"`
	StatusFrequency           *float64  `json:"statusFrequency,omitempty" yaml:"statusFrequency"`
	ChangeFrequency           *float64  `json:"changeFrequency,omitempty" yaml:"changeFrequency"`
	DeviceScanFrequency       *float64  `json:"deviceScanFrequency,omitempty" yaml:"deviceScanFrequency"`
	BluetoothEnabled          *bool     `json:"bluetoothEnabled,omitempty" yaml:"bluetoothEnabled"`
	WatchdogEnabled           *bool     `json:"watchdogEnabled,omitempty" yaml:"watchdogEnabled"`
	GpsMode                   *string   `yaml:"gpsMode,omitempty" json:"gpsMode,omitempty"`
	GpsScanFrequency          *float64  `yaml:"gpsScanFrequency,omitempty" json:"gpsScanFrequency,omitempty"`
	GpsDevice                 *string   `yaml:"gpsDevice,omitempty" json:"gpsDevice,omitempty"`
	EdgeGuardFrequency        *float64  `yaml:"edgeGuardFrequency,omitempty" json:"edgeGuardFrequency,omitempty"`
	AbstractedHardwareEnabled *bool     `json:"abstractedHardwareEnabled,omitempty" yaml:"abstractedHardwareEnabled"`
	IsSystem                  *bool     `json:"isSystem,omitempty" yaml:"-"` // Can't specify system agent using yaml file.
	UpstreamRouters           *[]string `json:"upstreamRouters,omitempty" yaml:"upstreamRouters,omitempty"`
	NetworkRouter             *string   `json:"networkRouter,omitempty" yaml:"networkRouter,omitempty"`
	Host                      *string   `json:"host,omitempty" yaml:"host,omitempty"`
	RouterConfig              `json:",omitempty" yaml:"routerConfig,omitempty"`
	LogLevel                  *string   `json:"logLevel,omitempty" yaml:"logLevel"`
	DockerPruningFrequency    *float64  `json:"dockerPruningFrequency,omitempty" yaml:"dockerPruningFrequency"`
	AvailableDiskThreshold    *float64  `json:"availableDiskThreshold,omitempty" yaml:"availableDiskThreshold"`
	TimeZone                  string    `json:"timeZone,omitempty" yaml:"timeZone"`
	UpstreamNatsServers       *[]string `json:"upstreamNatsServers,omitempty" yaml:"upstreamNatsServers,omitempty"`
	NatsConfig                `json:",omitempty" yaml:"natsConfig,omitempty"`
}

type AgentUpdateRequest struct {
	UUID        string    `json:"-"`
	Name        string    `json:"name,omitempty" yaml:"name"`
	Location    string    `json:"location,omitempty" yaml:"location"`
	Latitude    float64   `json:"latitude,omitempty" yaml:"latitude"`
	Longitude   float64   `json:"longitude,omitempty" yaml:"longitude"`
	Description string    `json:"description,omitempty" yaml:"description"`
	FogType     *int64    `json:"fogType,omitempty" yaml:"agentType"`
	Tags        *[]string `json:"tags,omitempty" yaml:"tags"`
	AgentConfiguration
}

type ListAgentsResponse struct {
	Agents []AgentInfo `json:"fogs"`
}

type AgentListFilter struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Condition string `json:"condition"`
}

type Router struct {
	RouterConfig
	Host string `json:"host"`
}

type UpdateConfigRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RouteListResponse is the response for listing routes.
//
// Deprecated: Controller no longer exposes /api/v3/routes. Use NATS for messaging.
type RouteListResponse struct {
	Routes []Route `json:"routes"`
}

// Route represents a route from one microservice to another.
//
// Deprecated: Controller no longer exposes routing endpoints. Use NATS for messaging.
type Route struct {
	Name        string `json:"name"`
	Application string `json:"application"`
	From        string `json:"from"`
	To          string `json:"to"`
}

// ApplicationRouteCreateRequest is the request body for creating a route.
//
// Deprecated: Controller no longer exposes /api/v3/routes. Use NATS for messaging.
type ApplicationRouteCreateRequest struct {
	Name string `json:"name"`
	From string `json:"from"`
	To   string `json:"to"`
}

type EdgeResourceDisplay struct {
	Name  string `json:"name,omitempty"`
	Icon  string `json:"icon,omitempty"`
	Color string `json:"color,omitempty"`
}

type EdgeResourceMetadata struct {
	Name              string               `json:"name,omitempty"`
	Description       string               `json:"description,omitempty"`
	Version           string               `json:"version,omitempty"`
	InterfaceProtocol string               `json:"interfaceProtocol,omitempty"`
	Display           *EdgeResourceDisplay `json:"display,omitempty"`
	Interface         HTTPEdgeResource     `json:"interface,omitempty"` // TODO: Make this generic
	OrchestrationTags []string             `json:"orchestrationTags,omitempty"`
	Custom            map[string]any       `json:"custom,omitempty"`
}

type HTTPEdgeResource struct {
	Endpoints []HTTPEndpoint `json:"endpoints,omitempty"`
}

type HTTPEndpoint struct {
	Name   string `json:"name,omitempty"`
	Method string `json:"method,omitempty"`
	URL    string `json:"url,omitempty"`
}

type LinkEdgeResourceRequest struct {
	AgentUUID           string `json:"uuid"`
	EdgeResourceName    string `json:"-"`
	EdgeResourceVersion string `json:"-"`
}

type ListEdgeResourceResponse struct {
	EdgeResources []EdgeResourceMetadata `json:"edgeResources"`
}

// Secrets
type SecretInfo struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Data      map[string]string `json:"data"`
	CreatedAt string            `json:"createdAt,omitempty"`
	UpdatedAt string            `json:"updatedAt,omitempty"`
}

type SecretCreateRequest struct {
	Name string            `json:"name"`
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

type SecretUpdateRequest struct {
	Name string            `json:"name,omitempty"`
	Data map[string]string `json:"data"`
}

type SecretListResponse struct {
	Secrets []SecretInfo `json:"secrets"`
}

// Services
type ServiceInfo struct {
	Tags               []string `json:"tags"`
	Name               string   `json:"name"`
	Type               string   `json:"type"`
	Resource           string   `json:"resource"`
	TargetPort         int      `json:"targetPort"`
	ServicePort        int      `json:"servicePort"`
	K8sType            string   `json:"k8sType"`
	BridgePort         int      `json:"bridgePort"`
	DefaultBridge      string   `json:"defaultBridge"`
	ServiceEndpoint    string   `json:"serviceEndpoint"`
	ProvisioningStatus string   `json:"provisioningStatus"`
	ProvisioningError  string   `json:"provisioningError"`
	CreatedAt          string   `json:"createdAt,omitempty"`
	UpdatedAt          string   `json:"updatedAt,omitempty"`
}

type ServiceCreateRequest struct {
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Resource      string   `json:"resource"`
	TargetPort    int      `json:"targetPort"`
	ServicePort   int      `json:"servicePort,omitempty"`
	K8sType       string   `json:"k8sType,omitempty"`
	DefaultBridge string   `json:"defaultBridge,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

type ServiceUpdateRequest struct {
	Name          string   `json:"name,omitempty"`
	Type          string   `json:"type"`
	Resource      string   `json:"resource"`
	TargetPort    int      `json:"targetPort"`
	ServicePort   int      `json:"servicePort,omitempty"`
	K8sType       string   `json:"k8sType,omitempty"`
	DefaultBridge string   `json:"defaultBridgePort,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

type ServiceListResponse struct {
	Services []ServiceInfo `json:"services"`
}

// ConfigMaps
type ConfigMapInfo struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Data      map[string]string `json:"data"`
	Immutable bool              `json:"immutable,omitempty"`
	CreatedAt string            `json:"createdAt,omitempty"`
	UpdatedAt string            `json:"updatedAt,omitempty"`
}

type ConfigMapCreateRequest struct {
	Name      string            `json:"name"`
	Data      map[string]string `json:"data"`
	Immutable bool              `json:"immutable,omitempty"`
}

type ConfigMapUpdateRequest struct {
	Name      string            `json:"name,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
	Immutable bool              `json:"immutable,omitempty"`
}

type ConfigMapListResponse struct {
	ConfigMaps []ConfigMapInfo `json:"configMaps"`
}

// VolumeMounts
type VolumeMountInfo struct {
	Name          string `json:"name"`
	UUID          string `json:"uuid,omitempty"`
	ConfigMapName string `json:"configMapName,omitempty"`
	SecretName    string `json:"secretName,omitempty"`
	Version       int    `json:"version"`
	CreatedAt     string `json:"createdAt,omitempty"`
	UpdatedAt     string `json:"updatedAt,omitempty"`
}

type VolumeMountCreateRequest struct {
	Name          string `json:"name"`
	ConfigMapName string `json:"configMapName,omitempty"`
	SecretName    string `json:"secretName,omitempty"`
}

type VolumeMountUpdateRequest struct {
	Name          string `json:"name,omitempty"`
	ConfigMapName string `json:"configMapName,omitempty"`
	SecretName    string `json:"secretName,omitempty"`
}

type VolumeMountListResponse struct {
	VolumeMounts []VolumeMountInfo `json:"volumeMounts"`
}

type VolumeMountLinkRequest struct {
	Name     string   `json:"name"`
	FogUUIDs []string `json:"fogUuids"`
}

type VolumeMountUnlinkRequest struct {
	Name     string   `json:"name"`
	FogUUIDs []string `json:"fogUuids"`
}

// Certificate Types
type CertificateCreateRequest struct {
	Name       string              `json:"name"`
	Subject    string              `json:"subject"`
	Hosts      string              `json:"hosts"`
	Expiration int                 `json:"expiration,omitempty"`
	CA         CertificateCreateCA `json:"ca"`
}

type CertificateCreateResponse struct {
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	Hosts     string    `json:"hosts"`
	ValidFrom time.Time `json:"validFrom"`
	ValidTo   time.Time `json:"validTo"`
	CAName    string    `json:"caName"`
}

type CertificateCACreateResponse struct {
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	Type      string    `json:"type"`
	ValidFrom time.Time `json:"validFrom"`
	ValidTo   time.Time `json:"validTo"`
}

type CertificateCreateCA struct {
	Type       string `json:"type"`
	SecretName string `json:"secretName,omitempty"`
}

type CACreateRequest struct {
	Name       string `json:"name"`
	Subject    string `json:"subject,omitempty"`
	Expiration int    `json:"expiration,omitempty"`
	Type       string `json:"type"`
	SecretName string `json:"secretName,omitempty"`
}

type CertificateInfo struct {
	Name             string                 `json:"name"`
	Subject          string                 `json:"subject"`
	Hosts            string                 `json:"hosts"`
	IsCA             bool                   `json:"isCA"`
	ValidFrom        time.Time              `json:"validFrom"`
	ValidTo          time.Time              `json:"validTo"`
	SerialNumber     string                 `json:"serialNumber"`
	CAName           *string                `json:"caName"`
	CertificateChain []CertificateChainItem `json:"certificateChain"`
	DaysRemaining    int                    `json:"daysRemaining"`
	IsExpired        bool                   `json:"isExpired"`
	Data             CertificateData        `json:"data"`
}

type CAInfo struct {
	Name         string          `json:"name"`
	Subject      string          `json:"subject"`
	IsCA         bool            `json:"isCA"`
	ValidFrom    time.Time       `json:"validFrom"`
	ValidTo      time.Time       `json:"validTo"`
	SerialNumber string          `json:"serialNumber"`
	Data         CertificateData `json:"data"`
}

type CertificateData struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

type CertificateChainItem struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
}

type CertificateListResponse struct {
	Certificates []CertificateInfo `json:"certificates"`
}

type CAListResponse struct {
	CAs []CAInfo `json:"cas"`
}

type AttachExecMicroserviceRequest struct {
	UUID string `json:"uuid"`
}

type DetachExecMicroserviceRequest struct {
	UUID string `json:"uuid"`
}

type AttachExecToAgentRequest struct {
	UUID  string  `json:"uuid"`
	Image *string `json:"image,omitempty"`
}

type DetachExecFromAgentRequest struct {
	UUID string `json:"uuid"`
}

// NATS API types (Controller /api/v3/nats/*)

// NatsOperatorResponse is the response for GET /nats/operator.
type NatsOperatorResponse struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	JWT       string `json:"jwt"`
}

// NatsHubResponse is the response for GET/PUT /nats/hub.
type NatsHubResponse struct {
	Host        string `json:"host"`
	ServerPort  int    `json:"serverPort"`
	ClusterPort int    `json:"clusterPort"`
	LeafPort    int    `json:"leafPort"`
	MqttPort    int    `json:"mqttPort"`
	HTTPPort    int    `json:"httpPort"`
}

// NatsHubRequest is the request body for PUT /nats/hub.
type NatsHubRequest struct {
	Host        *string `json:"host,omitempty"`
	ServerPort  *int    `json:"serverPort,omitempty"`
	ClusterPort *int    `json:"clusterPort,omitempty"`
	LeafPort    *int    `json:"leafPort,omitempty"`
	MqttPort    *int    `json:"mqttPort,omitempty"`
	HTTPPort    *int    `json:"httpPort,omitempty"`
}

// NatsBootstrapResponse is the response for GET /nats/bootstrap.
// Used by the K8s operator when the Controller runs on the Kubernetes control plane.
// Contains operator JWT/publicKey/seed and system account JWT/publicKey and system user creds (base64) for NATS bootstrap.
type NatsBootstrapResponse struct {
	OperatorJwt            string `json:"operatorJwt"`
	OperatorPublicKey      string `json:"operatorPublicKey"`
	OperatorSeed           string `json:"operatorSeed"`
	SystemAccountJwt       string `json:"systemAccountJwt"`
	SystemAccountPublicKey string `json:"systemAccountPublicKey"`
	SysUserCredsBase64     string `json:"sysUserCredsBase64"`
}

// NatsAccountInfo is a single NATS account (list or get by app).
type NatsAccountInfo struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	PublicKey     string `json:"publicKey"`
	JWT           string `json:"jwt"`
	IsSystem      bool   `json:"isSystem"`
	IsLeafSystem  bool   `json:"isLeafSystem,omitempty"`
	ApplicationID int    `json:"applicationId,omitempty"`
}

// NatsListAccountsResponse is the response for GET /nats/accounts.
type NatsListAccountsResponse struct {
	Accounts []NatsAccountInfo `json:"accounts"`
}

// NatsUserInfo is a single NATS user (list or create).
type NatsUserInfo struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	PublicKey        string `json:"publicKey"`
	JWT              string `json:"jwt"`
	IsBearer         bool   `json:"isBearer"`
	AccountID        int    `json:"accountId,omitempty"`
	AccountName      string `json:"accountName,omitempty"`
	ApplicationID    int    `json:"applicationId,omitempty"`
	ApplicationName  string `json:"applicationName,omitempty"`
	MicroserviceUUID string `json:"microserviceUuid,omitempty"`
}

// NatsListUsersResponse is the response for GET /nats/users.
type NatsListUsersResponse struct {
	Users []NatsUserInfo `json:"users"`
}

// NatsListAccountUsersResponse is the response for GET /nats/accounts/:appName/users.
type NatsListAccountUsersResponse struct {
	Users []NatsUserInfo `json:"users"`
}

// NatsCreateUserRequest is the request body for POST /nats/accounts/:appName/users.
type NatsCreateUserRequest struct {
	Name      string `json:"name"`
	ExpiresIn *int64 `json:"expiresIn,omitempty"`
	NatsRule  string `json:"natsRule,omitempty"`
}

// NatsEnsureAccountRequest is the request body for POST /nats/accounts/:appName (ensure account).
type NatsEnsureAccountRequest struct {
	NatsRule string `json:"natsRule,omitempty"`
}

// NatsUserCredsResponse is the response for GET /nats/accounts/:appName/users/:userName/creds.
type NatsUserCredsResponse struct {
	CredsBase64 string `json:"credsBase64"`
}

// NatsCreateMqttBearerRequest is the request body for POST /nats/accounts/:appName/mqtt-bearer.
type NatsCreateMqttBearerRequest struct {
	Name      string `json:"name"`
	ExpiresIn *int64 `json:"expiresIn,omitempty"`
	NatsRule  string `json:"natsRule,omitempty"`
}

// FlexInt64 unmarshals from JSON string or number (Controller may send e.g. "-1" as string for BIGINT).
type FlexInt64 int64

func (v *FlexInt64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*v = FlexInt64(n)
		return nil
	}
	var n int64
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*v = FlexInt64(n)
	return nil
}

func (v FlexInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(v))
}

// FlexInt unmarshals from JSON string or number (Controller may send e.g. "-1" as string).
type FlexInt int

func (v *FlexInt) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		n, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return err
		}
		*v = FlexInt(n)
		return nil
	}
	var n int
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*v = FlexInt(n)
	return nil
}

func (v FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(v))
}

// NatsRuleInfo is a generic NATS account or user rule (list responses).
// Rule payloads from the Controller can have many optional fields; use map or extend as needed.
type NatsRuleInfo struct {
	ID                     int        `json:"id"`
	Name                   string     `json:"name"`
	Description            string     `json:"description,omitempty"`
	IsSystem               bool       `json:"isSystem,omitempty"`
	InfoURL                string     `json:"infoUrl,omitempty"`
	MaxConnections         *FlexInt   `json:"maxConnections,omitempty"`
	MaxLeafNodeConnections *FlexInt   `json:"maxLeafNodeConnections,omitempty"`
	MaxData                *FlexInt64 `json:"maxData,omitempty"`
	MaxExports             *FlexInt   `json:"maxExports,omitempty"`
	MaxImports             *FlexInt   `json:"maxImports,omitempty"`
	MaxMsgPayload          *FlexInt64 `json:"maxMsgPayload,omitempty"`
	MaxSubscriptions       *FlexInt   `json:"maxSubscriptions,omitempty"`
	ExportsAllowWildcards  *bool      `json:"exportsAllowWildcards,omitempty"`
	DisallowBearer         *bool      `json:"disallowBearer,omitempty"`
	ResponsePermissions    any        `json:"responsePermissions,omitempty"`
	RespMax                *FlexInt   `json:"respMax,omitempty"`
	RespTTL                *FlexInt64 `json:"respTtl,omitempty"`
	Imports                any        `json:"imports,omitempty"`
	Exports                any        `json:"exports,omitempty"`
	MemStorage             *FlexInt64 `json:"memStorage,omitempty"`
	DiskStorage            *FlexInt64 `json:"diskStorage,omitempty"`
	Streams                any        `json:"streams,omitempty"`
	Consumer               any        `json:"consumer,omitempty"`
	MaxAckPending          *FlexInt   `json:"maxAckPending,omitempty"`
	MemMaxStreamBytes      *FlexInt64 `json:"memMaxStreamBytes,omitempty"`
	DiskMaxStreamBytes     *FlexInt64 `json:"diskMaxStreamBytes,omitempty"`
	MaxBytesRequired       *bool      `json:"maxBytesRequired,omitempty"`
	TieredLimits           any        `json:"tieredLimits,omitempty"`
	MaxPayload             *FlexInt64 `json:"maxPayload,omitempty"`
	BearerToken            *bool      `json:"bearerToken,omitempty"`
	ProxyRequired          *bool      `json:"proxyRequired,omitempty"`
	AllowedConnectionTypes any        `json:"allowedConnectionTypes,omitempty"`
	Src                    any        `json:"src,omitempty"`
	Times                  any        `json:"times,omitempty"`
	TimesLocation          string     `json:"timesLocation,omitempty"`
	PubAllow               any        `json:"pubAllow,omitempty"`
	PubDeny                any        `json:"pubDeny,omitempty"`
	SubAllow               any        `json:"subAllow,omitempty"`
	SubDeny                any        `json:"subDeny,omitempty"`
	Tags                   any        `json:"tags,omitempty"`
}

// NatsListAccountRulesResponse is the response for GET /nats/account-rules.
type NatsListAccountRulesResponse struct {
	Rules []NatsRuleInfo `json:"rules"`
}

// NatsListUserRulesResponse is the response for GET /nats/user-rules.
type NatsListUserRulesResponse struct {
	Rules []NatsRuleInfo `json:"rules"`
}

// NatsAccountRulePayload is the request body for creating/patching an account rule (JSON).
// Omit empty fields; Controller validates per schema.
type NatsAccountRulePayload map[string]any

// NatsUserRulePayload is the request body for creating/patching a user rule (JSON).
type NatsUserRulePayload map[string]any
