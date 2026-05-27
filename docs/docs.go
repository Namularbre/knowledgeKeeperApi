// Package docs is the target of `swag init` output. The file you are reading
// is a placeholder that lets the project build before the first generation.
//
// Run `./scripts/swag.sh` (or `swag init -g cmd/api/main.go -o docs`) to
// overwrite this file with the real generated docs.go, swagger.json and
// swagger.yaml.
package docs

import "github.com/swaggo/swag"

// init registers an empty spec so /swagger endpoints don't 500 before
// generation. After `swag init` runs, the generated file replaces this one
// and registers the real spec.
func init() {
	swag.Register(swag.Name, &emptySpec{})
}

type emptySpec struct{}

func (e *emptySpec) ReadDoc() string {
	return `{"swagger":"2.0","info":{"title":"knowledgeKeeperApi","version":"0.0.0","description":"Run scripts/swag.sh to generate the real spec."},"paths":{}}`
}
