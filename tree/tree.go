package tree

// Link
type Link struct {
	Text string
	Href string
}

// Base node
type Node struct {
	name         string
	path         string
	absolutePath string
	parent       *Node
}

// Returns node name
func (n *Node) Name() string {
	return n.name
}

// Returns node location
func (n *Node) Path() string {
	return n.path
}

// Returns node absolute path
func (n *Node) AbsolutePath() string {
	return n.absolutePath
}

// Returns parent node
func (n *Node) Parent() *Node {
	return n.parent
}

type File struct {
	*Node
	parent *Dir
}

// File constructor
func NewFile(parent *Node, name, path, absolutePath string) *File {
	return &File{
		Node: &Node{
			name:         name,
			path:         path,
			absolutePath: absolutePath,
			parent:       parent,
		},
	}
}

// Returns parent dir
func (f *File) Parent() *Dir {
	return f.parent
}

type Dir struct {
	*Node
	parent *Dir
	dirs   []*Dir
	files  []*File
}

// Dir constructor
func NewDir(parent *Node, name, path, absolutePath string) *Dir {
	return &Dir{
		Node: &Node{
			name:         name,
			path:         path,
			absolutePath: absolutePath,
			parent:       parent,
		},
	}
}

// Returns parent dir if dir has parent else nil
func (d *Dir) Parent() *Dir {
	return d.parent
}

// Returns true if dir has parent
func (d *Dir) HasParent() bool {
	return d.parent != nil
}

// Returns all dirs in current dir
func (d *Dir) Dirs() []*Dir {
	return d.dirs
}

// Returns all files in current dir
func (d *Dir) Files() []*File {
	return d.files
}

// Adds dir
func (d *Dir) AddDir(dir *Dir) {
	d.dirs = append(d.dirs, dir)
}

// Adds file
func (d *Dir) AddFile(file *File) {
	d.files = append(d.files, file)
}
