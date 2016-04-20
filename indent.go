//
// go.cli :: indent.go
//
//   Copyright (c) 2016 Akinori Hattori <hattya@gmail.com>
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

package cli

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"
)

func Dedent(s string) string {
	for _, newline := range []string{"\n", "\r\n"} {
		if strings.HasPrefix(s, newline) {
			s = s[len(newline):]
			break
		}
	}
	r := bufio.NewReader(strings.NewReader(s))

	isNewline := func(s string) bool {
		return s == "\n" || s == "\r\n"
	}

	var out []string
	var mgn string
	var eof bool
	n := 0
	for !eof {
		l, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if l == "" {
					break
				}
				eof = true
			} else {
				panic(err)
			}
		}

		i := strings.IndexFunc(l, func(r rune) bool {
			return r != '\t' && r != ' '
		})
		switch {
		case i == -1:
			continue
		case isNewline(l[i:]):
			out = append(out, l[i:])
			continue
		default:
			out = append(out, l)
			n++
		}

		ind := l[:i]
		switch {
		case n == 1:
			mgn = ind
		case !strings.HasPrefix(ind, mgn):
			i := 0
			for {
				r, w := utf8.DecodeRuneInString(ind[i:])
				if r == utf8.RuneError {
					if w == 1 {
						panic("invalid UTF-8")
					}
					break
				}
				if !strings.HasPrefix(mgn, ind[:i+w]) {
					break
				}
				i += w
			}
			mgn = mgn[:i]
		}
	}

	if mgn != "" {
		for i := range out {
			if !isNewline(out[i]) {
				out[i] = out[i][len(mgn):]
			}
		}
	}
	return strings.Join(out, "")
}
