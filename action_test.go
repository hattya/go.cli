//
// go.cli :: action_test.go
//
//   Copyright (c) 2014-2016 Akinori Hattori <hattya@gmail.com>
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
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hattya/go.cli"
)

func TestSubcommand(t *testing.T) {
	setup := func() (*cli.CLI, *cli.Command, *cli.Command) {
		app := cli.NewCLI()
		app.Stdout = ioutil.Discard
		app.Stderr = ioutil.Discard
		cmd := &cli.Command{
			Name:  []string{"cmd"},
			Flags: cli.NewFlagSet(),
		}
		app.Add(cmd)
		subcmd := &cli.Command{
			Name:  []string{"subcmd"},
			Flags: cli.NewFlagSet(),
		}
		cmd.Add(subcmd)
		return app, cmd, subcmd
	}

	app, cmd, _ := setup()
	args := []string{cmd.Name[0], "_"}
	if err := app.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(*cli.CommandError); !ok {
			t.Errorf("expected *cli.CommandError, got %T", err)
		}
		if !strings.Contains(err.Error(), "unknown") {
			t.Error("unexpected error:", err)
		}
	}

	app, cmd, _ = setup()
	args = []string{cmd.Name[0]}
	switch err := app.Run(args); {
	case err == nil:
		t.Error("expected error")
	case err != cli.ErrCommand:
		t.Error("unexpected error:", err)
	}

	app, cmd, subcmd := setup()
	args = []string{cmd.Name[0], subcmd.Name[0], "-g"}
	if err := app.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(cli.FlagError); !ok {
			t.Errorf("expected cli.FlagError, got %T", err)
		}
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	}

	app, cmd, subcmd = setup()
	args = []string{cmd.Name[0], subcmd.Name[0]}
	if err := app.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}

	app, cmd, subcmd = setup()
	subcmd.Action = func(*cli.Context) error {
		return nil
	}
	args = []string{cmd.Name[0], subcmd.Name[0]}
	if err := app.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}

	app, cmd, subcmd = setup()
	app.Flags.Bool("g", false, "")
	cmd.Flags.Bool("cmd", false, "")
	subcmd.Flags.Bool("subcmd", false, "")
	subcmd.Action = func(ctx *cli.Context) error {
		for _, n := range []string{"g", "cmd", "subcmd"} {
			if g, e := ctx.Bool(n), true; g != e {
				t.Errorf("expected %v, got %v", e, g)
			}
		}
		return nil
	}
	args = []string{"-g", cmd.Name[0], "-cmd", subcmd.Name[0], "-subcmd"}
	if err := app.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestSubcommandPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()

	app := cli.NewCLI()
	cmd := &cli.Command{Name: []string{"cmd"}}
	app.Add(cmd)
	cmd.Add(&cli.Command{
		Name:  []string{"subcmd"},
		Flags: cli.NewFlagSet(),
	})
	app.Run([]string{"cmd", "subcmd"})
}

func TestChain(t *testing.T) {
	setup := func() *cli.CLI {
		app := cli.NewCLI()
		app.Action = cli.Chain
		app.Stdout = ioutil.Discard
		app.Stderr = ioutil.Discard
		for _, n := range []string{"foo", "bar", "baz"} {
			cmd := &cli.Command{
				Name:  []string{n},
				Flags: cli.NewFlagSet(),
			}
			cmd.Flags.Bool(n, false, "")
			app.Add(cmd)
		}
		return app
	}

	app := setup()
	args := []string{"_"}
	if err := app.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(*cli.CommandError); !ok {
			t.Errorf("expected *cli.CommandError, got %T", err)
		}
		if !strings.Contains(err.Error(), "unknown") {
			t.Error("unexpected error:", err)
		}
	}

	app = setup()
	args = []string{}
	switch err := app.Run(args); {
	case err == nil:
		t.Error("expected error")
	case err != cli.ErrCommand:
		t.Error("unexpected error:", err)
	}

	app = setup()
	args = []string{app.Cmds[0].Name[0], "-chain"}
	if err := app.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(cli.FlagError); !ok {
			t.Errorf("expected cli.FlagError, got %T", err)
		}
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	}

	app = setup()
	args = nil
	for _, cmd := range app.Cmds {
		args = append(args, cmd.Name[0])
	}
	if err := app.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}

	app = setup()
	args = nil
	i := 0
	for _, cmd := range app.Cmds {
		args = append(args, cmd.Name[0], "-"+cmd.Name[0])
		cmd.Action = func(ctx *cli.Context) error {
			n := ctx.Stack[0].Name[0]
			if g, e := n, args[i]; g != e {
				t.Errorf("expected %v, got %v", e, g)
			}
			if g, e := ctx.Bool(n), true; g != e {
				t.Errorf("expected %v, got %v", e, g)
			}
			i += 2
			return nil
		}
	}
	if err := app.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestOption(t *testing.T) {
	setup := func() (*cli.CLI, *cli.Command) {
		app := cli.NewCLI()
		app.Action = cli.Option(func(*cli.Context) error {
			return nil
		})
		app.Stdout = ioutil.Discard
		app.Stderr = ioutil.Discard
		cmd := &cli.Command{
			Name:  []string{"cmd"},
			Flags: cli.NewFlagSet(),
		}
		app.Add(cmd)
		return app, cmd
	}

	app, _ := setup()
	args := []string{}
	if err := app.Run(args); err != nil {
		t.Error(err)
	}

	app, cmd := setup()
	args = []string{cmd.Name[0]}
	if err := app.Run(args); err != nil {
		t.Error(err)
	}
}

func TestSimple(t *testing.T) {
	app := cli.NewCLI()
	app.Action = cli.Simple(func(*cli.Context) error {
		return nil
	})
	app.Stdout = ioutil.Discard
	app.Stderr = ioutil.Discard

	args := []string{}
	if err := app.Run(args); err != nil {
		t.Error(err)
	}
}
