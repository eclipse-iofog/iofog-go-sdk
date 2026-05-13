package microservices

import (
	"encoding/json"
	"fmt"
)

type v3Envelope struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Error   struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details"`
	} `json:"error"`
}

func parseV3Envelope(body []byte, statusCode int) (map[string]interface{}, error) {
	var env v3Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		if statusCode >= 200 && statusCode < 300 {
			// Fallback for non-enveloped payloads.
			var payload map[string]interface{}
			if errMap := json.Unmarshal(body, &payload); errMap == nil {
				return payload, nil
			}
		}
		return nil, fmt.Errorf("failed to decode response envelope: %w", err)
	}
	if statusCode < 200 || statusCode >= 300 || !env.Success {
		return nil, &V3APIError{
			StatusCode: statusCode,
			Code:       env.Error.Code,
			Message:    env.Error.Message,
			Details:    env.Error.Details,
		}
	}
	if env.Data == nil {
		return map[string]interface{}{}, nil
	}
	return env.Data, nil
}
