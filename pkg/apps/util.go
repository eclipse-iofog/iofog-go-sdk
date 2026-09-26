package apps

import (
	"fmt"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
)

// validateRoutes validates route names against existing microservices.
// Deprecated: Controller routing is deprecated; use NATS for messaging.
func validateRoutes(routes []string, microserviceByName map[string]*client.MicroserviceInfo) (routesUUIDs []string, err error) { // nolint:deadcode,unused
	for _, route := range routes {
		msvc, foundTo := microserviceByName[route]
		if !foundTo {
			return routesUUIDs, NewNotFoundError(fmt.Sprintf("Could not find microservice [%s] required by a route", route))
		}
		routesUUIDs = append(routesUUIDs, msvc.UUID)
	}
	return routesUUIDs, nil
}

// createRoutes creates microservice routes via the Controller API.
// Deprecated: Controller no longer exposes route endpoints; use NATS for messaging.
// createRoutes returns client.ErrRoutesNotSupported.
func createRoutes(routes []Route, microserviceByName map[string]*client.MicroserviceInfo, clt *client.Client) error { // nolint:deadcode,unused
	return client.ErrRoutesNotSupported
}
