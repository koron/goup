package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/koron-go/subcmd"
	"github.com/koron/goup/internal/subcmd/install"
	"github.com/koron/goup/internal/subcmd/list"
	"github.com/koron/goup/internal/subcmd/remotelist"
	"github.com/koron/goup/internal/subcmd/uninstall"
)

type stderrWrapWriter struct{}

func (ww stderrWrapWriter) Write(b []byte) (int, error) {
	return os.Stderr.Write(b)
}

var logWrirter stderrWrapWriter

var (
	debugEnable = false
	debugLog    = log.New(logWrirter, "[DEBUG] ", log.LstdFlags)
	infoLog     = log.New(logWrirter, "", 0)
	warnLog     = log.New(logWrirter, "[WARN] ", log.LstdFlags)
)

func debugf(msg string, args ...interface{}) {
	if !debugEnable {
		return
	}
	debugLog.Printf(msg, args...)
}

func infof(msg string, args ...interface{}) {
	infoLog.Printf(msg, args...)
}

func warnf(msg string, args ...interface{}) {
	warnLog.Printf(msg, args...)
}

func envGoupRoot() string {
	s := os.Getenv("GOUP_ROOT")
	if s != "" {
		return s
	}
	// for comaptibility. this will be removed in future version.
	return os.Getenv("GODL_ROOT")
}

func envGoupLinkname() string {
	s := os.Getenv("GOUP_LINKNAME")
	if s != "" {
		return s
	}
	return "current"
}

var rootCommandSet = subcmd.DefineRootSet(
	remotelist.Command, // remotelist
	install.Command,    // install
	uninstall.Command,  // uninstall
	upgradeCommand,     // upgrade
	list.Command,       // list
	switchCommand,      // switch
	cleanCommand,       // clean
)

func main() {
	flag.BoolVar(&debugEnable, "debug", false, "enable debug log")
	flag.Parse()
	err := subcmd.Run(rootCommandSet, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %s\n", err)
		os.Exit(1)
	}
}
