package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/koron/goup/godlremote"
)

func remoteList(ctx context.Context, args []string) error {
	var all bool
	var match string
	fs := context2flagset(ctx)
	fs.BoolVar(&all, "all", false, "list all releases (archive and unstable)")
	fs.StringVar(&match, "match", "", "show only matched versions (regexp)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	var f func(godlremote.Release) bool
	if match != "" {
		rx, err := regexp.Compile(match)
		if err != nil {
			return fmt.Errorf("pattern error for -match: %w", err)
		}
		f = func(r godlremote.Release) bool {
			return rx.MatchString(r.Version)
		}
	}

	rels, err := godlremote.Download(ctx, all)
	if err != nil {
		return err
	}
	fmt.Println("Remote Version:")
	for _, r := range rels.Filter(f) {
		fmt.Printf("  %s\n", r.Version)
	}
	return nil
}
