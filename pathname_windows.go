package pathname

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CopyDirTo is demo wrapper of the windows xcopy command
func (o Pathname) CopyDirTo(toPath Pathname) error {
	if !o.IsDirectory() {
		return fmt.Errorf("source path not a directory")
	}
	if toPath.Exist() {
		return fmt.Errorf("destination path already exist")
	}
	cmd := exec.Command(`C:\Windows\System32\cmd.exe`, "/C", "xcopy", o.filePath, toPath.filePath, "/E/H/C/I")
	var buf = bytes.Buffer{}
	cmd.Stderr = &buf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s. %s", err.Error(), buf.Bytes())
	}
	return cmd.Run()
}
