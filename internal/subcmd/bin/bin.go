// Package bin provides "bin" sub command set of goup.
package bin

import (
	"context"
	"debug/buildinfo"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/koron-go/subcmd"
	"github.com/koron/goup/internal/bindir"
)

var Set = subcmd.DefineSet("bin", "operate executables in GOBIN dir",
	subcmd.DefineCommand("list", "list executables in GOBIN", ListCommand),
)

func ListCommand(ctx context.Context, args []string) error {
	//fs := subcmd.FlagSet(ctx)
	b, err := bindir.Open()
	if err != nil {
		return err
	}
	defer b.Close()
	fmt.Printf("%s\t%s\t%s\t%s\t%s\n", "name", "go", "path", "main.path", "main.version")
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
			slog.Warn("failed to read", "name", name, "err", err)
			continue
		}
		_ = bi
		//fmt.Printf("%s : go=%s path=%s main.path=%s main.version=%s\n", name, bi.GoVersion, bi.Path, bi.Main.Path, bi.Main.Version)
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", name, bi.GoVersion, bi.Path, bi.Main.Path, bi.Main.Version)
	}
	return nil
}
