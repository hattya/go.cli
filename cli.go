//
// go.cli :: cli.go
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

package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
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

	help    bool
	version bool
}

func NewCLI() *CLI {
	name := filepath.Base(os.Args[0])
	if runtime.GOOS == "windows" {
		name = name[:len(name)-len(filepath.Ext(name))]
	}
	return &CLI{
		Name:  name,
		Flags: NewFlagSet(),
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
	return ui.Action(ctx)
}

func (ui *CLI) Add(cmd *Command) {
	ui.Cmds = append(ui.Cmds, cmd)
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
	if f, ok := ui.Stdin.(*os.File); ok && terminal.IsTerminal(int(f.Fd())) {
		b, err := terminal.ReadPassword(int(f.Fd()))
		return string(b), err
	}
	return ui.readLine()
}

func (ui *CLI) readLine() (string, error) {
	b := make([]byte, 1024)
	var in []byte
	for {
		n, err := ui.Stdin.Read(b)
		if err != nil {
			return "", err
		}
		if n == 0 {
			if len(in) == 0 {
				return "", io.EOF
			}
			break
		}
		if i := bytes.IndexByte(b[:n], '\n'); i != -1 {
			n = i
		}
		in = append(in, b[:n]...)
		if n < len(b) {
			if 0 < len(in) && in[len(in)-1] == '\r' {
				in = in[:len(in)-1]
			}
			break
		}
	}
	return string(in), nil
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

func ErrorHandler(ctx *Context, err error) error {
	switch err := err.(type) {
	case nil:
	case FlagError:
		ctx.UI.Errorf("%v: %v\n", ctx.Name(), err)
		Help(ctx)
	case *CommandError:
		if len(err.List) == 0 {
			ctx.UI.Errorf("%v: %v\n", ctx.Name(), err)
			Help(ctx)
		} else {
			ctx.UI.Errorf("%v: command '%v' is ambiguous\n", ctx.Name(), err.Name)
			ctx.UI.Errorf("    %v\n", strings.Join(err.List, " "))
		}
	case *Abort:
		ctx.UI.Errorf("%v: %v\n", ctx.UI.Name, err)
		if err.Hint != "" {
			ctx.UI.Errorln(err.Hint)
		}
	default:
		if err == ErrCommand {
			Help(ctx)
		} else {
			ctx.UI.Errorf("%v: %v\n", ctx.Name(), err)
		}
	}
	return err
}
