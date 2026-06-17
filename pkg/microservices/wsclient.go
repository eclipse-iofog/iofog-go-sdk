package microservices

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	ws "github.com/gorilla/websocket"
)

type edgeletAPIWsClient struct {
	options ClientOptions
}

func newEdgeletAPIWsClient(options ClientOptions) *edgeletAPIWsClient {
	return &edgeletAPIWsClient{options: options}
}

func (client *edgeletAPIWsClient) connectToControlWs(signalChannel chan<- byte) {
	attempt := 0
	for {
		conn, err := client.dialControlWS()
		if err != nil {
			attempt++
			logger.Println(err.Error(), "Reconnecting to EdgeletAPI control ws...")
			sleepWithCap(attempt, client.options.WSReconnectBaseDelay, client.options.WSReconnectMaxDelay)
			continue
		}
		logger.Println("Control ws connection has been established.")
		attempt = 0
		if err := client.listenControl(conn, signalChannel); err != nil {
			logger.Println("Control ws reconnect after corruption:", err.Error())
		}
		_ = conn.Close()
	}
}

func (client *edgeletAPIWsClient) dialControlWS() (*ws.Conn, error) {
	dialer := ws.Dialer{
		HandshakeTimeout: client.options.WSHandshakeTimeout,
	}
	if client.options.UseTLS {
		rootCAs, err := loadRootCAs(client.options.CAPath)
		if err != nil {
			return nil, err
		}
		dialer.TLSClientConfig = buildTLSConfig(rootCAs)
	}
	token, err := readBearerToken(client.options.TokenPath)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)

	hosts := client.options.requestHosts()
	var lastErr error
	for idx, host := range hosts {
		url := fmt.Sprint(client.options.wsBaseURLForHost(host), URLGetControlWSV1)
		conn, _, dialErr := dialer.Dial(url, header)
		if dialErr != nil {
			lastErr = dialErr
			if idx < len(hosts)-1 && isRetriableHostError(dialErr) {
				logger.Println("WARN: failed to dial control ws host", host, "trying fallback host:", dialErr)
				continue
			}
			return nil, dialErr
		}
		setCustomPingHandler(conn)
		return conn, nil
	}
	return nil, lastErr
}

func (client *edgeletAPIWsClient) listenControl(conn *ws.Conn, signalChannel chan<- byte) error {
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if len(payload) == 0 {
			continue
		}
		if payload[0] != CODE_CONTROL_SIGNAL {
			continue
		}
		signalChannel <- payload[0]
		if err := conn.WriteMessage(ws.BinaryMessage, []byte{CODE_ACK}); err != nil {
			return err
		}
	}
}

func buildTLSConfig(rootCAs *x509.CertPool) *tls.Config {
	return &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
	}
}
