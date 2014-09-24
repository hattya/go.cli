//
// go.cli :: command_test.go
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
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hattya/go.cli"
)

func TestCommand(t *testing.T) {
	newCLI := func() (*cli.CLI, *cli.Command) {
		c := cli.NewCLI()
		c.Stdout = ioutil.Discard
		c.Stderr = ioutil.Discard
		c.Add(&cli.Command{
			Name:  []string{"cmd"},
			Flags: cli.NewFlagSet(),
		})
		return c, c.Cmds[0]
	}

	c, _ := newCLI()
	args := []string{"_"}
	if err := c.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(*cli.CommandError); !ok {
			t.Errorf("expected *cli.CommandError, got %T", err)
		}
		if !strings.Contains(err.Error(), "unknown") {
			t.Error("unexpected error:", err)
		}
	}

	c, _ = newCLI()
	args = []string{}
	switch err := c.Run(args); {
	case err == nil:
		t.Error("expected error")
	case err != cli.ErrCommand:
		t.Error("unexpected error:", err)
	}

	c, cmd := newCLI()
	args = []string{cmd.Name[0], "-cli"}
	if err := c.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(cli.FlagError); !ok {
			t.Errorf("expected cli.FlagError, got %T", err)
		}
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	}

	c, cmd = newCLI()
	args = []string{cmd.Name[0]}
	if err := c.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}

	c, cmd = newCLI()
	cmd.Action = func(*cli.Context) error {
		return nil
	}
	args = []string{cmd.Name[0]}
	if err := c.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}

	c, cmd = newCLI()
	c.Flags.Bool("cli", false, "")
	cmd.Flags.Bool("cmd", false, "")
	cmd.Action = func(ctx *cli.Context) error {
		for _, n := range []string{"cli", "cmd"} {
			if g, e := ctx.Bool(n), true; g != e {
				t.Errorf("expected %v, got %v", e, g)
			}
		}
		return nil
	}
	args = []string{"-cli", cmd.Name[0], "-cmd"}
	if err := c.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestSubcommand(t *testing.T) {
	newCLI := func() (*cli.CLI, *cli.Command, *cli.Command) {
		c := cli.NewCLI()
		c.Stdout = ioutil.Discard
		c.Stderr = ioutil.Discard
		c.Add(&cli.Command{
			Name:  []string{"cmd"},
			Flags: cli.NewFlagSet(),
		})
		c.Cmds[0].Add(&cli.Command{
			Name:  []string{"subcmd"},
			Flags: cli.NewFlagSet(),
		})
		return c, c.Cmds[0], c.Cmds[0].Cmds[0]
	}

	c, cmd, _ := newCLI()
	args := []string{cmd.Name[0], "_"}
	if err := c.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(*cli.CommandError); !ok {
			t.Errorf("expected *cli.CommandError, got %T", err)
		}
		if !strings.Contains(err.Error(), "unknown") {
			t.Error("unexpected error:", err)
		}
	}

	c, cmd, _ = newCLI()
	args = []string{cmd.Name[0]}
	switch err := c.Run(args); {
	case err == nil:
		t.Error("expected error")
	case err != cli.ErrCommand:
		t.Error("unexpected error:", err)
	}

	c, cmd, subcmd := newCLI()
	args = []string{cmd.Name[0], subcmd.Name[0], "-cli"}
	if err := c.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(cli.FlagError); !ok {
			t.Errorf("expected cli.FlagError, got %T", err)
		}
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	}

	c, cmd, subcmd = newCLI()
	args = []string{cmd.Name[0], subcmd.Name[0]}
	if err := c.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}

	c, cmd, subcmd = newCLI()
	subcmd.Action = func(*cli.Context) error {
		return nil
	}
	args = []string{cmd.Name[0], subcmd.Name[0]}
	if err := c.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}

	c, cmd, subcmd = newCLI()
	c.Flags.Bool("cli", false, "")
	cmd.Flags.Bool("cmd", false, "")
	subcmd.Flags.Bool("subcmd", false, "")
	subcmd.Action = func(ctx *cli.Context) error {
		for _, n := range []string{"cli", "cmd", "subcmd"} {
			if g, e := ctx.Bool(n), true; g != e {
				t.Errorf("expected %v, got %v", e, g)
			}
		}
		return nil
	}
	args = []string{"-cli", cmd.Name[0], "-cmd", subcmd.Name[0], "-subcmd"}
	if err := c.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestChain(t *testing.T) {
	newCLI := func() (*cli.CLI, []*cli.Command) {
		c := cli.NewCLI()
		c.Action = cli.Chain
		c.Stdout = ioutil.Discard
		c.Stderr = ioutil.Discard
		for _, n := range []string{"foo", "bar", "baz"} {
			cmd := &cli.Command{
				Name:  []string{n},
				Flags: cli.NewFlagSet(),
			}
			cmd.Flags.Bool(n, false, "")
			c.Add(cmd)
		}
		return c, c.Cmds
	}

	c, _ := newCLI()
	args := []string{"_"}
	if err := c.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(*cli.CommandError); !ok {
			t.Errorf("expected *cli.CommandError, got %T", err)
		}
		if !strings.Contains(err.Error(), "unknown") {
			t.Error("unexpected error:", err)
		}
	}

	c, _ = newCLI()
	args = []string{}
	switch err := c.Run(args); {
	case err == nil:
		t.Error("expected error")
	case err != cli.ErrCommand:
		t.Error("unexpected error:", err)
	}

	c, cmds := newCLI()
	args = []string{cmds[0].Name[0], "-chain"}
	if err := c.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(cli.FlagError); !ok {
			t.Errorf("expected cli.FlagError, got %T", err)
		}
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	}

	c, cmds = newCLI()
	args = nil
	for _, cmd := range cmds {
		args = append(args, cmd.Name[0])
	}
	if err := c.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}

	c, cmds = newCLI()
	args = nil
	i := 0
	for _, cmd := range cmds {
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
	if err := c.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestFindCommand(t *testing.T) {
	cmds := []*cli.Command{
		{Name: []string{"foo"}},
		{Name: []string{"bar"}},
		{Name: []string{"baz"}},
	}

	_, err := cli.FindCommand(cmds, "")
	switch {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), "unknown"):
		t.Fatal("unexpected error:", err)
	}

	cmd, err := cli.FindCommand(cmds, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if g, e := cmd.Name[0], "foo"; g != e {
		t.Errorf("expected %v, got %v", e, g)
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
		t.Errorf("expected %v, got %v", e, g)
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
