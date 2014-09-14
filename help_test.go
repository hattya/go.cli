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
	"strings"
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
		t.Fatal(err)
	}
	if err := testOut(b.String(), fmt.Sprintf(helpOut, c.Name)); err != nil {
		t.Error(err)
	}
}

type usageTest struct {
	usage  interface{}
	format string
}

var usageTests = []usageTest{
	{
		usage:  nil,
		format: "usage: %s",
	},
	{
		usage:  "<options>",
		format: "usage: %s <options>",
	},
	{
		usage:  []string{"", ""},
		format: "usage: %[1]s\n   or: %[1]s",
	},
}

func TestUsage(t *testing.T) {
	for _, tt := range usageTests {
		c := cli.NewCLI()
		c.Usage = tt.usage
		if err := testOut(strings.Join(cli.Usage(c), "\n"), fmt.Sprintf(tt.format, c.Name)); err != nil {
			t.Error(err)
		}
	}
}

func TestUsagePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()

	c := cli.NewCLI()
	c.Usage = 1
	cli.Usage(c)
}

type metaVarTest struct {
	name     string
	value    interface{}
	metaVar  string
	expected string
}

var metaVarTests = []metaVarTest{
	{
		name:     "b,bool",
		value:    false,
		metaVar:  "",
		expected: "",
	},
	{
		name:     "s,string",
		value:    "",
		metaVar:  "",
		expected: " <string>",
	},
	{
		name:     "i",
		value:    0,
		metaVar:  "",
		expected: " <i>",
	},
	{
		name:     "b,bool",
		value:    false,
		metaVar:  "=bool",
		expected: "=bool",
	},
}

func TestMetaVar(t *testing.T) {
	for _, tt := range metaVarTests {
		list := strings.Split(tt.name, ",")
		n := list[len(list)-1]

		flags := cli.NewFlagSet()
		switch v := tt.value.(type) {
		case bool:
			flags.Bool(tt.name, v, "")
		case int:
			flags.Int(tt.name, v, "")
		case string:
			flags.String(tt.name, v, "")
		}
		flags.MetaVar(n, tt.metaVar)
		if g, e := cli.MetaVar(flags.Lookup(n)), tt.expected; g != e {
			t.Errorf("expected %v, got %v", e, g)
		}
	}
}
