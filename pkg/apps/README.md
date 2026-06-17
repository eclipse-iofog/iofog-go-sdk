# Apps Package

This package contains executors to deploy ioFog applications and microservices using the `client` package. It is used by `iofogctl` and `iofog-operator` to apply YAML configuration to Controller.

Import path: `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/apps`

## apiVersion

Deploy YAML emitted by this package includes an `apiVersion` header. The default is **`iofog.org/v3`** (`apps.DefaultAPIVersion`).

For Datasance-flavored deploy YAML, pass `WithAPIVersion` at the call site:

```go
apps.DeployApplication(ctrl, app, name, apps.WithAPIVersion("datasance.com/v3"))
```

## Usage

```go
import (
	"os"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/apps"
	"gopkg.in/yaml.v2"
)

controller := apps.IofogController{
	Endpoint: "127.0.0.1:51121",
	Email:    "user@domain.com",
	Password: "password",
}

application := apps.Application{
	Name: "my-app",
	// ... remaining fields
}

// OR, read spec from a YAML file
yamlFile, err := os.ReadFile(filename)
if err != nil {
	return err
}
var header apps.IofogHeader
if err := yaml.Unmarshal(yamlFile, &header); err != nil {
	return err
}

// Deploy (default apiVersion: iofog.org/v3)
err = apps.DeployApplication(controller, application, "my-app")
```

Other entry points: `DeployMicroservice` and `DeployApplicationTemplate`. Each accepts optional `DeployOption` values (for example `WithAPIVersion`).
