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

import "time"

type Context struct {
	CLI *CLI
}

func NewContext(cli *CLI) *Context {
	return &Context{cli}
}

func (c *Context) Arg(i int) string { return c.CLI.Flags.Arg(i) }

func (c *Context) Args() []string { return c.CLI.Flags.Args() }

func (c *Context) NArg() int { return c.CLI.Flags.NArg() }

func (c *Context) Value(name string) interface{} {
	if f := c.CLI.Flags.Lookup(name); f != nil {
		return f.Get()
	}
	return nil
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
