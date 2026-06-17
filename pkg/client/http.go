package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	json "github.com/json-iterator/go"
)

type httpDo struct {
	timeout int
}

func (hd *httpDo) do(method, url string, headers map[string]string, requestBody any) ([]byte, error) {
	body, isIoReader := requestBody.(io.Reader)
	encodeType, ok := headers["Content-Type"]
	if ok && encodeType == "application/json" {
		if !isIoReader {
			jsonBody := ""
			if requestBody != nil {
				jsonBodyBytes, err := json.Marshal(requestBody)
				if err != nil {
					return nil, err
				}
				jsonBody = string(jsonBodyBytes)
			}

			Verbose(fmt.Sprintf("===> [%s] %s \nBody: %s\n", method, url, jsonBody))
			body = strings.NewReader(jsonBody)
		}
	} else {
		if !isIoReader {
			return nil, NewInternalError("Failed to convert request body to io.Reader")
		}
		Verbose(fmt.Sprintf("===> [%s] %s \nContent-Type: %s\n", method, url, encodeType))
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	request.Close = true

	for key, val := range headers {
		request.Header.Set(key, val)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402 -- Controller deployments commonly use self-signed TLS
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(hd.timeout),
	}

	httpResp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err = checkStatusCode(httpResp.StatusCode, method, url, httpResp.Body); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(httpResp.Body); err != nil {
		return nil, err
	}
	responseBody := buf.Bytes()
	Verbose(fmt.Sprintf("===> Response: %s\n\n", string(responseBody)))
	return responseBody, nil
}
