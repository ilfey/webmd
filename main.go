package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/gorilla/mux"
	mark "github.com/ilfey/webmd/internal/markdown"
	"github.com/ilfey/webmd/internal/node"

	"github.com/kyoto-framework/kyoto/v2"
	"github.com/kyoto-framework/zen/v2"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: mark.Hook,
	}

	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

type BaseState struct {
	Html   template.HTML
	Routes []*node.File
}

func setupPages(router *mux.Router, root *node.Node) {
	root.Execute(func(file *node.File) error {
		b, err := os.ReadFile(file.Path())
		if err != nil {
			return err
		}

		html := mdToHTML(b)

		route := "/" + file.Route()
		fmt.Printf("\n\n\n\nhtml: %v\n\n\n\n", string(html))

		fmt.Printf("register route: %v\n", route)

		router.HandleFunc(route, kyoto.HandlerPage(func(ctx *kyoto.Context) BaseState {
			// Setup rendering
			kyoto.Template(ctx, "base.html")

			return BaseState{
				Html:   template.HTML(html),
				Routes: root.AllFiles(),
			}
		}))

		return nil
	})
}

// setupAssets registers a static files handler.
func setupAssets(mux *mux.Router) {
	mux.PathPrefix("/dist/").Handler(
		http.StripPrefix("/dist/", http.FileServer(http.Dir(".dist"))),
	)
}

func main() {

	dirEntry, err := os.ReadDir(".md")
	if err != nil {
		panic(err)
	}

	root, err := node.NewFromEntries(dirEntry, ".md")
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	kyoto.TemplateConf.ParseGlob = "templates/*.html"

	// Setup kyoto
	kyoto.TemplateConf.FuncMap = kyoto.ComposeFuncMap(
		kyoto.FuncMap, zen.FuncMap,
	)

	setupAssets(router)
	setupPages(router, root)

	http.ListenAndServe(":8080", router)
}
