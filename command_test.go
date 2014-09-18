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
	case err != cli.ErrCmd:
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
	cmd.Action = func(*cli.Context) error { return nil }
	args = []string{cmd.Name[0]}
	if err := c.Run(args); err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestFindCmd(t *testing.T) {
	cmds := []*cli.Command{
		{Name: []string{"foo"}},
		{Name: []string{"bar"}},
		{Name: []string{"baz"}},
	}

	_, err := cli.FindCmd(cmds, "")
	switch {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), "unknown"):
		t.Fatal("unexpected error:", err)
	}

	cmd, err := cli.FindCmd(cmds, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if g, e := cmd.Name[0], "foo"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}

	_, err = cli.FindCmd(cmds, "b")
	switch {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), "ambiguous"):
		t.Fatal("unexpected error:", err)
	}

	cmds[1].Name = append(cmds[1].Name, "b")
	cmd, err = cli.FindCmd(cmds, "b")
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
