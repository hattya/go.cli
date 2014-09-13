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
	"strconv"
	"strings"
	"time"
)

type Flag struct {
	Name    []string
	Usage   string
	Value   flag.Getter
	MetaVar string
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
	f := &FlagSet{vars: make(map[string]*Flag)}
	f.fs.SetOutput(ioutil.Discard)
	return f
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
			for _, n := range f.Name {
				seen[n] = true
			}
			fn(f)
		}
	})
}

func (fs *FlagSet) VisitAll(fn func(*Flag)) {
	seen := make(map[string]bool)
	fs.fs.VisitAll(func(ff *flag.Flag) {
		if _, ok := seen[ff.Name]; !ok {
			f := fs.vars[ff.Name]
			for _, n := range f.Name {
				seen[n] = true
			}
			fn(f)
		}
	})
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

func (fs *FlagSet) Var(name string, value flag.Getter, usage string) {
	fs.each(name, func(n string) {
		fs.fs.Var(value, n, usage)
	})
}

func (fs *FlagSet) VarEnv(envVar, name string, value flag.Getter, usage string) {
	if s := os.Getenv(envVar); s != "" {
		value.Set(s)
	}
	fs.Var(name, value, usage)
}

func (fs *FlagSet) Bool(name string, value bool, usage string) {
	fs.each(name, func(n string) {
		fs.fs.Bool(n, value, usage)
	})
}

func (fs *FlagSet) BoolEnv(envVar, name string, value bool, usage string) {
	if s := os.Getenv(envVar); s != "" {
		if v, err := strconv.ParseBool(s); err == nil {
			value = v
		}
	}
	fs.Bool(name, value, usage)
}

func (fs *FlagSet) Duration(name string, value time.Duration, usage string) {
	fs.each(name, func(n string) {
		fs.fs.Duration(n, value, usage)
	})
}

func (fs *FlagSet) DurationEnv(envVar, name string, value time.Duration, usage string) {
	if s := os.Getenv(envVar); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			value = d
		}
	}
	fs.Duration(name, value, usage)
}

func (fs *FlagSet) Float64(name string, value float64, usage string) {
	fs.each(name, func(n string) {
		fs.fs.Float64(n, value, usage)
	})
}

func (fs *FlagSet) Float64Env(envVar, name string, value float64, usage string) {
	if s := os.Getenv(envVar); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			value = f
		}
	}
	fs.Float64(name, value, usage)
}

func (fs *FlagSet) Int(name string, value int, usage string) {
	fs.each(name, func(n string) {
		fs.fs.Int(n, value, usage)
	})
}

func (fs *FlagSet) IntEnv(envVar, name string, value int, usage string) {
	if s := os.Getenv(envVar); s != "" {
		if i, err := strconv.ParseInt(s, 10, 0); err == nil {
			value = int(i)
		}
	}
	fs.Int(name, value, usage)
}

func (fs *FlagSet) Int64(name string, value int64, usage string) {
	fs.each(name, func(n string) {
		fs.fs.Int64(n, value, usage)
	})
}

func (fs *FlagSet) Int64Env(envVar, name string, value int64, usage string) {
	if s := os.Getenv(envVar); s != "" {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			value = i
		}
	}
	fs.Int64(name, value, usage)
}

func (fs *FlagSet) String(name string, value string, usage string) {
	fs.each(name, func(n string) {
		fs.fs.String(n, value, usage)
	})
}

func (fs *FlagSet) StringEnv(envVar, name string, value string, usage string) {
	if s := os.Getenv(envVar); s != "" {
		value = s
	}
	fs.String(name, value, usage)
}

func (fs *FlagSet) Uint(name string, value uint, usage string) {
	fs.each(name, func(n string) {
		fs.fs.Uint(n, value, usage)
	})
}

func (fs *FlagSet) UintEnv(envVar, name string, value uint, usage string) {
	if s := os.Getenv(envVar); s != "" {
		if n, err := strconv.ParseUint(s, 10, 0); err == nil {
			value = uint(n)
		}
	}
	fs.Uint(name, value, usage)
}

func (fs *FlagSet) Uint64(name string, value uint64, usage string) {
	fs.each(name, func(n string) {
		fs.fs.Uint64(n, value, usage)
	})
}

func (fs *FlagSet) Uint64Env(envVar, name string, value uint64, usage string) {
	if s := os.Getenv(envVar); s != "" {
		if n, err := strconv.ParseUint(s, 10, 64); err == nil {
			value = n
		}
	}
	fs.Uint64(name, value, usage)
}

func (fs *FlagSet) each(name string, fn func(string)) {
	list := strings.Split(name, ",")
	for i := 0; i < len(list); i++ {
		n := strings.TrimSpace(list[i])
		list[i] = n
		fn(n)
	}

	sort.Sort(flagNames(list))
	f := &Flag{Name: list}
	for i, n := range list {
		ff := fs.fs.Lookup(n)
		if i == 0 {
			f.Usage = ff.Usage
			f.Value = ff.Value.(flag.Getter)
		} else {
			ff.Value = f.Value
		}
		fs.vars[n] = f
	}
	fs.list = append(fs.list, f)
}

type flagNames []string

func (p flagNames) Len() int           { return len(p) }
func (p flagNames) Less(i, j int) bool { return len(p[i]) < len(p[j]) || p[i] < p[j] }
func (p flagNames) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
