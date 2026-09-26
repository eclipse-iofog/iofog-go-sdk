package microservices

import (
	"encoding/json"
	"fmt"
)

type edgeletAPIErrorPayload struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

type edgeletAPIEnvelope struct {
	Success bool                   `json:"success"`
	Data    map[string]any         `json:"data"`
	Error   edgeletAPIErrorPayload `json:"error"`
}

func parseEdgeletAPIEnvelope(body []byte, statusCode int) (map[string]any, error) {
	var env edgeletAPIEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		if statusCode >= 200 && statusCode < 300 {
			// Fallback for non-enveloped payloads.
			var payload map[string]any
			if errMap := json.Unmarshal(body, &payload); errMap == nil {
				return payload, nil
			}
		}
		return nil, fmt.Errorf("failed to decode response envelope: %w", err)
	}
	if statusCode < 200 || statusCode >= 300 || !env.Success {
		return nil, &EdgeletAPIError{
			StatusCode: statusCode,
			Code:       env.Error.Code,
			Message:    env.Error.Message,
			Details:    env.Error.Details,
		}
	}
	if env.Data == nil {
		return map[string]any{}, nil
	}
	return env.Data, nil
}
