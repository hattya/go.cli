//
// go.cli :: cli_test.go
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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hattya/go.cli"
)

func TestCLI(t *testing.T) {
	app := cli.NewCLI()
	app.Stdout = ioutil.Discard
	app.Stderr = ioutil.Discard
	args := []string{"-g"}
	if err := app.Run(args); err == nil {
		t.Error("expected error")
	} else {
		if _, ok := err.(cli.FlagError); !ok {
			t.Errorf("expected cli.FlagError, got %T", err)
		}
		if !strings.Contains(err.Error(), "not defined") {
			t.Error("unexpected error:", err)
		}
	}

	app = cli.NewCLI()
	app.Flags.Bool("bool", false, "")
	app.Flags.Duration("duration", 0, "")
	app.Flags.Float64("float64", 0.0, "")
	app.Flags.Int("int", 0, "")
	app.Flags.Int64("int64", 0, "")
	app.Flags.String("string", "", "")
	app.Flags.Uint("uint", 0, "")
	app.Flags.Uint64("uint64", 0, "")
	app.Flags.Var("var", &value{}, "")
	args = strings.Fields("-bool -duration 1ms -float64 3.14 -int -1 -int64 -64 -string string -uint 1 -uint64 64 -var var 0 1")
	if err := app.Run(args); err != nil {
		t.Fatal(err)
	}
	ctx := cli.NewContext(app)
	for i := 0; i < len(ctx.Args); i++ {
		if g, e := ctx.Args[i], strconv.FormatInt(int64(i), 10); g != e {
			t.Errorf("Context.Args[%v] = %v, expected %v", i, g, e)
		}
	}
	if g, e := len(ctx.Args), 2; g != e {
		t.Errorf("len(Context.Args) = %v, expected %v", g, e)
	}
	if g := ctx.Value(""); g != nil {
		t.Errorf("Context.Value(%q) = %v, expected %v", "", g, nil)
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
			t.Errorf("Context.%s(%q) = %v, expected %v", strings.Title(tt.name), tt.name, g, e)
		}
	}
	n := "var"
	if g, e := ctx.Value(n), "var"; g != e {
		t.Errorf("Context.Value(%q) = %v, expected %v", n, g, e)
	}
}

func TestCLIOut(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := cli.NewCLI()
	app.Stdout = &stdout
	app.Stderr = &stderr

	app.Print("Print,")
	app.Println("Println")
	app.Printf("Printf")
	app.Error("Error,")
	app.Errorln("Errorln")
	app.Errorf("Errorf")

	if err := testOut(stdout.String(), "Print,Println\nPrintf"); err != nil {
		t.Error(err)
	}
	if err := testOut(stderr.String(), "Error,Errorln\nErrorf"); err != nil {
		t.Error(err)
	}
}

func TestTitle(t *testing.T) {
	app := cli.NewCLI()

	if err := app.Title(app.Name); err != nil {
		t.Error(err)
	}

	app.Stdout = os.Stdout
	if err := app.Title(app.Name); err != nil {
		t.Error(err)
	}
}

var promptTests = []struct {
	in, out string
	err     error
}{
	{
		in:  "input",
		out: "input",
	},
	{
		in:  "input\n",
		out: "input",
	},
	{
		in:  "input\r\n",
		out: "input",
	},
	{
		err: io.EOF,
	},
}

func TestPrompt(t *testing.T) {
	var stdin, stdout bytes.Buffer
	app := cli.NewCLI()
	app.Stdin = &stdin
	app.Stdout = &stdout
	prompt := ">> "

	for _, tt := range promptTests {
		stdin.Reset()
		stdout.Reset()
		stdin.WriteString(tt.in)
		l, err := app.Prompt(prompt)
		if tt.err == nil {
			if err != nil {
				t.Fatal(err)
			}
		} else {
			if g, e := err, tt.err; g != e {
				t.Errorf("expected %v, got %v", e, g)
			}
		}
		if g, e := l, tt.out; g != e {
			t.Errorf("expected %q, got %q", e, g)
		}
		if g, e := stdout.String(), prompt; g != e {
			t.Errorf("expected %q, got %q", e, g)
		}
	}
}

func TestPassword(t *testing.T) {
	var stdin, stdout bytes.Buffer
	app := cli.NewCLI()
	app.Stdin = &stdin
	app.Stdout = &stdout
	prompt := "Password: "

	stdin.WriteString("password")
	l, err := app.Password(prompt)
	if err != nil {
		t.Error(err)
	}
	if g, e := l, "password"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
	if g, e := stdout.String(), prompt+"\n"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestPrepare(t *testing.T) {
	app := cli.NewCLI()
	app.Stdout = ioutil.Discard
	app.Stderr = ioutil.Discard
	app.Prepare = func(ctx *cli.Context, cmd *cli.Command) error {
		cmd.Data = cmd.Data.(int) + 1
		return nil
	}
	app.Add(&cli.Command{
		Name: []string{"true"},
		Data: 0,
	})

	args := []string{"true"}
	if err := app.Run(args); err != nil {
		t.Fatal(err)
	}
	if g, e := app.Cmds[0].Data, 1; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
}

var errorHandlerTests = []struct {
	err error
	out string
}{
	{
		err: nil,
		out: "",
	},
	{
		err: cli.ErrCommand,
		out: cli.Dedent(`
			usage: %[1]v
		`),
	},
	{
		err: &cli.Abort{
			Err: fmt.Errorf("abort"),
		},
		out: cli.Dedent(`
			%[1]v: abort
		`),
	},
	{
		err: &cli.Abort{
			Err:  fmt.Errorf("abort"),
			Hint: "hint",
		},
		out: cli.Dedent(`
			%[1]v: abort
			hint
		`),
	},
	{
		err: cli.FlagError("flag error"),
		out: cli.Dedent(`
			%[1]v: flag error
			usage: %[1]v
		`),
	},
	{
		err: &cli.CommandError{
			Name: "cmd",
		},
		out: cli.Dedent(`
			%[1]v: unknown command 'cmd'
			usage: %[1]v
		`),
	},
	{
		err: &cli.CommandError{
			Name: "b",
			List: []string{"bar", "baz"},
		},
		out: cli.Dedent(`
			%[1]v: command 'b' is ambiguous
			    bar baz
		`),
	},
	{
		err: fmt.Errorf("error"),
		out: cli.Dedent(`
			%[1]v: error
		`),
	},
}

func TestErrorHandler(t *testing.T) {
	var b bytes.Buffer
	app := cli.NewCLI()
	app.Stdout = &b
	app.Stderr = &b
	ctx := cli.NewContext(app)
	for _, tt := range errorHandlerTests {
		b.Reset()
		cli.ErrorHandler(ctx, tt.err)
		var out string
		if tt.out != "" {
			out = fmt.Sprintf(tt.out, ctx.Name())
		}
		if err := testOut(b.String(), out); err != nil {
			t.Error(err)
		}
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
func (f *value) String() string   { return fmt.Sprintf("%v", f.s) }
