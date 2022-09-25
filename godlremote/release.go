package godlremote

import "github.com/koron/goup/internal/verutil"

// Release represents a release.
// It should be kept in sync with the dl code in golang/website/internal/dl.
//
// See https://github.com/golang/website/blob/d0b4462f2c677caac44e6f5cb06ea9fd3555f222/internal/dl/dl.go#L131-L137
type Release struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []File `json:"files"`
}

// Semver returns Version in semantic versioning notation.
func (r Release) Semver() string {
	return verutil.Regulate(r.Version)
}

// Releases is a collection of Release
type Releases []Release

// Filter returns only matched Releases.
func (rels Releases) Filter(f func(Release) bool) Releases {
	if f == nil {
		return rels
	}
	var matched Releases
	for _, r := range rels {
		if f(r) {
			matched = append(matched, r)
		}
	}
	return matched
}
