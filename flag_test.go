//
// go.cli :: flag_test.go
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

package cli_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hattya/go.cli"
)

type flagVar struct {
	s string
}

func (f *flagVar) Set(v string) error {
	f.s = v
	return nil
}

func (f *flagVar) Get() interface{} { return f.s }
func (f *flagVar) String() string   { return fmt.Sprintf("%s", f.s) }

func TestFlagSet(t *testing.T) {
	envVar := func(s string) string {
		return fmt.Sprintf("__CLI_%s__", strings.ToUpper(s))
	}
	get := func(flags *cli.FlagSet, name string) interface{} {
		return flags.Lookup(name).Value.Get()
	}

	values := map[string]interface{}{
		"var":      "var",
		"bool":     true,
		"duration": 1 * time.Millisecond,
		"float64":  3.14,
		"int":      -1,
		"int64":    int64(-64),
		"string":   "string",
		"uint":     uint(1),
		"uint64":   uint64(64),
	}
	environ := map[string]string{
		envVar("var"):      "var",
		envVar("bool"):     "true",
		envVar("duration"): "1ms",
		envVar("float64"):  "3.14",
		envVar("int"):      "-1",
		envVar("int64"):    "-64",
		envVar("string"):   "string",
		envVar("uint"):     "1",
		envVar("uint64"):   "64",
	}
	for k, v := range environ {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range environ {
			os.Setenv(k, "")
		}
	}()

	flags := cli.NewFlagSet()
	flags.VarEnv(envVar("var"), "var", &flagVar{}, "")
	flags.BoolEnv(envVar("bool"), "bool", false, "")
	flags.DurationEnv(envVar("duration"), "duration", 0, "")
	flags.Float64Env(envVar("float64"), "float64", 0.0, "")
	flags.IntEnv(envVar("int"), "int", 0, "")
	flags.Int64Env(envVar("int64"), "int64", 0, "")
	flags.StringEnv(envVar("string"), "string", "", "")
	flags.UintEnv(envVar("uint"), "uint", 0, "")
	flags.Uint64Env(envVar("uint64"), "uint64", 0, "")
	for k, v := range values {
		if g, e := get(flags, k), v; g != e {
			t.Errorf("expected %v, got %v", e, g)
		}
	}

	args := []string{"0", "1"}
	if err := flags.Parse(args); err != nil {
		t.Fatal(err)
	}
	if g, e := flags.Parsed(), true; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := flags.NFlag(), 0; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	for i := 0; i < flags.NArg(); i++ {
		if g, e := flags.Arg(i), strconv.FormatInt(int64(i), 10); g != e {
			t.Errorf("expected %v, got %v", e, g)
		}
	}
	if g, e := len(flags.Args()), len(args); g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	i := 0
	flags.Visit(func(*cli.Flag) {
		i++
	})
	if g, e := i, 0; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	i = 0
	flags.VisitAll(func(*cli.Flag) {
		i++
	})
	if g, e := i, len(environ); g != e {
		t.Errorf("expected %v, got %v", e, g)
	}

	flags.Set("var", "set")
	if g, e := get(flags, "var"), "set"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
}

func TestVisitFlags(t *testing.T) {
	flags := cli.NewFlagSet()
	flags.Bool("h, help", false, "")
	args := []string{"-h"}
	if err := flags.Parse(args); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"h", "help"} {
		if g, e := flags.Lookup(s).Value.Get().(bool), true; g != e {
			t.Errorf("expected %v, got %v", e, g)
		}
	}
	if g, e := flags.NFlag(), 1; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	i := 0
	flags.Visit(func(*cli.Flag) {
		i++
	})
	if g, e := i, 1; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	i = 0
	flags.VisitAll(func(*cli.Flag) {
		i++
	})
	if g, e := i, 1; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
}

func TestSortFlags(t *testing.T) {
	flags := cli.NewFlagSet()
	flags.Bool("help, h", false, "")
	f := flags.Lookup("help")
	if g, e := f.Name[0], "h"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := f.Name[1], "help"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
}
