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

func NewContext(ui *CLI) *Context {
	return &Context{
		CLI:   ui,
		Cmds:  ui.Cmds,
		Flags: ui.Flags,
		Args:  ui.Flags.Args(),
	}
}

func (ctx *Context) Name() string {
	if 0 < len(ctx.Stack) {
		var b bytes.Buffer
		b.WriteString(ctx.CLI.Name)
		for _, cmd := range ctx.Stack {
			b.WriteRune(' ')
			b.WriteString(cmd.Name[0])
		}
		return b.String()
	}
	return ctx.CLI.Name
}

func (ctx *Context) Command() (cmd *Command, err error) {
	switch {
	case len(ctx.Cmds) == 0:
	case len(ctx.Args) == 0:
		err = ErrCommand
	default:
		cmd, err = FindCommand(ctx.Cmds, ctx.Args[0])
		if err == nil {
			ctx.Cmds = cmd.Cmds
			ctx.Args = ctx.Args[1:]
		}
	}
	return
}

func (ctx *Context) Bool(name string) bool {
	return ctx.Flags.Get(name).(bool)
}

func (ctx *Context) Duration(name string) time.Duration {
	return ctx.Flags.Get(name).(time.Duration)
}

func (ctx *Context) Float64(name string) float64 {
	return ctx.Flags.Get(name).(float64)
}

func (ctx *Context) Int(name string) int {
	return ctx.Flags.Get(name).(int)
}

func (ctx *Context) Int64(name string) int64 {
	return ctx.Flags.Get(name).(int64)
}

func (ctx *Context) String(name string) string {
	return ctx.Flags.Get(name).(string)
}

func (ctx *Context) Uint(name string) uint {
	return ctx.Flags.Get(name).(uint)
}

func (ctx *Context) Uint64(name string) uint64 {
	return ctx.Flags.Get(name).(uint64)
}

func (ctx *Context) Value(name string) interface{} {
	return ctx.Flags.Get(name)
}

func (ctx *Context) Prepare(cmd *Command) error {
	if ctx.CLI.Prepare != nil {
		return ctx.CLI.Prepare(ctx, cmd)
	}
	return nil
}

func (ctx *Context) ErrorHandler(err error) error {
	return ctx.CLI.ErrorHandler(ctx, err)
}
