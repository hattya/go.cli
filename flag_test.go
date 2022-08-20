//
// go.cli :: flag_test.go
//
//   Copyright (c) 2014-2022 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package cli_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hattya/go.cli"
)

func TestFlagSet(t *testing.T) {
	envVar := func(s string) string { return fmt.Sprintf("__CLI_%v__", strings.ToUpper(s)) }
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
		t.Setenv(k, v)
	}

	flags := cli.NewFlagSet()
	flags.BoolEnv(envVar("bool"), "bool", false, "")
	flags.DurationEnv(envVar("duration"), "duration", 0, "")
	flags.Float64Env(envVar("float64"), "float64", 0.0, "")
	flags.IntEnv(envVar("int"), "int", 0, "")
	flags.Int64Env(envVar("int64"), "int64", 0, "")
	flags.StringEnv(envVar("string"), "string", "", "")
	flags.UintEnv(envVar("uint"), "uint", 0, "")
	flags.Uint64Env(envVar("uint64"), "uint64", 0, "")
	flags.VarEnv(envVar("var"), "var", new(value), "")
	for n, v := range values {
		if g, e := flags.Lookup(n).Value.Get(), v; g != e {
			t.Errorf("FlagSet.Lookup(%q).Value.Get() = %v, expected %v", n, g, e)
		}
	}

	if err := flags.Parse([]string{"0", "1"}); err != nil {
		t.Fatal(err)
	}
	if !flags.Parsed() {
		t.Error("FlagSet.Parsed() = false, expected true")
	}
	if g, e := flags.NFlag(), 0; g != e {
		t.Errorf("FlagSet.NFlag() = %v, expected %v", g, e)
	}
	if g, e := flags.NArg(), 2; g != e {
		t.Errorf("len(FlagSet.NArg()) = %v, expected %v", g, e)
	}
	for i := 0; i < flags.NArg(); i++ {
		if g, e := flags.Arg(i), strconv.FormatInt(int64(i), 10); g != e {
			t.Errorf("FlagSet.Arg(%v) = %v, expected %v", i, g, e)
		}
	}
	i := 0
	flags.Visit(func(*cli.Flag) { i++ })
	if g, e := i, 0; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	i = 0
	flags.VisitAll(func(*cli.Flag) { i++ })
	if g, e := i, len(environ); g != e {
		t.Errorf("expected %v, got %v", e, g)
	}

	if err := flags.MetaVar("", ""); err == nil {
		t.Error("expected error")
	}

	n := "var"
	v := "set"
	flags.Set(n, v)
	if g, e := flags.Lookup(n).Value.Get(), v; g != e {
		t.Errorf("FlagSet.Lookup(%q).Value.Get() = %q, expected %q", n, g, e)
	}
}

func TestAddFlags(t *testing.T) {
	f := &cli.Flag{
		Name:   []string{"var"},
		Value:  new(value),
		EnvVar: "__CLI_VAR__",
	}

	t.Setenv(f.EnvVar, "var")

	flags := cli.NewFlagSet()
	flags.Add(f)
	if g, e := f.Value.Get(), "var"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
	i := 0
	flags.VisitAll(func(*cli.Flag) { i++ })
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

		verb = "\n(default: %v)"
		f.Usage = usage + verb
		if g, e := f.Format("\t"), fmt.Sprintf(strings.Replace(s+verb, "\n", "\n\t", -1), f.Name[0], usage, f.Default); g != e {
			t.Errorf("expected %q, got %q", e, g)
		}
	})
}

func TestResetFlags(t *testing.T) {
	envVar := func(s string) string {
		return fmt.Sprintf("__CLI_%v__", strings.ToUpper(s))
	}

	t.Setenv(envVar("int"), "-2")

	type test struct {
		name  string
		value interface{}
	}

	flags := cli.NewFlagSet()
	flags.Bool("h, help", false, "")
	flags.IntEnv(envVar("int"), "int", -1, "")
	flags.UintEnv(envVar("uint"), "uint", 1, "")

	flags.Set("h", "true")
	flags.Reset()
	if flags.Parsed() {
		t.Error("FlagSet.Parsed() = true, expected false")
	}
	for _, tt := range []test{
		{"h", false},
		{"int", -2},
		{"uint", uint(1)},
	} {
		if g, e := flags.Get(tt.name), tt.value; g != e {
			t.Errorf("FlagSet.Get(%q) = %v, expected %v", tt.name, g, e)
		}
	}

	if err := flags.Parse(strings.Fields("-h -int 0 -uint 0")); err != nil {
		t.Fatal(err)
	}
	flags.Reset()
	if flags.Parsed() {
		t.Error("FlagSet.Parsed() = true, expected false")
	}
	for _, tt := range []test{
		{"h", false},
		{"int", -2},
		{"uint", uint(1)},
	} {
		if g, e := flags.Get(tt.name), tt.value; g != e {
			t.Errorf("FlagSet.Get(%q) = %v, expected %v", tt.name, g, e)
		}
	}
}

func TestVisitFlags(t *testing.T) {
	flags := cli.NewFlagSet()
	flags.Bool("h, help", false, "")
	flags.Bool("version", false, "")
	if err := flags.Parse([]string{"-h"}); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name  string
		value bool
	}{
		{"h", true},
		{"help", true},
		{"version", false},
	} {
		if g, e := flags.Lookup(tt.name).Value.Get(), tt.value; g != e {
			t.Errorf("FlagSet.Lookup(%q).Value.Get() = %v, expected %v", tt.name, g, e)
		}
	}
	if g, e := flags.NFlag(), 1; g != e {
		t.Errorf("FlagSet.NFlag() = %v, expected %v", g, e)
	}
	i := 0
	flags.Visit(func(*cli.Flag) { i++ })
	if g, e := i, 1; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	i = 0
	flags.VisitAll(func(*cli.Flag) { i++ })
	if g, e := i, 2; g != e {
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

	if err := flags.Parse([]string{"-c", "foo"}); err != nil {
		t.Fatal(err)
	}
	n := "c"
	if g, e := flags.Get(n), 1; g != e {
		t.Errorf("FlagSet.Get(%q) = %v, expected %v", n, g, e)
	}

	flags.Reset()
	switch err := flags.Parse([]string{"-c", ""}); {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), `choose from "bar", "baz", "foo" or "foobar"`):
		t.Error("unexpected error:", err)
	}

	flags.Reset()
	switch err := flags.Parse([]string{"-c", "b"}); {
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

	if err := flags.Parse([]string{"-c", "foo"}); err != nil {
		t.Fatal(err)
	}
	n := "c"
	if g, e := flags.Get(n), 1; g != e {
		t.Errorf("FlagSet.Get(%q) = %v, expected %v", n, g, e)
	}

	flags.Reset()
	switch err := flags.Parse([]string{"-c", ""}); {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), `choose from "bar", "baz", "foo" or "foobar"`):
		t.Error("unexpected error:", err)
	}

	flags.Reset()
	switch err := flags.Parse([]string{"-c", "b"}); {
	case err == nil:
		t.Fatal("expected error")
	case !strings.Contains(err.Error(), `choose from "bar" or "baz"`):
		t.Error("unexpected error:", err)
	}
}
