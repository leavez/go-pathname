//go:build linux || darwin
// +build linux darwin

package pathname

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CopyDirTo Just demo wrapper of the unix cp command with '-r'
// the toPath must not exist (while the toPath's parent must exist)
func (o Pathname) CopyDirTo(toPath Pathname) error {
	if !o.IsDirectory() {
		return fmt.Errorf("source path not a directory")
	}
	if toPath.Exist() {
		return fmt.Errorf("destination path already exist")
	}
	cmd := exec.Command("cp", "-r", o.filePath, toPath.filePath)
	var buf = bytes.Buffer{}
	cmd.Stderr = &buf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s. %s", err.Error(), buf.Bytes())
	}
	return nil
}
