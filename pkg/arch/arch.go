// Package arch defines architecture integer codes shared with Controller v3.8
// (Architectures table, archId on agents and catalog images) and Edgelet
// (ArchitectureCode in internal/config/arch.go).
package arch

const (
	Auto    = 0
	AMD64   = 1
	ARM64   = 2
	RISCV64 = 3
	ARM     = 4
)

// NameToID maps canonical architecture names to Controller/Edgelet integer codes.
var NameToID = map[string]int{
	"auto":    Auto,
	"amd64":   AMD64,
	"arm64":   ARM64,
	"riscv64": RISCV64,
	"arm":     ARM,
}

// IDToName maps Controller/Edgelet integer codes to canonical architecture names.
var IDToName = map[int]string{
	Auto:    "auto",
	AMD64:   "amd64",
	ARM64:   "arm64",
	RISCV64: "riscv64",
	ARM:     "arm",
}

// DeployArchIDs lists architecture IDs valid for deploy/catalog images (excludes auto).
var DeployArchIDs = []int{AMD64, ARM64, RISCV64, ARM}
