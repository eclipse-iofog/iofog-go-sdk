package container_sdk_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"errors"
)

type ioFogHttpClient struct {
	ssl                         bool
	host                        string
	port                        int
	url_base_rest               string
	url_get_config              string
	url_get_next_messages       string
	url_get_publishers_messages string
	url_post_message            string
	requestBodyId               []byte
}

func newIoFogHttpClient(id string, ssl bool, host string, port int) *ioFogHttpClient {
	client := ioFogHttpClient {
		host:host,
		port:port,
		ssl:ssl,
	}
	protocol_rest := HTTP
	if client.ssl {
		protocol_rest = HTTPS
	}
	client.url_base_rest = fmt.Sprintf("%s://%s:%d", protocol_rest, client.host, client.port)
	client.url_get_config = fmt.Sprint(client.url_base_rest, URL_GET_CONFIG)
	client.url_get_next_messages = fmt.Sprint(client.url_base_rest, URL_GET_NEXT_MESSAGES)
	client.url_get_publishers_messages = fmt.Sprint(client.url_base_rest, URL_GET_PUBLISHERS_MESSAGES)
	client.url_post_message = fmt.Sprint(client.url_base_rest, URL_POST_MESSAGE)
	client.requestBodyId, _ = json.Marshal(map[string]interface{}{
		ID : id,
	})
	return &client
}

func (client *ioFogHttpClient) getConfig() (map[string]interface{}, error) {
	resp, err := makePostRequest(client.url_get_config, APPLICATION_JSON, bytes.NewBuffer(client.requestBodyId))
	if err != nil {
		return nil, err
	}
	configResponse := new(getConfigResponse)
	config := make(map[string]interface{})
	if err := json.Unmarshal(resp, configResponse); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(configResponse.Config), &config); err != nil {
		return nil, err
	}
	return config, nil
}

func (client *ioFogHttpClient) getNextMessages() ([]IoMessage, error) {
	resp, err := makePostRequest(client.url_get_next_messages, APPLICATION_JSON, bytes.NewBuffer(client.requestBodyId))
	if err != nil {
		return nil, err
	}
	nextMessagesResponse := new(getNextMessagesResponse)
	if err := json.Unmarshal(resp, nextMessagesResponse); err != nil {
		return nil, err
	}
	return nextMessagesResponse.Messages, nil
}

func (client *ioFogHttpClient) postMessage(msg *IoMessage) (*PostMessageResponse, error) {
	requestBytes, _ := json.Marshal(msg)
	resp, err := makePostRequest(client.url_post_message, APPLICATION_JSON, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	postMessageResponse := new(PostMessageResponse)
	if err := json.Unmarshal(resp, postMessageResponse); err != nil {
		return nil, err
	}
	return postMessageResponse, nil
}

func (client *ioFogHttpClient) getMessagesFromPublishersWithinTimeFrame(query *MessagesQueryParameters) (*TimeFrameMessages, error) {
	requestBytes, _ := json.Marshal(query)
	resp, err := makePostRequest(client.url_get_publishers_messages, APPLICATION_JSON, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	nextMessagesResponse := new(getNextMessagesResponse)
	if err := json.Unmarshal(resp, nextMessagesResponse); err != nil {
		return nil, err
	}
	return (*TimeFrameMessages)(nextMessagesResponse), nil
}

func makePostRequest(url, bodyType string, body io.Reader) ([]byte, error) {
	resp, err := http.Post(url, bodyType, body)
	if err != nil {
		return nil, err
	}
	respBodyBytes := make([]byte, resp.ContentLength)
	resp.Body.Read(respBodyBytes)
	resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New(string(respBodyBytes))
	}
	return respBodyBytes, nil
}
