//
// go.cli :: help.go
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
	"fmt"
	"strings"
	"text/tabwriter"
	"text/template"
)

func NewHelpCommand() *Command {
	return &Command{
		Name:  []string{"help"},
		Usage: "[<command>]",
		Desc:  "show help for a specified command",
		Flags: NewFlagSet(),
		Action: func(ctx *Context) error {
			ctx.Stack = nil
			ctx.Cmds = ctx.CLI.Cmds
			for 0 < len(ctx.Args) {
				cmd, err := ctx.Command()
				switch {
				case err != nil:
					return &Abort{
						Err:  err,
						Hint: fmt.Sprintf("type '%v help' for usage", ctx.Name()),
					}
				case cmd == nil:
					return ErrArgs
				}
				ctx.Stack = append(ctx.Stack, cmd)
			}
			return Help(ctx)
		},
	}
}

var (
	Help    = ShowHelp
	Usage   = FormatUsage
	MetaVar = FormatMetaVar
)

func ShowHelp(ctx *Context) error {
	fm := template.FuncMap{
		"usage":  Usage,
		"cmd":    cmd,
		"cmds":   cmds,
		"format": format,
		"flags":  flags,
	}
	t := template.Must(template.New("help").Funcs(fm).Parse(helpTmpl))
	w := tabwriter.NewWriter(ctx.CLI.Stdout, 0, 8, 4, ' ', 0)
	defer w.Flush()
	return t.Execute(w, ctx)
}

const helpTmpl = `{{range usage .}}{{.}}
{{end}}{{with or (cmd .) .CLI}}{{if .Desc}}
{{.Desc}}
{{end}}{{range $i, $cmd := cmds .Cmds}}{{if eq $i 0}}
commands:

{{end}}  {{format $cmd "\t"}}
{{end}}{{$flags := flags .Flags}}{{range $i, $f := $flags}}{{if eq $i 0 }}
options:

{{end}}  {{$f.Format "\t"}}
{{end}}{{if .Epilog}}
{{.Epilog}}
{{else if or .Desc (lt 0 (len .Cmds)) (lt 0 (len $flags))}}
{{end}}{{end}}`

func cmd(ctx *Context) *Command {
	if len(ctx.Stack) == 0 {
		return nil
	}
	return ctx.Stack[len(ctx.Stack)-1]
}

func cmds(cmds []*Command) []*Command {
	list := make(CommandSlice, len(cmds))
	copy(list, cmds)
	list.Sort()
	return list
}

func format(cmd *Command, sep string) string {
	var b bytes.Buffer
	b.WriteString(cmd.Name[0])
	if cmd.Desc != "" {
		b.WriteString(sep)
		b.WriteString(strings.TrimSpace(strings.Split(cmd.Desc, "\n")[0]))
	}
	return b.String()
}

func flags(fs *FlagSet) []*Flag {
	var flags []*Flag
	if fs != nil {
		fs.VisitAll(func(f *Flag) {
			flags = append(flags, f)
		})
	}
	return flags
}

func FormatUsage(ctx *Context) []string {
	var cmd *Command
	var u interface{}
	if 0 < len(ctx.Stack) {
		cmd = ctx.Stack[len(ctx.Stack)-1]
		u = cmd.Usage
	} else {
		u = ctx.CLI.Usage
	}
	var usage []string
	switch v := u.(type) {
	case nil:
		usage = []string{""}
	case string:
		usage = []string{v}
	case []string:
		usage = make([]string, len(v))
		copy(usage, v)
	default:
		panic(fmt.Sprintf("unknown type '%T'", v))
	}

	var b bytes.Buffer
	for i, s := range usage {
		if i == 0 {
			b.WriteString("usage: ")
		} else {
			b.WriteString("   or: ")
		}
		b.WriteString(ctx.Name())
		if s != "" {
			b.WriteRune(' ')
			b.WriteString(s)
		}
		usage[i] = b.String()
		b.Reset()
	}
	if cmd != nil && 1 < len(cmd.Name) {
		usage = append(usage, "", "alias: "+strings.Join(cmd.Name[1:], ", "))
	}
	return usage
}

func FormatMetaVar(f *Flag) string {
	if f.IsBool() || f.MetaVar != "" {
		return f.MetaVar
	}
	s := f.Name[0]
	for _, n := range f.Name {
		if 1 < len(n) {
			s = n
			break
		}
	}
	return fmt.Sprintf(" <%v>", s)
}
