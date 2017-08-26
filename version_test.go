//
// go.cli :: version_test.go
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
	"testing"

	"github.com/hattya/go.cli"
)

var versionOut = "%v version %v\n"

var versionTests = []struct {
	in, out string
}{
	{"", "unknown"},
	{"1.0", "1.0"},
}

func TestVersionCommand(t *testing.T) {
	var b bytes.Buffer
	args := []string{"version"}
	for _, tt := range versionTests {
		b.Reset()
		app := cli.NewCLI()
		app.Version = tt.in
		app.Stdout = &b
		app.Add(cli.NewVersionCommand())
		if err := app.Run(args); err != nil {
			t.Fatal(err)
		}
		if err := testOut(b.String(), fmt.Sprintf(versionOut, app.Name, tt.out)); err != nil {
			t.Error(err)
		}
	}
}

func TestVersion(t *testing.T) {
	var b bytes.Buffer
	args := []string{"--version"}
	for _, tt := range versionTests {
		b.Reset()
		app := cli.NewCLI()
		app.Version = tt.in
		app.Stdout = &b
		if err := app.Run(args); err != nil {
			t.Fatal(err)
		}
		if err := testOut(b.String(), fmt.Sprintf(versionOut, app.Name, tt.out)); err != nil {
			t.Error(err)
		}
	}

	b.Reset()
	app := cli.NewCLI()
	app.Version = "1.0"
	app.Stdout = &b
	app.Add(&cli.Command{
		Name:  []string{"cmd"},
		Flags: cli.NewFlagSet(),
	})
	args = []string{"cmd", "--version"}
	if err := app.Run(args); err != nil {
		t.Fatal(err)
	}
	if err := testOut(b.String(), fmt.Sprintf(versionOut, app.Name, app.Version)); err != nil {
		t.Error(err)
	}
}
