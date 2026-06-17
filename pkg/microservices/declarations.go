package microservices

import (
	"log"
	"os"
	"time"
)

const (
	PortIoFog         = 54321
	MicroserviceUID   = "EDGELET_MICROSERVICE_UID"
	SSL               = "SSL"
	SSLDefault        = true
	HostDefault       = "edgelet.default.svc.bridge.local"
	FallbackHostLocal = "127.0.0.1"

	DefaultServiceAccountTokenPath = "/var/run/secrets/edgelet.iofog.org/serviceaccount/token"
	DefaultServiceAccountCAPath    = "/var/run/secrets/edgelet.iofog.org/serviceaccount/ca.crt"

	URLGetConfigV1    = "/v1/microservices/config"
	URLGetControlWSV1 = "/v1/microservices/control"
	ApplicationJSON   = "application/json"
	SchemeHTTP        = "http"
	SchemeHTTPS       = "https"
	SchemeWS          = "ws"
	SchemeWSS         = "wss"

	CODE_ACK            = 0xB
	CODE_CONTROL_SIGNAL = 0xC

	DefaultSignalBufferSize     = 5
	DefaultRequestTimeout       = 15 * time.Second
	DefaultWSHandshakeTimeout   = 10 * time.Second
	DefaultWSReconnectBaseDelay = time.Second
	DefaultWSReconnectMaxDelay  = 30 * time.Second
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
)
