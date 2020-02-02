// +build !windows

package symlink

// Dir creates symbolic link for dir.
func Dir(src, dst string) error {
	return errors.New("LinkDir not implemented yet")
}
