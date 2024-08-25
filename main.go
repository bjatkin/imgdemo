package main

import (
	"os"

	"github.com/bjatkin/imgdemo/cli"
	"github.com/bjatkin/imgdemo/cmd/find"
	"github.com/bjatkin/imgdemo/cmd/hide"
	"github.com/bjatkin/imgdemo/cmd/ishihara"
)

var Root = cli.Cmd[bool]{
	Name:        "imgdemo",
	Description: "a simple tool demoing what can be accomplished using the go standard library",
	Usage:       "imgdemo [COMMAND] [ARGS]",
	SubCmds: []cli.Runable{
		find.Cmd,
		hide.Cmd,
		ishihara.Cmd,
	},
}

func main() {
	code := Root.Run(os.Args[1:])
	if code != 0 {
		os.Exit(code)
	}
}
