// Package web provides the web UI and API server for PicoClaw.
package web

import (
	"embed"
)

//go:embed all:dist
var DistFS embed.FS
