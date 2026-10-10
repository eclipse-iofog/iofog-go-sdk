# Fleet resources

You will register a fleet model, link it to a fog node, and bind it into a microservice. Knowledge uses the same shape. RuntimeClasses link to fog nodes and are selected by the microservice `runtime` field.

Import `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client`. Names are the identity. UUIDs are server-assigned. Item names in a microservice catalog are fleet names, never a UUID or a host path.

`pkg/apps` can unmarshal `kind: Model`, `kind: Knowledge`, and `kind: RuntimeClass`. It does not deploy them.

## Prerequisites

- Logged-in Controller client.
- A registry id. `type` is `oci` or `hf`. Microservice image and catalog `registryId` values must refer to an `oci` registry.
- For RuntimeClass links, the fog `containerEngine` must be `edgelet`. Linking any other engine returns HTTP 400.

## Steps

### 1. Create a model

JSON:

```go
model, err := clt.CreateModel(&client.ModelCreateRequest{
	Name:       "llama",
	Repo:       "library/llama",
	RegistryID: registryID,
	Revision:   "v1", // empty means latest (OCI) or main (Hugging Face)
	Format:     "gguf",
})
```

YAML file (`multipart` field name `model`):

```go
f, err := os.Open("model.yaml")
if err != nil {
	return err
}
defer f.Close()
model, err = clt.CreateModelFromYAML(f)
```

Example document:

```yaml
apiVersion: iofog.org/v3
kind: Model
metadata:
  name: llama
spec:
  repo: library/llama
  revision: v1
  registryId: 1
  format: gguf
```

`UpsertModelFromYAML(name, file)` is `PUT /models/yaml/{name}`. `UpdateModel` is `PATCH /models/{name}`. Name is immutable and comes from the path.

`GetModel` does not include linked fog UUIDs. Call `GetModelLink`.

### 2. Link the model to fog nodes

```go
err = clt.LinkModel("llama", &client.FogLinkRequest{
	FogUUIDs: []string{fogUUID},
})
```

The body is only `fogUuids`. The name is the URL. `UnlinkModel` is `DELETE` on the same path. Unlink returns HTTP 409 when a microservice on that fog still binds the name. `DeleteModel` returns HTTP 409 when any microservice still binds the name.

### 3. Bind the model into a microservice

Catalog-only patch (`PATCH /microservices/{uuid}/models`). This does not reload the full microservice. The workload rebuilds when the catalog goes empty to non-empty or the reverse, or when `bindPath` or `permissions` change. Adding or removing items with the same bind path and permissions does not rebuild while the catalog stays non-empty.

```go
err = clt.PatchMicroserviceModels(uuid, client.MicroserviceCatalog{
	BindPath:    "/models",
	Permissions: "ro", // ro or rw
	Items: []client.MicroserviceCatalogItem{
		{Name: "llama"},
	},
})
```

The same fields exist on YAML as `spec.models`.

### 4. Knowledge

Replace the model methods with the Knowledge pair:

| Model | Knowledge | YAML multipart field |
| --- | --- | --- |
| `ListModels` | `ListKnowledge` | |
| `GetModel` | `GetKnowledge` | |
| `CreateModel` | `CreateKnowledge` | |
| `UpdateModel` | `UpdateKnowledge` | |
| `DeleteModel` | `DeleteKnowledge` (success HTTP 202) | |
| `CreateModelFromYAML` | `CreateKnowledgeFromYAML` | `knowledge` |
| `UpsertModelFromYAML` | `UpsertKnowledgeFromYAML` | `knowledge` |
| `GetModelLink` / `LinkModel` / `UnlinkModel` | `GetKnowledgeLink` / `LinkKnowledge` / `UnlinkKnowledge` | |

`Files` is Hugging Face only. The Controller rejects non-empty `files` on an OCI registry. `Format` is a hint; unknown values are stored as `unknown`.

Bind with `PatchMicroserviceKnowledge` (`PATCH /microservices/{uuid}/knowledge`, body `KnowledgeCatalog`, success HTTP 204). YAML field is `spec.knowledge`. This patch applies to user microservices.

Fog `GET` adds `modelStatus`, `activeModels`, `modelLastUpdate`, `knowledgeStatus`, `activeKnowledge`, and `knowledgeLastUpdate`. The status fields are JSON strings. Parse them in your program. `modelLastUpdate` and `knowledgeLastUpdate` are Unix milliseconds, or `0` when the managed list is empty. `activeModels` and `activeKnowledge` count managed fleet items only.

### 5. RuntimeClass

```go
_, err = clt.CreateRuntimeClass(&client.RuntimeClassCreateRequest{
	Name:    "nvidia",
	Handler: "nvidia",
})
err = clt.LinkRuntimeClass("nvidia", &client.FogLinkRequest{FogUUIDs: []string{fogUUID}})
```

YAML multipart field is `runtimeClass`. The handler is a top-level field, not under `spec`:

```yaml
apiVersion: iofog.org/v3
kind: RuntimeClass
metadata:
  name: nvidia
handler: nvidia
```

`DeleteRuntimeClass` and `UnlinkRuntimeClass` return HTTP 409 when a microservice still pins `runtime` to that name.

Fog `GET` returns `runtimeClasses` and `availableCdiDevices` as JSON strings.

## Verification

`GetModelLink("llama")` contains the fog UUID after link. `GetMicroserviceByID` shows `Models` with that item name. After the agent reports status, `ActiveModels` on the fog is at least 1 and `ModelLastUpdate` is non-zero.

## Troubleshooting

| HTTP | Meaning |
| --- | --- |
| 409 on delete or unlink | A microservice still references the name. |
| 400 on `LinkRuntimeClass` | Fog container engine is not `edgelet`. |
| Empty `GetModel` fog list | Expected. Links are a separate resource. |

## Related

- [Client reference](../reference/client.md)
- [Apps reference](../reference/apps.md)
