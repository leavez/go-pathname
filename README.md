# go-pathname

Managing paths conveniently. Inspired by ruby's `Pathname`. 


## TLDR

```go
path := pathname.New(p)

if path.Exist() {
	// do something
}

path.Append("abc.txt").AbsolutePath()
path.Children()

// ...
```

### Methods

```
New(pathString string) Pathname
Raw(pathString string) Pathname
(Pathname) String() string
(Pathname) Append(path string) Pathname
(Pathname) Parent() Pathname
(Pathname) Basename() string
(Pathname) Extension() string
(Pathname) IsAbsolute() bool
(Pathname) AbsolutePath() Pathname
(Pathname) Children() ([]Pathname, error)
(Pathname) Walk(block func(path Pathname, entry fs.DirEntry, err error) error) error
(Pathname) ChildrenNames() ([]string, error)
(Pathname) Mkpath() error
(Pathname) Rmtree() error
(Pathname) Glob(pattern string) []Pathname
(Pathname) AppendRandom() Pathname
(Pathname) Exist() bool
(Pathname) IsDirectory() bool
(Pathname) IsFile() bool
(Pathname) IsRegularFile() bool
(Pathname) IsSymlink() bool
(Pathname) Write(text string) error
(Pathname) Read() (string, error)
(Pathname) MoveTo(path Pathname) error
(Pathname) CopyTo(toPath Pathname) error
(Pathname) CopyFileTo(toPath Pathname) error
(Pathname) SymlinkTo(toPath Pathname) error
(Pathname) Perm() uint32
(Pathname) Chmod(perm uint32) error
(Pathname) Chown(uid int) error
Current() Pathname
TempDir() Pathname
HomeDir() Pathname
(Pathname) Chdir() error
(Pathname) ChdirWithin(block func(path Pathname)) error
ErrIgnore[T any](a T, err error) T
ErrDie[T any](a T, err error) T
```


## License

MIT
