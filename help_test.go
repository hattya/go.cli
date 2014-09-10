//
// go.cli :: help_test.go
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

package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hattya/go.cli"
)

var helpOut = `usage: %s

options:

  -h, --help    show help
  --version     show version information
`

func TestHelp(t *testing.T) {
	b := new(bytes.Buffer)
	args := []string{"--help"}

	c := cli.NewCLI()
	c.Stdout = b
	if err := c.Run(args); err != nil {
		t.Error(err)
	}
	if g, e := b.String(), fmt.Sprintf(helpOut, c.Name); g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestMetaVar(t *testing.T) {
	flags := cli.NewFlagSet()
	flags.Bool("b, bool", false, "")
	flags.String("s, string", "", "")
	flags.Int("i", 0, "")
	if g, e := cli.MetaVar(flags.Lookup("bool")), ""; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := cli.MetaVar(flags.Lookup("string")), " <string>"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := cli.MetaVar(flags.Lookup("i")), " <i>"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}

	flags.MetaVar("bool", "=BOOL")
	if g, e := cli.MetaVar(flags.Lookup("bool")), "=BOOL"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
}
