go.cli
======

A command line interface framework

.. image:: https://semaphoreci.com/api/v1/hattya/go-cli/branches/master/badge.svg
   :target: https://semaphoreci.com/hattya/go-cli

.. image:: https://ci.appveyor.com/api/projects/status/fwccrp8kt0g2mhik/branch/master?svg=true
   :target: https://ci.appveyor.com/project/hattya/go-cli

.. image:: https://codecov.io/gh/hattya/go.cli/branch/master/graph/badge.svg
   :target: https://codecov.io/gh/hattya/go.cli


Installation
------------

.. code:: console

   $ go get -u github.com/hattya/go.cli


Usage
-----

.. code:: go

   package main

   import (
   	"os"

   	"github.com/hattya/go.cli"
   )

   var app = cli.NewCLI()

   func main() {
   	app.Version = "1.0"
   	app.Usage = "<options> hello"
   	app.Add(&cli.Command{
   		Name: []string{"hello"},
   		Action: func(ctx *cli.Context) error {
   			ctx.UI.Println("Hello World!")
   			return nil
   		},
   	})

   	if err := app.Run(os.Args[1:]); err != nil {
   		if _, ok := err.(cli.FlagError); ok {
   			os.Exit(2)
   		}
   		os.Exit(1)
   	}
   }


License
-------

go.cli is distributed under the terms of the MIT License.
