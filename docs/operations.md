# Operations

## Module

```text
github.com/eclipse-iofog/iofog-go-sdk/v3
```

Go **1.26.10** or newer. The Datasance GitHub repository publishes this same module path. Do not change imports to a Datasance module path.

`go.mod` retracts a set of accidental and test tags (`v3.3.2` through parts of `v3.4`, and `v3.8.0-beta.0` / `beta.1`). Do not depend on those versions.

Current API level: Controller REST v3.9 and EdgeletAPI v1. History is [`CHANGELOG.md`](../CHANGELOG.md).

## Make targets

| Target | What it does |
| --- | --- |
| `make test` | `make gen`, `gofmt`, then `go test` on each package under `pkg/` (parallelism 1). Writes `reports/`. |
| `make gen` | Regenerates `pkg/apps/deepcopy_generated.go`. Uses BSD `sed` on Darwin and GNU `sed` otherwise. |
| `make gen-check` | Fails if generated deepcopy or `types.go` drift. |
| `make lint` | golangci-lint **v2.12.2** on `./pkg/...` (`.golangci.yaml`). |
| `make lint-fix` | Same, with `--fix`. |
| `make security-code` | gosec **v2.22.2** on `./pkg/...`, generated files excluded. |
| `make vulncheck` | govulncheck **v1.1.4** on `./pkg/...`, then `go mod verify`. |
| `make quality` | `lint`, `security-code`, `vulncheck`. |
| `make fmt` | `gofmt -s -w .` |

Pinned tool versions can be overridden with `GOLANGCI_LINT_VERSION`, `GOSEC_VERSION`, `GOVULNCHECK_VERSION`, and `DEEPCOPY_GEN_VERSION` (default `v0.32.0`, aligned with `k8s.io/apimachinery` v0.32.1).

CI (`.github/workflows/ci.yml`) runs on `develop`, `master`, `main`, `release/**`, and pull requests. The job is lint, `make security-code`, and `make vulncheck`. It does not run the integration scripts under `test/`.

## Documented gosec exceptions

| Rule | File | Why |
| --- | --- | --- |
| G402 | `pkg/client/http.go` | Default TLS config skips verify for private Controller certificates. Callers pass `Options.TLSConfig` in production. |
| G101 | `pkg/microservices/declarations.go` | The service-account mount path is a well-known path, not a secret. |

Any other `#nosec` fails `make security-code`.

## Integration scripts

`test/` contains shell helpers (`bootstrap.bash`, `setup.bash`, `run.bash`, `teardown.bash`) and sample YAML. They are not part of CI. Unit tests under `pkg/` use `httptest` and do not need a Controller.

## License

Eclipse Public License 2.0. See [`LICENSE`](../LICENSE), [`NOTICE`](../NOTICE), and [`SECURITY.md`](../SECURITY.md). Contributions require the Eclipse Contributor Agreement described in [`CONTRIBUTING`](../CONTRIBUTING).

## Related

- [Documentation index](README.md)
- [Packages](explanation/packages.md)
