package main

import (
	"context"
	"flag"
	"strings"

	"github.com/koron-go/subcmd"
)

func context2flagset(ctx context.Context) *flag.FlagSet {
	name := strings.Join(subcmd.Names(ctx), " ")
	return flag.NewFlagSet(name, flag.ExitOnError)
}

var remotelistCommand = subcmd.DefineCommand("remotelist", "list published releases", func(ctx context.Context, args []string) error {
	return remoteList(context2flagset(ctx), args)
})

var installCommand = subcmd.DefineCommand("install", "install Go releases", func(ctx context.Context, args []string) error {
	return installCmd(context2flagset(ctx), args)
})

var uninstallCommand = subcmd.DefineCommand("uninstall", "uninstall Go releases", func(ctx context.Context, args []string) error {
	return uninstallCmd(context2flagset(ctx), args)
})

var upgradeCommand = subcmd.DefineCommand("upgrade", "upgrade installed Go releases", func(ctx context.Context, args []string) error {
	return upgradeCmd(context2flagset(ctx), args)
})

var listCommand = subcmd.DefineCommand("list", "list installed releases", func(ctx context.Context, args []string) error {
	return localList(context2flagset(ctx), args)
})

var switchCommand = subcmd.DefineCommand("switch", "switch active Go release", func(ctx context.Context, args []string) error {
	return localSwitch(context2flagset(ctx), args)
})

var cleanCommand = subcmd.DefineCommand("clean", "clean download caches", func(ctx context.Context, args []string) error {
	return localClean(context2flagset(ctx), args)
})

var rootCommandSet = subcmd.DefineRootSet(
	remotelistCommand, // remotelist
	installCommand,    // install
	uninstallCommand,  // uninstall
	upgradeCommand,    // upgrade
	listCommand,       // list
	switchCommand,     // switch
	cleanCommand,      // clean
)
