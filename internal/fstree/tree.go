package fstree

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/rotisserie/eris"
)

type Link struct {
	text string
	href string
}

func NewLink(text, href string) *Link {
	return &Link{
		text: text,
		href: href,
	}
}

func (l *Link) Text() string {
	return l.text
}

func (l *Link) Href() string {
	return l.href
}

type File struct {
	parent *Dir
	name   string
	path   string
}

func (f *File) Name() string {
	return f.name
}

func (f *File) NameWithoutExtension() string {
	split := strings.Split(f.name, ".")
	split = split[:len(split)-1]

	return strings.Join(split, "")
}

func (f *File) Parent() *Dir {
	return f.parent
}

func (f *File) Link() *Link {
	return &Link{
		text: f.Name(),
		href: f.Route(),
	}
}

func (f *File) Links() []*Link {
	links := []*Link{
		f.Link(),
	}

	if f.parent == nil {
		return links
	}

	dir := f.parent

	for dir.parent != nil {
		links = append(links, dir.Link())
	}

	for i, j := 0, len(links)-1; i < j; i, j = i+1, j-1 {
		links[i], links[j] = links[j], links[i]
	}

	return links
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Route() string {
	split := strings.Split(f.path[4:], ".")
	split = split[:len(split)-1]

	return strings.Join(split, "")
}

type Dir struct {
	parent *Dir
	pwd    string
	files  []*File
	dirs   []*Dir
}

func NewFromEntries(entities []fs.DirEntry, pwd string) (*Dir, error) {
	root := &Dir{
		parent: nil,
		pwd:    pwd,
	}

	err := fillTree(root, entities)
	if err != nil {
		return nil, err
	}

	return root, nil
}

func (n *Dir) Parent() *Dir {
	return n.parent
}

func (n *Dir) HasParent() bool {
	return n.parent != nil
}

func fillTree(root *Dir, entries []fs.DirEntry) error {
	if len(entries) == 0 {
		return nil
	}

	entry := entries[0]

	if entry.IsDir() {
		node := &Dir{
			parent: root,
			pwd:    root.pwd + "/" + entry.Name(),
		}

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
		parent: root,
		name:   entry.Name(),
		path:   root.pwd + "/" + entry.Name(),
	}

	root.AddFile(file)

	// Next
	err := fillTree(root, entries[1:])
	if err != nil {
		return err
	}

	return nil
}

func (n *Dir) PWD() string {
	return n.pwd
}

func (n *Dir) Name() string {
	split := strings.Split(n.pwd, "/")

	return split[len(split)-1]
}

func (n *Dir) Route() string {
	split := strings.Split(n.pwd, "/")[1:]
	route := strings.Join(split, "/")

	return route
}

func (n *Dir) Link() *Link {
	return &Link{
		text: n.Name(),
		href: n.Route(),
	}
}

func (n *Dir) Links() []*Link {
	var links []*Link

	dir := n

	for dir.parent != nil {
		links = append(links, dir.Link())

		dir = dir.parent
	}

	for i, j := 0, len(links)-1; i < j; i, j = i+1, j-1 {
		links[i], links[j] = links[j], links[i]
	}

	return links
}

func (n *Dir) Files() []*File {
	return n.files
}

func (n *Dir) Dirs() []*Dir {
	return n.dirs
}

func (n *Dir) AddDir(dir *Dir) {
	n.dirs = append(n.dirs, dir)
}

func (n *Dir) AddFile(file *File) {
	n.files = append(n.files, file)
}

func (n *Dir) IsEmpty() bool {
	return len(n.dirs) == 0 && len(n.files) == 0
}

func (n *Dir) String() string {
	return fmt.Sprintf("PWD: %s\nFiles: %v\nDirs: %v\n", n.pwd, n.files, n.dirs)
}

func (n *Dir) getFiles() []*File {
	var files []*File

	for _, dir := range n.dirs {
		files = append(files, dir.getFiles()...)
	}

	return files
}

func (n *Dir) AllFiles() []*File {
	files := n.Files()

	for _, dir := range n.Dirs() {
		files = append(files, dir.getFiles()...)
	}

	return files
}

// func (n *Node) AllRoutes() []string {
// 	return zen.Map(n.AllFiles(), func(file *File) string {
// 		return file.Route()
// 	})
// }

func (n *Dir) ExecuteOnAllFiles(fn func(file *File) error) error {
	if n.IsEmpty() {
		return nil
	}

	for _, file := range n.files {
		err := fn(file)
		if err != nil {
			return err
		}
	}

	for _, dir := range n.dirs {
		err := dir.ExecuteOnAllFiles(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Dir) ExecuteOnAllDirs(fn func(dir *Dir) error) error {
	if n.IsEmpty() {
		return nil
	}

	for _, dir := range n.dirs {
		dir.ExecuteOnAllDirs(fn)
		err := fn(dir)
		if err != nil {
			return err
		}
	}

	return nil
}
