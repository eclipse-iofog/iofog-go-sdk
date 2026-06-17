# Client Package

This package provides an HTTP client for ioFog Controller's REST API **v3** (`/api/v3/*`).

Import path: `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client`

You can view the full REST API specification at [iofog.org](https://iofog.org/docs/1.3.0/controllers/rest-api.html).

## Usage

Create a client with the Controller base URL (including `/api/v3`), then log in with your credentials.

```go
import (
	"net/url"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
)

baseURL, err := url.Parse("http://localhost:51121/api/v3")
if err != nil {
	return err
}

clt, err := client.NewAndLogin(client.Options{BaseURL: baseURL}, "user@domain.com", "password")
if err != nil {
	return err
}
```

Call any method on the authenticated client:

```go
resp, err := clt.GetStatus()
if err != nil {
	return err
}
println(resp.Status)
```
