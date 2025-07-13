package assets

import (
	"embed"
)

//go:embed tty/logo.txt
var Logo string

//go:embed all:public all:templates
var FS embed.FS
