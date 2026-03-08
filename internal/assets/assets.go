package assets

import "embed"

// FS contains the generated in-repo game art.
//go:embed generated/*.png
var FS embed.FS
