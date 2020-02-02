// +build windows

package symlink

import "os/exec"

// Dir creates symbolic link for dir.
func Dir(src, dst string) error {
	cmd := exec.Command("cmd", "/C", "mklink", "/D", dst, src)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
