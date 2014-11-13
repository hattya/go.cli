//
// go.cli :: cli_unix.go
//
//   Copyright (c) 2014 Akinori Hattori <hattya@gmail.com>
//
//   Permission is hereby granted, free of charge, to any person
//   obtaining a copy of this software and associated documentation files
//   (the "Software"), to deal in the Software without restriction,
//   including without limitation the rights to use, copy, modify, merge,
//   publish, distribute, sublicense, and/or sell copies of the Software,
//   and to permit persons to whom the Software is furnished to do so,
//   subject to the following conditions:
//
//   The above copyright notice and this permission notice shall be
//   included in all copies or substantial portions of the Software.
//
//   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//   EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//   MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//   NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
//   BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
//   ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
//   CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//   SOFTWARE.
//

// +build !plan9,!windows

package cli

import (
	"bytes"
	"os"
	"regexp"
)

var xtermRx *regexp.Regexp

func init() {
	var b bytes.Buffer
	b.WriteString(`^(?:`)
	for i, s := range []string{
		"xterm",
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
	if xtermRx.MatchString(os.Getenv("TERM")) {
		ui.Printf("\x1b]2;%v\a", title)
	}
	return nil
}
