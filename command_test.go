//
// go.cli :: command_test.go
//
//   Copyright (c) 2014-2017 Akinori Hattori <hattya@gmail.com>
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

func TestCommand(t *testing.T) {
	setup := func() *cli.CLI {
		app := cli.NewCLI()
		app.Stdout = ioutil.Discard
		app.Stderr = ioutil.Discard
		app.Add(&cli.Command{
			Name:  []string{"cmd"},
			Flags: cli.NewFlagSet(),
		})
		return app
	}

	app := setup()
	if err := app.Run([]string{app.Cmds[0].Name[0]}); err != nil {
		t.Error("unexpected error:", err)
	}

	app = setup()
	app.Cmds[0].Action = func(*cli.Context) error { return nil }
	if err := app.Run([]string{app.Cmds[0].Name[0]}); err != nil {
		t.Error("unexpected error:", err)
	}

	app = setup()
	app.Flags.Bool("cli", false, "")
	app.Cmds[0].Flags.Bool("cmd", false, "")
	app.Cmds[0].Action = func(ctx *cli.Context) error {
		for _, n := range []string{"cli", "cmd"} {
			if !ctx.Bool(n) {
				t.Errorf("Context.Bool(%q) = false, expected true", n)
			}
		}
		return nil
	}
	if err := app.Run([]string{"-cli", app.Cmds[0].Name[0], "-cmd"}); err != nil {
		t.Error("unexpected error:", err)
	}
	// no command specified
	app = setup()
	if err := app.Run(nil); err != cli.ErrCommand {
		t.Errorf("expected ErrCommand, got %#v", err)
	}
	// unknown command
	app = setup()
	switch err := app.Run([]string{"_"}).(type) {
	case cli.CommandError:
		if !strings.Contains(err.Error(), "unknown") {
			t.Error("unexpected error:", err)
		}
	default:
		t.Errorf("expected CommandError, got %#v", err)
	}
	// flag error
	app = setup()
	switch err := app.Run([]string{app.Cmds[0].Name[0], "-cli"}).(type) {
	case cli.FlagError:
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	default:
		t.Errorf("expected FlagError, got %#v", err)
	}
}

func TestFindCommand(t *testing.T) {
	cmds := []*cli.Command{
		{Name: []string{"foo"}},
		{Name: []string{"bar"}},
		{Name: []string{"baz"}},
	}

	cmd, err := cli.FindCommand(cmds, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if g, e := cmd.Name[0], "foo"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}

	_, err = cli.FindCommand(cmds, "")
	switch {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), "unknown"):
		t.Fatal("unexpected error:", err)
	}

	_, err = cli.FindCommand(cmds, "b")
	switch {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), "ambiguous"):
		t.Fatal("unexpected error:", err)
	}

	cmds[1].Name = append(cmds[1].Name, "b")
	cmd, err = cli.FindCommand(cmds, "b")
	if err != nil {
		t.Fatal(err)
	}
	if g, e := cmd.Name[0], "bar"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestSortCommands(t *testing.T) {
	cmds := []*cli.Command{
		{Name: []string{"2"}},
		{Name: []string{"3"}},
		{Name: []string{"1"}},
	}
	cli.CommandSlice(cmds).Sort()
	if err := testStrings(func(i int) string { return cmds[i].Name[0] }, []string{"1", "2", "3"}); err != nil {
		t.Error(err)
	}
}
