//
// go.cli :: cli_unix.go
//
//   Copyright (c) 2014-2021 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

// +build !plan9,!windows

package cli

import (
	"os"
	"regexp"
	"strings"

	"golang.org/x/term"
)

var xtermRx *regexp.Regexp

func init() {
	var b strings.Builder
	b.WriteString(`^(?:`)
	for i, s := range []string{
		"xterm",
		"putty",
		"rxvt",
		"screen",
	} {
		if i > 0 {
			b.WriteRune('|')
		}
		b.WriteString(regexp.QuoteMeta(s))
	}
	b.WriteString(`)`)
	xtermRx = regexp.MustCompile(b.String())
}

func (ui *CLI) title(title string) error {
	if f, ok := ui.Stdout.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		if xtermRx.MatchString(os.Getenv("TERM")) {
			ui.Printf("\x1b]2;%v\a", title)
		}
	}
	return nil
}
