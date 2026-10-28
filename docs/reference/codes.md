# Architecture codes and registry ids

## `pkg/arch`

Import `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/arch`.

These integers are the Controller architectures table, `archId` on agents and catalog images, and the Edgelet architecture code.

| Name | ID | Constant |
| --- | --- | --- |
| `auto` | 0 | `arch.Auto` |
| `amd64` | 1 | `arch.AMD64` |
| `arm64` | 2 | `arch.ARM64` |
| `riscv64` | 3 | `arch.RISCV64` |
| `arm` | 4 | `arch.ARM` |

`arch.NameToID` and `arch.IDToName` are the two maps. `arch.DeployArchIDs` is `AMD64`, `ARM64`, `RISCV64`, `ARM` (`auto` is excluded).

`client.ArchNameToID` and `client.ArchIDToName` are the same maps.

YAML image keys are `amd64`, `arm64`, `riscv64`, and `arm`. The old `x86` key is not decoded.

`client.GetArchitectures` returns the Controller's `/architectures/` list (`Architecture` objects). Use that list when you need the server's names. Use `pkg/arch` when you need the stable integer.

## Registry id aliases

These aliases are for YAML `registry` fields. They are not the registry `type` (`oci` or `hf`).

| Alias | ID |
| --- | --- |
| `remote` | 1 |
| `local` | 2 |

Maps: `client.RegistryTypeRegistryTypeIDDict`, `client.RegistryTypeIDRegistryTypeDict`.

`apps.RegistryRef` accepts the alias, a number, or a numeric string. Empty and `0` marshal as `1`.

## Ports

| Constant | Value |
| --- | --- |
| `client.ControllerPort` | `51121` |
| `client.ControllerPortString` | `"51121"` |
| `microservices.PortEdgeletAPI` | `54321` |

## Related

- [Apps](apps.md)
- [Client](client.md)
