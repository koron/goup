// Package verutil provides version utilities.
package verutil

import (
	"bytes"
	"regexp"
)

var rxGoVer = regexp.MustCompile(`^go(?P<major>\d+)(?:\.(?P<minor>\d+)(?:\.(?P<patch>\d+))?)?(?P<pr>[A-Za-z][-.0-9A-Za-z]*)?`)

func regnum(s string) string {
	if s == "" {
		return "0"
	}
	return s
}

// Regulate regulates Go version string as semantic versioning.
//
// Examples:
//
//	go1.19    -> v1.19.0
//	go1.18.6  -> v1.18.6
//	go1.20rc1 -> v1.20.0-rc1
func Regulate(s string) string {
	m := rxGoVer.FindStringSubmatch(s)
	if m == nil {
		return ""
	}
	bb := &bytes.Buffer{}
	bb.WriteRune('v')
	bb.WriteString(regnum(m[1]))
	bb.WriteRune('.')
	bb.WriteString(regnum(m[2]))
	bb.WriteRune('.')
	bb.WriteString(regnum(m[3]))
	if m[4] != "" {
		bb.WriteRune('-')
		bb.WriteString(m[4])
	}
	return bb.String()
}
