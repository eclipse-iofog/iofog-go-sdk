package client

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"
)

type controllerStatus struct {
	version         string
	versionNoSuffix string
	versionNums     []string
}

type Client struct {
	baseURL      *url.URL
	accessToken  string
	refreshToken string
	retries      Retries
	status       controllerStatus
	timeout      int
}

type Options struct {
	BaseURL *url.URL
	Retries *Retries
	Timeout int
}

func New(opt Options) *Client {
	if opt.Timeout == 0 {
		opt.Timeout = 10
	}
	retries := GlobalRetriesPolicy
	if opt.Retries != nil {
		retries = *opt.Retries
	}
	client := &Client{
		retries: retries,
		baseURL: opt.BaseURL,
		timeout: opt.Timeout,
	}
	if client.baseURL.Scheme == "" {
		client.baseURL.Path = "http"
	}
	if client.baseURL.Path == "" {
		client.baseURL.Path = "api/v3"
	}
	// Get Controller version
	if status, err := client.GetStatus(); err == nil {
		versionNoSuffix := before(status.Versions.Controller, "-")
		versionNums := strings.Split(versionNoSuffix, ".")
		client.status = controllerStatus{
			version:         status.Versions.Controller,
			versionNoSuffix: versionNoSuffix,
			versionNums:     versionNums,
		}
	}
	return client
}

func NewAndLogin(opt Options, email, password string) (*Client, error) {
	clt := New(opt)
	if err := clt.Login(LoginRequest{Email: email, Password: password}); err != nil {
		return nil, err
	}
	return clt, nil
}

func SessionLogin(opt Options, token, email, password string) (clt *Client, err error) {
	clt = New(opt)

	// Attempt to login using the refresh token
	if err = clt.Refresh(RefreshTokenRequest{RefreshToken: token}); err != nil {
		// If session login fails, fall back to normal login with email and password
		clt, err = NewAndLogin(opt, email, password)
		if err != nil {
			return nil, fmt.Errorf("fallback login failed: %w", err)
		}
	}

	return clt, nil
}

func RefreshUserSubscriptionKey(opt Options, refreshToken, email, password string) (clt *Client, subscriptionKey string, err error) {
	// Attempt session login using the refresh token
	clt, err = SessionLogin(opt, refreshToken, email, password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to login: %w", err)
	}

	// Get the access token and fetch user profile
	accessToken := clt.GetAccessToken()
	userResponse, err := clt.Profile(WithTokenRequest{AccessToken: accessToken})
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch profile: %w", err)
	}

	return clt, userResponse.SubscriptionKey, nil
}

func NewWithToken(opt Options, token string) (*Client, error) {
	clt := New(opt)
	clt.SetAccessToken(token)
	return clt, nil
}

func NewWithRefreshToken(opt Options, refreshToken string) (*Client, error) {
	clt := New(opt)
	clt.SetAccessToken(refreshToken)
	return clt, nil
}

func (clt *Client) GetBaseURL() string {
	return clt.baseURL.String()
}

func (clt *Client) GetRetries() Retries {
	return clt.retries
}

func (clt *Client) SetRetries(retries Retries) {
	clt.retries = retries
}

func (clt *Client) GetAccessToken() string {
	return clt.accessToken
}

func (clt *Client) SetAccessToken(token string) {
	clt.accessToken = token
}

func (clt *Client) GetRefreshToken() string {
	return clt.refreshToken
}

func (clt *Client) SetRefreshToken(token string) {
	clt.refreshToken = token
}

func (clt *Client) doRequestWithRetries(currentRetries Retries, method, requestURL string, headers map[string]string, request any) ([]byte, error) {
	// Send request
	httpDo := httpDo{timeout: clt.timeout}
	bytes, err := httpDo.do(method, requestURL, headers, request)
	if err != nil {
		httpErr := &HTTPError{}
		ok := errors.As(err, &httpErr)
		// If HTTP Error
		if ok {
			if httpErr.Code == 408 { // HTTP Timeout
				if currentRetries.Timeout < clt.retries.Timeout {
					currentRetries.Timeout++
					time.Sleep(time.Duration(currentRetries.Timeout) * time.Second)
					return clt.doRequestWithRetries(currentRetries, method, requestURL, headers, request)
				}
				return bytes, err
			}
		}
		// If custom retries defined
		if clt.retries.CustomMessage != nil {
			for message, allowedRetries := range clt.retries.CustomMessage {
				if strings.Contains(err.Error(), message) {
					if currentRetries.CustomMessage[message] < allowedRetries {
						currentRetries.CustomMessage[message]++
						time.Sleep(time.Duration(currentRetries.CustomMessage[message]) * time.Second)
						return clt.doRequestWithRetries(currentRetries, method, requestURL, headers, request)
					}
					return bytes, err
				}
			}
		}
	}
	return bytes, err
}

func (clt *Client) doRequestWithHeaders(method, requestPath string, request any, headers map[string]string) ([]byte, error) {
	// Copy the base URL
	requestURL, err := url.Parse(clt.baseURL.String())
	if err != nil {
		return nil, err
	}
	// Get query params
	qpSplit := strings.Split(requestPath, "?")
	switch len(qpSplit) {
	case 1:
		requestURL.Path = path.Join(requestURL.Path, requestPath)
	case 2:
		requestURL.Path = path.Join(requestURL.Path, qpSplit[0])
		requestURL.RawQuery = qpSplit[1]
	default:
		return nil, fmt.Errorf("failed to parse request URL %s", requestPath)
	}

	// Set auth header
	headers["Authorization"] = "Bearer " + clt.accessToken

	currentRetries := Retries{CustomMessage: make(map[string]int)}
	if clt.retries.CustomMessage != nil {
		for message := range clt.retries.CustomMessage {
			currentRetries.CustomMessage[message] = 0
		}
	}

	return clt.doRequestWithRetries(currentRetries, method, requestURL.String(), headers, request)
}

func (clt *Client) doRequest(method, requestPath string, request any) ([]byte, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	return clt.doRequestWithHeaders(method, requestPath, request, headers)
}

func (clt *Client) isLoggedIn() bool {
	return clt.accessToken != ""
}
