package main

import (
	"debug/buildinfo"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/koron/goup/internal/bindir"
)

func binaryList(fs *flag.FlagSet, args []string) error {
	b, err := bindir.Open()
	if err != nil {
		return err
	}
	defer b.Close()
	for {
		name, err := b.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		bi, err := buildinfo.ReadFile(name)
		if err != nil {
			fmt.Printf("%s : failed to read: %s\n", name, err)
			continue
		}
		_ = bi
		fmt.Printf("%s : go=%s path=%s main.path=%s main.version=%s\n", name, bi.GoVersion, bi.Path, bi.Main.Path, bi.Main.Version)
	}
	return nil
}
