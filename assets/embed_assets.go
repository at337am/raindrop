package assets

import _ "embed"

//go:embed help_en.txt
var HelpEN string

//go:embed help_zh.txt
var HelpZH string

//go:embed tty_logo.txt
var TTYLogo string
