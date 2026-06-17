package client

// ListRoutes returns all routes. Deprecated: Controller no longer exposes /api/v3/routes.
// Use NATS for messaging. ListRoutes returns ErrRoutesNotSupported.
func (clt *Client) ListRoutes() (response RouteListResponse, err error) {
	return RouteListResponse{}, ErrRoutesNotSupported
}

// GetRoute returns a single route by application and name.
// Deprecated: Controller no longer exposes /api/v3/routes. Use NATS for messaging.
// GetRoute returns ErrRoutesNotSupported.
func (clt *Client) GetRoute(appName, name string) (route Route, err error) {
	return Route{}, ErrRoutesNotSupported
}

// CreateRoute creates a route. Deprecated: Controller no longer exposes /api/v3/routes.
// Use NATS for messaging. CreateRoute returns ErrRoutesNotSupported.
func (clt *Client) CreateRoute(route *Route) (err error) {
	return ErrRoutesNotSupported
}

// UpdateRoute creates or updates a route. Deprecated: Controller no longer exposes /api/v3/routes.
// Use NATS for messaging. UpdateRoute returns ErrRoutesNotSupported.
func (clt *Client) UpdateRoute(route *Route) (err error) {
	return ErrRoutesNotSupported
}

// PatchRoute updates a route by application and name.
// Deprecated: Controller no longer exposes /api/v3/routes. Use NATS for messaging.
// PatchRoute returns ErrRoutesNotSupported.
func (clt *Client) PatchRoute(appName, name string, route *Route) (err error) {
	return ErrRoutesNotSupported
}

// DeleteRoute deletes a route. Deprecated: Controller no longer exposes /api/v3/routes.
// Use NATS for messaging. DeleteRoute returns ErrRoutesNotSupported.
func (clt *Client) DeleteRoute(appName, name string) (err error) {
	return ErrRoutesNotSupported
}
