package pathname

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	assert.Equal(t, New("./"), Pathname{filePath: "."})
	assert.Equal(t, New("demo/b/c/.").String(), "demo/b/c")
	assert.Equal(t, New("demo/b/c/./").String(), "demo/b/c")
	assert.Equal(t, Raw("demo/b/c/./").String(), "demo/b/c/./")
}

func TestPathname_Append(t *testing.T) {
	var p = New("/")
	var p1 = p.Append("bbb")
	assert.Equal(t, p1, Pathname{filePath: "/bbb"})
	var p2 = p.Append("/bbb/ccc")
	assert.Equal(t, p2, Pathname{filePath: "/bbb/ccc"})
}

func TestPathname_AbsolutePath(t *testing.T) {
	var p = New("demo/b/c")
	var abs = p.AbsolutePath()
	assert.Equal(t, Current().Append(p.filePath), abs)

	var p2 = New("/demo/b/c")
	assert.Equal(t, p2.AbsolutePath(), p2)
}

func TestPathname_Basename(t *testing.T) {
	assert.Equal(t, New("/abc/efg").Basename(), "efg")
	assert.Equal(t, New("/abc/efg/demo.123").Basename(), "demo.123")
	assert.Equal(t, New(".").Basename(), ".")
}

func TestPathname_Chdir(t *testing.T) {
	var oldPwd, _ = os.Getwd()
	_ = New("/").Chdir()
	var current, _ = os.Getwd()
	assert.Equal(t, current, "/")

	_ = New(oldPwd).Chdir()
	var current2, _ = os.Getwd()
	assert.Equal(t, current2, oldPwd)

	var err = New("/abcdefg").Chdir()
	assert.Error(t, err, os.PathError{})
}

func TestPathname_ChdirWithin(t *testing.T) {
	var oldPwd, _ = os.Getwd()
	var err = New("/").ChdirWithin(func(path Pathname) {
		var pwd, _ = os.Getwd()
		assert.Equal(t, pwd, "/")
	})
	assert.NoError(t, err)
	var newPwd, _ = os.Getwd()
	assert.Equal(t, oldPwd, newPwd)

	var err2 = New("/123/gjfo/adfle").ChdirWithin(func(path Pathname) {
	})
	assert.Error(t, err2)
}

func TestPathname_Children(t *testing.T) {
	var paths, _ = New("/").Children()
	var pathStrings = make([]string, 0)
	for _, p := range paths {
		pathStrings = append(pathStrings, p.Basename())
	}
	assert.Contains(t, pathStrings, "usr")

	var _, err = New("/asdfljlsajf").Children()
	assert.Error(t, err)
}

func TestCurrent(t *testing.T) {
	var pwd, _ = os.Getwd()
	assert.Equal(t, Current().filePath, pwd)
}

func TestPathname_Exist(t *testing.T) {
	var bashPath = New("/bin/bash")
	assert.True(t, bashPath.Exist())
	assert.False(t, New("/bin/111").Exist())
	assert.False(t, New("/bin/bash/111").Exist())
	assert.False(t, New("./bin/bash/111").Exist())
}

func TestPathname_Extension(t *testing.T) {
	assert.Equal(t, New("demo.123").Extension(), ".123")
	assert.Equal(t, New("/demo").Extension(), "")
}

func TestPathname_Glob(t *testing.T) {
	var base = New("tmp" + fmt.Sprintf("%d", rand.Int()))
	base.Mkpath()
	base.ChdirWithin(func(_ Pathname) {
		var base = New(".")
		var f1 = base.Append("1.txt")
		var f2 = base.Append("sub/2")
		var f3 = base.Append("sub2/sub22/3.txt")

		for _, f := range []Pathname{f1, f2, f3} {
			f.Parent().Mkpath()
			f.Write("")
		}

		assert.Equal(t, ToStrings(base.Glob("*")), []string{"1.txt", "sub", "sub2"})
		assert.Equal(t, ToStrings(base.Glob("*.txt")), []string{"1.txt"})
		assert.Equal(t, ToStrings(base.Glob("*/*")), []string{"sub/2", "sub2/sub22"})
		assert.Equal(t, ToStrings(base.Glob("**/**/*.txt")), []string{"sub2/sub22/3.txt"})
		assert.Equal(t, ToStrings(base.Glob("*/*/*.txt")), []string{"sub2/sub22/3.txt"})
	})
	defer func() {
		base.Rmtree()
	}()
}

func TestPathname_IsAbsolute(t *testing.T) {
	assert.Equal(t, New("/abc").IsAbsolute(), true)
	assert.Equal(t, New("../abc").IsAbsolute(), false)
	assert.Equal(t, New(".123").IsAbsolute(), false)
}

func TestPathname_IsDirectory(t *testing.T) {
	assert.Equal(t, New("/").IsDirectory(), true)
	assert.Equal(t, New("/bin/bash").IsDirectory(), false)
	assert.Equal(t, New("/1/2/3/4/5/6/7/8").IsDirectory(), false)
}

func TestPathname_IsFile(t *testing.T) {
	assert.Equal(t, New("/").IsFile(), false)
	assert.Equal(t, New("/bin/bash").IsFile(), true)
	assert.Equal(t, New("/1/2/3/4/5/6/7/8").IsFile(), false)
}

