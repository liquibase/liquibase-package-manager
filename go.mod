module package-manager

go 1.25.12

require (
	github.com/google/go-github/v55 v55.0.0
	github.com/hashicorp/go-version v1.9.0
	github.com/spf13/cobra v1.10.2
	github.com/vifraa/gopom v1.0.0
	golang.org/x/oauth2 v0.36.0
)

require (
	github.com/ProtonMail/go-crypto v0.0.0-20230217124315-7d5c6f04bbb8 // indirect
	github.com/cloudflare/circl v1.6.3 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/telemetry v0.0.0-20260708182218-49f421fb7959 // indirect
	golang.org/x/tools v0.48.0 // indirect
	golang.org/x/vuln v1.6.0 // indirect
)

// govulncheck is tracked as a tool dependency (not `go install …@version` in CI)
// so its version AND transitive build deps are pinned and checksum-verified via
// go.sum — a reproducible, lock-file-enforced scanner install. Required by the
// enterprise SHA-pinning policy and SonarCloud rule githubactions:S8545 ("Go
// dependencies should be locked to verified versions"). Run it with `go tool
// govulncheck`. Build-only: not compiled into the shipped lpm binary.
tool golang.org/x/vuln/cmd/govulncheck
