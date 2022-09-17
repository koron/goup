package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/koron/goup/godlremote"
	"golang.org/x/mod/semver"
)

type upgradePlan struct {
	local  installedGo
	remote latestRelease
	curr   bool
}

// upgrade upgrades installed Go versions.
func upgrade(fs *flag.FlagSet, args []string) error {
	var root string
	var linkname string
	var dryrun bool
	var all bool
	fs.StringVar(&root, "root", envGoupRoot(), "root dir to install")
	fs.StringVar(&linkname, "linkname", "current", "name of symbolic link to switch")
	fs.BoolVar(&dryrun, "dryrun", false, "don't switch, just test")
	fs.BoolVar(&all, "all", false, "clean all caches")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if root == "" {
		fs.PrintDefaults()
		return errors.New("required -root")
	}

	ctx := context.Background()

	// list local versions.
	installed, err := listInstalledGo(root)
	if err != nil {
		return fmt.Errorf("failed to list installed Go: %w", err)
	}

	// list remote versions (with considering -all option)
	rels, err := godlremote.Download(ctx, all)
	if err != nil {
		return fmt.Errorf("failed to list remote Go: %w", err)
	}
	latests := groupReleases(rels)

	currName, err := localCurrent(filepath.Join(root, linkname))
	if err != nil {
		return fmt.Errorf("failed to determine \"current\" version: %w", err)
	}

	// compare versions, determine versions to be upgraded
	upgrades := make([]upgradePlan, 0, len(installed))
	for _, local := range installed {
		latest, ok := latests[shrinkVersion(local.semver)]
		if !ok || semver.Compare(latest.semver, local.semver) <= 0 {
			continue
		}
		upgrades = append(upgrades, upgradePlan{
			local:  local,
			remote: latest,
			curr:   currName != "" && local.name == currName,
		})
	}

	// repeat versions to be upgraded
	// TODO: consider dryrun
	for _, target := range upgrades {
		// install new version
		f, err := installer{
			releases: godlremote.Releases{target.remote.origin},
			rootdir:  root,
			force:    false,
			goos:     target.local.os,
			goarch:   target.local.arch,
		}.install2(ctx, target.remote.origin.Version)
		if err != nil {
			return fmt.Errorf("failed to install Go %s: %w", target.remote.origin.Version, err)
		}
		installedName := f.Name()
		// switch "current" version if needed
		if target.curr {
			err := switchGo(root, linkname, installedName)
			if err != nil {
				return fmt.Errorf("failed to switch Go %s: %w", installedName, err)
			}
		}
		// clean old version
		err = uninstaller{
			rootdir: root,
			goos:    target.local.os,
			goarch:  target.local.arch,
			clean:   false,
		}.uninstall(ctx, target.local.version)
		if err != nil {
			return fmt.Errorf("failed to uinstall Go %s: %w", target.local.name, err)
		}
	}

	return errors.New("not implemented yet")
}

type latestRelease struct {
	semver string
	origin godlremote.Release
}

func shrinkVersion(ver string) string {
	majorMinor := semver.MajorMinor(ver)
	if semver.Prerelease(ver) == "" {
		return majorMinor
	}
	return majorMinor + "-pre"
}

func groupReleases(releases godlremote.Releases) map[string]latestRelease {
	m := map[string]latestRelease{}
	for _, r := range releases {
		v := r.Semver()
		k := shrinkVersion(v)
		curr, ok := m[k]
		if !ok || semver.Compare(v, curr.semver) > 0 {
			m[k] = latestRelease{
				semver: v,
				origin: r,
			}
		}
	}
	return m
}
