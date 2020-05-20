//
// go.cli :: version.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package cli

func NewVersionCommand() *Command {
	return &Command{
		Name:  []string{"version"},
		Desc:  "show version information",
		Flags: NewFlagSet(),
		Action: func(ctx *Context) error {
			return Version(ctx)
		},
	}
}

var Version = ShowVersion

func ShowVersion(ctx *Context) error {
	version := ctx.UI.Version
	if version == "" {
		version = "unknown"
	}
	ctx.UI.Printf("%v version %v\n", ctx.UI.Name, version)
	return nil
}
