// Package version exposes build metadata for the knowledgeKeeperApi service.
//
// Version is declared as a var (not a const) so it can be overridden at build
// time via linker flags, which is the idiomatic way to stamp a binary from
// CI/CD. The default value mirrors the contents of the repo-root VERSION file
// so local builds (`go build`/`go run`) still report a sensible value.
//
// Example CI override:
//
//	VERSION=$(cat VERSION)
//	go build -ldflags "-X github.com/Namularbre/knowledgeKeeperApi/internal/version.Version=${VERSION}" ./cmd/api
package version

// Version is the semantic version of the API. Overridable via -ldflags.
var Version = "1.0.0"

// APIName is the public name of the service.
const APIName = "knowledgeKeeperApi"
