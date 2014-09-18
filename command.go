//
// go.cli :: command.go
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
	"fmt"
	"sort"
	"strings"
)

type Command struct {
	Name   []string
	Usage  interface{}
	Desc   string
	Epilog string
	Flags  *FlagSet
	Action func(*Context) error
}

func (c *Command) Run(ctx *Context) error {
	if c.Flags != nil {
		ctx.Flags = NewFlagSet()
		for _, fs := range []*FlagSet{ctx.CLI.Flags, c.Flags} {
			fs.VisitAll(func(f *Flag) {
				ctx.Flags.Add(f)
			})
		}
		if err := ctx.Flags.Parse(ctx.Args); err != nil {
			return err
		}
		ctx.Args = ctx.Flags.Args()
	}
	if c.Action == nil {
		return c.action(ctx)
	}
	return c.Action(ctx)
}

func (c *Command) action(ctx *Context) error {
	switch {
	case ctx.CLI.help && ctx.Bool("help"):
		return Help(ctx, nil)
	case ctx.CLI.version && ctx.Bool("version"):
		return Version(ctx)
	}
	return nil
}

func FindCmd(cmds []*Command, name string) (cmd *Command, err error) {
	set := make(map[string]*Command)
L:
	for _, c := range cmds {
		// exact match
		for _, n := range c.Name {
			if n == name {
				set[n] = c
				continue L
			}
		}
		// prefix match
		if name != "" {
			for _, n := range c.Name {
				if strings.HasPrefix(n, name) {
					set[n] = c
					continue L
				}
			}
		}
	}

	switch len(set) {
	case 0:
		err = &CommandError{Name: name}
	case 1:
		for _, cmd = range set {
		}
	default:
		if c, ok := set[name]; ok {
			cmd = c
		} else {
			list := make([]string, len(set))
			i := 0
			for n := range set {
				list[i] = n
				i++
			}
			sort.Strings(list)
			err = &CommandError{name, list}
		}
	}
	return
}

type CommandError struct {
	Name string
	List []string
}

func (e *CommandError) Error() string {
	if len(e.List) == 0 {
		return fmt.Sprintf("unknown command '%s'", e.Name)
	}
	return fmt.Sprintf("command '%s' is ambiguous (%s)", e.Name, strings.Join(e.List, ", "))
}

type CommandSlice []*Command

func (p CommandSlice) Len() int           { return len(p) }
func (p CommandSlice) Less(i, j int) bool { return p[i].Name[0] < p[j].Name[0] }
func (p CommandSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p CommandSlice) Sort() { sort.Sort(p) }
