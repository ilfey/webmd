package main

import (
	"io/fs"
	"strings"

	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/ilfey/webmd/tree"

	mark "github.com/ilfey/webmd/markdown"
	"github.com/sirupsen/logrus"
)

func DirPage(dir *tree.Dir, logger *logrus.Logger) gin.HandlerFunc {
	logger.Infof("register dir: %s on %s", dir.Name(), dir.Path())

	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "dir.html", gin.H{
			"Title":   dir.Name(),
			"Links":   parseLinks(dir.Node),
			"Dirs":    dir.Dirs(),
			"IsEmpty": len(dir.Files()) == 0 && len(dir.Dirs()) == 0,
			"Files":   dir.Files(),
		})
	}
}

func FilePage(file *tree.File, logger *logrus.Logger) gin.HandlerFunc {

	// Read file
	b, err := os.ReadFile(file.AbsolutePath())
	if err != nil {
		logger.Panic(eris.Wrapf(err, "failed to read file %s", file.AbsolutePath()))
	}

	// Create parser
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.OrderedListStart | parser.SuperSubscript | parser.NoEmptyLineBeforeBlock | parser.EmptyLinesBreakList
	mdParser := parser.NewWithExtensions(extensions)

	// Create renderer
	opts := html.RendererOptions{
		Flags:          html.CommonFlags | html.HrefTargetBlank | html.SmartypantsAngledQuotes | html.SmartypantsQuotesNBSP,
		RenderNodeHook: mark.RenderHook,
	}
	mdRenderer := html.NewRenderer(opts)

	// Parse
	doc := mdParser.Parse(b)

	// Render
	html := template.HTML(markdown.Render(doc, mdRenderer))

	// Get nav
	nav := mark.GetNav(doc)

	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "file.html", gin.H{
			"Title": file.Name(),
			"Links": parseLinks(file.Node),
			"Html":  html,
			"Nav":   nav,
		})
	}
}

func parseLinks(root *tree.Node) []*tree.Link {
	links := []*tree.Link{
		{
			Text: root.Name(),
			Href: root.Path(),
		},
	}

	parent := root.Parent()

	for parent != nil {
		links = append(links, &tree.Link{
			Text: parent.Name(),
			Href: parent.Path(),
		})

		parent = parent.Parent()
	}

	for i, j := 0, len(links)-1; i < j; i, j = i+1, j-1 {
		links[i], links[j] = links[j], links[i]
	}

	return links
}

func bindDirs(engine *gin.Engine, root *tree.Dir, logger *logrus.Logger) {
	for _, d := range root.Dirs() {
		bindDirs(engine, d, logger)
	}

	engine.Handle(http.MethodGet, root.Path(), DirPage(root, logger))
}

func bindFiles(engine *gin.Engine, root *tree.Dir, logger *logrus.Logger) {
	for _, d := range root.Dirs() {
		bindFiles(engine, d, logger)
	}

	for _, f := range root.Files() {
		engine.Handle(http.MethodGet, f.Path(), FilePage(f, logger))
	}
}

// setupPages registers
func setupPages(engine *gin.Engine, root *tree.Dir, logger *logrus.Logger) {
	bindDirs(engine, root, logger)
	bindFiles(engine, root, logger)
}

func addNextDirs(root *tree.Dir, entries []fs.DirEntry) error {
	for _, de := range entries {
		name := de.Name()
		absolutePath := root.AbsolutePath() + "/" + name

		var nodename string

		if strings.HasSuffix(root.Path(), "/") {
			nodename = root.Path() + name
		} else {
			nodename = root.Path() + "/" + name
		}

		if de.IsDir() {

			next := tree.NewDir(root.Node, name, nodename, absolutePath)

			nextPath := root.AbsolutePath() + "/" + name
			nextEntries, err := os.ReadDir(nextPath)
			if err != nil {
				return eris.Wrapf(err, "failed to read dir %s", nextPath)
			}

			addNextDirs(next, nextEntries)

			root.AddDir(next)

			continue
		}

		file := tree.NewFile(root.Node, name, nodename, absolutePath)

		root.AddFile(file)
	}

	return nil
}

func GetFS(root *tree.Dir, path string) (*tree.Dir, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to read dir %s", path)
	}

	if root == nil {
		root = tree.NewDir(nil, "root", "/", path)
	}

	err = addNextDirs(root, dirEntries)
	if err != nil {
		return nil, err
	}

	return root, nil
}

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01/02 15:04:05",
	}

	root, err := GetFS(nil, ".md")
	if err != nil {
		logger.Panic(err)
	}

	engine := gin.New()

	engine.LoadHTMLGlob("templates/*.html")
	engine.StaticFS("/dist", http.Dir(".dist"))

	engine.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	setupPages(engine, root, logger)

	http.ListenAndServe(":8080", engine)
}
