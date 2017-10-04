//
// go.cli :: action.go
//
//   Copyright (c) 2014-2017 Akinori Hattori <hattya@gmail.com>
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
	if 0 < len(ctx.Stack) {
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
		if 0 < len(ctx.Args) {
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
