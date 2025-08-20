package common

import "os"

// GoupRoot retrieves GOUP_ROOT environment variable's value
func GoupRoot() string {
	s := os.Getenv("GOUP_ROOT")
	if s != "" {
		return s
	}
	// for comaptibility. this will be removed in future version.
	return os.Getenv("GODL_ROOT")
}

// GoupLinkname retrieves GOUP_LINKNAME environment variable's value
func GoupLinkname() string {
	s := os.Getenv("GOUP_LINKNAME")
	if s != "" {
		return s
	}
	return "current"
}
