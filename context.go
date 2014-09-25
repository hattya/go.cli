//
// go.cli :: context.go
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
	"bytes"
	"time"
)

type Context struct {
	CLI   *CLI
	Stack []*Command

	Cmds  []*Command
	Flags *FlagSet
	Args  []string
	Data  interface{}
}

func NewContext(cli *CLI) *Context {
	return &Context{
		CLI:   cli,
		Cmds:  cli.Cmds,
		Flags: cli.Flags,
		Args:  cli.Flags.Args(),
	}
}

func (c *Context) Name() string {
	if 0 < len(c.Stack) {
		var b bytes.Buffer
		b.WriteString(c.CLI.Name)
		for _, cmd := range c.Stack {
			b.WriteRune(' ')
			b.WriteString(cmd.Name[0])
		}
		return b.String()
	}
	return c.CLI.Name
}

func (c *Context) Command() (cmd *Command, err error) {
	switch {
	case len(c.Cmds) == 0:
	case len(c.Args) == 0:
		err = ErrCommand
	default:
		cmd, err = FindCommand(c.Cmds, c.Args[0])
		if err == nil {
			c.Cmds = cmd.Cmds
			c.Args = c.Args[1:]
		}
	}
	return
}

func (c *Context) Bool(name string) bool {
	return c.Value(name).(bool)
}

func (c *Context) Duration(name string) time.Duration {
	return c.Value(name).(time.Duration)
}

func (c *Context) Float64(name string) float64 {
	return c.Value(name).(float64)
}

func (c *Context) Int(name string) int {
	return c.Value(name).(int)
}

func (c *Context) Int64(name string) int64 {
	return c.Value(name).(int64)
}

func (c *Context) String(name string) string {
	return c.Value(name).(string)
}

func (c *Context) Uint(name string) uint {
	return c.Value(name).(uint)
}

func (c *Context) Uint64(name string) uint64 {
	return c.Value(name).(uint64)
}

func (c *Context) Value(name string) interface{} {
	if f := c.Flags.Lookup(name); f != nil {
		return f.Value.Get()
	}
	return nil
}

func (c *Context) ErrorHandler(err error) error {
	return c.CLI.ErrorHandler(c, err)
}
