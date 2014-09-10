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
	"fmt"
	"text/tabwriter"
	"text/template"
)

const helpTmpl = `usage: {{.Name}}{{if .Usage}} {{.Usage}}{{end}}
{{range $i, $flag := flags .Flags}}{{if eq $i 0 }}
options:
{{end}}
  {{$flag.Format "\t"}}{{end}}
`

var Help = PrintHelp

func PrintHelp(ctx *Context) error {
	t := template.New("help")
	t.Funcs(template.FuncMap{
		"flags": flags,
	})
	template.Must(t.Parse(helpTmpl))
	w := tabwriter.NewWriter(ctx.CLI.Stdout, 0, 8, 4, ' ', 0)
	defer w.Flush()
	return t.Execute(w, ctx.CLI)
}

func flags(fs *FlagSet) []*Flag {
	var flags []*Flag
	fs.VisitAll(func(f *Flag) {
		flags = append(flags, f)
	})
	return flags
}

var MetaVar = func(f *Flag) string {
	if f.MetaVar != "" {
		return f.MetaVar
	}
	if b, ok := f.Value.(boolFlag); ok && b.IsBoolFlag() {
		return ""
	}
	s := f.Name[0]
	for _, n := range f.Name {
		if 1 < len(n) {
			s = n
			break
		}
	}
	return fmt.Sprintf(" <%s>", s)
}
