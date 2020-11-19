//
// go.cli :: cli.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/term"
)

type CLI struct {
	Name    string
	Version string
	Usage   interface{}
	Desc    string
	Epilog  string
	Cmds    []*Command
	Flags   *FlagSet

	Prepare      func(*Context, *Command) error
	Action       Action
	ErrorHandler func(*Context, error) error

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	ctx     context.Context
	cancel  context.CancelFunc
	help    bool
	version bool
}

func NewCLI() *CLI {
	name := filepath.Base(os.Args[0])
	if runtime.GOOS == "windows" {
		name = name[:len(name)-len(filepath.Ext(name))]
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &CLI{
		Name:   name,
		Flags:  NewFlagSet(),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (ui *CLI) Run(args []string) error {
	if ui.Action == nil {
		ui.Action = DefaultAction
	}
	if ui.ErrorHandler == nil {
		ui.ErrorHandler = ErrorHandler
	}

	if ui.Stdin == nil {
		ui.Stdin = os.Stdin
	}
	if ui.Stdout == nil {
		ui.Stdout = os.Stdout
	}
	if ui.Stderr == nil {
		ui.Stderr = os.Stderr
	}

	if ui.Flags.Lookup("h") == nil && ui.Flags.Lookup("help") == nil {
		ui.Flags.Bool("h, help", false, "show help")
		ui.help = true
	}
	if ui.Flags.Lookup("version") == nil {
		ui.Flags.Bool("version", false, "show version information")
		ui.version = true
	}

	ctx := NewContext(ui)
	select {
	case <-ui.ctx.Done():
		return ctx.ErrorHandler(Interrupt{})
	default:
	}
	if err := ui.Flags.Parse(args); err != nil {
		return ctx.ErrorHandler(err)
	}
	ctx.Args = ui.Flags.Args()
	switch {
	case ui.help && ctx.Bool("help"):
		return Help(ctx)
	case ui.version && ctx.Bool("version"):
		return Version(ctx)
	}
	err := ui.Action(ctx)
	select {
	case <-ui.ctx.Done():
		err = Interrupt{}
	default:
	}
	return ctx.ErrorHandler(err)
}

func (ui *CLI) Add(cmd *Command) {
	ui.Cmds = append(ui.Cmds, cmd)
}

func (ui *CLI) Context() context.Context {
	return ui.ctx
}

func (ui *CLI) Interrupt() {
	ui.cancel()
}

func (ui *CLI) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(ui.Stdout, a...)
}

func (ui *CLI) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(ui.Stdout, a...)
}

func (ui *CLI) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(ui.Stdout, format, a...)
}

func (ui *CLI) Error(a ...interface{}) (int, error) {
	return fmt.Fprint(ui.Stderr, a...)
}

func (ui *CLI) Errorln(a ...interface{}) (int, error) {
	return fmt.Fprintln(ui.Stderr, a...)
}

func (ui *CLI) Errorf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(ui.Stderr, format, a...)
}

func (ui *CLI) Title(title string) error {
	return ui.title(title)
}

func (ui *CLI) Prompt(prompt string) (string, error) {
	ui.Print(prompt)
	return ui.readLine()
}

func (ui *CLI) Password(prompt string) (string, error) {
	ui.Print(prompt)
	defer ui.Println()
	if f, ok := ui.Stdin.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		b, err := term.ReadPassword(int(f.Fd()))
		return string(b), err
	}
	return ui.readLine()
}

func (ui *CLI) readLine() (string, error) {
	b := make([]byte, 1)
	var in []byte
	for {
		n, err := ui.Stdin.Read(b)
		if 0 < n {
			if b[0] == '\n' {
				if 0 < len(in) && in[len(in)-1] == '\r' {
					in = in[:len(in)-1]
				}
				return string(in), nil
			}
			in = append(in, b...)
		}
		if err != nil {
			if err == io.EOF && 0 < len(in) {
				err = nil
			}
			return string(in), err
		}
	}
}

var (
	ErrCommand = errors.New("cli: command required")
	ErrFlags   = errors.New("cli: flag parsing is disabled")
	ErrArgs    = errors.New("invalid arguments")
)

type Abort struct {
	Err  error
	Hint string
}

func (e Abort) Error() string { return e.Err.Error() }

type Interrupt struct{}

func (e Interrupt) Error() string { return "interrupted" }

func ErrorHandler(ctx *Context, err error) error {
	switch err := err.(type) {
	case nil:
	case Abort:
		ctx.UI.Errorf("%v: %v\n", ctx.UI.Name, err)
		if err.Hint != "" {
			ctx.UI.Errorln(err.Hint)
		}
	case CommandError:
		if len(err.List) == 0 {
			ctx.UI.Errorf("%v: %v\n", ctx.Name(), err)
			Help(ctx)
		} else {
			ctx.UI.Errorf("%v: command '%v' is ambiguous\n", ctx.Name(), err.Name)
			ctx.UI.Errorf("    %v\n", strings.Join(err.List, " "))
		}
	case FlagError:
		ctx.UI.Errorf("%v: %v\n", ctx.Name(), err)
		Help(ctx)
	default:
		if err == ErrCommand {
			Help(ctx)
		} else {
			ctx.UI.Errorf("%v: %v\n", ctx.Name(), err)
		}
	}
	return err
}
