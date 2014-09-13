//
// go.cli :: flag.go
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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

type Flag struct {
	Name    []string
	Usage   string
	Value   interface{}
	MetaVar string
	EnvVar  string

	value flag.Getter
}

func (f *Flag) Get() interface{} {
	return f.value.Get()
}

func (f *Flag) IsBool() bool {
	if b, ok := f.value.(boolFlag); ok {
		return b.IsBoolFlag()
	}
	return false
}

func (f *Flag) Format(sep string) string {
	var b bytes.Buffer
	for i, n := range f.Name {
		if 0 < i {
			b.WriteString(", ")
		}
		if len(n) == 1 {
			b.WriteRune('-')
		} else {
			b.WriteString("--")
		}
		b.WriteString(n)
	}
	b.WriteString(MetaVar(f))
	if f.Usage != "" {
		b.WriteString(sep)
		b.WriteString(f.Usage)
	}
	return b.String()
}

type boolFlag interface {
	flag.Getter
	IsBoolFlag() bool
}

type FlagSet struct {
	fs   flag.FlagSet
	vars map[string]*Flag
	list []*Flag
}

func NewFlagSet() *FlagSet {
	fs := &FlagSet{vars: make(map[string]*Flag)}
	fs.fs.SetOutput(ioutil.Discard)
	return fs
}

func (fs *FlagSet) Parse(args []string) error { return fs.fs.Parse(args) }

func (fs *FlagSet) Lookup(name string) *Flag { return fs.vars[name] }

func (fs *FlagSet) MetaVar(name, metaVar string) { fs.vars[name].MetaVar = metaVar }

func (fs *FlagSet) Set(name, value string) error { return fs.fs.Set(name, value) }

func (fs *FlagSet) Visit(fn func(*Flag)) {
	seen := make(map[string]bool)
	fs.fs.Visit(func(ff *flag.Flag) {
		if _, ok := seen[ff.Name]; !ok {
			f := fs.vars[ff.Name]
			fn(f)
			for _, n := range f.Name {
				seen[n] = true
			}
		}
	})
}

func (fs *FlagSet) VisitAll(fn func(*Flag)) {
	list := make(sort.StringSlice, len(fs.list))
	for i, f := range fs.list {
		list[i] = f.Name[0]
	}
	sort.Sort(list)
	flags := make([]*Flag, len(list))
	for i, n := range list {
		flags[i] = fs.vars[n]
	}

	for _, f := range flags {
		fn(f)
	}
}

func (fs *FlagSet) NFlag() int {
	n := 0
	fs.Visit(func(*Flag) {
		n++
	})
	return n
}

func (fs *FlagSet) Arg(i int) string { return fs.fs.Arg(i) }

func (fs *FlagSet) Args() []string { return fs.fs.Args() }

func (fs *FlagSet) NArg() int { return fs.fs.NArg() }

func (fs *FlagSet) Parsed() bool { return fs.fs.Parsed() }

func (fs *FlagSet) Bool(name string, value bool, usage string) *Flag {
	return fs.Var(name, value, usage)
}

func (fs *FlagSet) BoolEnv(envVar, name string, value bool, usage string) *Flag {
	return fs.VarEnv(envVar, name, value, usage)
}

func (fs *FlagSet) Duration(name string, value time.Duration, usage string) *Flag {
	return fs.Var(name, value, usage)
}

func (fs *FlagSet) DurationEnv(envVar, name string, value time.Duration, usage string) *Flag {
	return fs.VarEnv(envVar, name, value, usage)
}

func (fs *FlagSet) Float64(name string, value float64, usage string) *Flag {
	return fs.Var(name, value, usage)
}

func (fs *FlagSet) Float64Env(envVar, name string, value float64, usage string) *Flag {
	return fs.VarEnv(envVar, name, value, usage)
}

func (fs *FlagSet) Int(name string, value int, usage string) *Flag {
	return fs.Var(name, value, usage)
}

