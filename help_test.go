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

var options = `options:

  -h, --help    show help
  --version     show version information`

type helpTest struct {
	usage  interface{}
	desc   string
	epilog string
	cmds   []*cli.Command
	out    string
}

var helpTests = []helpTest{
	{
		out: `usage: %[1]s

%[2]s

`,
	},
	{
		usage: "<options>",
		out: `usage: %[1]s <options>

%[2]s

`,
	},
	{
		usage: []string{
			"add <path>...",
			"rm <path>...",
		},
		out: `usage: %[1]s add <path>...
   or: %[1]s rm <path>...

%[2]s

`,
	},
	{
		desc: "    desc",
		out: `usage: %[1]s

    desc

%[2]s

`,
	},
	{
		epilog: "epilog",
		out: `usage: %[1]s

%[2]s

epilog
`,
	},
	{
		desc:   "    desc",
		epilog: "epilog",
		out: `usage: %[1]s

    desc

%[2]s

epilog
`,
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
			},
		},
		out: `usage: %[1]s

commands:

  cmd

%[2]s

`,
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
				Desc: "desc",
			},
		},
		out: `usage: %[1]s

commands:

  cmd    desc

%[2]s

`,
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
				Desc: " desc \n",
			},
		},
		out: `usage: %[1]s

commands:

  cmd    desc

%[2]s

`,
	},
}

func TestHelp(t *testing.T) {
	var b bytes.Buffer
	args := []string{"--help"}
	for _, tt := range helpTests {
		b.Reset()
		c := cli.NewCLI()
		c.Usage = tt.usage
		c.Desc = tt.desc
		c.Epilog = tt.epilog
		c.Cmds = tt.cmds
		c.Stdout = &b
		if err := c.Run(args); err != nil {
			t.Fatal(err)
		}
		if err := testOut(b.String(), fmt.Sprintf(tt.out, c.Name, options)); err != nil {
			t.Error(err)
		}
	}
}

type cmdHelpTest struct {
	usage  interface{}
	desc   string
	epilog string
	cmds   []*cli.Command
	out    string
}

var cmdHelpTests = []cmdHelpTest{
	{
		out: `usage: %[1]s %[2]s
`,
	},
	{
		desc: "    desc",
		out: `usage: %[1]s %[2]s

    desc

`,
	},
	{
		epilog: "epilog",
		out: `usage: %[1]s %[2]s

epilog
`,
	},
	{
		desc:   "    desc",
		epilog: "epilog",
		out: `usage: %[1]s %[2]s

    desc

epilog
`,
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
				Desc: "desc",
			},
		},
		out: `usage: %[1]s %[2]s

commands:

  cmd    desc

`,
	},
}

func TestCmdHelp(t *testing.T) {
	b := new(bytes.Buffer)
	name := []string{"cmd"}
	args := []string{name[0], "--help"}
	for _, tt := range cmdHelpTests {
		b.Reset()
		c := cli.NewCLI()
		c.Add(&cli.Command{
			Name:   name,
			Usage:  tt.usage,
			Desc:   tt.desc,
			Epilog: tt.epilog,
			Cmds:   tt.cmds,
			Flags:  cli.NewFlagSet(),
		})
		c.Stdout = b
		if err := c.Run(args); err != nil {
			t.Fatal(err)
		}
		if err := testOut(b.String(), fmt.Sprintf(tt.out, c.Name, name[0])); err != nil {
			t.Error(err)
		}
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
		if err := testOut(strings.Join(cli.Usage(cli.NewContext(c)), "\n"), fmt.Sprintf(tt.format, c.Name)); err != nil {
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
	cli.Usage(cli.NewContext(c))
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
