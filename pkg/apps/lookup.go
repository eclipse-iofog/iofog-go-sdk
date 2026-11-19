package apps

import (
	"errors"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
)

// lookupErrorAllowNotFound returns nil when the resource is absent (NotFound).
// Any other lookup error is returned as-is.
func lookupErrorAllowNotFound(err error) error {
	if err == nil {
		return nil
	}
	var notFoundError *client.NotFoundError
	if errors.As(err, &notFoundError) {
		return nil
	}
	return err
}
