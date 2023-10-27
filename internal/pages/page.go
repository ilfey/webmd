package pages

import (
	"html/template"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/ilfey/webmd/internal/fstree"
	mark "github.com/ilfey/webmd/internal/markdown"
	"github.com/kyoto-framework/kyoto/v2"
	"github.com/rotisserie/eris"
)

type PPageState struct {
	*fstree.File
	Html template.HTML
	Nav  []*fstree.Link
}

func PPage(file *fstree.File) (kyoto.Component[*PPageState], error) {

	// Read file
	b, err := os.ReadFile(file.Path())
	if err != nil {
		return nil, eris.Wrapf(err, "failed to read file %s", file.Path())
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
	html := markdown.Render(doc, mdRenderer)

	// Get nav
	nav := mark.GetNav(doc)

	return func(ctx *kyoto.Context) *PPageState {
		kyoto.Template(ctx, "page.html")

		return &PPageState{
			File: file,
			Html: template.HTML(html),
			Nav:  nav,
		}
	}, nil
}
