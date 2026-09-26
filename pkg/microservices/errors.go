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

// EdgeletAPIError describes a structured EdgeletAPI v1 error envelope.
type EdgeletAPIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    map[string]any
}

func (e *EdgeletAPIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code == "" {
		return fmt.Sprintf("edgeletapi request failed (%d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("edgeletapi %s (%d): %s", e.Code, e.StatusCode, e.Message)
}
