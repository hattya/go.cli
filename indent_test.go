//
// go.cli :: indent_test.go
//
//   Copyright (c) 2016-2017 Akinori Hattori <hattya@gmail.com>
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

package cli_test

import (
	"testing"

	"github.com/hattya/go.cli"
)

var dedentTests = []struct {
	in, out string
}{
	// no margin
	{
		in:  "1\n2\n3",
		out: "1\n2\n3",
	},
	{
		in:  "1\n 2\n  3",
		out: "1\n 2\n  3",
	},
	{
		in:  "1\n\t2\n\t\t3",
		out: "1\n\t2\n\t\t3",
	},
	{
		in:  "1\n\n3\n4",
		out: "1\n\n3\n4",
	},
	{
		in:  "1\n \n  3\n   4",
		out: "1\n\n  3\n   4",
	},
	{
		in:  "1\n\t\n\t\t3\n\t\t\t4",
		out: "1\n\n\t\t3\n\t\t\t4",
	},
	// indent >>>
	{
		in:  " 1\n  2\n   3",
		out: "1\n 2\n  3",
	},
	{
		in:  "\t1\n\t\t2\n\t\t\t3",
		out: "1\n\t2\n\t\t3",
	},
	{
		in:  " 1\n\n   3\n    4",
		out: "1\n\n  3\n   4",
	},
	{
		in:  "\t1\n\n\t\t\t3\n\t\t\t\t4",
		out: "1\n\n\t\t3\n\t\t\t4",
	},
	{
		in:  " 1\n  \n   3\n    4",
		out: "1\n\n  3\n   4",
	},
	{
		in:  "\t1\n\t\t\n\t\t\t3\n\t\t\t\t4",
		out: "1\n\n\t\t3\n\t\t\t4",
	},
	// indent <<<
	{
		in:  "   3\n  2\n 1",
		out: "  3\n 2\n1",
	},
	{
		in:  "\t\t\t3\n\t\t2\n\t1",
		out: "\t\t3\n\t2\n1",
	},
	{
		in:  "    4\n\n  2\n 1",
		out: "   4\n\n 2\n1",
	},
	{
		in:  "\t\t\t\t4\n\n\t\t2\n\t1",
		out: "\t\t\t4\n\n\t2\n1",
	},
	{
		in:  "    4\n   \n  2\n 1",
		out: "   4\n\n 2\n1",
	},
	{
		in:  "\t\t\t\t4\n\t\t\t\n\t\t2\n\t1",
		out: "\t\t\t4\n\n\t2\n1",
	},
	// mix
	{
		in:  "",
		out: "",
	},
	{
		in:  "\t ",
		out: "",
	},
	{
		in:  "\t  1\n\t \t2\n\t  3\n",
		out: " 1\n\t2\n 3\n",
	},
	{
		in:  " 1\n\t2\n 3\n",
		out: " 1\n\t2\n 3\n",
	},
	{
		in:  "\n\n\t 1\n\t\t2\n\t 3\n\t ",
		out: "\n 1\n\t2\n 3\n",
	},
	{
		in:  "\r\n\r\n\t 1\r\n\t\t2\r\n\t 3\r\n\t ",
		out: "\r\n 1\r\n\t2\r\n 3\r\n",
	},
}

func TestDedent(t *testing.T) {
	for _, tt := range dedentTests {
		if g, e := cli.Dedent(tt.in), tt.out; g != e {
			t.Fatalf("expected %q, got %q", e, g)
		}
	}
}
