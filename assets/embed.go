package assets

import (
	"embed"
)

//go:embed tty/help_en.txt
var HelpEN string

//go:embed tty/help_zh.txt
var HelpZH string

//go:embed tty/logo.txt
var Logo string

//go:embed all:public all:templates
var FS embed.FS
