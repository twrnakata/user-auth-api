package swagger

import "embed"

//go:embed index.html openapi.yaml
var Files embed.FS
