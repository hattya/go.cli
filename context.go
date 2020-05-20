//
// go.cli :: context.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package cli

import (
	"bytes"
	"context"
	"time"
)

type Context struct {
	UI    *CLI
	Stack []*Command

	Cmds  []*Command
	Flags *FlagSet
	Args  []string
	Data  interface{}
}

func NewContext(ui *CLI) *Context {
	return &Context{
		UI:    ui,
		Cmds:  ui.Cmds,
		Flags: ui.Flags,
		Args:  ui.Flags.Args(),
	}
}

func (ctx *Context) Name() string {
	if 0 < len(ctx.Stack) {
		var b bytes.Buffer
		b.WriteString(ctx.UI.Name)
		for _, cmd := range ctx.Stack {
			b.WriteRune(' ')
			b.WriteString(cmd.Name[0])
		}
		return b.String()
	}
	return ctx.UI.Name
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
	if ctx.UI.Prepare != nil {
		return ctx.UI.Prepare(ctx, cmd)
	}
	return nil
}

func (ctx *Context) ErrorHandler(err error) error {
	return ctx.UI.ErrorHandler(ctx, err)
}

func (ctx *Context) Context() context.Context {
	return ctx.UI.Context()
}

func (ctx *Context) Interrupt() {
	ctx.UI.Interrupt()
}
