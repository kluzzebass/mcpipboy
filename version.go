// Package version provides version information for mcpipboy
package version

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var versionContent string

// Version is the current version of mcpipboy
var Version = strings.TrimSpace(versionContent)
