package web

import "embed"

// Embedded holds the Vue production build (web/dist).
//
//go:embed all:dist
var Embedded embed.FS
