# Security Policy

## Reporting a vulnerability

If you believe you have found a security issue in the ioFog Go SDK:

1. **Do not** open a public GitHub issue for exploitable vulnerabilities.
2. Report through the Eclipse ioFog project security process or contact the maintainers privately with a description, impact, and reproduction steps.

For non-security bugs, use the public issue tracker.

## Security gates (maintainers)

Before release tags, run:

```bash
make security-code   # gosec on ./pkg/...
make vulncheck       # govulncheck + go mod verify
make lint            # golangci-lint v2 (gosec intentionally excluded)
```

- **gosec** runs via `make security-code` — **not** inside golangci-lint (edgelet pattern).
- **govulncheck** scans `./pkg/...` and verifies module integrity.

## Documented gosec exceptions

| Rule | Location | Rationale |
|------|----------|-----------|
| G402 | `pkg/client/http.go` | Controller deployments commonly use self-signed TLS |
| G101 | `pkg/microservices/declarations.go` | Well-known service-account mount path constant, not a credential |
| G115 | `pkg/microservices/util.go` | Minimal big-endian encoding of non-negative protocol integers |

Undocumented `#nosec` findings **fail** `make security-code`.
