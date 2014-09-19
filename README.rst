go.cli
======

A command line interface framework


Install
-------

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

   func main() {
   	app := cli.NewCLI()
   	app.Version = "1.0"
   	app.Usage = "<options> hello"
   	app.Add(&cli.Command{
   		Name: []string{"hello"},
   		Action: func(ctx *cli.Context) error {
   			ctx.CLI.Println("Hello World!")
   			return nil
   		},
   	})

   	if err := app.Run(os.Args[1:]); err != nil {
   		switch err.(type) {
   		case cli.FlagError:
   			os.Exit(2)
   		case *cli.CommandError:
   		}
   		os.Exit(1)
   	}
   }


License
-------

go.cli is distributed under the terms of the MIT License
