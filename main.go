package main

import (
	"fmt"
	"os"

	"github.com/ilfey/webmd/node"
	"github.com/kyoto-framework/kyoto/v2"
	"github.com/kyoto-framework/zen/v2"

	"github.com/gomarkdown/markdown"
	// "github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func main() {
	dirEntry, err := os.ReadDir("md")
	if err != nil {
		// TODO: handle error
		panic(err)
	}

	root, err := node.NewFromEntries(dirEntry, "md")
	if err != nil {
		panic(err)
	}

	kyoto.TemplateConf.FuncMap = kyoto.ComposeFuncMap(
		kyoto.FuncMap, zen.FuncMap,
	)

	// Setup pages

	root.Execute(func(path string) error {
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		html := mdToHTML(b)

		fmt.Printf("path: %v\n", path)

		kyoto.HandlePage("/" + path, func(ctx *kyoto.Context) struct {
			HTML string
		} {
			kyoto.Template(ctx, "base.html")

			return struct {
				HTML string
			}{
				HTML: string(html),
			}
		})

		return nil
	})

	kyoto.Serve(":8080")
}
