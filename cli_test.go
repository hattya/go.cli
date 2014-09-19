//
// go.cli :: cli_test.go
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
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hattya/go.cli"
)

func TestCLI(t *testing.T) {
	c := cli.NewCLI()
	c.Stdout = ioutil.Discard
	c.Stderr = ioutil.Discard
	args := []string{"-cli"}
	if err := c.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(cli.FlagError); !ok {
			t.Errorf("expected cli.FlagError, got %T", err)
		}
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	}

	c = cli.NewCLI()
	c.Flags.Bool("bool", false, "")
	c.Flags.Duration("duration", 0, "")
	c.Flags.Float64("float64", 0.0, "")
	c.Flags.Int("int", 0, "")
	c.Flags.Int64("int64", 0, "")
	c.Flags.String("string", "", "")
	c.Flags.Uint("uint", 0, "")
	c.Flags.Uint64("uint64", 0, "")
	c.Flags.Var("var", &value{}, "")
	args = strings.Fields("-bool -duration 1ms -float64 3.14 -int -1 -int64 -64 -string string -uint 1 -uint64 64 -var var 0 1")
	if err := c.Run(args); err != nil {
		t.Fatal(err)
	}
	ctx := cli.NewContext(c)
	for i := 0; i < len(ctx.Args); i++ {
		if g, e := ctx.Args[i], strconv.FormatInt(int64(i), 10); g != e {
			t.Errorf("expected %v, got %v", e, g)
		}
	}
	if g, e := len(ctx.Args), 2; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g := ctx.Value(""); g != nil {
		t.Errorf("expected %v, got %v", nil, g)
	}
	for _, tt := range []struct {
		name string
		fn   reflect.Value
		val  interface{}
	}{
		{"bool", reflect.ValueOf(ctx.Bool), true},
		{"duration", reflect.ValueOf(ctx.Duration), 1 * time.Millisecond},
		{"float64", reflect.ValueOf(ctx.Float64), 3.14},
		{"int", reflect.ValueOf(ctx.Int), -1},
		{"int64", reflect.ValueOf(ctx.Int64), int64(-64)},
		{"string", reflect.ValueOf(ctx.String), "string"},
		{"uint", reflect.ValueOf(ctx.Uint), uint(1)},
		{"uint64", reflect.ValueOf(ctx.Uint64), uint64(64)},
	} {
		rv := tt.fn.Call([]reflect.Value{reflect.ValueOf(tt.name)})
		if g, e := rv[0].Interface(), tt.val; g != e {
			t.Errorf("expected %v, got %v", e, g)
		}
	}
	if g, e := ctx.Value("var"), "var"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
}

func TestCLIOut(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := cli.NewCLI()
	c.Stdout = &stdout
	c.Stderr = &stderr

	c.Print("Print,")
	c.Println("Println")
	c.Printf("Printf")
	c.Error("Error,")
	c.Errorln("Errorln")
	c.Errorf("Errorf")

	if err := testOut(stdout.String(), "Print,Println\nPrintf"); err != nil {
		t.Error(err)
	}
	if err := testOut(stderr.String(), "Error,Errorln\nErrorf"); err != nil {
		t.Error(err)
	}
}

func testOut(g, e string) error {
	if g != e {
		return fmt.Errorf("output differ\nexpected: %q\n     got: %q", e, g)
	}
	return nil
}

func testStrings(get func(int) string, e []string) (err error) {
	g := make([]string, len(e))
	i := 0
	for ; i < len(e); i++ {
		g[i] = get(i)
		if g[i] != e[i] {
			break
		}

	}
	if i < len(e) {
		for ; i < len(e); i++ {
			g[i] = get(i)
		}
		err = fmt.Errorf("expected %v, got %v", e, g)
	}
	return
}

type value struct {
	s string
}

func (f *value) Set(v string) error {
	f.s = v
	return nil
}

func (f *value) Get() interface{} { return f.s }
func (f *value) String() string   { return fmt.Sprintf("%s", f.s) }
