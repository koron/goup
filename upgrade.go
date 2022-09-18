package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/koron/goup/godlremote"
	"golang.org/x/mod/semver"
)

type upgradePlan struct {
	local  installedGo
	remote latestRelease
	curr   bool
}

type upgrader interface {
	install(ctx context.Context, ins installer, ver string) error
	setCurrent(root, linkname, installedName string) error
	uninstall(ctx context.Context, uni uninstaller, ver string) error
}

type upgraderActual struct{}

func (ua upgraderActual) install(ctx context.Context, ins installer, ver string) error {
	return ins.install(ctx, ver)
}

func (ua upgraderActual) setCurrent(root, linkname, installedName string) error {
	return switchGo(root, linkname, installedName)
}

func (ua upgraderActual) uninstall(ctx context.Context, uni uninstaller, ver string) error {
	return uni.uninstall(ctx, ver)
}

type upgraderRehearsal struct{}

func (ur upgraderRehearsal) install(ctx context.Context, ins installer, ver string) error {
	debugf("DRYRUN: install Go %s", ver)
	return nil
}

func (ur upgraderRehearsal) setCurrent(root, linkname, installedName string) error {
	debugf("DRYRUN: set current \"%s\" as %s in %s", linkname, installedName, root)
	return nil
}

func (ur upgraderRehearsal) uninstall(ctx context.Context, uni uninstaller, ver string) error {
	debugf("DRYRUN: uninstall Go %s", ver)
	return nil
}

type debugInstalledGos installedGos

func (d debugInstalledGos) String() string {
	bb := &bytes.Buffer{}
	for _, g := range d {
		bb.WriteString("\n\t" + g.name)
	}
	return bb.String()
}

// upgradeCmd upgrades installed Go version.
func upgradeCmd(fs *flag.FlagSet, args []string) error {
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
	return upgrade(ctx, root, linkname, dryrun, all)
}

// upgrade upgrades installed Go versions.
func upgrade(ctx context.Context, root, linkname string, dryrun, all bool) error {
	debugf("upgrade processing...")

	// list local versions.
	installed, err := listInstalledGo(root)
	if err != nil {
		return fmt.Errorf("failed to list installed Go: %w", err)
	}
	debugf("detect installed Go: %s", debugInstalledGos(installed))

	// list remote versions (with considering -all option)
	rels, err := godlremote.Download(ctx, all)
	if err != nil {
		return fmt.Errorf("failed to list remote Go: %w", err)
	}
	latests := groupReleases(rels)
	debugf("found releases: %s", debugLatestReleases(latests))

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
	if len(upgrades) == 0 {
		debugf("no upgrades")
		return nil
	}

	var upg upgrader = upgraderActual{}
	if dryrun {
		upg = upgraderRehearsal{}
	}

	// repeat versions to be upgraded
	for _, target := range upgrades {
		debugf("upgrading Go %s", target.local.name)

		// install new version
		ins := installer{
			releases: godlremote.Releases{target.remote.origin},
			rootdir:  root,
			force:    false,
			goos:     target.local.os,
			goarch:   target.local.arch,
		}
		ver := target.remote.origin.Version
		archiveFile, ok := ins.archiveFile(ver)
		if !ok {
			warnf("no archive files found for version=%s os=%s arch=%s, skipped", ver, ins.goos, ins.goarch)
			continue
		}
		err := upg.install(ctx, ins, ver)
		if err != nil {
			return fmt.Errorf("failed to install Go %s: %w", ver, err)
		}
		installedName := archiveFile.Name()

		// switch "current" version if needed
		if target.curr {
			err := upg.setCurrent(root, linkname, installedName)
			if err != nil {
				return fmt.Errorf("failed to switch current Go as %s: %w", installedName, err)
			}
		}

		// clean old version
		uni := uninstaller{
			rootdir: root,
			goos:    target.local.os,
			goarch:  target.local.arch,
			clean:   false,
		}
		err = upg.uninstall(ctx, uni, target.local.version)
		if err != nil {
			return fmt.Errorf("failed to uinstall Go %s: %w", target.local.name, err)
		}
		infof("upgraded Go %s to %s", target.local.name, installedName)
	}

	return nil
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

type debugLatestReleases map[string]latestRelease

func (rels debugLatestReleases) String() string {
	vers := make([]string, 0, len(rels))
	for _, r := range rels {
		vers = append(vers, r.semver)
	}
	sort.Slice(vers, func(i, j int) bool {
		return semver.Compare(vers[i], vers[j]) > 0
	})
	bb := &bytes.Buffer{}
	for _, v := range vers {
		bb.WriteString("\n\t" + v)
	}
	return bb.String()
}
