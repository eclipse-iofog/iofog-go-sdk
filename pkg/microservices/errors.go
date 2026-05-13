package microservices

import "fmt"

// AuthMaterialError indicates a token/CA read or parse failure.
type AuthMaterialError struct {
	Kind string
	Path string
	Err  error
}

func (e *AuthMaterialError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s material error at %s: %v", e.Kind, e.Path, e.Err)
}

func (e *AuthMaterialError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// V3APIError describes a structured LocalAPI v3 API error envelope.
type V3APIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    map[string]interface{}
}

func (e *V3APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code == "" {
		return fmt.Sprintf("localapi request failed (%d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("localapi %s (%d): %s", e.Code, e.StatusCode, e.Message)
}
