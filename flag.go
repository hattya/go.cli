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
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

type Flag struct {
	Name    []string
	Usage   string
	Value   flag.Getter
	MetaVar string
	EnvVar  string
}

func (f *Flag) IsBool() bool {
	if b, ok := f.Value.(boolFlag); ok {
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

func (fs *FlagSet) Parse(args []string) error { return fs.error(fs.fs.Parse(args)) }

func (fs *FlagSet) Lookup(name string) *Flag { return fs.vars[name] }

func (fs *FlagSet) MetaVar(name, metaVar string) { fs.vars[name].MetaVar = metaVar }

func (fs *FlagSet) Set(name, value string) error { return fs.error(fs.fs.Set(name, value)) }

func (fs *FlagSet) error(err error) error {
	if err != nil {
		return FlagError(err.Error())
	}
	return nil
}

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

func (fs *FlagSet) Add(f *Flag) {
	sort.Sort(flagName(f.Name))
	for _, n := range f.Name {
		fs.fs.Var(f.Value, n, f.Usage)
		fs.vars[n] = f
	}
	fs.list = append(fs.list, f)

	if f.EnvVar != "" {
		if s := os.Getenv(f.EnvVar); s != "" {
			f.Value.Set(s)
		}
	}
}

func (fs *FlagSet) Bool(name string, value bool, usage string) *Flag {
	return fs.BoolEnv("", name, value, usage)
}

func (fs *FlagSet) BoolEnv(envVar, name string, value bool, usage string) *Flag {
	return fs.each(name, envVar, func(n string) {
		fs.fs.BoolVar(&value, n, value, usage)
	})
}

func (fs *FlagSet) Duration(name string, value time.Duration, usage string) *Flag {
	return fs.DurationEnv("", name, value, usage)
}

func (fs *FlagSet) DurationEnv(envVar, name string, value time.Duration, usage string) *Flag {
	return fs.each(name, envVar, func(n string) {
		fs.fs.DurationVar(&value, n, value, usage)
	})
}

func (fs *FlagSet) Float64(name string, value float64, usage string) *Flag {
	return fs.Float64Env("", name, value, usage)
}

func (fs *FlagSet) Float64Env(envVar, name string, value float64, usage string) *Flag {
	return fs.each(name, envVar, func(n string) {
		fs.fs.Float64Var(&value, n, value, usage)
	})
}

func (fs *FlagSet) Int(name string, value int, usage string) *Flag {
	return fs.IntEnv("", name, value, usage)
}

func (fs *FlagSet) IntEnv(envVar, name string, value int, usage string) *Flag {
	return fs.each(name, envVar, func(n string) {
		fs.fs.IntVar(&value, n, value, usage)
	})
}

func (fs *FlagSet) Int64(name string, value int64, usage string) *Flag {
	return fs.Int64Env("", name, value, usage)
}

func (fs *FlagSet) Int64Env(envVar, name string, value int64, usage string) *Flag {
	return fs.each(name, envVar, func(n string) {
		fs.fs.Int64Var(&value, n, value, usage)
	})
}

func (fs *FlagSet) String(name string, value string, usage string) *Flag {
	return fs.StringEnv("", name, value, usage)
}

func (fs *FlagSet) StringEnv(envVar, name string, value string, usage string) *Flag {
	return fs.each(name, envVar, func(n string) {
		fs.fs.StringVar(&value, n, value, usage)
	})
}

func (fs *FlagSet) Uint(name string, value uint, usage string) *Flag {
	return fs.UintEnv("", name, value, usage)
}

func (fs *FlagSet) UintEnv(envVar, name string, value uint, usage string) *Flag {
	return fs.each(name, envVar, func(n string) {
		fs.fs.UintVar(&value, n, value, usage)
	})
}

func (fs *FlagSet) Uint64(name string, value uint64, usage string) *Flag {
	return fs.Uint64Env("", name, value, usage)
}

func (fs *FlagSet) Uint64Env(envVar, name string, value uint64, usage string) *Flag {
	return fs.each(name, envVar, func(n string) {
		fs.fs.Uint64Var(&value, n, value, usage)
	})
}

func (fs *FlagSet) Var(name string, value flag.Getter, usage string) *Flag {
	return fs.VarEnv("", name, value, usage)
}

func (fs *FlagSet) VarEnv(envVar, name string, value flag.Getter, usage string) *Flag {
	return fs.each(name, envVar, func(n string) {
		fs.fs.Var(value, n, usage)
	})
}

func (fs *FlagSet) each(name, envVar string, fn func(string)) *Flag {
	list := strings.Split(name, ",")
	for i, s := range list {
		list[i] = strings.TrimSpace(s)
	}

	f := &Flag{
		Name:   list,
		EnvVar: envVar,
	}
	sort.Sort(flagName(f.Name))
	for _, n := range f.Name {
		fn(n)
		fs.vars[n] = f
	}
	fs.list = append(fs.list, f)

	ff := fs.fs.Lookup(f.Name[0])
	f.Usage = ff.Usage
	f.Value = ff.Value.(flag.Getter)
	if f.EnvVar != "" {
		if s := os.Getenv(f.EnvVar); s != "" {
			f.Value.Set(s)
		}
	}
	return f
}

type flagName []string

func (p flagName) Len() int           { return len(p) }
func (p flagName) Less(i, j int) bool { return len(p[i]) < len(p[j]) || p[i] < p[j] }
func (p flagName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type FlagError string

func (e FlagError) Error() string { return string(e) }
