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

	isEmpty := func(s string) bool {
		return s == "" || s == "\r"
	}

	out := strings.Split(s, "\n")
	var mgn string
	n := 0
	for i, l := range out {
		j := strings.IndexFunc(l, func(r rune) bool {
			return r != '\t' && r != ' '
		})
		switch {
		case j == -1:
			out[i] = ""
			continue
		case isEmpty(l[j:]):
			out[i] = l[j:]
			continue
		}

		ind := l[:j]
		switch {
		case n == 0:
			mgn = ind
		case !strings.HasPrefix(ind, mgn):
			j := 0
			for {
				r, w := utf8.DecodeRuneInString(ind[j:])
				if r == utf8.RuneError {
					if w == 1 {
						panic("invalid UTF-8")
					}
					break
				}
				if !strings.HasPrefix(mgn, ind[:j+w]) {
					break
				}
				j += w
			}
			mgn = mgn[:j]
		}
		n++
	}

	if mgn != "" {
		for i, l := range out {
			if !isEmpty(l) {
				out[i] = l[len(mgn):]
			}
		}
	}
	return strings.Join(out, "\n")
}
