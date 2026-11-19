# Deploy YAML

You will create or update an application on the Controller and start it. The executor is `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/apps`.

## Prerequisites

- A Controller base URL that `url.Parse` accepts, including scheme and `/api/v3`.
- Either a bearer `Token`, or `RefreshToken` plus `Email` and `Password`.
- An application spec. You can pass an `apps.Application` value or any value that YAML-marshals as `spec`.

`DeployApplication` and `DeployMicroservice` read credentials from `apps.IofogController` and parse `Endpoint` as the base URL.

`DeployApplicationTemplate` and `DeployMicroserviceTemplate` take that base URL as a separate `*url.URL` argument. `IofogController.Endpoint` is not used for those two calls.

## Steps

### 1. Describe the Controller

```go
ctrl := apps.IofogController{
	Endpoint: "https://controller.example:51121/api/v3",
	Token:    accessToken,
}
```

### 2. Deploy an application

```go
app := apps.Application{
	Name: "sensor-app",
	Microservices: []apps.Microservice{
		{
			Name: "reader",
			Agent: apps.MicroserviceAgent{Name: "edge-1"},
			Images: &apps.MicroserviceImages{
				AMD64: "example/reader:1.0",
			},
		},
	},
}

err := apps.DeployApplication(ctrl, app, "sensor-app")
```

What the executor does:

1. Log in (`NewWithToken` when `Token` is set, otherwise `SessionLogin`).
2. `GetApplicationByName`. A not-found error means create. Any other lookup error stops the deploy.
3. Marshal `apiVersion`, `kind: Application`, `metadata.name`, and `spec`, then `POST /application/yaml` or `PUT /application/yaml/{name}`.
4. `StartApplication` (`PATCH` with `isActivated: true`).

The default header is `apiVersion: iofog.org/v3`. For a Datasance Controller header:

```go
err = apps.DeployApplication(ctrl, app, "sensor-app", apps.WithAPIVersion("datasance.com/v3"))
```

### 3. Deploy one microservice

`appName` is required. The metadata name sent to the Controller is `appName/name`.

```go
ms := apps.Microservice{
	Name:  "reader",
	Agent: apps.MicroserviceAgent{Name: "edge-1"},
	Images: &apps.MicroserviceImages{
		AMD64: "example/reader:1.0",
	},
}
err = apps.DeployMicroservice(ctrl, ms, "sensor-app", "reader")
```

If the application is a system application, the executor updates an existing microservice and returns an error on create (`cannot create system microservice`).

### 4. Deploy templates

```go
base, err := url.Parse(ctrl.Endpoint)
if err != nil {
	return err
}
err = apps.DeployApplicationTemplate(ctrl, base, templateSpec, "sensor-template")
err = apps.DeployMicroserviceTemplate(ctrl, base, msTemplateSpec, "reader-template")
```

Template deploy creates when missing and updates when present. It does not start an application.

## Kinds this package deploys

| Kind constant | Deploy function |
| --- | --- |
| `ApplicationKind` (`Application`) | `DeployApplication` |
| `MicroserviceKind` (`Microservice`) | `DeployMicroservice` |
| `ApplicationTemplateKind` | `DeployApplicationTemplate` |
| `MicroserviceTemplateKind` | `DeployMicroserviceTemplate` |

`ModelKind`, `KnowledgeKind`, and `RuntimeClassKind` unmarshal in this package. Upload them with the client. See [Fleet resources](fleet-resources.md).

## Verification

`GetApplicationByName` on a `client.Client` returns `IsActivated: true` after `DeployApplication` succeeds. The new or updated microservice appears on `GetMicroservicesByApplication`.

## Troubleshooting

| Symptom | Cause |
| --- | --- |
| `failed to parse Controller endpoint as URL` | `Endpoint` is not an absolute URL. |
| Prompt for email or OTP | `Token` was empty, so `SessionLogin` fell through to interactive `Login`. |
| `Application name missing` | `DeployMicroservice` was called with an empty `appName`. |
| `cannot create system microservice` | The name matched a system application and the microservice did not already exist. |
| Registry field rejected | `images.registry` must be an integer, a numeric string, or `remote` (1) / `local` (2). Zero and empty become `1`. |

`apps.RegistryRef` marshals as an unquoted integer.

## Related

- [Apps reference](../reference/apps.md)
- [Client reference](../reference/client.md)
