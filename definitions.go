package container_sdk_go

import (
	"time"
)

const IOFOG = "iofog"
const PORT_IOFOG = 54321
const SELFNAME = "SELFNAME"
const SSL = "SSL"
const SSL_DEFAULT = false
const HOST_DEFAULT = "127.0.0.1"

const URL_GET_CONFIG = "/v2/config/get"
const URL_GET_NEXT_MESSAGES = "/v2/messages/next"
const URL_GET_PUBLISHERS_MESSAGES = "/v2/messages/query"
const URL_POST_MESSAGE = "/v2/messages/new"
const URL_GET_CONTROL_WS = "/v2/control/socket/id/"
const URL_GET_MESSAGE_WS = "/v2/message/socket/id/"

const APPLICATION_JSON = "application/json"
const HTTP = "http"
const HTTPS = "https"
const WS = "ws"
const WSS = "wss"

const CODE_ACK = 0xB;
const CODE_CONTROL_SIGNAL = 0xC;
const CODE_MSG = 0xD;
const CODE_RECEIPT = 0xE;

const wsAttemptLimit = 5;
const wsConnectTimeout = time.Second;

const ID = "id"

type getConfigResponse struct {
	Config string        `json:"config"`
}

type PostMessageResponse struct {
	ID        string      `json:"id"`
	Timestamp int         `json:"timestamp"`
}

type MessagesQueryParameters struct {
	ID             string             `json:"id"`
	TimeFrameStart int                `json:"timeframestart"`
	TimeFrameEnd   int                `json:"timeframeend"`
	Publishers     []string           `json:"publishers"`
}

type getNextMessagesResponse struct {
	TimeFrameStart int                `json:"timeframestart"`
	TimeFrameEnd   int                `json:"timeframeend"`
	Messages       []IoMessage        `json:"messages"`
}

type TimeFrameMessages getNextMessagesResponse