func (fs *FlagSet) IntEnv(envVar, name string, value int, usage string) *Flag {
	return fs.VarEnv(envVar, name, value, usage)
}

func (fs *FlagSet) Int64(name string, value int64, usage string) *Flag {
	return fs.Var(name, value, usage)
}

func (fs *FlagSet) Int64Env(envVar, name string, value int64, usage string) *Flag {
	return fs.VarEnv(envVar, name, value, usage)
}

func (fs *FlagSet) String(name string, value string, usage string) *Flag {
	return fs.Var(name, value, usage)
}

func (fs *FlagSet) StringEnv(envVar, name string, value string, usage string) *Flag {
	return fs.VarEnv(envVar, name, value, usage)
}

func (fs *FlagSet) Uint(name string, value uint, usage string) *Flag {
	return fs.Var(name, value, usage)
}

func (fs *FlagSet) UintEnv(envVar, name string, value uint, usage string) *Flag {
	return fs.VarEnv(envVar, name, value, usage)
}

func (fs *FlagSet) Uint64(name string, value uint64, usage string) *Flag {
	return fs.Var(name, value, usage)
}

func (fs *FlagSet) Uint64Env(envVar, name string, value uint64, usage string) *Flag {
	return fs.VarEnv(envVar, name, value, usage)
}

func (fs *FlagSet) Var(name string, value interface{}, usage string) *Flag {
	return fs.VarEnv("", name, value, usage)
}

func (fs *FlagSet) VarEnv(envVar, name string, value interface{}, usage string) *Flag {
	list := strings.Split(name, ",")
	for i, s := range list {
		list[i] = strings.TrimSpace(s)
	}
	sort.Sort(flagName(list))

	f := &Flag{
		Name:   list,
		Usage:  usage,
		Value:  value,
		EnvVar: envVar,
	}
	fs.Add(f)
	return f
}

func (fs *FlagSet) Add(f *Flag) {
	switch v := f.Value.(type) {
	case bool:
		fs.each(f, func(n string) {
			fs.fs.BoolVar(&v, n, v, f.Usage)
		})
	case time.Duration:
		fs.each(f, func(n string) {
			fs.fs.DurationVar(&v, n, v, f.Usage)
		})
	case float64:
		fs.each(f, func(n string) {
			fs.fs.Float64Var(&v, n, v, f.Usage)
		})
	case int:
		fs.each(f, func(n string) {
			fs.fs.IntVar(&v, n, v, f.Usage)
		})
	case int64:
		fs.each(f, func(n string) {
			fs.fs.Int64Var(&v, n, v, f.Usage)
		})
	case string:
		fs.each(f, func(n string) {
			fs.fs.StringVar(&v, n, v, f.Usage)
		})
	case uint:
		fs.each(f, func(n string) {
			fs.fs.UintVar(&v, n, v, f.Usage)
		})
	case uint64:
		fs.each(f, func(n string) {
			fs.fs.Uint64Var(&v, n, v, f.Usage)
		})
	case flag.Getter:
		fs.each(f, func(n string) {
			fs.fs.Var(v, n, f.Usage)
		})
	default:
		panic(fmt.Sprintf("unknown type '%T'", v))
	}
}

func (fs *FlagSet) each(f *Flag, fn func(string)) {
	for _, n := range f.Name {
		fn(n)
		fs.vars[n] = f
	}
	fs.list = append(fs.list, f)

	f.value = fs.fs.Lookup(f.Name[0]).Value.(flag.Getter)
	if f.EnvVar != "" {
		if s := os.Getenv(f.EnvVar); s != "" {
			f.value.Set(s)
		}
	}
}

type flagName []string

func (p flagName) Len() int           { return len(p) }
func (p flagName) Less(i, j int) bool { return len(p[i]) < len(p[j]) || p[i] < p[j] }
func (p flagName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
