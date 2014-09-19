//
// go.cli :: cli.go
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

package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

var ErrCmd = errors.New("cli: command required")

type CLI struct {
	Name    string
	Version string
	Usage   interface{}
	Desc    string
	Epilog  string
	Cmds    []*Command
	Flags   *FlagSet
	Action  func(*Context) error

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
		Name:   name,
		Flags:  NewFlagSet(),
		Action: Action,
	}
}

func (c *CLI) Run(args []string) error {
	if c.Stdin == nil {
		c.Stdin = os.Stdin
	}
	if c.Stdout == nil {
		c.Stdout = os.Stdout
	}
	if c.Stderr == nil {
		c.Stderr = os.Stderr
	}

	if c.Flags.Lookup("h") == nil && c.Flags.Lookup("help") == nil {
		c.Flags.Bool("h, help", false, "show help")
		c.help = true
	}
	if c.Flags.Lookup("version") == nil {
		c.Flags.Bool("version", false, "show version information")
		c.version = true
	}

	ctx := NewContext(c)
	if err := c.Flags.Parse(args); err != nil {
		Help(ctx, err)
		return err
	}
	ctx.Args = c.Flags.Args()
	return c.Action(ctx)
}

func (c *CLI) Add(cmd *Command) {
	c.Cmds = append(c.Cmds, cmd)
}

func (c *CLI) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(c.Stdout, a...)
}

func (c *CLI) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(c.Stdout, a...)
}

func (c *CLI) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(c.Stdout, format, a...)
}

func (c *CLI) Error(a ...interface{}) (int, error) {
	return fmt.Fprint(c.Stderr, a...)
}

func (c *CLI) Errorln(a ...interface{}) (int, error) {
	return fmt.Fprintln(c.Stderr, a...)
}

func (c *CLI) Errorf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(c.Stderr, format, a...)
}

var Action = DefaultAction

func DefaultAction(ctx *Context) error {
	switch {
	case ctx.CLI.help && ctx.Bool("help"):
		return Help(ctx, nil)
	case ctx.CLI.version && ctx.Bool("version"):
		return Version(ctx)
	}

	cmd, err := ctx.Command()
	if cmd != nil {
		ctx.Stack = append(ctx.Stack, cmd)
		err = cmd.Run(ctx)
	}
	if err != nil {
		Help(ctx, err)
	}
	return err
}
