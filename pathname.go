package pathname

import (
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"path/filepath"
)

// Pathname is convenient tool to manipulate filepath
// inspired by Ruby's Pathname
type Pathname struct {
	filePath string
}

func New(pathString string) Pathname {
	return Pathname{filePath: path.Clean(pathString)}
}

// Raw create demo Pathname without any modification
func Raw(pathString string) Pathname {
	return Pathname{filePath: pathString}
}

// -- path string manipulation (without interaction with filesystem)

// return the string representation
func (o Pathname) String() string {
	return o.filePath
}

func (o Pathname) Append(path string) Pathname {
	return Raw(filepath.Join(o.filePath, path))
}

func (o Pathname) Parent() Pathname {
	return Raw(filepath.Dir(o.filePath))
}

func (o Pathname) Basename() string {
	return filepath.Base(o.filePath)
}

func (o Pathname) Extension() string {
	return filepath.Ext(o.filePath)
}

func (o Pathname) IsAbsolute() bool {
	return filepath.IsAbs(o.filePath)
}

// -- fs tree

// AbsolutePath returns an absolute representation of path.
// If the path is not absolute it will be joined with the current
func (o Pathname) AbsolutePath() Pathname {
	var path, _ = filepath.Abs(o.filePath)
	return Raw(path)
}

func (o Pathname) Children() ([]Pathname, error) {
	files, err := os.ReadDir(o.filePath)
	if err != nil {
		return nil, err
	}

	var paths = make([]Pathname, len(files))
	for i, file := range files {
		var p = filepath.Join(o.filePath, file.Name())
		paths[i] = Raw(p)
	}
	return paths, nil
}

// Walk is just a wrapper of filepath.WalkDir
func (o Pathname) Walk(block func(path Pathname, entry fs.DirEntry, err error) error) error {
	return filepath.WalkDir(o.filePath, func(path string, info fs.DirEntry, err error) error {
		return block(Raw(path), info, err)
	})
}

func (o Pathname) ChildrenNames() []string {
	files, err := os.ReadDir(o.filePath)
	if err != nil {
		return nil
	}
	var paths = make([]string, len(files))
	for i, file := range files {
		paths[i] = file.Name()
	}
	return paths
}

func (o Pathname) Mkpath() error {
	return os.MkdirAll(o.filePath, 0755)
}

func (o Pathname) Mkpath2(perm uint32) error {
	return os.MkdirAll(o.filePath, os.FileMode(perm))
}

func (o Pathname) Rmtree() error {
	return os.RemoveAll(o.filePath)
}

// Glob is the shell style glob
// Doesn't support '**' as any nested path
func (o Pathname) Glob(pattern string) []Pathname {
	var fullPattern = filepath.Join(o.filePath, pattern)
	var m, _ = filepath.Glob(fullPattern)

	var paths = make([]Pathname, len(m))
	for i, s := range m {
		paths[i] = Raw(s)
	}
	return paths
}

func (o Pathname) AppendRandom() Pathname {
	var name = fmt.Sprintf("%8x", rand.Int())
	var t = o.Append(name)
	if t.Exist() {
		return o.AppendRandom()
	}
	return t
}

// -- file manipulation

func (o Pathname) Exist() bool {
	var _, err = os.Stat(o.filePath)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (o Pathname) IsDirectory() bool {
	fileInfo, err := os.Stat(o.filePath)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func (o Pathname) IsFile() bool {
	fileInfo, err := os.Stat(o.filePath)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func (o Pathname) IsRegularFile() bool {
	fileInfo, err := os.Lstat(o.filePath)
	if err != nil {
		return false
	}
	return fileInfo.Mode().IsRegular()
}

func (o Pathname) IsSymlink() bool {
	fileInfo, err := os.Lstat(o.filePath)
	if err != nil {
		return false
	}
	return fileInfo.Mode()&os.ModeSymlink != 0
}

func (o Pathname) Write(text string) error {
	var bytes = []byte(text)
	return os.WriteFile(o.filePath, bytes, 0644)
}

func (o Pathname) Read() (string, error) {
	var data, err = os.ReadFile(o.filePath)
	if err != nil {
		return "", err
	}
	var text = string(data)
	return text, err
}

func (o Pathname) MoveTo(path Pathname) error {
	return os.Rename(o.filePath, path.filePath)
}

func (o Pathname) CopyTo(toPath Pathname) error {
	if toPath.IsDirectory() {
		return o.CopyDirTo(toPath)
	} else {
		return o.CopyFileTo(toPath)
	}
}

func (o Pathname) CopyFileTo(toPath Pathname) error {
	fileInfo, err := os.Stat(o.filePath)
	if err != nil {
		return err
	}
	if !fileInfo.Mode().IsRegular() {
		if fileInfo.IsDir() {
			return fmt.Errorf("source is dir (not a file)")
		} else {
			return fmt.Errorf("source not a regular file")
		}
	}
	if toPath.Exist() {
		return fmt.Errorf("target already exist")
	}

	source, err := os.Open(o.filePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.OpenFile(toPath.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileInfo.Mode().Perm())
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func (o Pathname) SymlinkTo(toPath Pathname) error {
	return os.Symlink(o.filePath, toPath.filePath)
}

func (o Pathname) Perm() uint32 {
	fileInfo, err := os.Stat(o.filePath)
	if err != nil {
		return 0
	}
	return uint32(fileInfo.Mode().Perm())
}

func (o Pathname) Chmod(perm uint32) error {
	return os.Chmod(o.filePath, os.FileMode(perm))
}

func (o Pathname) Chown(uid int) error {
	return os.Chown(o.filePath, uid, -1)
}

// -- change working path

func Current() Pathname {
	var path, _ = os.Getwd()
	return New(path)
}

func TempDir() Pathname {
	return Raw(os.TempDir())
}

func HomeDir() Pathname {
	var path, _ = os.UserHomeDir()
	return Raw(path)
}

func (o Pathname) Chdir() error {
	return os.Chdir(o.filePath)
}

func (o Pathname) ChdirWithin(block func(path Pathname)) error {
	var oldDir = Current()
	var error = os.Chdir(o.filePath)
	if error != nil {
		return error
	}
	block(o)
	return os.Chdir(oldDir.filePath)
}

// --- others

func ErrIgnore[T any](a T, err error) T {
	return a
}

func ErrDie[T any](a T, err error) T {
	if err != nil {
		panic(err)
	}
	return a
}
