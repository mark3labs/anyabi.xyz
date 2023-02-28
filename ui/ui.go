// Package ui handles the PocketBase Admin frontend embedding.
package ui

import (
	"embed"

	"github.com/labstack/echo/v5"
)

//go:embed all:build
var buildDir embed.FS

// BuildDirFS contains the embedded dist directory files (without the "build" prefix)
var BuildDirFS = echo.MustSubFS(buildDir, "build")
