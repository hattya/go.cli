//
// go.cli :: action_test.go
//
//   Copyright (c) 2014-2022 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package cli_test

import (
	"io"
	"strings"
	"testing"

	"github.com/hattya/go.cli"
)

func TestSubcommand(t *testing.T) {
	setup := func() *cli.CLI {
		app := cli.NewCLI()
		app.Stdout = io.Discard
		app.Stderr = io.Discard
		app.Add(&cli.Command{
			Name:  []string{"cmd"},
			Flags: cli.NewFlagSet(),
		})
		app.Cmds[0].Add(&cli.Command{
			Name:  []string{"subcmd"},
			Flags: cli.NewFlagSet(),
		})
		return app
	}

	app := setup()
	if err := app.Run([]string{app.Cmds[0].Name[0], app.Cmds[0].Cmds[0].Name[0]}); err != nil {
		t.Error("unexpected error:", err)
	}

	app = setup()
	app.Cmds[0].Cmds[0].Action = func(*cli.Context) error { return nil }
	if err := app.Run([]string{app.Cmds[0].Name[0], app.Cmds[0].Cmds[0].Name[0]}); err != nil {
		t.Error("unexpected error:", err)
	}

	app = setup()
	app.Flags.Bool("cli", false, "")
	app.Cmds[0].Flags.Bool("cmd", false, "")
	app.Cmds[0].Cmds[0].Flags.Bool("subcmd", false, "")
	app.Cmds[0].Cmds[0].Action = func(ctx *cli.Context) error {
		for _, n := range []string{"cli", "cmd", "subcmd"} {
			if !ctx.Bool(n) {
				t.Errorf("Context.Bool(%q) = false, expected true", n)
			}
		}
		return nil
	}
	if err := app.Run([]string{"-cli", app.Cmds[0].Name[0], "-cmd", app.Cmds[0].Cmds[0].Name[0], "-subcmd"}); err != nil {
		t.Error("unexpected error:", err)
	}
	// no subcommand specified
	app = setup()
	if err := app.Run([]string{app.Cmds[0].Name[0]}); err != cli.ErrCommand {
		t.Errorf("expected ErrCommand, got %#v", err)
	}
	// unknown subcommand
	app = setup()
	switch err := app.Run([]string{app.Cmds[0].Name[0], "_"}).(type) {
	case cli.CommandError:
		if !strings.Contains(err.Error(), "unknown") {
			t.Error("unexpected error:", err)
		}
	default:
		t.Errorf("expected CommandError, got %#v", err)
	}
	// flag error
	app = setup()
	switch err := app.Run([]string{app.Cmds[0].Name[0], app.Cmds[0].Cmds[0].Name[0], "-subcmd"}).(type) {
	case cli.FlagError:
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	default:
		t.Errorf("expected FlagError, got %#v", err)
	}
}

func TestSubcommandPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()

	app := cli.NewCLI()
	app.Add(&cli.Command{
		Name: []string{"cmd"},
	})
	app.Cmds[0].Add(&cli.Command{
		Name:  []string{"subcmd"},
		Flags: cli.NewFlagSet(),
	})
	app.Run([]string{app.Cmds[0].Name[0], app.Cmds[0].Cmds[0].Name[0]})
}

func TestChain(t *testing.T) {
	setup := func() *cli.CLI {
		app := cli.NewCLI()
		app.Action = cli.Chain
		app.Stdout = io.Discard
		app.Stderr = io.Discard
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
	args := make([]string, len(app.Cmds))
	for i, cmd := range app.Cmds {
		args[i] = cmd.Name[0]
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
				t.Errorf("expected %q, got %q", e, g)
			}
			if !ctx.Bool(n) {
				t.Errorf("Context.Bool(%q) = false, expected true", n)
			}
			i += 2
			return nil
		}
	}
	if err := app.Run(args); err != nil {
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
	switch err := app.Run([]string{app.Cmds[0].Name[0], "-chain"}).(type) {
	case cli.FlagError:
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	default:
		t.Errorf("expected FlagError, got %#v", err)
	}
}

func TestChainInterrupt(t *testing.T) {
	app := cli.NewCLI()
	app.Action = cli.Chain
	app.Stdout = io.Discard
	app.Stderr = io.Discard
	for _, n := range []string{"foo", "bar", "baz"} {
		app.Add(&cli.Command{
			Name: []string{n},
			Action: func(ctx *cli.Context) error {
				ctx.Interrupt()
				return nil
			},
		})
	}

	args := make([]string, len(app.Cmds))
	for i, cmd := range app.Cmds {
		args[i] = cmd.Name[0]
	}
	switch err := app.Run(args).(type) {
	case cli.Interrupt:
	default:
		t.Errorf("expected Interrupt, got %#v", err)
	}
}

func TestOption(t *testing.T) {
	setup := func() *cli.CLI {
		app := cli.NewCLI()
		app.Action = cli.Option(func(*cli.Context) error { return nil })
		app.Stdout = io.Discard
		app.Stderr = io.Discard
		app.Add(&cli.Command{
			Name:  []string{"cmd"},
			Flags: cli.NewFlagSet(),
		})
		return app
	}

	app := setup()
	if err := app.Run(nil); err != nil {
		t.Error(err)
	}

	app = setup()
	if err := app.Run([]string{app.Cmds[0].Name[0]}); err != nil {
		t.Error(err)
	}
}

func TestSimple(t *testing.T) {
	app := cli.NewCLI()
	app.Action = cli.Simple(func(*cli.Context) error { return nil })
	app.Stdout = io.Discard
	app.Stderr = io.Discard

	if err := app.Run(nil); err != nil {
		t.Error(err)
	}
}
