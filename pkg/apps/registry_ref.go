package apps

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
)

// RegistryRef is a registry id for YAML images.registry and catalog registry fields.
// Unmarshal accepts int, numeric string, or aliases remote (1) and local (2).
// Marshal emits an unquoted integer.
type RegistryRef int

const defaultRegistryRef RegistryRef = 1

// Int returns the numeric registry id. Zero resolves to the default remote registry (1).
func (r RegistryRef) Int() int {
	if r == 0 {
		return int(defaultRegistryRef)
	}
	return int(r)
}

func parseRegistryRef(v any) (RegistryRef, error) {
	switch x := v.(type) {
	case int:
		if x == 0 {
			return defaultRegistryRef, nil
		}
		return RegistryRef(x), nil
	case int64:
		if x == 0 {
			return defaultRegistryRef, nil
		}
		return RegistryRef(x), nil
	case float64:
		if x == 0 {
			return defaultRegistryRef, nil
		}
		return RegistryRef(int(x)), nil
	case string:
		if x == "" {
			return defaultRegistryRef, nil
		}
		if id, ok := client.RegistryTypeRegistryTypeIDDict[x]; ok {
			return RegistryRef(id), nil
		}
		n, err := strconv.Atoi(x)
		if err != nil {
			return 0, fmt.Errorf("invalid registry %q: expected integer or remote|local", x)
		}
		if n == 0 {
			return defaultRegistryRef, nil
		}
		return RegistryRef(n), nil
	default:
		return 0, fmt.Errorf("invalid registry type %T", v)
	}
}

// UnmarshalYAML implements yaml.v2 unmarshaling for RegistryRef.
func (r *RegistryRef) UnmarshalYAML(unmarshal func(any) error) error {
	var raw any
	if err := unmarshal(&raw); err != nil {
		return err
	}
	if raw == nil {
		*r = defaultRegistryRef
		return nil
	}
	parsed, err := parseRegistryRef(raw)
	if err != nil {
		return err
	}
	*r = parsed
	return nil
}

// MarshalYAML implements yaml.v2 marshaling for RegistryRef.
func (r RegistryRef) MarshalYAML() (any, error) {
	return r.Int(), nil
}

// UnmarshalJSON implements JSON unmarshaling for RegistryRef.
func (r *RegistryRef) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*r = 0
		return nil
	}
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed, err := parseRegistryRef(raw)
	if err != nil {
		return err
	}
	*r = parsed
	return nil
}

// MarshalJSON implements JSON marshaling for RegistryRef.
func (r RegistryRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Int())
}
