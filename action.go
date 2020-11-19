//
// go.cli :: action.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package cli

type Action func(*Context) error

var DefaultAction = Subcommand

func Subcommand(ctx *Context) error {
	cmd, err := ctx.Command()
	if cmd != nil {
		if err = ctx.Prepare(cmd); err == nil {
			ctx.Stack = append(ctx.Stack, cmd)
			err = cmd.Run(ctx)
		}
	}
	return err
}

func Chain(ctx *Context) error {
	if len(ctx.Stack) > 0 {
		return nil
	}

	for {
		select {
		case <-ctx.Context().Done():
			return ctx.Context().Err()
		default:
		}

		cmd, err := ctx.Command()
		if cmd != nil {
			if err = ctx.Prepare(cmd); err == nil {
				ctx.Cmds = ctx.UI.Cmds
				if len(ctx.Stack) == 0 {
					ctx.Stack = []*Command{cmd}
				} else {
					ctx.Stack[0] = cmd
				}
				err = cmd.Run(ctx)
			}
		}
		switch {
		case err != nil:
			return err
		case len(ctx.Args) == 0:
			return nil
		}
	}
}

func Option(action Action) Action {
	return func(ctx *Context) error {
		if len(ctx.Args) > 0 {
			return DefaultAction(ctx)
		}
		err := ctx.Prepare(nil)
		if err == nil {
			err = action(ctx)
		}
		return err
	}
}

func Simple(action Action) Action {
	return func(ctx *Context) error {
		err := ctx.Prepare(nil)
		if err == nil {
			err = action(ctx)
		}
		return err
	}
}
