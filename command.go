//
// go.cli :: command.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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
	Cmds   []*Command
	Flags  *FlagSet
	Action func(*Context) error
	Data   interface{}
}

func (c *Command) Run(ctx *Context) error {
	if c.Flags != nil {
		ctx.Flags = NewFlagSet()
		ctx.UI.Flags.VisitAll(ctx.Flags.Add)
		for _, cmd := range ctx.Stack {
			if cmd.Flags == nil {
				panic(ErrFlags)
			}
			cmd.Flags.VisitAll(ctx.Flags.Add)
		}
		if err := ctx.Flags.Parse(ctx.Args); err != nil {
			return err
		}
		ctx.Args = ctx.Flags.Args()
		switch {
		case ctx.UI.help && ctx.Bool("help"):
			return Help(ctx)
		case ctx.UI.version && ctx.Bool("version"):
			return Version(ctx)
		}
	}
	if c.Action == nil {
		return ctx.UI.Action(ctx)
	}
	return c.Action(ctx)
}

func (c *Command) Add(cmd *Command) {
	c.Cmds = append(c.Cmds, cmd)
}

func FindCommand(cmds []*Command, name string) (cmd *Command, err error) {
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
		err = CommandError{Name: name}
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
			err = CommandError{
				Name: name,
				List: list,
			}
		}
	}
	return
}

type CommandError struct {
	Name string
	List []string
}

func (e CommandError) Error() string {
	if len(e.List) == 0 {
		return fmt.Sprintf("unknown command '%v'", e.Name)
	}
	return fmt.Sprintf("command '%v' is ambiguous (%v)", e.Name, strings.Join(e.List, ", "))
}

type CommandSlice []*Command

func (p CommandSlice) Len() int           { return len(p) }
func (p CommandSlice) Less(i, j int) bool { return p[i].Name[0] < p[j].Name[0] }
func (p CommandSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p CommandSlice) Sort() { sort.Sort(p) }
