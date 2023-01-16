package main

import (
	"fmt"
	"os"

	"github.com/koron-go/subcmd"
	"github.com/koron/goup/internal/subcmd/clean"
	"github.com/koron/goup/internal/subcmd/install"
	"github.com/koron/goup/internal/subcmd/list"
	"github.com/koron/goup/internal/subcmd/remotelist"
	"github.com/koron/goup/internal/subcmd/switchgo"
	"github.com/koron/goup/internal/subcmd/uninstall"
	"github.com/koron/goup/internal/subcmd/upgrade"
)

var rootCommandSet = subcmd.DefineRootSet(
	remotelist.Command, // remotelist
	install.Command,    // install
	uninstall.Command,  // uninstall
	upgrade.Command,    // upgrade
	list.Command,       // list
	switchgo.Command,   // switch
	clean.Command,      // clean
	binSet,             // bin
)

func main() {
	err := subcmd.Run(rootCommandSet, os.Args[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %s\n", err)
		os.Exit(1)
	}
}
