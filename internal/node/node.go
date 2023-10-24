package node

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/kyoto-framework/zen/v2"
	"github.com/rotisserie/eris"
)

type File struct {
	name string
	path string
}

func (f *File) Name() string {
	return f.name
}

func (f *File) NameWithoutExtension() string {
	split := strings.Split(f.name, ".")
	split = split[:len(split)-1]

	return strings.Join(split, "")
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Route() string {
	split := strings.Split(f.path[4:], ".")
	split = split[:len(split)-1]

	return strings.Join(split, "")
}

type Node struct {
	pwd   string
	files []*File
	dirs  []*Node
}

type Options struct {
	Pwd   string
	Files []*File
	Dirs  []*Node
}

func New(opt Options) *Node {
	return &Node{
		pwd:   opt.Pwd,
		files: opt.Files,
		dirs:  opt.Dirs,
	}
}

func NewFromEntries(entities []fs.DirEntry, pwd string) (*Node, error) {
	root := New(Options{
		Pwd:   pwd,
		Files: []*File{},
		Dirs:  []*Node{},
	})

	err := fillTree(root, entities)
	if err != nil {
		return nil, err
	}

	return root, nil
}

func fillTree(root *Node, entries []fs.DirEntry) error {
	if len(entries) == 0 {
		return nil
	}

	entry := entries[0]

	if entry.IsDir() {
		node := New(Options{
			Pwd:   root.pwd + "/" + entry.Name(),
			Files: []*File{},
			Dirs:  []*Node{},
		})

		nodeEntries, err := os.ReadDir(root.pwd + "/" + entry.Name())
		if err != nil {
			return eris.Wrapf(err, "failed to read dir %s", root.pwd+"/"+entry.Name())
		}

		err = fillTree(node, nodeEntries)
		if err != nil {
			return err
		}

		root.AddDir(node)

		// Next
		err = fillTree(root, entries[1:])
		if err != nil {
			return err
		}

		return nil
	}

	file := &File{
		name: entry.Name(),
		path: root.pwd + "/" + entry.Name(),
	}

	root.AddFile(file)

	// Next
	err := fillTree(root, entries[1:])
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) PWD() string {
	return n.pwd
}

func (n *Node) Files() []*File {
	return n.files
}

func (n *Node) Dirs() []*Node {
	return n.dirs
}

func (n *Node) AddDir(dir *Node) {
	n.dirs = append(n.dirs, dir)
}

func (n *Node) AddFile(file *File) {
	n.files = append(n.files, file)
}

func (n *Node) IsEmpty() bool {
	return len(n.dirs) == 0 && len(n.files) == 0
}

func (n *Node) String() string {
	return fmt.Sprintf("PWD: %s\nFiles: %v\nDirs: %v\n", n.pwd, n.files, n.dirs)
}

func (n *Node) getFiles() []*File {
	var files []*File

	for _, dir := range n.dirs {
		files = append(files, dir.getFiles()...)
	}

	return files
}

func (n *Node) AllFiles() []*File {
	files := n.Files()

	for _, dir := range n.Dirs() {
		files = append(files, dir.getFiles()...)
	}

	return files
}

func (n *Node) AllRoutes() []string {
	return zen.Map(n.AllFiles(), func(file *File) string {
		return file.Route()
	})
}

func (n *Node) Execute(fn func(file *File) error) error {
	for _, file := range n.Files() {
		err := fn(file)
		if err != nil {
			return err
		}
	}

	return nil
}
