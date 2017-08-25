//
// go.cli :: flag_test.go
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

func TestFlagSet(t *testing.T) {
	envVar := func(s string) string {
		return fmt.Sprintf("__CLI_%v__", strings.ToUpper(s))
	}

	values := map[string]interface{}{
		"bool":     true,
		"duration": 1 * time.Millisecond,
		"float64":  3.14,
		"int":      -1,
		"int64":    int64(-64),
		"string":   "string",
		"uint":     uint(1),
		"uint64":   uint64(64),
		"var":      "var",
	}
	environ := map[string]string{
		envVar("bool"):     "true",
		envVar("duration"): "1ms",
		envVar("float64"):  "3.14",
		envVar("int"):      "-1",
		envVar("int64"):    "-64",
		envVar("string"):   "string",
		envVar("uint"):     "1",
		envVar("uint64"):   "64",
		envVar("var"):      "var",
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
	flags.BoolEnv(envVar("bool"), "bool", false, "")
	flags.DurationEnv(envVar("duration"), "duration", 0, "")
	flags.Float64Env(envVar("float64"), "float64", 0.0, "")
	flags.IntEnv(envVar("int"), "int", 0, "")
	flags.Int64Env(envVar("int64"), "int64", 0, "")
	flags.StringEnv(envVar("string"), "string", "", "")
	flags.UintEnv(envVar("uint"), "uint", 0, "")
	flags.Uint64Env(envVar("uint64"), "uint64", 0, "")
	flags.VarEnv(envVar("var"), "var", &value{}, "")
	for k, v := range values {
		if g, e := flags.Lookup(k).Value.Get(), v; g != e {
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

	if err := flags.MetaVar("", ""); err == nil {
		t.Error("expected error")
	}

	flags.Set("var", "set")
	if g, e := flags.Lookup("var").Value.Get(), "set"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
}

func TestAddFlags(t *testing.T) {
	f := &cli.Flag{
		Name:   []string{"var"},
		Value:  &value{},
		EnvVar: "__CLI_VAR__",
	}

	os.Setenv(f.EnvVar, "var")
	defer os.Setenv(f.EnvVar, "")

	flags := cli.NewFlagSet()
	flags.Add(f)
	if g, e := f.Value.Get(), "var"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	i := 0
	flags.VisitAll(func(*cli.Flag) {
		i++
	})
	if g, e := i, 1; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
}

func TestFormatFlags(t *testing.T) {
	usage := "usage"

	flags := cli.NewFlagSet()
	flags.Bool("bool", true, usage)
	flags.Duration("duration", 1*time.Millisecond, usage)
	flags.Float64("float64", 3.14, usage)
	flags.Int("int", -1, usage)
	flags.String("string", "string", usage)
	flags.Uint("uint", 1, usage)

	flags.VisitAll(func(f *cli.Flag) {
		s := "--%v"
		if f.Name[0] != "bool" {
			s += " <%[1]v>"
		}
		s += "\t%v"

		if g, e := f.Format("\t"), fmt.Sprintf(s, f.Name[0], usage); g != e {
			t.Errorf("expected %q, got %q", e, g)
		}

		verb := " %%"
		f.Usage = usage + verb
		if g, e := f.Format("\t"), fmt.Sprintf(s+verb, f.Name[0], usage); g != e {
			t.Errorf("expected %q, got %q", e, g)
		}

		verb = " (default: %v)"
		f.Usage = usage + verb
		if g, e := f.Format("\t"), fmt.Sprintf(s+verb, f.Name[0], usage, f.Default); g != e {
			t.Errorf("expected %q, got %q", e, g)
		}
	})
}

func TestResetFlags(t *testing.T) {
	envVar := func(s string) string {
		return fmt.Sprintf("__CLI_%v__", strings.ToUpper(s))
	}

	os.Setenv(envVar("int"), "-2")
	defer os.Setenv(envVar("int"), "")

	flags := cli.NewFlagSet()
	flags.Bool("h, help", false, "")
	flags.IntEnv(envVar("int"), "int", 0, "")
	flags.UintEnv(envVar("uint"), "uint", 0, "")

	flags.Set("h", "true")
	flags.Reset()
	if g, e := flags.Parsed(), false; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := flags.Get("h"), false; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := flags.Get("int"), -2; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := flags.Get("uint"), uint(0); g != e {
		t.Errorf("expected %v, got %v", e, g)
	}

	args := []string{"-h", "-int", "-1", "-uint", "1"}
	if err := flags.Parse(args); err != nil {
		t.Fatal(err)
	}
	flags.Reset()
	if g, e := flags.Parsed(), false; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := flags.Get("h"), false; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := flags.Get("int"), -2; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := flags.Get("uint"), uint(0); g != e {
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
	if err := testStrings(func(i int) string { return f.Name[i] }, []string{"h", "help"}); err != nil {
		t.Error(err)
	}
}

func TestChoiceFlag(t *testing.T) {
	flags := cli.NewFlagSet()
	flags.Choice("c, choice", 0, map[string]interface{}{
		"foo":    1,
		"bar":    2,
		"baz":    3,
		"foobar": 4,
	}, "")

	args := []string{"-c", "foo"}
	if err := flags.Parse(args); err != nil {
		t.Fatal(err)
	}
	if g, e := flags.Get("c"), 1; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}

	flags.Reset()
	args = []string{"-c", ""}
	switch err := flags.Parse(args); {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), `choose from "bar", "baz", "foo" or "foobar"`):
		t.Error("unexpected error:", err)
	}

	flags.Reset()
	args = []string{"-c", "b"}
	switch err := flags.Parse(args); {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), `choose from "bar", "baz", "foo" or "foobar"`):
		t.Error("unexpected error:", err)
	}
}

func TestPrefixChoiceFlag(t *testing.T) {
	flags := cli.NewFlagSet()
	flags.PrefixChoice("c, choice", 0, map[string]interface{}{
		"foo":    1,
		"bar":    2,
		"baz":    3,
		"foobar": 4,
	}, "")

	args := []string{"-c", "foo"}
	if err := flags.Parse(args); err != nil {
		t.Fatal(err)
	}
	if g, e := flags.Get("c"), 1; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}

	flags.Reset()
	args = []string{"-c", ""}
	switch err := flags.Parse(args); {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), `choose from "bar", "baz", "foo" or "foobar"`):
		t.Error("unexpected error:", err)
	}

	flags.Reset()
	args = []string{"-c", "b"}
	switch err := flags.Parse(args); {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), `choose from "bar" or "baz"`):
		t.Error("unexpected error:", err)
	}
}
