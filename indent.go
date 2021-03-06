//
// go.cli :: indent.go
//
//   Copyright (c) 2016-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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
