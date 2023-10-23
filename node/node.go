package node

import (
	"fmt"
	"io/fs"
	"os"
)

type Node struct {
	pwd   string
	files []string
	dirs  []*Node
}

type Options struct {
	Pwd   string
	Files []string
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
		Files: []string{},
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
			Pwd:   root.PWD() + "/" + entry.Name(),
			Files: []string{},
			Dirs:  []*Node{},
		})

		nodeEntries, err := os.ReadDir(root.PWD() + "/" + entry.Name())
		if err != nil {
			return err
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

	root.AddFile(entry.Name())

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

func (n *Node) Files() []string {
	return n.files
}

func (n *Node) Dirs() []*Node {
	return n.dirs
}

func (n *Node) AddDir(dir *Node) {
	n.dirs = append(n.dirs, dir)
}

func (n *Node) AddFile(file string) {
	n.files = append(n.files, n.PWD()+"/"+file)
}

func (n *Node) IsEmpty() bool {
	return len(n.dirs) == 0 && len(n.files) == 0
}

func (n *Node) String() string {
	return fmt.Sprintf("PWD: %s\nFiles: %v\nDirs: %v\n", n.pwd, n.files, n.dirs)
}

func (n *Node) Execute(fn func(path string) error) error {
	for _, file := range n.Files() {
		err := fn(file)
		if err != nil {
			return err
		}
	}

	return nil
}
