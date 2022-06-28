# koron/goup

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/goup)](https://pkg.go.dev/github.com/koron/goup)
[![Actions/Go](https://github.com/koron/goup/workflows/Go/badge.svg)](https://github.com/koron/goup/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/goup)](https://goreportcard.com/report/github.com/koron/goup)

Utility to download, extract and switch released go versions on Windows.

## Getting started

Check stable go releases.

```
> goup remotelist
Remote Version:
  go1.13.7
  go1.12.16
```

Check all go releases.

```
> goup remotelist -all
Remote Version:
  go1.13.7
  go1.12.16
  go1.13.6
  go1.13.5
  go1.13.4
  (...snip...)
  go1.3.1
  go1.3
  go1.2.2
  go1
  go1.14beta1
```

Install (=download and extract) a go archive.

```
> goup install -root C:\golang go1.13.7
```

List installed go.

```
> goup list -root C:\golang
Local Version:
  go1.12.16.windows-amd64
  go1.13.7.windows-amd64
```

Switch go version.  This requires symbolic link to directory, so you may need
to install ["Windows Developer mode"][devmode].

```
> goup switch -root C:\golang go1.13.7
go1.13.7.windows-amd64

> dir C:\golang
 Volume in drive C has no label.
 Volume Serial Number is DEAD-BEEF

 Directory of C:\golang

2020/02/03  00:42    <DIR>          .
2020/02/03  00:42    <DIR>          ..
2020/02/03  00:42    <SYMLINKD>     current [go1.13.7.windows-amd64]
2020/02/02  18:48    <DIR>          dl
2020/02/02  18:48    <DIR>          go1.12.16.windows-amd64
2020/02/02  18:47    <DIR>          go1.13.7.windows-amd64
               0 File(s)              0 bytes
               6 Dir(s)  999,999,999,999 bytes free
```

Now `C:\golang\current` is a symblic link to go1.13.7's dir.
You can add `C:\golang\current\bin` to `PATH` env, to run `go` command.

`GOUP_ROOT` env works as instead of `-root C:\golang` option.

## Environment variables

*   `GOUP_ROOT` - root dir to install.

    `GODL_ROOT` can be used for same purpose but it is obsoleted.
    It will be removed at future release.

*   `GOUP_LINKNAME` - name of symbolic link for active version.

    Default is `current` but you can override it by this.

[devmode]:https://docs.microsoft.com/en-us/windows/uwp/get-started/enable-your-device-for-development
