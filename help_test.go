//
// go.cli :: help_test.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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
			{
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
			case tt.err && err == nil:
				t.Fatal("expected error")
			case !tt.err && err != nil:
				t.Fatal(err)
			}
			if err := testOut(strings.SplitN(b.String(), "\n", 2)[0], fmt.Sprintf(tt.out, app.Name)); err != nil {
				t.Error(err)
			}
		}
	}
}

var options = strings.TrimSpace(cli.Dedent(`
	options:

	  -h, --help    show help
	  --version     show version information
`))

var helpTests = []struct {
	usage  interface{}
	desc   string
	epilog string
	cmds   []*cli.Command
	out    string
}{
	{
		out: cli.Dedent(`
			usage: %[1]v

			%[2]v

		`),
	},
	{
		usage: "<options>",
		out: cli.Dedent(`
			usage: %[1]v <options>

			%[2]v

		`),
	},
	{
		usage: []string{
			"add <path>...",
			"rm <path>...",
		},
		out: cli.Dedent(`
			usage: %[1]v add <path>...
			   or: %[1]v rm <path>...

			%[2]v

		`),
	},
	{
		desc: "    desc",
		out: cli.Dedent(`
			usage: %[1]v

			    desc

			%[2]v

		`),
	},
	{
		epilog: "epilog",
		out: cli.Dedent(`
			usage: %[1]v

			%[2]v

			epilog
		`),
	},
	{
		desc:   "    desc",
		epilog: "epilog",
		out: cli.Dedent(`
			usage: %[1]v

			    desc

			%[2]v

			epilog
		`),
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
			},
		},
		out: cli.Dedent(`
			usage: %[1]v

			commands:

			  cmd

			%[2]v

		`),
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
				Desc: "desc",
			},
		},
		out: cli.Dedent(`
			usage: %[1]v

			commands:

			  cmd    desc

			%[2]v

		`),
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"cmd"},
				Desc: " desc \n",
			},
		},
		out: cli.Dedent(`
			usage: %[1]v

			commands:

			  cmd    desc

			%[2]v

		`),
	},
}

func TestHelp(t *testing.T) {
	for _, tt := range helpTests {
		var b bytes.Buffer
		app := cli.NewCLI()
		app.Usage = tt.usage
		app.Desc = tt.desc
		app.Epilog = tt.epilog
		app.Cmds = tt.cmds
		app.Stdout = &b
		if err := app.Run([]string{"--help"}); err != nil {
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
		out: cli.Dedent(`
			usage: %[1]v %[2]v
		`),
	},
	{
		alias: []string{"alias"},
		out: cli.Dedent(`
			usage: %[1]v %[2]v

			alias: alias
		`),
	},
	{
		desc: "    desc",
		out: cli.Dedent(`
			usage: %[1]v %[2]v

			    desc

		`),
	},
	{
		epilog: "epilog",
		out: cli.Dedent(`
			usage: %[1]v %[2]v

			epilog
		`),
	},
	{
		alias:  []string{"alias"},
		desc:   "    desc",
		epilog: "epilog",
		out: cli.Dedent(`
			usage: %[1]v %[2]v

			alias: alias

			    desc

			epilog
		`),
	},
	{
		cmds: []*cli.Command{
			{
				Name: []string{"subcmd"},
				Desc: "desc",
			},
		},
		out: cli.Dedent(`
			usage: %[1]v %[2]v

			commands:

			  subcmd    desc

		`),
	},
}

func TestCommandHelp(t *testing.T) {
	for _, tt := range commandHelpTests {
		var b bytes.Buffer
		app := cli.NewCLI()
		app.Stdout = &b
		app.Add(&cli.Command{
			Name:   append([]string{"cmd"}, tt.alias...),
			Usage:  tt.usage,
			Desc:   tt.desc,
			Epilog: tt.epilog,
			Cmds:   tt.cmds,
			Flags:  cli.NewFlagSet(),
		})
		if err := app.Run([]string{app.Cmds[0].Name[0], "--help"}); err != nil {
			t.Fatal(err)
		}
		if err := testOut(b.String(), fmt.Sprintf(tt.out, app.Name, app.Cmds[0].Name[0])); err != nil {
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
		if recover() == nil {
			t.Error("expected panic")
		}
	}()

	app := cli.NewCLI()
	app.Usage = 1
	cli.Usage(cli.NewContext(app))
}

var metaVarTests = []struct {
	name    string
	value   interface{}
	metaVar string
	out     string
}{
	{
		name:    "b, bool",
		value:   false,
		metaVar: "",
		out:     "",
	},
	{
		name:    "s, string",
		value:   "",
		metaVar: "",
		out:     " <string>",
	},
	{
		name:    "i",
		value:   0,
		metaVar: "",
		out:     " <i>",
	},
	{
		name:    "b, bool",
		value:   false,
		metaVar: "=bool",
		out:     "=bool",
	},
}

func TestMetaVar(t *testing.T) {
	for _, tt := range metaVarTests {
		list := strings.Split(tt.name, ",")
		n := strings.TrimSpace(list[len(list)-1])

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
		if g, e := cli.MetaVar(flags.Lookup(n)), tt.out; g != e {
			t.Errorf("expected %q, got %q", e, g)
		}
	}
}
