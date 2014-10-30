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

var helpCommandTests = []struct {
	args []string
	cmds []*cli.Command
	err  bool
	out  string
}{
	{
		args: []string{"help"},
		out:  "usage: %v",
	},
	{
		args: []string{"help", "help"},
		out:  "usage: %v help [<command>]",
	},
	{
		args: []string{"help", "_"},
		err:  true,
		out:  "%v: unknown command '_'",
	},
	{
		args: []string{"help", "help", "_"},
		err:  true,
		out:  "%v help: " + cli.ErrArgs.Error(),
	},
	{
		args: []string{"cmd", "help"},
		out:  "usage: %v",
		cmds: []*cli.Command{
			&cli.Command{
				Name: []string{"cmd"},
				Cmds: []*cli.Command{
					cli.NewHelpCommand(),
				},
				Flags: cli.NewFlagSet(),
			},
		},
	},
}

func TestHelpCommand(t *testing.T) {
	for _, tt := range helpCommandTests {
		for _, action := range []cli.Action{cli.Subcommand, cli.Chain} {
			var b bytes.Buffer
			app := cli.NewCLI()
			app.Cmds = tt.cmds
			app.Action = action
			app.Stdout = &b
			app.Stderr = &b
			app.Add(cli.NewHelpCommand())
			switch err := app.Run(tt.args); {
			case !tt.err && err != nil:
				t.Fatal(err)
			case tt.err && err == nil:
				t.Fatal("expected error")
			}
			if err := testOut(strings.SplitN(b.String(), "\n", 2)[0], fmt.Sprintf(tt.out, app.Name)); err != nil {
				t.Error(err)
			}
		}
	}
}

var options = `options:

  -h, --help    show help
  --version     show version information`

var helpTests = []struct {
	usage  interface{}
	desc   string
	epilog string
	cmds   []*cli.Command
	out    string
}{
	{
		out: `usage: %[1]v

%[2]v

`,
	},
	{
		usage: "<options>",
		out: `usage: %[1]v <options>

%[2]v

`,
	},
	{
		usage: []string{
			"add <path>...",
			"rm <path>...",
		},
		out: `usage: %[1]v add <path>...
   or: %[1]v rm <path>...

%[2]v

`,
	},
	{
		desc: "    desc",
		out: `usage: %[1]v

    desc

%[2]v

`,
	},
	{
		epilog: "epilog",
		out: `usage: %[1]v

%[2]v

epilog
`,
	},
	{
		desc:   "    desc",
		epilog: "epilog",
		out: `usage: %[1]v

    desc

%[2]v

epilog
`,
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
			},
		},
		out: `usage: %[1]v

commands:

  cmd

%[2]v

`,
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
				Desc: "desc",
			},
		},
		out: `usage: %[1]v

commands:

  cmd    desc

%[2]v

`,
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
				Desc: " desc \n",
			},
		},
		out: `usage: %[1]v

commands:

  cmd    desc

%[2]v

`,
	},
}

func TestHelp(t *testing.T) {
	var b bytes.Buffer
	args := []string{"--help"}
	for _, tt := range helpTests {
		b.Reset()
		app := cli.NewCLI()
		app.Usage = tt.usage
		app.Desc = tt.desc
		app.Epilog = tt.epilog
		app.Cmds = tt.cmds
		app.Stdout = &b
		if err := app.Run(args); err != nil {
			t.Fatal(err)
		}
		if err := testOut(b.String(), fmt.Sprintf(tt.out, app.Name, options)); err != nil {
			t.Error(err)
		}
	}
}

var commandHelpTests = []struct {
	alias  []string
	usage  interface{}
	desc   string
	epilog string
	cmds   []*cli.Command
	out    string
}{
	{
		out: `usage: %[1]v %[2]v
`,
	},
	{
		alias: []string{"alias"},
		out: `usage: %[1]v %[2]v

alias: alias
`,
	},
	{
		desc: "    desc",
		out: `usage: %[1]v %[2]v

    desc

`,
	},
	{
		epilog: "epilog",
		out: `usage: %[1]v %[2]v

epilog
`,
	},
	{
		alias:  []string{"alias"},
		desc:   "    desc",
		epilog: "epilog",
		out: `usage: %[1]v %[2]v

alias: alias

    desc

epilog
`,
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"subcmd"},
				Desc: "desc",
			},
		},
		out: `usage: %[1]v %[2]v

commands:

  subcmd    desc

`,
	},
}

func TestCommandHelp(t *testing.T) {
	b := new(bytes.Buffer)
	name := []string{"cmd"}
	args := []string{name[0], "--help"}
	for _, tt := range commandHelpTests {
		b.Reset()
		app := cli.NewCLI()
		app.Stdout = b
		app.Add(&cli.Command{
			Name:   append(append([]string{}, name...), tt.alias...),
			Usage:  tt.usage,
			Desc:   tt.desc,
			Epilog: tt.epilog,
			Cmds:   tt.cmds,
			Flags:  cli.NewFlagSet(),
		})
		if err := app.Run(args); err != nil {
			t.Fatal(err)
		}
		if err := testOut(b.String(), fmt.Sprintf(tt.out, app.Name, name[0])); err != nil {
			t.Error(err)
		}
	}
}

var usageTests = []struct {
	usage  interface{}
	format string
}{
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
		format: "usage: %[1]v\n   or: %[1]v",
	},
}

func TestUsage(t *testing.T) {
	for _, tt := range usageTests {
		app := cli.NewCLI()
		app.Usage = tt.usage
		if err := testOut(strings.Join(cli.Usage(cli.NewContext(app)), "\n"), fmt.Sprintf(tt.format, app.Name)); err != nil {
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

	app := cli.NewCLI()
	app.Usage = 1
	cli.Usage(cli.NewContext(app))
}

var metaVarTests = []struct {
	name     string
	value    interface{}
	metaVar  string
	expected string
}{
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