func TestPathname_Mkpath(t *testing.T) {
	var base = Current().AppendRandom()
	base.Mkpath()
	defer base.Rmtree()
	assert.Equal(t, base.IsDirectory(), true)

	var p2 = base.Append("1").Append("2").Append("3")
	var err = p2.Mkpath()
	defer p2.Rmtree()
	assert.True(t, p2.IsDirectory())
	assert.NoError(t, err)

	// if the target is exist, it doesn't throw error
	var p3 = base.Append("1").Append("2").Append("3")
	var err3 = p3.Mkpath()
	assert.True(t, p3.IsDirectory())
	assert.NoError(t, err3)
}

func TestPathname_Rmtree(t *testing.T) {
	var base = Current().AppendRandom()
	base.Mkpath()
	assert.Equal(t, base.Exist(), true)
	base.Rmtree()
	assert.Equal(t, base.Exist(), false)

	base.Mkpath()
	var p2 = base.Append("1").Append("2").Append("3")
	var p3 = base.Append("1").Append("file")
	p2.Mkpath()
	p3.Write("")
	base.Rmtree()
	assert.Equal(t, base.Exist(), false)
	assert.Equal(t, p3.Exist(), false)
	assert.Equal(t, p2.Exist(), false)
}

func TestPathname_Parent(t *testing.T) {
	assert.Equal(t, New("demo/b/c").Parent().String(), "demo/b")
	assert.Equal(t, New(".").Parent().String(), ".")
	assert.Equal(t, New("/").Parent().String(), "/")
}

func TestPathname_RandomChildren(t *testing.T) {
	var base = New("demo/b")
	var r = base.AppendRandom()
	assert.Equal(t, r.Parent(), base)
}

func TestPathname_Read(t *testing.T) {
	for _, text := range []string{"123", "你好", "123🥳123"} {
		var bytes = []byte(text)
		var path = Current().AppendRandom()
		ioutil.WriteFile(path.String(), bytes, 0744)
		var readText, err = path.Read()
		assert.Equal(t, readText, text)
		assert.NoError(t, err)
		path.Rmtree()
	}
}

func TestPathname_Write(t *testing.T) {
	for _, text := range []string{"123", "你好", "123🥳123"} {
		var path = Current().AppendRandom()
		var err = path.Write(text)
		var readText, _ = path.Read()
		assert.Equal(t, readText, text)
		assert.NoError(t, err)
		path.Rmtree()
	}

	func() {
		var path = Current().AppendRandom()
		var err = path.Write("1")
		assert.NoError(t, err)
		var err2 = path.Write("2")
		assert.NoError(t, err2)
		readValue, _ := path.Read()
		assert.Equal(t, readValue, "2")
		path.Rmtree()
	}()
	func() {
		var path = TempDir().AppendRandom()
		_ = path.Mkpath()
		var err = path.Write("123 by ut")
		assert.Error(t, err)
		_ = path.Rmtree()
	}()
}

func TestTempDir(t *testing.T) {
	if runtime.GOOS == "darwin" {
		if os.Getuid() == 0 {
			assert.Equal(t, TempDir().String(), "/tmp")
		} else {
			assert.True(t, strings.HasPrefix(TempDir().String(), "/var/folders/"), TempDir().String())
		}
	}
}

func TestHomeDir(t *testing.T) {
	if runtime.GOOS == "darwin" {
		var h = os.Getenv("HOME")
		assert.Equal(t, HomeDir().String(), h)

		filepath.WalkDir(h, func(path string, info fs.DirEntry, err error) error {
			return nil
		})
	}
}

func TestPathname_CopyFileTo(t *testing.T) {
	var path = TempDir().AppendRandom()
	path.Write("123")
	defer path.Rmtree()
	var toPath = TempDir().AppendRandom()
	var copyErr = path.CopyFileTo(toPath)
	defer toPath.Rmtree()
	assert.NoError(t, copyErr)
	var content, err = toPath.Read()
	assert.Equal(t, content, "123")
	assert.NoError(t, err)

	// test error
	func() {
		var path = Current().AppendRandom()
		var err = path.CopyFileTo(New("./"))
		assert.Error(t, err)
	}()
	func() {
		// not a file
		var err = Current().CopyFileTo(TempDir().AppendRandom())
		assert.Error(t, err)
	}()
}

func TestPathname_CopyDirTo(t *testing.T) {
	var path = Current().AppendRandom()
	_ = path.Mkpath()
	defer path.Rmtree()

	file := path.AppendRandom()
	file.Write("1234")
	var toPath = Current().AppendRandom()
	defer toPath.Rmtree()

	var copyErr = path.CopyDirTo(toPath)
	assert.NoError(t, copyErr)

	assert.True(t, toPath.IsDirectory())
	assert.True(t, toPath.Append(file.Basename()).Exist())
}

func TestPathname_MoveTo(t *testing.T) {
	var path = Current().AppendRandom()
	path.Write("123")
	defer func() {
		path.Rmtree()
		if path.Exist() {
		}
	}()
	var toPath = Current().AppendRandom()
	var moveErr = path.MoveTo(toPath)
	defer func() {
		if toPath.Exist() {
			toPath.Rmtree()
		}
	}()
	assert.NoError(t, moveErr)
	var content, err = toPath.Read()
	assert.Equal(t, content, "123")
	assert.NoError(t, err)

	// test error
	func() {
		var path = Current().AppendRandom()
		var err = path.MoveTo(New("./"))
		assert.Error(t, err)
	}()
}

func TestPathname_String(t *testing.T) {
	assert.Equal(t, New("demo/b/c").String(), "demo/b/c")
}

func ToStrings(paths []Pathname) []string {
	var strings = make([]string, len(paths))
	for i, p := range paths {
		strings[i] = p.String()
	}
	return strings
}
