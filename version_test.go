//
// go.cli :: version_test.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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
	for _, tt := range versionTests {
		var b bytes.Buffer
		app := cli.NewCLI()
		app.Version = tt.in
		app.Stdout = &b
		app.Add(cli.NewVersionCommand())
		if err := app.Run([]string{"version"}); err != nil {
			t.Fatal(err)
		}
		if err := testOut(b.String(), fmt.Sprintf(versionOut, app.Name, tt.out)); err != nil {
			t.Error(err)
		}
	}
}

func TestVersion(t *testing.T) {
	for _, tt := range versionTests {
		var b bytes.Buffer
		app := cli.NewCLI()
		app.Version = tt.in
		app.Stdout = &b
		if err := app.Run([]string{"--version"}); err != nil {
			t.Fatal(err)
		}
		if err := testOut(b.String(), fmt.Sprintf(versionOut, app.Name, tt.out)); err != nil {
			t.Error(err)
		}
	}

	var b bytes.Buffer
	app := cli.NewCLI()
	app.Version = "1.0"
	app.Stdout = &b
	app.Add(&cli.Command{
		Name:  []string{"cmd"},
		Flags: cli.NewFlagSet(),
	})
	if err := app.Run([]string{"cmd", "--version"}); err != nil {
		t.Fatal(err)
	}
	if err := testOut(b.String(), fmt.Sprintf(versionOut, app.Name, app.Version)); err != nil {
		t.Error(err)
	}
}
