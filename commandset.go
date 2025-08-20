package main

import (
	"github.com/koron-go/subcmd"
)

var remotelistCommand = subcmd.DefineCommand("remotelist", "list published releases", remoteList)

var installCommand = subcmd.DefineCommand("install", "install Go releases", installCmd)

var uninstallCommand = subcmd.DefineCommand("uninstall", "uninstall Go releases", uninstallCmd)

var upgradeCommand = subcmd.DefineCommand("upgrade", "upgrade installed Go releases", upgradeCmd)

var listCommand = subcmd.DefineCommand("list", "list installed releases", localList)

var switchCommand = subcmd.DefineCommand("switch", "switch active Go release", localSwitch)

var cleanCommand = subcmd.DefineCommand("clean", "clean download caches", localClean)

var rootCommandSet = subcmd.DefineRootSet(
	remotelistCommand, // remotelist
	installCommand,    // install
	uninstallCommand,  // uninstall
	upgradeCommand,    // upgrade
	listCommand,       // list
	switchCommand,     // switch
	cleanCommand,      // clean
)
