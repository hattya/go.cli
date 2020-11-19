//
// go.cli :: cli_unix.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

// +build !plan9,!windows

package cli

import (
	"bytes"
	"os"
	"regexp"

	"golang.org/x/term"
)

var xtermRx *regexp.Regexp

func init() {
	var b bytes.Buffer
	b.WriteString(`^(?:`)
	for i, s := range []string{
		"xterm",
		"putty",
		"rxvt",
		"screen",
	} {
		if 0 < i {
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
