package list

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/koron-go/subcmd"
	"github.com/koron/goup/internal/common"
)

var Command = subcmd.DefineCommand("list", "list installed releases", func(ctx context.Context, args []string) error {
	var root string
	var linkname string
	fs := subcmd.FlagSet(ctx)
	fs.StringVar(&root, "root", common.GoupRoot(), "root dir to install")
	fs.StringVar(&linkname, "linkname", common.GoupLinkname(), "name of symbolic link to switch")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if root == "" {
		fs.PrintDefaults()
		return errors.New("required -root")
	}

	curr, err := common.LocalCurrent(filepath.Join(root, linkname))
	if err != nil {
		return err
	}

	list, err := common.ListInstalledGo(root)
	if err != nil {
		return err
	}
	fmt.Println("Local Version:")
	for _, g := range list {
		if curr != "" && g.Name == curr {
			fmt.Printf("  %s (%s)\n", g.Name, linkname)
			continue
		}
		fmt.Printf("  %s\n", g.Name)
	}
	return nil
})
