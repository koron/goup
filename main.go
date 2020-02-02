package main

import (
	"flag"
	"log"
	"os"

	"github.com/koron/go-subcmd"
)

var (
	debugEnable = false
	debugLog    = log.New(os.Stderr, "[DEBUG]", log.LstdFlags)
	infoLog     = log.New(os.Stderr, "[INFO]", log.LstdFlags)
	warnLog     = log.New(os.Stderr, "[WARN]", log.LstdFlags)
	errorLog    = log.New(os.Stderr, "[ERROR]", log.LstdFlags)
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

func errorf(msg string, args ...interface{}) {
	errorLog.Printf(msg, args...)
}

func main() {
	flag.BoolVar(&debugEnable, "debug", false, "enable debug log")
	flag.Parse()
	err := cmds.Run(flag.Args())
	if err != nil {
		errorf("failed: %s", err)
		os.Exit(1)
	}
}

var cmds = subcmd.Subcmds{
	"remotelist": subcmd.Main2(remoteList),
	"install":    subcmd.Main2(install),
	"list":       subcmd.Main2(localList),
	"switch":     subcmd.Main2(localSwitch),
}
