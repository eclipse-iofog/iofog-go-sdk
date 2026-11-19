package microservices

import (
	"errors"
	"os"
	"strconv"
)

// EdgeletAPIClient talks to the EdgeletAPI v1 microservices endpoints.
type EdgeletAPIClient struct {
	id         string
	options    ClientOptions
	httpClient *edgeletAPIHttpClient
	wsClient   *edgeletAPIWsClient
}

func (client *EdgeletAPIClient) initClient(options ClientOptions) {
	client.options = options
	client.httpClient = newEdgeletAPIHttpClient(options)
	client.wsClient = newEdgeletAPIWsClient(options)
}

// NewEdgeletAPIClient creates an EdgeletAPI v1 client with explicit options.
func NewEdgeletAPIClient(id string, opts ...ClientOption) (*EdgeletAPIClient, error) {
	if id == "" {
		return nil, errors.New("cannot create client with empty id")
	}
	client := EdgeletAPIClient{id: id}
	client.initClient(applyClientOptions(defaultClientOptions(), opts...))
	return &client, nil
}

// NewDefaultEdgeletAPIClient creates an EdgeletAPI v1 client using mounted token/CA defaults.
func NewDefaultEdgeletAPIClient() (*EdgeletAPIClient, error) {
	microserviceUID := os.Getenv(MicroserviceUID)
	if microserviceUID == "" {
		return nil, errors.New("cannot create client with empty id: " + MicroserviceUID + " environment variable is not set")
	}
	ssl, err := strconv.ParseBool(os.Getenv(SSL))
	if err != nil {
		logger.Println("Empty or malformed", SSL, "environment variable. Using default value of", SSLDefault)
		ssl = SSLDefault
	}
	return NewEdgeletAPIClient(microserviceUID,
		WithTLS(ssl),
		WithHost(HostDefault),
		WithPort(PortEdgeletAPI),
	)
}

func (client *EdgeletAPIClient) GetConfig() (map[string]any, error) {
	return client.httpClient.getConfig()
}

func (client *EdgeletAPIClient) GetConfigIntoStruct(config any) error {
	return client.httpClient.getConfigIntoStruct(config)
}

func (client *EdgeletAPIClient) EstablishControlWsConnection(signalBufSize int) <-chan byte {
	if signalBufSize == 0 {
		signalBufSize = DefaultSignalBufferSize
	}
	signalChannel := make(chan byte, signalBufSize)
	go client.wsClient.connectToControlWs(signalChannel)
	return signalChannel
}
