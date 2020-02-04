// +build !windows

package symlink

import "os"

// Dir creates symbolic link for dir.
func Dir(src, dst string) error {
	return os.Symlink(src, dst)
}
